package sys

/// Abort execution with the given code and message. The code is recorded in the receipt, the
/// message is for debugging only.
///
/// # Errors
///
/// None. This function doesn't return.
//go:wasm-module vm
//export abort
func vmAbort(code uint32, msg uintptr, msgLen uint32) uint32

func Abort(code uint32, msg string) {
	strPtr, strLen := GetStringPointerAndLen(msg)
	_ = vmAbort(code, strPtr, strLen)
}
