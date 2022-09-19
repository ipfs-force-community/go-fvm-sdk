//go:build !simulated
// +build !simulated

package sys

import (
	"unsafe"

	"github.com/filecoin-project/go-state-types/abi"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
)

// BaseFee gets the base fee for the current epoch.
func BaseFee() (*abi.TokenAmount, error) {
	result := new(abi.TokenAmount)
	code := networkBaseFee(uintptr(unsafe.Pointer(result)))
	if code != 0 {
		return nil, ferrors.NewFvmError(ferrors.ExitCode(code), "failed to get base fee")
	}

	return result, nil
}

// TotalFilCircSupply gets the circulating supply.
func TotalFilCircSupply() (*abi.TokenAmount, error) {
	result := new(abi.TokenAmount)
	code := networkTotalFilCircSupply(uintptr(unsafe.Pointer(result)))
	if code != 0 {
		return nil, ferrors.NewFvmError(ferrors.ExitCode(code), "failed to get circulating supply")
	}

	return result, nil
}
