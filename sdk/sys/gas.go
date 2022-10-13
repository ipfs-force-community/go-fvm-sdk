//go:build !simulate
// +build !simulate

package sys

import (
	"context"
	"fmt"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
)

// Charge charge gas for the operation identified by name.
func Charge(ctx context.Context, name string, compute uint64) error {
	nameBufPtr, nameBufLen := GetStringPointerAndLen(name)
	code := gasCharge(nameBufPtr, nameBufLen, compute)
	if code != 0 {
		return ferrors.NewFvmError(ferrors.ExitCode(code), fmt.Sprintf("charge gas to %s", name))
	}
	return nil
}
