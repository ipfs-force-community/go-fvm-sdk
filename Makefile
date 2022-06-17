current_dir = $(shell pwd)

build-tool:
ifeq (,$(wildcard ${current_dir}/bin/go-fvm-sdk-tools))
	cargo build -p go-fvm-sdk-tools --out-dir ${current_dir}/bin --release -Z unstable-options
endif
	
build: build-tool

gen:
	cd ./sdk/gen && go run main.go

clean:
	rm -rf ./bin/*

gen-case:
	cd ./sdk/cases && ./gen.sh ${current_dir}

test: build gen-case
	${current_dir}/bin/go-fvm-sdk-tools test -- ./sdk/cases
