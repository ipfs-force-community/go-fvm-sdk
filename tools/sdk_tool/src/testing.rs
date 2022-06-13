use anyhow::Result;
use clap::Parser;
use colored::*;
use fvm::executor::ApplyFailure;
use fvm::executor::ApplyRet;
use fvm::executor::{ApplyKind, Executor};
use fvm_integration_tests::tester::Tester;
use fvm_ipld_blockstore::MemoryBlockstore;
use fvm_ipld_encoding::tuple::*;
use fvm_shared::address::Address;
use fvm_shared::bigint::BigInt;
use fvm_shared::econ::TokenAmount;
use fvm_shared::message::Message;
use fvm_shared::state::StateTreeVersion;
use fvm_shared::version::NetworkVersion;
use hex::FromHex;
use libsecp256k1::SecretKey;
use serde::{Deserialize, Serialize};
use std::fs;
use std::iter::Iterator;
use std::path::PathBuf;

#[derive(Parser, Debug)]
#[clap(author, version, about, long_about = None)]
pub struct TestConfig {
    #[clap(short, long)]
    path: String,
}

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct InitAccount {
    priv_key: String,
    balance: u64,
}

#[derive(Serialize, Deserialize, Clone, Debug)]
struct TestJson {
    accounts: Vec<InitAccount>,
    cases: Vec<TestCase>,
}

#[derive(Serialize, Deserialize, Clone, Debug)]
struct TestCase {
    name: String,
    binary: String,
    method_num: u64,
    actor_balance: u64,
    send_value: u64,
    params: String,

    expect_code: u32,
    expect_message: String,
}

pub fn run_testing(cfg: &TestConfig) {
    let case_meta_path: PathBuf = [cfg.path.clone(), "test.json".to_string()].iter().collect();
    let buf = fs::read(case_meta_path).unwrap();
    let test_json: TestJson = serde_json::from_slice(&buf).unwrap();

    test_json.cases.iter().for_each(|test_case| {
        let path: PathBuf = [cfg.path.clone(), test_case.binary.to_owned()]
            .iter()
            .collect();
        let buf = fs::read(path.clone())
            .unwrap_or_else(|_| panic!("path {} not found", path.to_str().unwrap()));
        let ret = exec(
            &buf,
            &test_json.accounts,
            test_case.method_num,
            test_case.actor_balance,
            test_case.send_value,
        )
        .unwrap();
        if ret.msg_receipt.exit_code.value() != test_case.expect_code {
            if let Some(fail_info) = ret.failure_info {
                panic!(
                    "case {} expect exit code {} but got {} {}",
                    test_case.name, test_case.expect_code, ret.msg_receipt.exit_code, fail_info
                )
            } else {
                panic!(
                    "case {} expect exit code {} but got {}",
                    test_case.name, test_case.expect_code, ret.msg_receipt.exit_code
                )
            }
        }

        if test_case.expect_code != 0 {
            if let ApplyFailure::MessageBacktrace(mut trace) = ret.failure_info.unwrap() {
                let abort_msg = trace.frames.pop().unwrap().message;
                if abort_msg != test_case.expect_message {
                    panic!(
                        "case {} expect messcage {} but got {}",
                        test_case.name, test_case.expect_message, abort_msg
                    )
                }
            }
        }
        println! {"{}: case {}", "passed".green(), test_case.name}
    });
}

#[derive(Serialize_tuple, Deserialize_tuple, Clone, Debug)]
struct State {
    empty: bool,
}

pub fn exec(
    wasm_bin: &[u8],
    init_accounts: &[InitAccount],
    method_num: u64,
    actor_balance: u64,
    send_value: u64,
) -> Result<ApplyRet> {
    // Instantiate tester
    let mut tester = Tester::new(
        NetworkVersion::V15,
        StateTreeVersion::V4,
        MemoryBlockstore::default(),
    )
    .unwrap();
    let mut accounts = vec![];
    for init_account in init_accounts {
        let priv_key = SecretKey::parse(&<[u8; 32]>::from_hex(init_account.priv_key.clone())?)?;
        let account =
            tester.make_secp256k1_account(priv_key, TokenAmount::from(init_account.balance))?;
        accounts.push(account);
    }
    // Get wasm bin
    //  let wasm_bin = wat2wasm(wat).unwrap();

    // Set actor state
    let actor_state = State { empty: true };
    let state_cid = tester.set_state(&actor_state)?;

    // Set actor

    let actor_address = Address::new_id(10000);
    tester
        .set_actor_from_bin(
            wasm_bin,
            state_cid,
            actor_address,
            BigInt::from(actor_balance),
        )
        .unwrap();

    // Instantiate machine
    tester.instantiate_machine()?;

    // Send message
    let message = Message {
        from: accounts[0].1,
        to: actor_address,
        gas_limit: 1000000000000,
        method_num,
        value: BigInt::from(send_value),
        ..Message::default()
    };

    tester
        .executor
        .unwrap()
        .execute_message(message, ApplyKind::Explicit, 100)
}
