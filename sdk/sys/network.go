//go:build !simulate
// +build !simulate

package sys

import (
	"context"
	"unsafe"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/ipfs/go-cid"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
)

// BaseFee gets the base fee for the current epoch.
func BaseFee(_ context.Context) (abi.TokenAmount, error) {
	result := new(fvmTokenAmount)
	code := networkBaseFee(uintptr(unsafe.Pointer(result)))
	if code != 0 {
		return abi.TokenAmount{}, ferrors.NewSysCallError(ferrors.ErrorNumber(code), "failed to get base fee")
	}
	return *result.TokenAmount(), nil
}

// TotalFilCircSupply gets the circulating supply.
func TotalFilCircSupply(_ context.Context) (abi.TokenAmount, error) {
	result := new(fvmTokenAmount)
	code := networkTotalFilCircSupply(uintptr(unsafe.Pointer(result)))
	if code != 0 {
		return abi.TokenAmount{}, ferrors.NewSysCallError(ferrors.ErrorNumber(code), "failed to get circulating supply")
	}

	return *result.TokenAmount(), nil
}

func TipsetTimestamp(_ context.Context) (uint64, error) {
	var timestamp uint64
	code := networkTipsetTimestamp(uintptr(unsafe.Pointer(&timestamp)))
	if code != 0 {
		return 0, ferrors.NewSysCallError(ferrors.ErrorNumber(code), "failed to get timestamp")
	}

	return timestamp, nil
}

func TipsetCid(_ context.Context, epoch uint64) (*cid.Cid, error) {
	buf := make([]byte, types.MaxCidLen)
	bufPtr, bufLen := GetSlicePointerAndLen(buf)
	var result uint32
	code := networkTipsetCid(uintptr(unsafe.Pointer(&result)), epoch, bufPtr, bufLen)
	if code != 0 {
		return nil, ferrors.NewSysCallError(ferrors.ErrorNumber(code), "unexpected cid resolution failure: "+ferrors.EnToString(code))

	}
	if result > 0 {
		_, result, err := cid.CidFromBytes(buf)
		if err != nil {
			return nil, err
		}
		return &result, nil
	} else {
		return nil, nil
	}
}
