//go:build !simulated
// +build !simulated

package sys

import (
	"unsafe"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
)

func VMContext() (*types.InvocationContext, error) {
	var result types.InvocationContext
	code := vmContext(uintptr(unsafe.Pointer(&result)))
	if code != 0 {
		return nil, ferrors.NewFvmError(ferrors.ExitCode(code), "unable to get invocation context")
	}
	return &result, nil
}
