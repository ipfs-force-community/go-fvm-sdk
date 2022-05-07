# Go-Fvm-SDK

Convert go contract to wasm and help wasm code run on fvm

1. convert go code to wasm
2. implement fvm system call interface
3. pure modify library for running on fvm

## dependency

### install modified tinygo

modify fmt to remove os refrence

```go
git clone git@github.com:ipfs-force-community/tinygo.git
git checkout fvm
```

[build doc](https://tinygo.org/docs/guides/build/)

## dependency analyze

```
go list -json -deps > deps.json 
```

### build and run hello-world contract
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
