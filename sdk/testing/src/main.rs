use clap::Parser;
use fvm::executor::ApplyFailure;
use fvm::executor::ApplyRet;
use fvm::executor::{ApplyKind, Executor};
use fvm_integration_tests::tester::{Account, Tester};
use fvm_ipld_encoding::tuple::*;
use fvm_shared::address::Address;
use fvm_shared::bigint::BigInt;
use fvm_shared::message::Message;
use fvm_shared::state::StateTreeVersion;
use fvm_shared::version::NetworkVersion;
use num_traits::Zero;
use serde::{Deserialize, Serialize};
use std::fs;
use std::iter::Iterator;
use std::path::PathBuf;

#[derive(Parser, Debug)]
#[clap(author, version, about, long_about = None)]
struct Args {
    #[clap(short, long)]
    path: String,
}

#[derive(Serialize, Deserialize, Clone, Debug)]
struct TestCase {
    name: String,
    binary: String,
    expect_code: u32,
    expect_message: String,
    method_num: u64,
}

fn main() {
    let args = Args::parse();
    let case_meta_path: PathBuf = [args.path.clone().to_owned(), "cases.json".to_string()]
        .iter()
        .collect();
    let buf = fs::read(case_meta_path).unwrap();
    let test_cases: Vec<TestCase> = serde_json::from_slice(&buf).unwrap();

    test_cases.iter().for_each(|test_case| {
        let path: PathBuf = [args.path.clone().to_owned(), test_case.binary.to_owned()]
            .iter()
            .collect();
        let buf = fs::read(path.clone())
            .expect(format!("path {} not found", path.clone().to_str().unwrap()).as_str());
        let ret = exec(&buf, test_case.method_num);
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
        /*
        let ret_msg = std::str::from_utf8(ret.msg_receipt.return_data.bytes()).uwrap();
        if ret_msg.clone() != test_case.expect_message {
            panic!("case {} expect messcage {} but got {}", test_case.name, test_case.expect_message, ret_msg)
        }
        */
        if test_case.expect_code != 0 {
            match ret.failure_info.unwrap() {
                ApplyFailure::MessageBacktrace(mut trace) => {
                    let abort_msg = trace.frames.pop().unwrap().message;
                    if abort_msg.clone() != test_case.expect_message {
                        panic!(
                            "case {} expect messcage {} but got {}",
                            test_case.name, test_case.expect_message, abort_msg
                        )
                    }
                }
                _ => {}
            }
        }
    });
}

#[derive(Serialize_tuple, Deserialize_tuple, Clone, Debug)]
struct State {
    empty: bool,
}

pub fn exec(wasm_bin: &[u8], method_num: u64) -> ApplyRet {
    // Instantiate tester
    let mut tester = Tester::new(NetworkVersion::V15, StateTreeVersion::V4).unwrap();

    let sender: [Account; 1] = tester.create_accounts().unwrap();

    // Get wasm bin
    //  let wasm_bin = wat2wasm(wat).unwrap();

    // Set actor state
    let actor_state = State { empty: true };
    let state_cid = tester.set_state(&actor_state).unwrap();

    // Set actor
    let actor_address = Address::new_id(10000);

    tester
        .set_actor_from_bin(&wasm_bin, state_cid, actor_address, BigInt::zero())
        .unwrap();

    // Instantiate machine
    tester.instantiate_machine().unwrap();

    // Send message
    let message = Message {
        from: sender[0].1,
        to: actor_address,
        gas_limit: 1000000000,
        method_num: method_num,
        value: BigInt::from(10),
        ..Message::default()
    };

    let exec_result = tester
        .executor
        .unwrap()
        .execute_message(message, ApplyKind::Explicit, 100);
    exec_result.unwrap()
}
