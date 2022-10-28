//go:build !simulate
// +build !simulate

package sys

import (
	"context"
	"unsafe"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
)

func Enabled(_ context.Context) (bool, error) {
	var result int32
	code := debugEnabled(uintptr(unsafe.Pointer(&result)))
	if code != 0 {
		return false, ferrors.NewFvmError(ferrors.ExitCode(code), "failed to get debug-enabled")
	}

	return result == 0, nil
}

func Log(_ context.Context, msg string) error {
	msgBufPtr, msgBufLen := GetStringPointerAndLen(msg)
	code := debugLog(msgBufPtr, msgBufLen)
	if code != 0 {
		return ferrors.NewFvmError(ferrors.ExitCode(code), "failed to record debug log")
	}

	return nil
}

func StoreArtifact(_ context.Context, name string, data string) error {
	nameBufPtr, nameBufLen := GetStringPointerAndLen(name)
	dataBufPtr, dataBufLen := GetStringPointerAndLen(data)
	code := debugStoreArtifact(nameBufPtr, nameBufLen, dataBufPtr, dataBufLen)
	if code != 0 {
		return ferrors.NewFvmError(ferrors.ExitCode(code), "failed to record debug log")
	}

	return nil
}
