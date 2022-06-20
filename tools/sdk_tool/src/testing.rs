use anyhow::Result;
use cid::Cid;
use clap::Parser;
use colored::*;
use fvm::executor::ApplyFailure;
use fvm::executor::ApplyRet;
use fvm::executor::{ApplyKind, Executor};
use fvm::init_actor::INIT_ACTOR_ADDR;
use fvm::machine::Machine;
use fvm_integration_tests::tester::{Account, Tester};
use fvm_ipld_blockstore::MemoryBlockstore;
use fvm_ipld_encoding::tuple::*;
use fvm_ipld_encoding::{Cbor, RawBytes};
use fvm_shared::address::Address;
use fvm_shared::bigint::BigInt;
use fvm_shared::bigint::Zero;
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
    #[clap(last = true)]
    path: String,

    #[clap(short, long)]
    name: Option<String>,
}

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct InitAccount {
    priv_key: String,
    balance: u64,
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

pub fn run_testing(cfg: &TestConfig) {
    let case_meta_path: PathBuf = [cfg.path.clone(), "test.json".to_string()].iter().collect();
    let buf = fs::read(case_meta_path).unwrap();

    let test_json: TestJson = serde_json::from_slice(&buf).unwrap();

    if let Some(cases) = test_json.cases {
        cases
            .iter()
            .filter(|v| cfg.name.is_none() || v.name.eq(cfg.name.as_ref().unwrap()))
            .for_each(|test_case| {
                run_signle_wasm(cfg.path.clone(), &test_json.accounts, test_case);
            });
    }

    if let Some(contracts) = &test_json.contracts {
        contracts
            .iter()
            .filter(|v| cfg.name.is_none() || v.name.eq(cfg.name.as_ref().unwrap()))
            .for_each(|group_case| {
                run_action_group(cfg.path.clone(), &test_json.accounts, group_case).unwrap();
            });
    }
}

pub fn run_action_group(
    root_path: String,
    accounts_cfg: &[InitAccount],
    contract_case: &ContractCase,
) -> Result<()> {
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
        let install_params = InstallParams {
            code: RawBytes::from(buf),
        }
        .marshal_cbor()?;
        let install_message = Message {
            from: accounts[contract_case.owner_account].1,
            to: INIT_ACTOR_ADDR,
            gas_limit: 1000000000000,
            method_num: 3,
            value: BigInt::zero(),
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
        );
        InstallReturn::unmarshal_cbor(ret.msg_receipt.return_data.as_slice())
    }?;

    println!("code cid {}", install_return.code_cid.clone());
    //create
    let create_return: ExecReturn = {
        let constructor_params = hex::decode(contract_case.constructor.clone())?;
        let exec_params = ExecParams {
            code_cid: install_return.code_cid,
            constructor_params: RawBytes::from(constructor_params),
        }
        .marshal_cbor()?;

        let create_message = Message {
            from: accounts[contract_case.owner_account].1,
            to: INIT_ACTOR_ADDR,
            gas_limit: 1000000000000,
            method_num: 2,
            value: BigInt::zero(),
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
        );
        ExecReturn::unmarshal_cbor(ret.msg_receipt.return_data.as_slice())
    }?;
    //invoke
    println!("actor cid {}", create_return.id_address);
    for wasm_case in &contract_case.cases {
        let from_addr = accounts[wasm_case.send_from].1;
        let actor = executor.state_tree().get_actor(&from_addr)?.unwrap();
        let message = Message {
            from: from_addr,
            sequence: actor.sequence,
            to: create_return.id_address,
            gas_limit: 1000000000000,
            method_num: wasm_case.method_num,
            params: RawBytes::from(hex::decode(&wasm_case.params)?),
            value: BigInt::from(wasm_case.send_value),
            ..Message::default()
        };

        let ret = executor.execute_message(message, ApplyKind::Explicit, 100)?;
        check_message_receipt(
            contract_case.name.clone() + "_" + wasm_case.name.as_str(),
            &ret,
            wasm_case.expect_code,
            wasm_case.expect_message.clone(),
            hex::decode(wasm_case.return_data.clone())?,
        );
    }

    Ok(())
}

pub fn run_signle_wasm(root_path: String, accounts: &[InitAccount], wasm_case: &WasmCase) {
    let path: PathBuf = [root_path, wasm_case.binary.to_owned()].iter().collect();
    let buf = fs::read(path.clone())
        .unwrap_or_else(|_| panic!("path {} not found", path.to_str().unwrap()));
    let ret = exec(&buf, accounts, wasm_case).unwrap();
    check_message_receipt(
        wasm_case.name.clone(),
        &ret,
        wasm_case.expect_code,
        wasm_case.expect_message.clone(),
        vec![],
    );
}

fn check_message_receipt(
    name: String,
    ret: &ApplyRet,
    expect_code: u32,
    expect_message: String,
    expect_receipt: Vec<u8>,
) {
    if ret.msg_receipt.exit_code.value() != expect_code {
        if let Some(fail_info) = &ret.failure_info {
            panic!(
                "{}:case {} expect exit code {} but got {} {}",
                "failed".red(),
                name,
                expect_code,
                ret.msg_receipt.exit_code,
                fail_info
            )
        } else {
            panic!(
                "{}: case {} expect exit code {} but got {}",
                "failed".red(),
                name,
                expect_code,
                ret.msg_receipt.exit_code
            )
        }
    }

    if expect_code != 0 {
        if let Some(ApplyFailure::MessageBacktrace(trace)) = &ret.failure_info {
            let abort_msg = trace.frames.iter().last().unwrap().to_string();
            if !abort_msg.contains(expect_message.as_str()) {
                panic!(
                    "{}: case {} expect messcage `{}` but got `{}`",
                    "failed".red(),
                    name,
                    expect_message,
                    abort_msg
                )
            }
        }
    }

    if expect_code == 0 && !expect_receipt.is_empty() {
        let return_data = ret.msg_receipt.return_data.bytes().to_vec();
        if return_data != expect_receipt {
            panic!(
                "{}: case {} expect return data {} but got {}",
                "failed".red(),
                name,
                hex::encode(expect_receipt),
                hex::encode(return_data)
            )
        }
    }
    println! {"{}: case {}", "passed".green(), name}
}

#[derive(Serialize_tuple, Deserialize_tuple, Clone, Debug)]
struct State {
    empty: bool,
}

pub fn new_tester(
    accounts_cfg: &[InitAccount],
) -> Result<(Tester<MemoryBlockstore>, Vec<Account>)> {
    let mut tester = Tester::new(
        NetworkVersion::V15,
        StateTreeVersion::V4,
        MemoryBlockstore::default(),
    )
    .unwrap();
    let mut accounts: Vec<Account> = vec![];
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
    wasm_case: &WasmCase,
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
            BigInt::from(wasm_case.actor_balance),
        )
        .unwrap();

    // Instantiate machine
    tester.instantiate_machine()?;

    // Send message
    let message = Message {
        from: accounts[wasm_case.send_from].1,
        to: actor_address,
        gas_limit: 1000000000000,
        method_num: wasm_case.method_num,
        value: BigInt::from(wasm_case.send_value),
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
