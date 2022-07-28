use xshell::{cmd, Shell};

use bitflags::bitflags;

bitflags! {
    struct Check: u32 {
        const FORMAT = 0b00000001;
        const CLIPPY = 0b00000010;
        const TEST = 0b00000100;
        const COMPILE_CHECK = 0b00001000;
        const COMPILE_EXAMPLE = 0b00010000;
    }
}

const CLIPPY_FLAGS: [&str; 8] = [
    "-Aclippy::type_complexity",
    "-Wclippy::doc_markdown",
    "-Wclippy::redundant_else",
    "-Wclippy::match_same_arms",
    "-Wclippy::semicolon_if_nothing_returned",
    "-Wclippy::explicit_iter_loop",
    "-Wclippy::map_flatten",
    "-Dwarnings",
];

fn main() {
    // When run locally, results may differ from actual CI runs triggered by
    // .github/workflows/ci.yml
    // - Official CI runs latest stable
    // - Local runs use whatever the default Rust is locally

    let what_to_run = match std::env::args().nth(1).as_deref() {
        Some("format") => Check::FORMAT,
        Some("clippy") => Check::CLIPPY,
        Some("test") => Check::TEST,
        Some("lints") => Check::FORMAT | Check::CLIPPY,
        Some("compile") => Check::COMPILE_CHECK,
        Some("example") => Check::COMPILE_EXAMPLE,
        _ => Check::all(),
    };

    let sh = Shell::new().unwrap();

    if what_to_run.contains(Check::FORMAT) {
        // See if any code needs to be formatted
        cmd!(sh, "cargo fmt --all -- --check")
            .run()
            .expect("Please run 'cargo fmt --all' to format your code.");
    }

    if what_to_run.contains(Check::CLIPPY) {
        // See if clippy has any complaints.
        // - Type complexity must be ignored because we use huge templates for queries
        cmd!(
            sh,
            "cargo clippy --workspace --all-targets --all-features -- {CLIPPY_FLAGS...}"
        )
        .run()
        .expect("Please fix clippy errors in output above.");
    }

    if what_to_run.contains(Check::TEST) {
        // Run tests (except doc tests and without building examples)
        cmd!(sh, "cargo test  --bins --tests")
            .run()
            .expect("Please fix failing tests in output above.");
    }

    if what_to_run.contains(Check::COMPILE_CHECK) {
        // Build examples and check they compile
        cmd!(sh, "cargo check --workspace")
            .run()
            .expect("Please fix failing doc-tests in output above.");
    }

    if what_to_run.contains(Check::COMPILE_EXAMPLE) {
        // Build examples and check they compile
        sh.change_dir("./examples/hellocontract");
        cmd!(sh, "go-fvm-sdk-tools build")
            .run()
            .expect("Please fix hellcontract example.");
    }
}
