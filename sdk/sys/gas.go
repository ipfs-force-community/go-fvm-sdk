//go:build !simulate
// +build !simulate

package sys

import (
	"context"
	"fmt"
	"unsafe"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
)

// Charge charge gas for the operation identified by name.
func Charge(_ context.Context, name string, compute uint64) error {
	nameBufPtr, nameBufLen := GetStringPointerAndLen(name)
	code := gasCharge(nameBufPtr, nameBufLen, compute)
	if code != 0 {
		return ferrors.NewFvmError(ferrors.ExitCode(code), fmt.Sprintf("charge gas to %s", name))
	}
	return nil
}

// Available Returns the amount of gas remaining.
func Available(_ context.Context) (uint64, error) {
	var retptr uint32
	code := gasAvailable(uintptr(unsafe.Pointer(&retptr)))
	if code != 0 {
		return 0, ferrors.NewFvmError(ferrors.ExitCode(code), "Available is fail")
	}
	return uint64(retptr), nil
}
