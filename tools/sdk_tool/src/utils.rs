use anyhow::{anyhow, Result};
use std::collections::HashMap;
use std::fs::File;
use std::io;
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
                Err(anyhow!("unable to found tinygo, please install tinygo https://tinygo.org/getting-started/install"))
            } else {
                Err(anyhow!("fail exec tinygo version {}", e))
            }
        }
    }
}

pub fn get_tinygo_env() -> Result<HashMap<String, String>> {
    match Command::new("tinygo").arg("env").output() {
        Ok(output) => Ok(HashMap::from_iter(
            String::from_utf8(output.stdout)?
                .split('\n')
                .map(|v| v.trim())
                .filter(|v| !v.is_empty())
                .map(|v| {
                    let key_pare: Vec<&str> = v.split('=').collect();
                    (
                        key_pare[0].to_string(),
                        key_pare[1].trim_matches('\"').to_string(),
                    )
                }),
        )),
        Err(e) => {
            if let ErrorKind::NotFound = e.kind() {
                Err(anyhow!("unable to found tinygo, please install tinygo https://tinygo.org/getting-started/install"))
            } else {
                Err(anyhow!("fail exec tinygo env {}", e))
            }
        }
    }
}

pub fn get_tinygo_version() -> Result<String> {
    match Command::new("tinygo").arg("version").output() {
        Ok(output) => Ok(String::from_utf8(output.stdout)?),
        Err(e) => {
            if let ErrorKind::NotFound = e.kind() {
                Err(anyhow!("unable to found tinygo, please install tinygo https://tinygo.org/getting-started/install"))
            } else {
                Err(anyhow!("fail exec tinygo version {}", e))
            }
        }
    }
}

pub fn get_patch_version() -> Result<String> {
    match Command::new("patch").arg("-version").output() {
        Ok(output) => Ok(String::from_utf8(output.stdout)?),
        Err(e) => {
            if let ErrorKind::NotFound = e.kind() {
                Err(anyhow!("unable to found patch tool"))
            } else {
                Err(anyhow!("fail exec patch -version {}", e))
            }
        }
    }
}

pub fn check_fvm_tool_install() -> Result<()> {
    match Command::new("fvm_go_sdk")
        .stdout(Stdio::null())
        .arg("--help")
        .status()
    {
        Ok(_) => Ok(()),
        Err(e) => {
            if let ErrorKind::NotFound = e.kind() {
                Err(anyhow!("unable to found fvm_go_sdk, please install this tool in https://github.com/ipfs-force-community/go-fvm-sdk/releases"))
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
            if version_str.contains("go1.16.")
                || version_str.contains("go1.17.")
                || version_str.contains("go1.18.")
                || version_str.contains("go1.19")
            {
                Ok(true)
            } else {
                Err(anyhow!(
                    "uncorect go version must be go 1.16.x/go1.17.x/go1.18.x/go1.19 but got {}",
                    version_str
                ))
            }
        }
        Err(e) => {
            if let ErrorKind::NotFound = e.kind() {
                Err(anyhow!(
                    "unable to found go, please install this tool in https://go.dev/dl"
                ))
            } else {
                Err(anyhow!("check err {}", e))
            }
        }
    }
}

pub fn download_file(path: &str, file_path: &str) -> Result<()> {
    let mut resp = reqwest::blocking::get(path)?;
    if resp.status().is_success() {
        let mut out = File::create(file_path)?;
        io::copy(&mut resp, &mut out)?;
        Ok(())
    } else {
        Err(anyhow!(
            "download {} fail status {}",
            file_path,
            resp.status()
        ))
    }
}
