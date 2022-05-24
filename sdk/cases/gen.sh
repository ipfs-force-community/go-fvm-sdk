#/bin/bash.

root_dir=$1
echo ${root_dir}
for i in *.go; do
    filename=$(echo $i | cut -d . -f1)
    target="${filename}.wasm"
    tinygo build -target wasi -no-debug -panic trap -o $target $i
    ${root_dir}/bin/go-fvm-sdk-tools --input $target --output $target --wat
done
