# go-fvm-sdk

`go-fvm-sdk` enables GO developers to write [FVM](https://fvm.filecoin.io/) contracts in TinyGo. To learn more about FVM, please check out resources on [fvm forum](https://fvm-forum.filecoin.io/).

Key features of `go-fvm-sdk` are...

- Compiles GO contracts into WASM
- Full compatibility with FVM system call interface
- unit/integration test support
- Deploy native WASM actor on FVM!

## Why choose TinyGO
- Support most Go syntax and libraries. 
- Smaller footprint which is great for use cases on the blockchain.
- Nimble management of GC and Memory, which is great for contracts with  runtime constraints of tens to hundreds of milliseconds and memory constraints of a few megabytes.

## Install

Install [Go](https://go.dev/doc/install)v1.17 or above and [TinyGo](https://tinygo.org/getting-started/install/)v0.25 or above.

*Note: latest Go or TinyGo versions may not be tested.*

### Use Binary

Go to [releases](https://github.com/ipfs-force-community/go-fvm-sdk/releases) and download the executables.

Rename your executable to `fvm_go_sdk` and add execution permission.

```bash
$ mv fvm_go_sdk_vX.Y.Z_xxx fvm_go_sdk
$ chmod +x fvm_go_sdk
```

### Build from source

Install [WASM tools](https://github.com/WebAssembly/binaryen) and [Rust](https://www.rust-lang.org/tools/install).

Clone `fvm_go_sdk` repo and make.

```bash
$ git clone git@github.com:ipfs-force-community/go-fvm-sdk.git
$ make
```

*Note: it may take a while before you finish buiding from source.*

Check if `fvm_go_sdk` is operational.

```bash
$ fvm_go_sdk -h
```

Once you have `fvm_go_sdk` ready, `Go` and `TinyGo` need to be patched in order to have the compiled WASM contract compatible with FVM.

```bash
$ fvm_go_sdk patch
```

Successful patching of `go-fvm-sdk` will give the following output. And you are done with installation! 

```
patching file src/reflect/value.go
Hunk #1 succeeded at 754 (offset 3 lines).
patching file targets/wasi.json
Hunk #2 succeeded at 10 with fuzz 1.
```

## Create your first Contract!

It is recommended that you create your first project using our template. 

```bash
$ fvm_go_sdk new -- <YOUR_ACTOR_NAME>
```

`fvm_go_sdk` will generate the following scaffolding for you.

```bash
.
├── README.md 
├── actor
│   ├── actor.go # Write your contract here
│   └── cbor_gen.go
├── client
│   └── client.go # client for installation and deployment of the contract
├── entry.go # Main entry point of the contract. If you change interfaces of your actor, please build again to renew the entry point file
├── gen
│   ├── go.mod
│   ├── go.sum
│   └── main.go
├── go.mod
├── go.sum
├── <YOUR_ACTOR_NAME>.wasm # WASM contract compiled from actor.go
└── test.json
```

Compile your contract. And you now have your very first contract ready to be depolyed!

```bash
$ fvm_go_sdk build
```

## Deploy your Contract

Follow this [instruction](https://lotus.filecoin.io/developers/local-network/) to setup local devnet with fvm branch `experimental/fvm-m2`. 

Install actor.

```bash
$ lotus chain install-actor <path-to-wasm-bytecode>
```

Instantiate actor.

```bash
$ lotus chain create-actor <code-cid> <encoded-params>
```

Invoke actor.

```bash
$ lotus chain invoke <address> <method_num>
```

## Tips

- Do not use ASM code.
- Avoid using the reflect library. If you have to, make sure that `TinyGo` supports it.
- Use `go list` command to check actor dependency issue.

```bash
$ go list -json -deps > deps.json 
```