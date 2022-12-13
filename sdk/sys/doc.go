// Package sys is used to interact with fvm
// and in this package non-fvm file used to index for ide, this can be convenient for development
// file with _call suffix defines the wasm export function
// file with _simulate used to simulate system call, useful for actor unit test
// others file was wrap for _call file, make the type and style more closely match the standard of go
package sys
