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

// TotalFilCircSupply gets the circulating supply.
func TotalFilCircSupply(_ context.Context) (abi.TokenAmount, error) {
	result := new(fvmTokenAmount)
	code := networkTotalFilCircSupply(uintptr(unsafe.Pointer(result)))
	if code != 0 {
		return abi.TokenAmount{}, ferrors.NewSysCallError(ferrors.ErrorNumber(code), "failed to get circulating supply")
	}

	return *result.TokenAmount(), nil
}

func TipsetCid(_ context.Context, epoch abi.ChainEpoch) (*cid.Cid, error) {
	buf := make([]byte, types.MaxCidLen)
	bufPtr, bufLen := GetSlicePointerAndLen(buf)
	var result uint32
	code := networkTipsetCid(uintptr(unsafe.Pointer(&result)), int64(epoch), bufPtr, bufLen)
	if code != 0 {
		return nil, ferrors.NewSysCallError(ferrors.ErrorNumber(code), "unexpected cid resolution failure: ")

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

func NetworkContext(_ context.Context) (*types.NetworkContext, error) {
	var result networkContext_
	code := networkContext(uintptr(unsafe.Pointer(&result)))
	if code != 0 {
		return nil, ferrors.NewSysCallError(ferrors.ErrorNumber(code), "unable to get invocation context")
	}
	return &types.NetworkContext{
		Epoch:          result.Epoch,
		Timestamp:      result.Timestamp,
		BaseFee:        *result.BaseFee.TokenAmount(),
		NetworkVersion: result.NetworkVersion,
	}, nil
}
