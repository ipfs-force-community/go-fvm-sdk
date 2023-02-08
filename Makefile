build:
ifeq (,$(wildcard ${current_dir}/bin/fvm_go_sdk))
	BUILD_FIL_NETWORK=devnet-wasm cargo build -p fvm_go_sdk --release
endif

build-client-example:
	cd examples/erc20/client_example && go build
	cd examples/frc46token/client_example && go build
	cd examples/hellocontract/client_example && go build

install:build
	cp -f ./target/release/fvm_go_sdk /usr/local/bin

code-gen:
	cd ./sdk/gen && go run main.go
	cd ./sdk/cases && ./gen.sh
	cd ./examples/hellocontract/gen && go run main.go
	cd ./examples/hellocontract && fvm_go_sdk build
	cd ./examples/erc20/gen && go run main.go
	cd ./examples/erc20 && fvm_go_sdk build
	cd ./examples/frc46token/gen && go run main.go
	cd ./examples/frc46token && fvm_go_sdk build

clean:
	rm -rf ./bin/*

lint:
	golangci-lint run
	cargo run -p ci -- lint

test: build code-gen
	cd ./sdk/cases && fvm_go_sdk test
	cd ./examples/hellocontract && fvm_go_sdk test
	cd ./examples/erc20 && fvm_go_sdk test
	cd ./examples/erc20/contract && go test --tags simulate
	cd ./examples/frc46token && fvm_go_sdk test
	cd ./examples/frc46token/contract && go test --tags simulate
	cd ./examples/hellocontract/contract && go test --tags simulate


check: code-gen lint build-client-example test
