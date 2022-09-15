//go:build !simulate
// +build !simulate

package sys

import (
	"unsafe"

	"github.com/filecoin-project/go-state-types/big"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
)

// BaseFee gets the base fee for the current epoch.
func BaseFee() (*big.Int, error) {
	result := new(types.TokenAmount)
	code := networkBaseFee(uintptr(unsafe.Pointer(result)))
	if code != 0 {
		return nil, ferrors.NewFvmError(ferrors.ExitCode(code), "failed to get base fee")
	}

	return result.Big(), nil
}

// TotalFilCircSupply gets the circulating supply.
func TotalFilCircSupply() (*big.Int, error) {
	result := new(types.TokenAmount)
	code := networkTotalFilCircSupply(uintptr(unsafe.Pointer(result)))
	if code != 0 {
		return nil, ferrors.NewFvmError(ferrors.ExitCode(code), "failed to get circulating supply")
	}

	return result.Big(), nil
}
