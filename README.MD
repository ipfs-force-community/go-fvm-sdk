# Go-Fvm-SDK

Convert go contract to wasm and help wasm code run on fvm

1. convert go code to wasm
2. implement fvm system call interface
3. pure modify library for running on fvm

## Install

1. download the latest release version from [wasm tools](https://github.com/WebAssembly/binaryen), and configure the corresponding bin path in the $PATH
2. install [go](https://go.dev/doc/install)  require v1.17+
3. install [tinygo](https://tinygo.org/getting-started/install/) require v0.26+
4. build go-fvm-sdk tool, tool will in path ``` target/release/go-fvm-sdk-tools```
  ```
  git clone git@github.com:ipfs-force-community/go-fvm-sdk.git
  make
  ```
5. use go-fvm-sdk-tools tool to patch your environment(go/tinygo) 
  ```
    go-fvm-sdk patch
  ```

  this log meaning operate successfully
  ```log
  patching file src/reflect/value.go
  Hunk #1 succeeded at 754 (offset 3 lines).
  patching file targets/wasi.json
  patch unexpectedly ends in middle of line
  Hunk #2 succeeded at 10 with fuzz 1.
  ```

### quick start this hello world actor
* Compile hellocontract to Wasm
    ```
        git clone git@github.com:ipfs-force-community/go-fvm-sdk.git
        cd hellocontract
        tinygo build -target wasi -no-debug -panic trap -scheduler none -o helloworld.wasm ./helloworld.go
    ```
* Set up a Lotus devnet
  * Experimental branch: experimental/fvm-m2
  * Instructions: https://lotus.filecoin.io/developers/local-network/
* Install Actor
  * lotus chain install-actor <path-to-wasm-bytecode>
* Instantiate Actor
  * lotus chain create-actor <code-cid> <encoded-params>
* Invoke actor
  * lotus chain invoke <address> <method_num>

## Suggestion

1. Do not use asm code
2. Try not to use the reflect library, if you must use it, make sure that tinygo supports it.
3. use go list command to check actor dependency graph, check which package
    ```
    go list -json -deps > deps.json 
    ```