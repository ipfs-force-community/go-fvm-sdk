use crate::utils;
use anyhow::{anyhow, Result};
use clap::Parser;
use regex::Regex;
use std::collections::HashMap;
use std::env;
use std::path::Path;
use xshell::Shell;

#[derive(Parser, Debug)]
#[clap(author, version, about, long_about = None)]
pub struct PatchConfig {}

pub fn apply_patch(_: &PatchConfig) -> Result<()> {
    utils::check_tinygo_install()?;
    utils::check_go_install()?;
    let patch_version = utils::get_patch_version()?;
    println!("patch version {}", patch_version);
    let envs = utils::get_tinygo_env()?;
    println!("{:?}", envs);

    let mut go_patch_map: HashMap<String, String> = HashMap::new();
    go_patch_map.insert("1.16.x".to_string(), "go_v1.16.x.patch".to_string());
    go_patch_map.insert("1.17.x".to_string(), "go_v1.17.x.patch".to_string());
    go_patch_map.insert("1.18.x".to_string(), "go_v1.18.x.patch".to_string());
    go_patch_map.insert("1.19".to_string(), "go_v1.19.patch".to_string());

    let mut tinygo_patch_map: HashMap<String, String> = HashMap::new();
    tinygo_patch_map.insert("0.24.x".to_string(), "tinygo_v0.24.x.patch".to_string());
    tinygo_patch_map.insert("0.25.x".to_string(), "tinygo_v0.25.x.patch".to_string());

    let version_str = utils::get_tinygo_version()?;
    let re = Regex::new(r"\d+\.\d+(\.\d+)?").unwrap();
    let version_arr: Vec<String> = re
        .captures_iter(version_str.as_str())
        .map(|c| c[0].to_string())
        .collect();
    let tinygo_version = version_arr.get(0).unwrap();
    let go_version = version_arr.get(1).unwrap();
    println!(
        "go version {} tinygo version {}",
        go_version, tinygo_version
    );

    let dir = env::current_dir()?;
    let current_dir = dir.as_os_str().to_str().unwrap();

    {
        let go_root_path = envs.get("GOROOT").expect("unable to locate GOROOT");
        println!("go root path {}", go_root_path);
        let default_patch_name = &format!("go_v{}.patch", default_version(go_version));
        let patch_name = go_patch_map.get(go_version).unwrap_or(default_patch_name);

        let patch_url = format!(
            "https://raw.githubusercontent.com/ipfs-force-community/go_tinygo_patch/main/patchs/{}",
            patch_name
        );
        println!("download go patch from {}", patch_url);
        utils::download_file(&patch_url, patch_name)?;

        let sh = Shell::new()?;
        sh.change_dir(Path::new(&go_root_path));
        sh.cmd("patch")
            .arg("-p1")
            .arg("-f")
            .arg("-i")
            .arg(format!("{}/{}", current_dir, patch_name))
            .run()
            .map_err(|e| anyhow!("unable to apply patch for go {}", e))?;
        std::fs::remove_file(patch_name)?;
    }
    {
        let tinygo_root_path = envs.get("TINYGOROOT").expect("unable to locate TINYGOROOT");
        println!("tinygo root path {}", tinygo_root_path);
        let default_patch_name = &format!("tinygo_v{}.patch", default_version(tinygo_version));
        let patch_name = tinygo_patch_map
            .get(tinygo_version)
            .unwrap_or(default_patch_name);
        let patch_url = format!(
            "https://raw.githubusercontent.com/ipfs-force-community/go_tinygo_patch/main/patchs/{}",
            patch_name
        );
        println!("download tinygo patch from {}", patch_url);
        utils::download_file(&patch_url, patch_name)?;

        let sh = Shell::new()?;
        sh.change_dir(Path::new(&tinygo_root_path));
        sh.cmd("patch")
            .arg("-p1")
            .arg("-f")
            .arg("-i")
            .arg(format!("{}/{}", current_dir, patch_name))
            .run()
            .map_err(|e| anyhow!("unable to apply patch for tinygo {}", e))?;
        std::fs::remove_file(patch_name)?;
    }

    Ok(())
}

fn default_version(str: &str) -> String {
    let mut version_seq: Vec<String> = str.split('.').map(|s| s.to_string()).collect();
    if version_seq.len() == 3 {
        //pop minor version number
        version_seq.pop().unwrap();
        version_seq.push("x".to_owned());
    }
    version_seq.join(".")
}
