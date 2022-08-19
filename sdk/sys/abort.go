//go:build !simulate
//+build !tinygo.wasm
package sys

func Abort(code uint32, msg string) {
	strPtr, strLen := GetStringPointerAndLen(msg)
	_ = vmAbort(code, strPtr, strLen)
}
