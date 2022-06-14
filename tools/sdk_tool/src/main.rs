mod testing;
mod wasmprocess;

use clap::Parser;
use clap::Subcommand;

#[derive(Parser)]
#[clap(author, version, about, long_about = None)]
#[clap(propagate_version = true)]
struct Cli {
    #[clap(subcommand)]
    command: Commands,
}

#[derive(Subcommand)]
enum Commands {
    /// build and process wasm
    Build(wasmprocess::BuildCLiConfig),
    /// test wasm on fvm
    Test(testing::TestConfig),
}

fn main() {
    let cli = Cli::parse();
    match &cli.command {
        Commands::Build(cfg) => {
            if let Err(e) = wasmprocess::run_process(cfg) {
                println!("run build command fail {}", e);
                std::process::exit(1);
            }
        }
        Commands::Test(cfg) => testing::run_testing(cfg),
    }
}
