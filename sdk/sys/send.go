//go:build !simulate
// +build !simulate

package sys

import (
	"context"
	"unsafe"

	"github.com/filecoin-project/go-state-types/abi"

	"github.com/filecoin-project/go-address"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
)

func Send(_ context.Context, to address.Address, method abi.MethodNum, params uint32, value abi.TokenAmount, gasLimit uint64, flag uint64) (*types.SendResult, error) {
	fvmTokenAmount := FromBig(&value)
	send := new(types.SendResult)
	addrBufPtr, addrBufLen := GetSlicePointerAndLen(to.Bytes())
	code := sysSend(uintptr(unsafe.Pointer(send)), addrBufPtr, addrBufLen, uint64(method), params, fvmTokenAmount.Hi, fvmTokenAmount.Lo, gasLimit, flag)
	if code != 0 {
		return nil, ferrors.NewSysCallError(ferrors.ErrorNumber(code), "failed to send")
	}

	return send, nil
}
