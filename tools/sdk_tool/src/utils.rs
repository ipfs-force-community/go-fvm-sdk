use anyhow::{anyhow, Result};
use std::io::ErrorKind;
use std::process::{Command, Stdio};

pub fn check_tinygo_install() -> Result<()> {
    match Command::new("tinygo")
        .stdout(Stdio::null())
        .arg("version")
        .status()
    {
        Ok(_) => Ok(()),
        Err(e) => {
            if let ErrorKind::NotFound = e.kind() {
                Err(anyhow!("unable to found tinygo(fvm), please install this tool in https://github.com/ipfs-force-community/tinygo/releases"))
            } else {
                Err(anyhow!("fvm-tinygo not install, please install  err {}", e))
            }
        }
    }
}

pub fn check_fvm_tool_install() -> Result<()> {
    match Command::new("go-fvm-sdk-tools")
        .stdout(Stdio::null())
        .arg("--help")
        .status()
    {
        Ok(_) => Ok(()),
        Err(e) => {
            if let ErrorKind::NotFound = e.kind() {
                Err(anyhow!("unable to found go-fvm-sdk-tools(fvm), please install this tool in https://github.com/ipfs-force-community/go-fvm-sdk/releases"))
            } else {
                Err(anyhow!("check err {}", e))
            }
        }
    }
}

pub fn check_go_install() -> Result<bool> {
    match Command::new("go")
        .stdout(Stdio::piped())
        .arg("version")
        .spawn()
    {
        Ok(child) => {
            let output = child.wait_with_output()?;
            let version_str = String::from_utf8(output.stdout)?;
            for s in SUPPORT_GO_VERSIONS.iter() {
                if version_str.contains(s) {
                    return Ok(true);
                }
            }
            Err(anyhow!("incorrect go version:{},must in {:?}", version_str, SUPPORT_GO_VERSIONS))
        },
        Err(e) => {
            if let ErrorKind::NotFound = e.kind() {
                Err(anyhow!("unable to found go-fvm-sdk-tools(fvm), please install this tool in https://go.dev/dl"))
            } else {
                Err(anyhow!("check err {}", e))
            }
        }
    }
}

const SUPPORT_GO_VERSIONS: &'static [&'static str] = &["go1.16.", "go1.17.", "go1.18."];
