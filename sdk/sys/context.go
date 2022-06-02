package sys

import (
	"unsafe"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
)

/// Returns the details about this invocation.
///
/// # Errors
///
/// None
//go:wasm-module vm
//export context
func vmContext(ret uintptr) uint32

func VmContext() (*types.InvocationContext, error) {
	var result types.InvocationContext
	code := vmContext(uintptr(unsafe.Pointer(&result)))
	if code != 0 {
		return nil, ferrors.NewFvmError(ferrors.ExitCode(code), "unable to get invocation context")
	}
	return &result, nil
}
