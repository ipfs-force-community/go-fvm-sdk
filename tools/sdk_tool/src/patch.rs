use crate::utils;
use anyhow::{anyhow, Result};
use clap::Parser;
use regex::Regex;
use std::env;
use std::path::Path;
use xshell::Shell;

#[derive(Parser, Debug)]
#[clap(author, version, about, long_about = None)]
pub struct PatchConfig {}

pub fn apply_patch(_: &PatchConfig) -> Result<()> {
    utils::check_tinygo_install()?;
    utils::check_go_install()?;
    let envs = utils::get_tinygo_env()?;
    println!("{:?}", envs);

    let version_str = utils::get_tinygo_version()?;
    let re = Regex::new(r"\d+\.\d+\.").unwrap();
    let version_arr: Vec<String> = re
        .captures_iter(version_str.as_str())
        .map(|c| c[0].to_string())
        .collect();
    let tinyo_version = version_arr.get(0).unwrap();
    let go_version = version_arr.get(1).unwrap();
    println!(
        "go version {}.x tinygo version {}.x",
        go_version, tinyo_version
    );

    let dir = env::current_dir()?;
    let current_dir = dir.as_os_str().to_str().unwrap();

    {
        let go_root_path = envs.get("GOROOT").expect("unable to locate GOROOT");
        println!("go root path {}", go_root_path);
        utils::download_file(format!("https://raw.githubusercontent.com/ipfs-force-community/go_tinygo_patch/main/patchs/fmt_v{}x.patch", go_version).as_str(),
                             format!("fmt_v{}x.patch", go_version).as_str())?;

        let sh = Shell::new()?;
        sh.change_dir(Path::new(&go_root_path));
        sh.cmd("patch")
            .arg("-p1")
            .arg("-i")
            .arg(format!("{}/fmt_v{}x.patch", current_dir, go_version))
            .run()
            .map_err(|e| anyhow!("unable to apply patch for go {}", e))?;
    }
    {
        let tinygo_root_path = envs.get("TINYGOROOT").expect("unable to locate TINYGOROOT");
        println!("tinygo root path {}", tinygo_root_path);
        utils::download_file(format!("https://raw.githubusercontent.com/ipfs-force-community/go_tinygo_patch/main/patchs/tinygo_0.24.0_reflect.patch").as_str(),
                             format!("tinygo_0.24.0_reflect.patch").as_str())?;

        let sh = Shell::new()?;
        sh.change_dir(Path::new(&tinygo_root_path));
        sh.cmd("patch")
            .arg("-p1")
            .arg("-i")
            .arg(format!("{}/tinygo_0.24.0_reflect.patch", current_dir))
            .run()
            .map_err(|e| anyhow!("unable to apply patch for tinygo {}", e))?;
    }

    Ok(())
}
