//go:build !simulated
// +build !simulated

package sys

import (
	"context"
	"unsafe"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
)

func Enabled(ctx context.Context) (bool, error) {
	if env, ok := isSimulatedEnv(ctx); ok {
		return env.Enabled()
	}

	var result int32
	code := debugEnabled(uintptr(unsafe.Pointer(&result)))
	if code != 0 {
		return false, ferrors.NewFvmError(ferrors.ExitCode(code), "failed to get debug-enabled")
	}

	return result == 0, nil
}

func Log(ctx context.Context, msg string) error {
	if env, ok := isSimulatedEnv(ctx); ok {
		return env.Log(msg)
	}
	msgBufPtr, msgBufLen := GetStringPointerAndLen(msg)
	code := debugLog(msgBufPtr, msgBufLen)
	if code != 0 {
		return ferrors.NewFvmError(ferrors.ExitCode(code), "failed to record debug log")
	}

	return nil
}
