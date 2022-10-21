//go:build !simulate
// +build !simulate

package sys

import (
	"context"
	"unsafe"

	"github.com/filecoin-project/go-state-types/abi"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
)

// BaseFee gets the base fee for the current epoch.
func BaseFee(_ context.Context) (abi.TokenAmount, error) {
	result := new(fvmTokenAmount)
	code := networkBaseFee(uintptr(unsafe.Pointer(result)))
	if code != 0 {
		return abi.TokenAmount{}, ferrors.NewFvmError(ferrors.ExitCode(code), "failed to get base fee")
	}
	return result.TokenAmount(), nil
}

// TotalFilCircSupply gets the circulating supply.
func TotalFilCircSupply(_ context.Context) (abi.TokenAmount, error) {
	result := new(fvmTokenAmount)
	code := networkTotalFilCircSupply(uintptr(unsafe.Pointer(result)))
	if code != 0 {
		return abi.TokenAmount{}, ferrors.NewFvmError(ferrors.ExitCode(code), "failed to get circulating supply")
	}

	return result.TokenAmount(), nil
}
