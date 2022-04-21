# Go-Fvm-SDK

Convert go contract to wasm and help wasm code run on fvm

1. convert go code to wasm
2. implement fvm system call interface
3. pure modify library for running on fvm

## dependency

### modify tinygo

modify fmt to remove os refrence

```go
git clone git@github.com:ipfs-force-community/tinygo.git
git checkout venus
```

[build doc](https://tinygo.org/docs/guides/build/)