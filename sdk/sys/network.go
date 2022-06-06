package sys

import (
	"unsafe"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
)

func BaseFee() (*types.TokenAmount, error) {
	result := new(types.TokenAmount)
	code := networkBaseFee(uintptr(unsafe.Pointer(result)))
	if code != 0 {
		return nil, ferrors.NewFvmError(ferrors.ExitCode(code), "failed to get base fee")
	}

	return result, nil
}

func TotalFilCircSupply() (*types.TokenAmount, error) {
	result := new(types.TokenAmount)
	code := networkTotalFilCircSupply(uintptr(unsafe.Pointer(result)))
	if code != 0 {
		return nil, ferrors.NewFvmError(ferrors.ExitCode(code), "failed to get circulating supply")
	}

	return result, nil
}
