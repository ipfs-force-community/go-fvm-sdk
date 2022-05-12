package sys

import (
	"fmt"
	"unsafe"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
)

/// Returns if we're in debug mode. A zero or positive return value means
/// yes, a negative return value means no.
//go:wasm-module debug
//export enabled
func debugEnabled(ret uintptr) uint32

/// Logs a message on the node.
//go:wasm-module debug
//export log
func debugLog(message uintptr, message_len uint32) uint32

func Enabled() (bool, error) {
	var result int32
	code := debugEnabled(uintptr(unsafe.Pointer(&result)))
	if code != 0 {
		return false, ferrors.NewFvmError(ferrors.ExitCode(code), fmt.Sprintf("failed to get debug-enabled"))
	}

	return result >= 0, nil
}

func Log(msg string) error {
	msgBufPtr, msgBufLen := GetStringPointerAndLen(msg)
	code := debugLog(msgBufPtr, msgBufLen)
	if code != 0 {
		return ferrors.NewFvmError(ferrors.ExitCode(code), fmt.Sprintf("failed to record debug log"))
	}

	return nil
}
