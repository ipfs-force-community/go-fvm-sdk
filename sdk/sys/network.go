package sys

import (
	"unsafe"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
)

//! Syscalls for network metadata.
/// Gets the current epoch.
///
/// # Errors
///
/// None
//go:wasm-module network
//export curr_epoch
func networkCurrEpoch(ret uintptr) uint32

/// Gets the network version.
///
/// # Errors
///
/// None
//go:wasm-module network
//export version
func networkVersion(ret uintptr) uint32

/// Gets the base fee for the current epoch.
///
/// # Errors
///
/// None
//go:wasm-module network
//export base_fee
func networkBaseFee(ret uintptr) uint32

/// Gets the circulating supply.
///
/// # Errors
///
/// None
//go:wasm-module network
//export total_fil_circ_supply
func networkTotalFilCircSupply(ret uintptr) uint32

func CurrEpoch() (abi.ChainEpoch, error) {
	var result int64
	code := networkCurrEpoch(uintptr(unsafe.Pointer(&result)))
	if code != 0 {
		return 0, ferrors.NewFvmError(ferrors.ExitCode(code), "failed to get current epoch")
	}

	return abi.ChainEpoch(result), nil
}

func Version() (uint32, error) {
	var result uint32
	code := networkVersion(uintptr(unsafe.Pointer(&result)))
	if code != 0 {
		return 0, ferrors.NewFvmError(ferrors.ExitCode(code), "failed to get network version")
	}

	return result, nil
}

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
