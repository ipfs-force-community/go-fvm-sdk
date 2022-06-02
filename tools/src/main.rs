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
    /// Adds files to myapp
    Process(wasmprocess::ProcessConfig),
    Test(testing::TestConfig),
}

fn main() {
    let cli = Cli::parse();
    match &cli.command {
        Commands::Process(cfg ) => {
           wasmprocess::run_process(cfg).unwrap();
        }
        Commands::Test ( cfg )=> {
           testing::run_testing(cfg)
        }
    }
}
