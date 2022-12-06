//go:build tinygo.wasm
// +build tinygo.wasm

package sys

// Returns the details about this invocation.
// /
// # Errors
// /
// None
//
//go:wasm-module vm
//export message_context
func vmMessageContext(ret uintptr) uint32

// Abort execution with the given code and message. The code is recorded in the receipt, the
// message is for debugging only.
// /
// # Errors
// /
// None. This function doesn't return.
//
//go:wasm-module vm
//export exit
func vmExit(code uint32, blkId uint32, msgOff uintptr, msgLen uint32) uint32
