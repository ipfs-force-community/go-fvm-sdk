current_dir = $(shell pwd)

build-tool:
	cd ./tools/go-fvm-sdk-tools && cargo build --out-dir ${current_dir}/bin --release -Z unstable-options

build-testing:
	cd ./tools/testing && cargo build --out-dir ${current_dir}/bin --release -Z unstable-options

build: build-tool build-testing

gen-case:
	cd ./sdk/cases && ./gen.sh ${current_dir}

test: build gen-case
	./tools/testing/target/release/testing --path ./sdk/cases
