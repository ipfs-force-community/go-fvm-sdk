use crate::utils;
use anyhow::{anyhow, Result};
use clap::Parser;
use std::env::current_dir;
use std::fs;
use xshell::{cmd, Shell};

#[derive(Parser, Debug)]
#[clap(author, version, about, long_about = None)]
pub struct NewTemplateConfig {
    //must be a validate go module package name
    #[clap(last = true)]
    pub name: Option<String>,
}

pub fn new_template_project(cfg: &NewTemplateConfig) -> Result<()> {
    utils::check_tinygo_install()?;
    utils::check_go_install()?;
    utils::check_fvm_tool_install()?;
    let template_name = "gofvm-counter";
    //market
    let sh = Shell::new()?;
    cmd!(
        sh,
        "git clone https://github.com/ipfs-force-community/gofvm-counter.git"
    )
    .run()
    .map_err(|e| anyhow!("unable to checkout template project {}", e))?;

    let mut old_tmp_dir = current_dir()?;
    old_tmp_dir.push(template_name);
    sh.change_dir(&old_tmp_dir);
    cmd!(sh, "rm -rf .git")
        .run()
        .map_err(|e| anyhow!("unable to remove git {}", e))?;

    sh.change_dir(current_dir()?);
    let mut module_name = template_name.to_string();
    if let Some(new_module_name) = &cfg.name {
        sh.change_dir(current_dir()?);
        sh.cmd("mv")
            .args([template_name, new_module_name])
            .run()
            .map_err(|e| {
                anyhow!(
                    "unable to rename template project to {} {}",
                    new_module_name,
                    e
                )
            })?;
        module_name = new_module_name.to_string();
    }

    //replace module name
    let mut new_cur_dir = current_dir()?;
    new_cur_dir.push(module_name.clone());
    if let Some(new_module_name) = &cfg.name {
        for file in walkdir::WalkDir::new(&new_cur_dir)
            .into_iter()
            .filter_map(|file| file.ok())
            .filter(|file| {
                if let Ok(meta) = file.metadata() {
                    return meta.is_file();
                }
                false
            })
        {
            let file_content = fs::read_to_string(file.path())?;
            let file_content = file_content.replace(template_name, new_module_name);
            fs::write(file.path(), file_content)?;
        }
        println!("module....");
    }

    sh.change_dir(new_cur_dir.clone());
    cmd!(sh, "mkdir client")
        .run()
        .map_err(|e| anyhow!("unable to create client dir {}", e))?;

    let mut gen_dir = new_cur_dir.clone();
    gen_dir.push("gen");
    sh.change_dir(gen_dir);

    cmd!(sh, "go mod tidy")
        .run()
        .map_err(|e| anyhow!("unable to create client dir {}", e))?;
    cmd!(sh, "go run main.go")
        .run()
        .map_err(|e| anyhow!("unable to create run gen tool {}", e))?;

    sh.change_dir(&new_cur_dir);
    cmd!(sh, "go mod tidy")
        .run()
        .map_err(|e| anyhow!("unable to build template in template project {}", e))?;

    sh.cmd("go-fvm-sdk-tools")
        .args(["build", "-o", &(module_name + ".wasm")])
        .run()
        .map_err(|e| anyhow!("unable to build template in template project {}", e))?;

    cmd!(sh, "go-fvm-sdk-tools test -- ./tests")
        .run()
        .map_err(|e| anyhow!("unable to run test in template project {}", e))?;

    Ok(())
}
