use anyhow::{anyhow, Result};
use clap::Parser;
use colored::*;
use fvm::executor::ApplyFailure;
use fvm::executor::ApplyRet;
use fvm::executor::{ApplyKind, Executor};
use fvm_integration_tests::tester::{Account, Tester};
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
use fvm::init_actor::INIT_ACTOR_ADDR;
use libsecp256k1::PublicKeyFormat::Raw;
use fvm_ipld_encoding::tuple::*;
use fvm_ipld_encoding::{Cbor, RawBytes};
use cid::Cid;
use fvm_shared::bigint::Zero;

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
pub struct TestJson {
    accounts: Vec<InitAccount>,
    cases: Vec<WasmCase>,
    contracts : Vec<ContractCase>,
}

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct WasmCase {
    name: String,
    binary: String,
    method_num: u64,
    actor_balance: u64,
    send_value: u64,
    params: String,
    expect_code: u32,
    expect_message: String,
}

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct ContractCase {
    binary: String,
    constructor: String,
    cases: Vec<WasmCase>
}


pub fn run_testing(cfg: &TestConfig) {
    let case_meta_path: PathBuf = [cfg.path.clone(), "test.json".to_string()].iter().collect();
    let buf = fs::read(case_meta_path).unwrap();

    let test_json: TestJson = serde_json::from_slice(&buf).unwrap();

    test_json.cases.iter().for_each(|test_case| { run_signle_wasm(cfg.path.clone(), &test_json.accounts, test_case)});
    test_json.contracts.iter().for_each(|group_case| {
        run_action_group(cfg.path.clone(), &test_json.accounts,  group_case).unwrap()
    })
}

pub fn run_action_group(root_path: String, accounts_cfg: &Vec<InitAccount>, contract_case: &ContractCase) -> Result<()>{
    let path: PathBuf = [root_path, contract_case.binary.to_owned()]
        .iter()
        .collect();
    let buf = fs::read(path.clone())
        .unwrap_or_else(|_| panic!("path {} not found", path.to_str().unwrap()));
    // Instantiate tester
    let (mut tester, accounts) = new_tester(accounts_cfg)?;
    // Instantiate machine
    tester.instantiate_machine()?;

    let mut executor = tester.executor.unwrap();
    //install
    // Send message
    let install_return: InstallReturn = {
        let install_message = Message {
            from: accounts[0].1,
            to: INIT_ACTOR_ADDR,
            gas_limit: 1000000000000,
            method_num:3,
            value: BigInt::zero(),
            params: RawBytes::from(buf),
            ..Message::default()
        };

        let result =  executor.execute_message(install_message, ApplyKind::Explicit, 100)?;
        if result.msg_receipt.exit_code.value() != 0 {
            return Err(anyhow!("failed to install code"));
        }
        InstallReturn::unmarshal_cbor(result.msg_receipt.return_data.as_slice())
    }?;

    //create
    let create_return:ExecReturn = {
        let constructor_params =  hex::decode(contract_case.constructor.clone())?;
        let exec_params = ExecParams{
            code_cid: install_return.code_cid,
            constructor_params: RawBytes::from(constructor_params) ,
        }.marshal_cbor()?;

        let create_message = Message {
            from: accounts[0].1,
            to: INIT_ACTOR_ADDR,
            gas_limit: 1000000000000,
            method_num:2,
            value: BigInt::zero(),
            params:  RawBytes::from(exec_params),
            ..Message::default()
        };
        let result = executor.execute_message(create_message, ApplyKind::Explicit, 100)?;
        if result.msg_receipt.exit_code.value() != 0 {
            return Err(anyhow!("failed to install code"));
        }
        ExecReturn::unmarshal_cbor(result.msg_receipt.return_data.as_slice())
    }?;
    //invoke

    for xx in &contract_case.cases {

    }

    Ok(())
}

pub fn run_signle_wasm(root_path: String, account: &Vec<InitAccount>, wasm_case: &WasmCase) {
    let path: PathBuf = [root_path, wasm_case.binary.to_owned()]
        .iter()
        .collect();
    let buf = fs::read(path.clone())
        .unwrap_or_else(|_| panic!("path {} not found", path.to_str().unwrap()));
    let ret = exec(
        &buf,
        account,
        wasm_case.method_num,
        wasm_case.actor_balance,
        wasm_case.send_value,
    )
        .unwrap();
    if ret.msg_receipt.exit_code.value() != wasm_case.expect_code {
        if let Some(fail_info) = ret.failure_info {
            panic!(
                "case {} expect exit code {} but got {} {}",
                wasm_case.name, wasm_case.expect_code, ret.msg_receipt.exit_code, fail_info
            )
        } else {
            panic!(
                "case {} expect exit code {} but got {}",
                wasm_case.name, wasm_case.expect_code, ret.msg_receipt.exit_code
            )
        }
    }

    if wasm_case.expect_code != 0 {
        if let ApplyFailure::MessageBacktrace(mut trace) = ret.failure_info.unwrap() {
            let abort_msg = trace.frames.pop().unwrap().message;
            if abort_msg != wasm_case.expect_message {
                panic!(
                    "case {} expect messcage {} but got {}",
                    wasm_case.name, wasm_case.expect_message, abort_msg
                )
            }
        }
    }
    println! {"{}: case {}", "passed".green(), wasm_case.name}
}

#[derive(Serialize_tuple, Deserialize_tuple, Clone, Debug)]
struct State {
    empty: bool,
}

pub fn new_tester(accounts_cfg: &[InitAccount]) -> Result<(Tester<MemoryBlockstore>, Vec<Account>)>{
    let mut tester = Tester::new(
        NetworkVersion::V15,
        StateTreeVersion::V4,
        MemoryBlockstore::default(),
    ).unwrap();
    let mut accounts:Vec<Account> = vec![];
    for init_account in accounts_cfg {
        let priv_key = SecretKey::parse(&<[u8; 32]>::from_hex(init_account.priv_key.clone())?)?;
        let account =
            tester.make_secp256k1_account(priv_key, TokenAmount::from(init_account.balance))?;
        accounts.push(account);
    }
    Ok((tester, accounts))
}

pub fn exec(
    wasm_bin: &[u8],
    init_accounts: &[InitAccount],
    method_num: u64,
    actor_balance: u64,
    send_value: u64,
) -> Result<ApplyRet> {
    // Instantiate tester
    let (mut tester, accounts) = new_tester(init_accounts)?;
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

//todo remove ExecParams/ExecReturn/InstallReturn/InstallParams after v8 release and  upgrade actor version  in integration test
/// Init actor Exec Params
#[derive(Serialize_tuple, Deserialize_tuple)]
pub struct ExecParams {
    pub code_cid: Cid,
    pub constructor_params: RawBytes,
}

/// Init actor Exec Return value
#[derive(Serialize_tuple, Deserialize_tuple)]
pub struct ExecReturn {
    /// ID based address for created actor
    pub id_address: Address,
    /// Reorg safe address for actor
    pub robust_address: Address,
}

impl Cbor for ExecReturn {}
impl Cbor for ExecParams {}

/// Init actor Install Params
#[derive(Serialize_tuple, Deserialize_tuple)]
pub struct InstallParams {
    pub code: RawBytes,
}

/// Init actor Install Return value
#[derive(Serialize_tuple, Deserialize_tuple)]
pub struct InstallReturn {
    pub code_cid: Cid,
    pub installed: bool,
}

impl Cbor for InstallParams {}
impl Cbor for InstallReturn {}
