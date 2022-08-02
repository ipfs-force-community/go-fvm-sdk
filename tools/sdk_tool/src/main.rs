mod patch;
mod template;
mod testing;
mod utils;
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
    /// create new template project by module name
    New(template::NewTemplateConfig),
    /// apply path for go/tinygo
    /// if your go and tinygo install in user home directory, just run./go-fvm-sdk-tools patch
    /// if you go and tinygo is installed in /usr/local/go, use sudo ./go-fvm-sdk-tools patch
    /// if you are in china and need proxy, exec sudo(opt) https_proxy=<proxy> http_proxy=<> ./go-fvm-sdk-tools patch
    /// if want to install manual or know more detail refer https://github.com/ipfs-force-community/go_tinygo_patch
    Patch(patch::PatchConfig),
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
        Commands::Test(cfg) => {
            if let Err(e) = testing::run_testing(cfg) {
                println!("run test command fail {}", e);
                std::process::exit(1);
            }
        }
        Commands::New(cfg) => {
            if let Err(e) = template::new_template_project(cfg) {
                println!("run new template command fail {}", e);
                std::process::exit(1);
            }
        }
        Commands::Patch(cfg) => {
            if let Err(e) = patch::apply_patch(cfg) {
                println!("apply patch command fail {}", e);
                std::process::exit(1);
            }
        }
    }
}
