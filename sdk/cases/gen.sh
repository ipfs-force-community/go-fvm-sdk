#!/bin/bash

for i in *.go; do
    filename=$(echo $i | cut -d . -f1)
    target="${filename}.wasm"
    fvm_go_sdk build --output $target --wat -- $i
    echo "Build ${target} successfully"
done
