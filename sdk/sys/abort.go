//go:build !simulated
// +build !simulated

// Package sys ...
package sys

import "context"

func Abort(ctx context.Context, code uint32, msg string) {
	strPtr, strLen := GetStringPointerAndLen(msg)
	_ = vmAbort(code, strPtr, strLen)
}
