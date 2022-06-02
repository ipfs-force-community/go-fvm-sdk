current_dir = $(shell pwd)

build-tool:
ifeq (,$(wildcard ${current_dir}/bin/go-fvm-sdk-tools))
	cd ./tools && cargo build --out-dir ${current_dir}/bin --release -Z unstable-options
endif
	
build: build-tool

gen-case:
	cd ./sdk/cases && ./gen.sh ${current_dir}

test: build gen-case
	${current_dir}/bin/go-fvm-sdk-tools test --path ./sdk/cases
