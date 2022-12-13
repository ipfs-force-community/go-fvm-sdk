use crate::utils;
use anyhow::{anyhow, Result};
use clap::Parser;
use parity_wasm::elements::Type::Function;
use parity_wasm::elements::{
    BlockType, External, Func, FuncBody, FunctionType, ImportCountType, ImportEntry, Instruction,
    Instructions, Module, Type, ValueType,
};
use path_absolutize::*;
use std::env;
use std::io::{self, Write};
use std::path::{Path, PathBuf};
use std::process::Command;

#[derive(Parser, Debug)]
pub struct BuildCliConfig {
    #[clap(last = true)]
    pub input: Option<String>,
    #[clap(short, long)]
    pub output: Option<String>,
    #[clap(short, long)]
    pub wat: bool,
    #[clap(short, long)]
    pub debug: bool,
}

pub struct BuildOptions {
    pub code_path: String,
    pub target_dir: String,
    pub target_name: String,
    pub output_wasm_path: String,
    pub output_wat_path: String,
    pub debug: bool,
}

impl BuildOptions {
    pub fn new(pwd_path: PathBuf, cfg: &BuildCliConfig) -> Result<Self> {
        let code_path = if let Some(input) = &cfg.input {
            pwd_path
                .join(Path::new(input))
                .absolutize()
                .map(|v| v.into_owned())?
        } else {
            pwd_path.clone()
        };

        let (target_dir, target_name): (PathBuf, String) = if let Some(o_path) = &cfg.output {
            let abs_output_path = pwd_path
                .join(Path::new(o_path))
                .absolutize()
                .map(|v| v.into_owned())?;
            if abs_output_path.extension().is_none() {
                let target_name = code_path
                    .with_extension("")
                    .file_name()
                    .ok_or_else(|| anyhow!("get file name from {:?}", code_path))
                    .map(|v| v.to_str().unwrap().to_string())?;
                (abs_output_path, target_name)
            } else {
                let target_name = abs_output_path
                    .with_extension("")
                    .file_name()
                    .ok_or_else(|| anyhow!("get file name from {:?}", abs_output_path))
                    .map(|v| v.to_str().unwrap().to_string())?;
                let target_dir = abs_output_path
                    .parent()
                    .ok_or_else(|| anyhow!("get parent path for {:?}", abs_output_path))
                    .map(|v| v.to_path_buf())?;
                (target_dir, target_name)
            }
        } else {
            //output is current dir if not specify output
            let target_name = code_path
                .with_extension("")
                .file_name()
                .ok_or_else(|| anyhow!("get file name from {:?}", code_path))
                .map(|v| v.to_str().unwrap().to_string())?;
            (pwd_path, target_name)
        };

        let output_wasm_path = Path::new(&target_dir)
            .join(target_name.clone() + ".wasm")
            .try_to_string()?;
        let output_wat_path = Path::new(&target_dir)
            .join(target_name.clone() + ".wat")
            .try_to_string()?;
        Ok(BuildOptions {
            code_path: code_path.try_to_string()?,
            target_name,
            target_dir: target_dir.try_to_string()?,
            output_wasm_path,
            output_wat_path,
            debug: cfg.debug,
        })
    }
}

pub fn run_process(cfg: &BuildCliConfig) -> Result<()> {
    let parent = env::current_dir().unwrap();
    let build_opts = BuildOptions::new(parent, cfg)?;
    let result = GoFvmBinProcessor::new(&build_opts)
        .build()?
        .append_init_to_invoke()?
        .replace_fd_write()?
        .get_binary()?;

    let wat_str = wasmprinter::print_bytes(result)?;

    if cfg.wat {
        std::fs::write(build_opts.output_wat_path, wat_str.clone())?;
    }

    let mut features = wabt::Features::new();
    features.set_annotations_enabled(true);
    features.set_bulk_memory_enabled(true);
    features.set_sat_float_to_int_enabled(true);
    features.set_sign_extension_enabled(true);
    let wat_bin = wabt::wat2wasm_with_features(wat_str, features)?;
    std::fs::write(build_opts.output_wasm_path, wat_bin)?;
    Ok(())
}

pub struct GoFvmBinProcessor<'a> {
    module: Module,
    build_cfg: &'a BuildOptions,
}

impl<'a> GoFvmBinProcessor<'a> {
    pub fn new(build_cfg: &'a BuildOptions) -> Self {
        GoFvmBinProcessor {
            module: Module::default(),
            build_cfg,
        }
    }

    pub fn build(&mut self) -> Result<&mut Self> {
        utils::check_tinygo_install()?;
        let mut args = vec!["build", "-target", "wasi", "-no-debug"];
        if !self.build_cfg.debug {
            args.append(&mut vec!["-panic", "trap"]);
        }

        args.append(&mut vec![
            "-o",
            &self.build_cfg.output_wasm_path,
            &self.build_cfg.code_path,
        ]);
        let output = Command::new("tinygo")
            .args(args)
            .spawn()?
            .wait_with_output()
            .expect("unable to get output");
        if !output.status.success() {
            return Err(anyhow!(format!("run tinygo command failed err {output:?}")));
        }

        io::stdout().write_all(&output.stdout).unwrap();
        io::stderr().write_all(&output.stderr).unwrap();
        let module = parity_wasm::deserialize_file(&self.build_cfg.output_wasm_path)?
            .parse_names()
            .map_err(|_| anyhow!("parser names in wasm"))?;
        self.module = module;
        Ok(self)
    }

    pub fn get_binary(&self) -> Result<Vec<u8>> {
        parity_wasm::serialize(self.module.clone())
            .map_err(|e| anyhow!("convert module to binary {}", e))
    }

    pub fn replace_fd_write(&mut self) -> Result<&mut Self> {
        //探测需不需要插入fd_write
        if !self.build_cfg.debug || !self.has_fd_write() {
            return Ok(self);
        }

        let fd_write_index = self
            .get_import_func_index("wasi_snapshot_preview1", "fd_write")
            .expect("unable to find fd_write import");

        {
            //fd_write import 替换成debug.log
            let debug_type_index = self
                .get_debug_type()
                .or_else(|| {
                    let types = self.module.type_section_mut().unwrap().types_mut();
                    types.push(Type::Function(FunctionType::new(
                        vec![ValueType::I32, ValueType::I32],
                        vec![ValueType::I32],
                    )));
                    Some(types.len() as u32)
                })
                .unwrap();

            //插入debug log到import部分和function部分
            let imports = self.module.import_section_mut().unwrap().entries_mut();
            imports[fd_write_index] = ImportEntry::new(
                "debug".to_owned(),
                "log".to_owned(),
                External::Function(debug_type_index),
            );
        };

        {
            let func_names_map = self
                .module
                .names_section_mut()
                .unwrap()
                .functions_mut()
                .as_mut()
                .unwrap()
                .names_mut();
            let fd_name = func_names_map.get(fd_write_index as u32).unwrap();
            func_names_map.insert(func_names_map.len() as u32, fd_name.clone());
            func_names_map.insert(fd_write_index as u32, "main.debugLogCp".to_string());
        }

        //插入fd_write到function部分
        {
            let fd_write_type_index = self
                .get_fd_write_type()
                .or_else(|| {
                    let types = self.module.type_section_mut().unwrap().types_mut();
                    types.push(Type::Function(FunctionType::new(
                        vec![
                            ValueType::I32,
                            ValueType::I32,
                            ValueType::I32,
                            ValueType::I32,
                        ],
                        vec![ValueType::I32],
                    )));
                    Some(types.len() as u32)
                })
                .unwrap();
            let functions = self.module.function_section_mut().unwrap().entries_mut();
            functions.push(Func::new(fd_write_type_index));
        }

        //编译所有的函数体，改变其中所有call/callindirect指令参数。
        {
            let new_fd_write_index = self.module.functions_space() - 1;
            if let Some(m) = self.module.code_section_mut() {
                for body in m.bodies_mut() {
                    for ins in body.code_mut().elements_mut() {
                        if let Instruction::Call(func_index) = ins {
                            if *func_index == fd_write_index as u32 {
                                *func_index = new_fd_write_index as u32;
                            }
                        }
                    }
                }
            }
        }

        //插入fd_write函数体到function body部分。
        {
            /*
            func fd_write(id uint32, iovs *__wasi_iovec_t, iovs_len uint, nwritten *uint)  uint {
                //only support println in fvm
                errno := debugLog(uintptr(iovs.buf), uint32(iovs.bufLen) )
                return uint(errno)
            }

             (func $runtime.fd_write (type 1) (param i32 i32 i32 i32) (result i32)
                block  ;; label = @1
                  local.get 1
                  br_if 0 (;@1;)
                  unreachable
                  unreachable
                end
                local.get 1
                i32.load
                local.get 1
                i32.load offset=4
                call $main.debugLog)
             */
            let new_insert_debug_index = self
                .get_import_func_index("debug", "log")
                .expect("unable to get debug log") as u32;
            let codes = self.module.code_section_mut().unwrap().bodies_mut();
            let fd_write_code = FuncBody::new(
                vec![],
                Instructions::new(vec![
                    Instruction::Block(BlockType::NoResult),
                    Instruction::GetLocal(1),
                    Instruction::BrIf(0),
                    Instruction::Unreachable,
                    Instruction::Unreachable,
                    Instruction::End,
                    Instruction::GetLocal(1),
                    Instruction::I32Load(2, 0), //why -1
                    Instruction::GetLocal(1),
                    Instruction::I32Load(2, 4),
                    Instruction::Call(new_insert_debug_index),
                ]),
            );
            codes.push(fd_write_code);
        }

        Ok(self)
    }

    pub fn append_init_to_invoke(&mut self) -> Result<&mut Self> {
        let import_func_count = self.module.import_count(ImportCountType::Function);
        if let Some(invoke_index) = self.get_func_index("invoke") {
            if let Some(start_func_index) = self.get_func_index("_start") {
                let invoke_body: &mut FuncBody = self
                    .module
                    .code_section_mut()
                    .unwrap()
                    .bodies_mut()
                    .get_mut(invoke_index)
                    .unwrap();
                invoke_body.code_mut().elements_mut().insert(
                    0,
                    Instruction::Call(start_func_index as u32 + import_func_count as u32),
                );
            }
            Ok(self)
        } else {
            Err(anyhow!("unable to find invoke function"))
        }
    }

    fn get_import_func_index(&self, module: &str, filed: &str) -> Option<usize> {
        if let Some(import_section) = self.module.import_section() {
            for (index, import) in import_section.entries().iter().enumerate() {
                if import.module() == module && import.field() == filed {
                    return Some(index);
                }
            }
        }
        None
    }

    fn get_func_index(&self, func_name: &str) -> Option<usize> {
        let import_func_count = self.module.import_count(ImportCountType::Function);
        let names_map = self
            .module
            .names_section()
            .unwrap()
            .functions()
            .unwrap()
            .names();
        if let Some(function_section) = self.module.function_section() {
            for (index, _) in function_section.entries().iter().enumerate() {
                if let Some(name) = names_map.get(index as u32 + import_func_count as u32) {
                    if name == func_name {
                        return Some(index);
                    }
                }
            }
        }
        None
    }

    // (import "wasi_snapshot_preview1" "fd_write" (func $runtime.fd_write (type 6)))
    fn has_fd_write(&self) -> bool {
        if let Some(import_section) = self.module.import_section() {
            for import in import_section.entries() {
                if import.module() == "wasi_snapshot_preview1" && import.field() == "fd_write" {
                    return true;
                }
            }
        }
        false
    }

    //(import "debug" "log" (func $main.debugLog (type 0)))
    #[allow(dead_code)]
    fn has_debug_import(&self) -> bool {
        if let Some(import_section) = self.module.import_section() {
            for import in import_section.entries() {
                if import.module() == "debug" && import.field() == "log" {
                    return true;
                }
            }
        }
        false
    }

    // (type (;0;) (func (param i32 i32) (result i32)))
    // (import "debug" "log" (func $main.debugLog (type 0)))
    fn get_debug_type(&self) -> Option<u32> {
        if let Some(type_section) = self.module.type_section() {
            for (i, wtype) in type_section.types().iter().enumerate() {
                match wtype {
                    Function(f) => {
                        if f.params().len() == 2
                            && f.params()[0] == ValueType::I32
                            && f.params()[1] == ValueType::I32
                            && f.results().len() == 1
                            && f.results()[0] == ValueType::I32
                        {
                            return Some(i as u32);
                        }
                    }
                }
            }
        }
        None
    }

    //(type (;1;) (func (param i32 i32 i32 i32) (result i32)))
    fn get_fd_write_type(&self) -> Option<u32> {
        if let Some(type_section) = self.module.type_section() {
            for (i, wtype) in type_section.types().iter().enumerate() {
                match wtype {
                    Function(f) => {
                        if f.params().len() == 4
                            && f.params()[0] == ValueType::I32
                            && f.params()[1] == ValueType::I32
                            && f.params()[0] == ValueType::I32
                            && f.params()[1] == ValueType::I32
                            && f.results().len() == 1
                            && f.results()[0] == ValueType::I32
                        {
                            return Some(i as u32);
                        }
                    }
                }
            }
        }
        None
    }
}

trait TryString {
    fn try_to_string(&self) -> Result<String>;
}

impl TryString for PathBuf {
    fn try_to_string(&self) -> Result<String> {
        self.to_str()
            .ok_or_else(|| anyhow!("unable to get string from pathbuf"))
            .map(|v| v.to_string())
    }
}

#[cfg(test)]
mod tests {
    use crate::wasmprocess::{BuildCliConfig, BuildOptions};
    use std::path::Path;

    #[test]
    fn no_input_output() {
        let cli_cfg = BuildCliConfig {
            input: None,
            output: None,
            wat: false,
            debug: false,
        };
        let build_opt =
            BuildOptions::new(Path::new("/foo").to_path_buf(), &cli_cfg).expect("build opt");
        assert_eq!(build_opt.code_path, "/foo");
        assert_eq!(build_opt.target_name, "foo");
        assert_eq!(build_opt.output_wat_path, "/foo/foo.wat");
        assert_eq!(build_opt.output_wasm_path, "/foo/foo.wasm");
        assert_eq!(build_opt.target_dir, "/foo");
    }

    #[test]
    fn abs_input_no_output() {
        let cli_cfg = BuildCliConfig {
            input: Some("/ggg".to_owned()),
            output: None,
            wat: false,
            debug: false,
        };
        let build_opt =
            BuildOptions::new(Path::new("/foo").to_path_buf(), &cli_cfg).expect("build opt");
        assert_eq!(build_opt.code_path, "/ggg");
        assert_eq!(build_opt.target_name, "ggg");
        assert_eq!(build_opt.output_wat_path, "/foo/ggg.wat");
        assert_eq!(build_opt.output_wasm_path, "/foo/ggg.wasm");
        assert_eq!(build_opt.target_dir, "/foo");
    }

    #[test]
    fn rel_input_no_output() {
        let cli_cfg = BuildCliConfig {
            input: Some("../mm".to_owned()),
            output: None,
            wat: false,
            debug: false,
        };
        let build_opt =
            BuildOptions::new(Path::new("/foo").to_path_buf(), &cli_cfg).expect("build opt");
        assert_eq!(build_opt.code_path, "/mm");
        assert_eq!(build_opt.target_name, "mm");
        assert_eq!(build_opt.target_dir, "/foo");
        assert_eq!(build_opt.output_wat_path, "/foo/mm.wat");
        assert_eq!(build_opt.output_wasm_path, "/foo/mm.wasm");
    }

    #[test]
    fn go_file_input_no_output() {
        let cli_cfg = BuildCliConfig {
            input: Some("../mm/nn.go".to_owned()),
            output: None,
            wat: false,
            debug: false,
        };
        let build_opt =
            BuildOptions::new(Path::new("/foo").to_path_buf(), &cli_cfg).expect("build opt");
        assert_eq!(build_opt.code_path, "/mm/nn.go");
        assert_eq!(build_opt.target_name, "nn");
        assert_eq!(build_opt.target_dir, "/foo");
        assert_eq!(build_opt.output_wat_path, "/foo/nn.wat");
        assert_eq!(build_opt.output_wasm_path, "/foo/nn.wasm");
    }

    #[test]
    fn no_input_abs_output() {
        let cli_cfg = BuildCliConfig {
            input: None,
            output: Some("/mmm".to_owned()),
            wat: false,
            debug: false,
        };
        let build_opt =
            BuildOptions::new(Path::new("/foo").to_path_buf(), &cli_cfg).expect("build opt");
        assert_eq!(build_opt.code_path, "/foo");
        assert_eq!(build_opt.target_name, "foo");
        assert_eq!(build_opt.target_dir, "/mmm");
        assert_eq!(build_opt.output_wat_path, "/mmm/foo.wat");
        assert_eq!(build_opt.output_wasm_path, "/mmm/foo.wasm");
    }

    #[test]
    fn no_input_rel_output() {
        let cli_cfg = BuildCliConfig {
            input: None,
            output: Some("../mmm".to_owned()),
            wat: false,
            debug: false,
        };
        let build_opt =
            BuildOptions::new(Path::new("/foo/ppp").to_path_buf(), &cli_cfg).expect("build opt");
        assert_eq!(build_opt.code_path, "/foo/ppp");
        assert_eq!(build_opt.target_name, "ppp");
        assert_eq!(build_opt.target_dir, "/foo/mmm");
        assert_eq!(build_opt.output_wat_path, "/foo/mmm/ppp.wat");
        assert_eq!(build_opt.output_wasm_path, "/foo/mmm/ppp.wasm");
    }

    #[test]
    fn no_input_go_output() {
        let cli_cfg = BuildCliConfig {
            input: None,
            output: Some("../mmm/main.go".to_owned()),
            wat: false,
            debug: false,
        };
        let build_opt =
            BuildOptions::new(Path::new("/foo/ppp").to_path_buf(), &cli_cfg).expect("build opt");
        assert_eq!(build_opt.code_path, "/foo/ppp");
        assert_eq!(build_opt.target_name, "main");
        assert_eq!(build_opt.target_dir, "/foo/mmm");
        assert_eq!(build_opt.output_wat_path, "/foo/mmm/main.wat");
        assert_eq!(build_opt.output_wasm_path, "/foo/mmm/main.wasm");
    }

    #[test]
    fn abs_input_rel_output() {
        let cli_cfg = BuildCliConfig {
            input: Some("/lll".to_owned()),
            output: Some("../mmm".to_owned()),
            wat: false,
            debug: false,
        };
        let build_opt =
            BuildOptions::new(Path::new("/foo/ppp").to_path_buf(), &cli_cfg).expect("build opt");
        assert_eq!(build_opt.code_path, "/lll");
        assert_eq!(build_opt.target_name, "lll");
        assert_eq!(build_opt.target_dir, "/foo/mmm");
        assert_eq!(build_opt.output_wat_path, "/foo/mmm/lll.wat");
        assert_eq!(build_opt.output_wasm_path, "/foo/mmm/lll.wasm");
    }

    #[test]
    fn rel_input_rel_output() {
        let cli_cfg = BuildCliConfig {
            input: Some("../lll".to_owned()),
            output: Some("../mmm".to_owned()),
            wat: false,
            debug: false,
        };
        let build_opt =
            BuildOptions::new(Path::new("/foo/ppp").to_path_buf(), &cli_cfg).expect("build opt");
        assert_eq!(build_opt.code_path, "/foo/lll");
        assert_eq!(build_opt.target_name, "lll");
        assert_eq!(build_opt.target_dir, "/foo/mmm");
        assert_eq!(build_opt.output_wat_path, "/foo/mmm/lll.wat");
        assert_eq!(build_opt.output_wasm_path, "/foo/mmm/lll.wasm");
    }

    #[test]
    fn rel_go_input_rel_output() {
        let cli_cfg = BuildCliConfig {
            input: Some("../lll/main.go".to_owned()),
            output: Some("../mmm".to_owned()),
            wat: false,
            debug: false,
        };
        let build_opt =
            BuildOptions::new(Path::new("/foo/ppp").to_path_buf(), &cli_cfg).expect("build opt");
        assert_eq!(build_opt.code_path, "/foo/lll/main.go");
        assert_eq!(build_opt.target_name, "main");
        assert_eq!(build_opt.target_dir, "/foo/mmm");
        assert_eq!(build_opt.output_wat_path, "/foo/mmm/main.wat");
        assert_eq!(build_opt.output_wasm_path, "/foo/mmm/main.wasm");
    }
}
