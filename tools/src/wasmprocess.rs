use anyhow::{anyhow, Result};
use parity_wasm::elements::Type::Function;
use parity_wasm::elements::{BlockType, External, Func, FuncBody, FunctionType, ImportCountType, ImportEntry, Instruction, Instructions, Internal, Module, Type, ValueType};
use std::collections::HashMap;
use clap::Parser;

#[derive(Parser, Debug)]
#[clap(author, version, about, long_about = None)]
pub struct ProcessConfig {
    #[clap(short, long)]
    pub  input: String,
    #[clap(short, long)]
    pub output: String,
    #[clap(short, long)]
    pub wat: bool,
}

pub fn run_process(cfg: &ProcessConfig) -> Result<()> {
    let result = GoFvmBinProcessor::new(cfg.input.as_str())
        .append_init_to_invoke()?
        // .replace_fd_write()?
        .get_binary()?;

    let wat_str = wasmprinter::print_bytes(result)?;
    if cfg.wat {

        std::fs::write(
            cfg.output.clone() + ".wat",
            wat_str.clone(),
        )?;
    }


    let mut features = wabt::Features::new();
    features.set_annotations_enabled(true);
    let wat_bin = wabt::wat2wasm_with_features(wat_str, features)?;
    std::fs::write(
        cfg.output.clone(),
        wat_bin,
    )?;
    Ok(())
}

pub struct GoFvmBinProcessor {
    module: Module,
}

impl GoFvmBinProcessor {
    pub fn new(path: &str) -> Self {
        let module = parity_wasm::deserialize_file(path)
            .expect("Should be deserialized")
            .parse_names()
            .expect("Names to be parsed");
        GoFvmBinProcessor { module }
    }

    pub fn get_binary(&self) -> Result<Vec<u8>> {
        return parity_wasm::serialize(self.module.clone())
            .map_err(|e| anyhow!("convert module to binary {}", e));
    }

    #[allow(dead_code)]
    //暂不可用，因为这里需要无法替换call_indirect函数的参数，动态调用无法做。暂时搁置，只能在底层直接使用黑科技改了。以后需要从go的ir层面进行修改才行，或许可以使用nestmodule来修改。
    pub fn replace_fd_write(&mut self) -> Result<&mut Self> {
        //探测需不需要插入fd_write
        if !self.has_fd_write() {
            return Ok(self);
        }

        let import_func_count = self.module.import_count(ImportCountType::Function);
        let func_count = self.module.function_section().unwrap().entries().len();
        let fd_write_index = self
            .get_import_func_index("wasi_snapshot_preview1", "fd_write")
            .expect("unable to find fs_write inport");
        //探测debug log是否已经存在，
        let has_debug_import = self.has_debug_import();
        let mut namevec = vec![];
        let mut func_index_map = HashMap::new();
        {
            let mut debug_offset: i32 = 0;
            let mut fd_write_offset: i32 = 0;
            let func_names_map = self
                .module
                .names_section_mut()
                .unwrap()
                .functions_mut()
                .as_mut()
                .unwrap()
                .names_mut();
            for index in 0..(import_func_count + func_count) {
                if index == fd_write_index {
                    fd_write_offset = -1;
                }


                func_index_map.insert(
                    index,
                    (index as i32 + fd_write_offset + debug_offset) as usize,
                );

                if index == fd_write_index {
                    let new_fd_write_index = if has_debug_import {
                        //（l1+l2)-fd
                        (import_func_count + func_count - 1) as i32
                    } else {
                        //老的调用这个位置的，全部调用到末尾的函数
                        (import_func_count + func_count) as i32
                    };
                    func_index_map.insert(index, new_fd_write_index as usize);
                }

                if !has_debug_import&&index == import_func_count -1 {
                    debug_offset = 1;
                    namevec.push("main.debugLog".to_string());
                    continue
                }

                if index != fd_write_index {
                    namevec.push(func_names_map.get(index as u32).unwrap().clone());
                }
            }
            namevec.push("runtime.fd_write".to_owned());
        }

        {
            //删除fd_write import 部分
            let imports = self.module.import_section_mut().unwrap().entries_mut();
            imports.remove(fd_write_index);
        }

        //插入debug log
        {
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
            if !has_debug_import {
                //插入debug log到import部分和function部分
                let imports = self.module.import_section_mut().unwrap().entries_mut();
                imports.push(ImportEntry::new(
                    "debug".to_owned(),
                    "log".to_owned(),
                    External::Function(debug_type_index),
                ));
            }
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
            for m in self.module.code_section_mut() {
                for body in m.bodies_mut() {
                    for ins in body.code_mut().elements_mut() {
                        match ins {
                            Instruction::Call(func_index) => {
                                if let Some(new_func_index) =
                                func_index_map.get(&(*func_index as usize))
                                {
                                    *func_index = *new_func_index as u32;
                                }
                            }
                            Instruction::CallIndirect(typer, index) => {
                                //unable todo
                                println!("todo support CallIndirect {}  {}", typer, index);
                            }
                            _ => {}
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
            let new_insert_debug_index =  self.get_import_func_index("debug", "log")
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
                    Instruction::Call(new_insert_debug_index as u32),
                ]),
            );
            codes.push(fd_write_code)
        }

        //重建namemap
        {
            let func_names_map = self
                .module
                .names_section_mut()
                .unwrap()
                .functions_mut()
                .as_mut()
                .unwrap()
                .names_mut();
            func_names_map.clear();
            for (i, val) in namevec.iter_mut().enumerate() {
                func_names_map.insert(i as u32, val.clone());
            }
        }

        //重建export
        {
            let  exports = self.module.export_section_mut().expect("unable to get export section");
            for export in  exports.entries_mut() {
                if let Internal::Function(func_index_ref) = export.internal_mut() {
                    if let Some(new_func_index) =
                    func_index_map.get(&(*func_index_ref as usize))
                    {
                        *func_index_ref = *new_func_index as u32;
                    }
                }
            }
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
                )
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
        return None;
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
        return None;
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
        return false;
    }

    //(import "debug" "log" (func $main.debugLog (type 0)))
    fn has_debug_import(&self) -> bool {
        if let Some(import_section) = self.module.import_section() {
            for import in import_section.entries() {
                if import.module() == "debug" && import.field() == "log" {
                    return true;
                }
            }
        }
        return false;
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
        return None;
    }

    //(type (;1;) (func (param i32 i32 i32 i32) (result i32)))
    pub fn get_fd_write_type(&self) -> Option<u32> {
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
        return None;
    }
}
