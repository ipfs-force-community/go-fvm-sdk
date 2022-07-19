#!/bin/bash

for i in *.go; do
    filename=$(echo $i | cut -d . -f1)
    target="${filename}.wasm"
    go-fvm-sdk-tools build --output $target --wat -- $i
    echo "Build ${target} successfully"
done
