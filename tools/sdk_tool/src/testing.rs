use anyhow::{anyhow, Result};
use cid::Cid;
use clap::Parser;
use colored::*;
use fvm::executor::ApplyFailure;
use fvm::executor::ApplyRet;
use fvm::executor::{ApplyKind, Executor};
use fvm::init_actor::INIT_ACTOR_ID;
use fvm::machine::Machine;
use fvm_integration_tests::bundle;
use fvm_integration_tests::dummy::DummyExterns;
use fvm_integration_tests::tester::{Account, Tester};
use fvm_ipld_blockstore::MemoryBlockstore;
use fvm_ipld_encoding::tuple::*;
use fvm_ipld_encoding::{from_slice, to_vec, RawBytes};
use fvm_shared::address;
use fvm_shared::address::Address;
use fvm_shared::bigint::BigInt;
use fvm_shared::bigint::Zero;
use fvm_shared::econ::TokenAmount;
use fvm_shared::message::Message;
use fvm_shared::state::StateTreeVersion;
use fvm_shared::version::NetworkVersion;
use hex::FromHex;
use libsecp256k1::SecretKey;
use path_absolutize::Absolutize;
use serde::{Deserialize, Serialize};
use std::env::current_dir;
use std::fs;
use std::iter::Iterator;
use std::path::{Path, PathBuf};

#[derive(Parser, Debug)]
pub struct TestConfig {
    //specify test file path
    #[clap(last = true)]
    path: Option<String>,

    //specify which test case to run
    #[clap(short, long)]
    name: Option<String>,
}

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct InitAccount {
    priv_key: Option<String>,
    balance: u64,

    subaddress: Option<String>,
    #[serde(default)]
    isf4: bool,
}

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct TestJson {
    accounts: Vec<InitAccount>,
    #[serde(default)]
    cases: Option<Vec<WasmCase>>,
    #[serde(default)]
    contracts: Option<Vec<ContractCase>>,
}

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct WasmCase {
    name: String,
    method_num: u64,
    #[serde(default)]
    expect_code: u32,
    #[serde(default)]
    expect_message: String,
    #[serde(default)]
    send_from: usize,
    #[serde(default)]
    params: String,
    #[serde(default)]
    send_value: u64,

    //contract case specify
    #[serde(default)]
    return_data: String,

    //signle case specify
    #[serde(default)]
    binary: String,
    #[serde(default)]
    actor_balance: u64,
}

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct ContractCase {
    name: String,
    binary: String,
    #[serde(default)]
    constructor: String,
    #[serde(default)]
    owner_account: usize,
    cases: Vec<WasmCase>,
}

pub fn run_testing(cfg: &TestConfig) -> Result<()> {
    let mut root_dir = if let Some(dir) = &cfg.path {
        Path::new(dir).to_path_buf()
    } else {
        current_dir()?
    };

    if root_dir.extension().is_none() {
        root_dir.push("test.json");
    }
    let buf = fs::read(root_dir).unwrap();

    let test_json: TestJson = serde_json::from_slice(&buf).unwrap();

    if let Some(cases) = test_json.cases {
        cases
            .iter()
            .filter(|v| cfg.name.is_none() || v.name.eq(cfg.name.as_ref().unwrap()))
            .for_each(|test_case| {
                if let Err(e) = run_signle_wasm(&test_json.accounts, test_case) {
                    panic!(
                        "{}:case {} run failed {}",
                        "failed".red(),
                        test_case.name,
                        e
                    )
                }
            });
    }

    if let Some(contracts) = &test_json.contracts {
        contracts
            .iter()
            .filter(|v| cfg.name.is_none() || v.name.eq(cfg.name.as_ref().unwrap()))
            .for_each(|group_case| {
                if let Err(e) = run_action_group(&test_json.accounts, group_case) {
                    panic!(
                        "{}:case {} run failed {}",
                        "failed".red(),
                        group_case.name,
                        e
                    )
                }
            });
    }

    Ok(())
}

pub fn run_action_group(accounts_cfg: &[InitAccount], contract_case: &ContractCase) -> Result<()> {
    let path: PathBuf = [
        current_dir()?,
        Path::new(&contract_case.binary).to_path_buf(),
    ]
    .iter()
    .collect::<PathBuf>()
    .absolutize()
    .map(|v| v.into_owned())
    .expect("get binary path");
    let buf = fs::read(path.clone())
        .unwrap_or_else(|_| panic!("path {} not found", path.to_str().get_or_insert("unknown")));
    // Instantiate tester
    let (mut tester, accounts) = new_tester(accounts_cfg)?;
    // Instantiate machine
    tester.instantiate_machine(DummyExterns)?;

    let mut executor = tester.executor.expect("unable to get executor");

    let init_actor_addr = address::Address::new_id(INIT_ACTOR_ID);
    //install
    // Send message
    let install_return: InstallReturn = {
        let install_params = to_vec(&InstallParams {
            code: RawBytes::from(buf),
        })?;
        let install_message = Message {
            from: accounts[contract_case.owner_account].1,
            to: init_actor_addr,
            gas_limit: 1000000000000,
            method_num: 4,
            value: TokenAmount::zero(),
            params: RawBytes::from(install_params),
            sequence: 0,
            ..Message::default()
        };

        let ret = executor.execute_message(install_message, ApplyKind::Explicit, 100)?;
        check_message_receipt(
            contract_case.name.clone() + "_install_code",
            &ret,
            0,
            "".to_owned(),
            vec![],
        )?;
        from_slice(ret.msg_receipt.return_data.as_slice())
    }?;

    println!("code cid {}", install_return.code_cid.clone());
    //create
    let create_return: ExecReturn = {
        let constructor_params = hex::decode(contract_case.constructor.clone())?;
        let exec_params = to_vec(&ExecParams {
            code_cid: install_return.code_cid,
            constructor_params: RawBytes::from(constructor_params),
        })?;

        let create_message = Message {
            from: accounts[contract_case.owner_account].1,
            to: init_actor_addr,
            gas_limit: 1000000000000,
            method_num: 2,
            value: TokenAmount::zero(),
            params: RawBytes::from(exec_params),
            sequence: 1,
            ..Message::default()
        };
        let ret = executor.execute_message(create_message, ApplyKind::Explicit, 100)?;
        check_message_receipt(
            contract_case.name.clone() + "_create_actor",
            &ret,
            0,
            "".to_owned(),
            vec![],
        )?;
        from_slice(ret.msg_receipt.return_data.as_slice())
    }?;
    //invoke
    println!("actor cid {}", create_return.id_address);
    for wasm_case in &contract_case.cases {
        let from = accounts[wasm_case.send_from];

        let actor = executor.state_tree().get_actor(from.0)?.unwrap();
        let send_value = BigInt::from(wasm_case.send_value);
        let message = Message {
            from: from.1,
            sequence: actor.sequence,
            to: create_return.id_address,
            gas_limit: 1000000000000,
            method_num: wasm_case.method_num,
            params: RawBytes::from(hex::decode(&wasm_case.params)?),
            value: TokenAmount::from_whole(send_value),
            ..Message::default()
        };

        let ret = executor.execute_message(message, ApplyKind::Explicit, 100)?;
        check_message_receipt(
            contract_case.name.clone() + "_" + wasm_case.name.as_str(),
            &ret,
            wasm_case.expect_code,
            wasm_case.expect_message.clone(),
            hex::decode(wasm_case.return_data.clone())?,
        )?;
    }

    Ok(())
}

pub fn run_signle_wasm(accounts: &[InitAccount], wasm_case: &WasmCase) -> Result<()> {
    let path: PathBuf = [current_dir()?, Path::new(&wasm_case.binary).to_path_buf()]
        .iter()
        .collect::<PathBuf>()
        .absolutize()
        .map(|v| v.into_owned())
        .expect("get binary path");
    let buf = fs::read(path.clone()).map_err(|e| {
        anyhow!(
            "path {} not found {}",
            path.to_str().get_or_insert("known"),
            e
        )
    })?;

    let ret = exec(&buf, accounts, wasm_case)?;
    check_message_receipt(
        wasm_case.name.clone(),
        &ret,
        wasm_case.expect_code,
        wasm_case.expect_message.clone(),
        vec![],
    )
}

fn check_message_receipt(
    name: String,
    ret: &ApplyRet,
    expect_code: u32,
    expect_message: String,
    expect_receipt: Vec<u8>,
) -> Result<()> {
    if ret.msg_receipt.exit_code.value() != expect_code {
        if let Some(fail_info) = &ret.failure_info {
            println!(
                "{}:case {} expect exit code\n\t{} \nbut got \n\t{} {}",
                "failed".red(),
                name,
                expect_code,
                ret.msg_receipt.exit_code,
                fail_info
            );
            return Ok(());
        }
        println!(
            "{}: case {} expect exit code \n\t{} but got \n\t{}",
            "failed".red(),
            name,
            expect_code,
            ret.msg_receipt.exit_code
        );
        return Ok(());
    }

    if expect_code != 0 {
        if let Some(ApplyFailure::MessageBacktrace(trace)) = &ret.failure_info {
            let abort_msg = trace.frames.iter().last().unwrap().to_string();
            if !abort_msg.contains(expect_message.as_str()) {
                println!(
                    "{}: case {} expect message \n\t`{}` \nbut got \n\t`{}`",
                    "failed".red(),
                    name,
                    expect_message,
                    abort_msg
                );
                return Ok(());
            }
        }
    }

    if expect_code == 0 && !expect_receipt.is_empty() {
        let return_data = ret.msg_receipt.return_data.bytes();
        if return_data != expect_receipt {
            println!(
                "{}: case {} expect return data \n\t{} \n but got \n\t{}",
                "failed".red(),
                name,
                hex::encode(expect_receipt),
                hex::encode(return_data)
            );
            return Ok(());
        }
    }
    println! {"{}: case {}", "passed".green(), name};
    Ok(())
}

#[derive(Serialize_tuple, Deserialize_tuple, Clone, Debug)]
struct State {
    empty: bool,
}

pub fn new_tester(
    accounts_cfg: &[InitAccount],
) -> Result<(Tester<MemoryBlockstore, DummyExterns>, Vec<Account>)> {
    let bs = MemoryBlockstore::default();
    let bundle_root = bundle::import_bundle(&bs, actors_v10::BUNDLE_CAR).unwrap();
    let mut tester =
        Tester::new(NetworkVersion::V18, StateTreeVersion::V5, bundle_root, bs).unwrap();
    let mut accounts: Vec<Account> = vec![];
    for init_account in accounts_cfg {
        let balance = BigInt::from(init_account.balance);
        if init_account.isf4 {
            let sub_address =
                Address::new_delegated(10, init_account.subaddress.as_ref().unwrap().as_bytes())?;
            accounts
                .push(tester.create_placeholder(&sub_address, TokenAmount::from_whole(balance))?);
        } else {
            let priv_key = SecretKey::parse(&<[u8; 32]>::from_hex(
                init_account
                    .priv_key
                    .as_ref()
                    .expect("must specific priv key for secp address"),
            )?)?;
            let account =
                tester.make_secp256k1_account(priv_key, TokenAmount::from_whole(balance))?;
            accounts.push(account);
        }
    }
    Ok((tester, accounts))
}

pub fn exec(
    wasm_bin: &[u8],
    init_accounts: &[InitAccount],
    wasm_case: &WasmCase,
) -> Result<ApplyRet> {
    // Instantiate tester
    let (mut tester, accounts) = new_tester(init_accounts)?;
    // Set actor state
    let actor_state = State { empty: true };
    let state_cid = tester.set_state(&actor_state)?;

    // Set actor
    let actor_address = Address::new_id(10000);
    let actor_balance = BigInt::from(wasm_case.actor_balance);
    tester
        .set_actor_from_bin(
            wasm_bin,
            state_cid,
            actor_address,
            TokenAmount::from_whole(actor_balance),
        )
        .unwrap();

    // Instantiate machine
    tester.instantiate_machine(DummyExterns)?;

    // Send message
    let send_value = BigInt::from(wasm_case.send_value);
    let message = Message {
        from: accounts[wasm_case.send_from].1,
        to: actor_address,
        gas_limit: 1000000000000,
        method_num: wasm_case.method_num,
        value: TokenAmount::from_whole(send_value),
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
