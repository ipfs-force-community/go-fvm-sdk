#/bin/bash.

for i in *.go; do
    filename=$(echo $i | cut -d . -f1)
    target="${filename}.wasm"
    tinygo build -target wasi -no-debug -panic trap -o $target $i
done
