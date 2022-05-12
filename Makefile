build-testing:
	cd ./sdk/testing && cargo build --release

gen-case:
	cd ./sdk/cases && ./gen.sh

test: build-testing gen-case
	./sdk/testing/target/release/testing --path ./sdk/cases