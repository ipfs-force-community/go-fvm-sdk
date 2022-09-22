//go:build !simulated
// +build !simulated

package sys

func Abort(code uint32, msg string) {
	strPtr, strLen := GetStringPointerAndLen(msg)
	_ = vmAbort(code, strPtr, strLen)
}
