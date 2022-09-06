build:
ifeq (,$(wildcard ${current_dir}/bin/go-fvm-sdk-tools))
	BUILD_FIL_NETWORK=devnet-wasm cargo build -p go-fvm-sdk-tools --release
endif

build-example:
	cd examples/erc20/client_example && go build
	cd examples/hellocontract/client_example && go build

install:build
	cp -f ./target/release/go-fvm-sdk-tools /usr/local/bin

code-gen:
	cd ./sdk/gen && go run main.go
	cd ./sdk/cases && ./gen.sh
	cd ./examples/hellocontract/gen && go run main.go
	cd ./examples/hellocontract && go-fvm-sdk-tools build
	cd ./examples/erc20/gen && go run main.go
	cd ./examples/erc20 && go-fvm-sdk-tools build

clean:
	rm -rf ./bin/*

lint:
	golangci-lint run
	cargo run -p ci -- lint

test: build code-gen
	cd ./sdk/cases && go-fvm-sdk-tools test
	cd ./examples/hellocontract && go-fvm-sdk-tools test
	cd ./examples/erc20 && go-fvm-sdk-tools test
	cd ./examples/erc20/contract && go test -tags simulate 

