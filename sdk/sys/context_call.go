//go:build tinygo.wasm
// +build tinygo.wasm

package sys

// / Returns the details about this invocation.
// /
// / # Errors
// /
// / None
//
//go:wasm-module vm
//export context
func vmContext(ret uintptr) uint32
