#!/bin/bash

root_dir=$1
echo ${root_dir}
for i in *.go; do
    filename=$(echo $i | cut -d . -f1)
    target="${filename}.wasm"
    ${root_dir}/bin/go-fvm-sdk-tools build --output $target --wat -- $i
    echo "Build ${target} successfully"
done
