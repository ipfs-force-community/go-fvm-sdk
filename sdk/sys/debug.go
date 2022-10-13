//go:build !simulate
// +build !simulate

package sys

import (
	"context"
	"unsafe"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
)

func Enabled(ctx context.Context) (bool, error) {
	var result int32
	code := debugEnabled(uintptr(unsafe.Pointer(&result)))
	if code != 0 {
		return false, ferrors.NewFvmError(ferrors.ExitCode(code), "failed to get debug-enabled")
	}

	return result == 0, nil
}

func Log(ctx context.Context, msg string) error {
	msgBufPtr, msgBufLen := GetStringPointerAndLen(msg)
	code := debugLog(msgBufPtr, msgBufLen)
	if code != 0 {
		return ferrors.NewFvmError(ferrors.ExitCode(code), "failed to record debug log")
	}

	return nil
}
