//go:build !simulate
// +build !simulate

// Package sys ...
package sys

import "context"

func Abort(_ context.Context, code uint32, msg string) {
	strPtr, strLen := GetStringPointerAndLen(msg)
	_ = vmAbort(code, strPtr, strLen)
}
