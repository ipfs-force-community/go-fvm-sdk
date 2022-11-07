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

func Send(_ context.Context, to address.Address, method uint64, params uint32, value abi.TokenAmount) (*types.Send, error) {
	fvmTokenAmount := FromBig(&value)
	send := new(types.Send)
	addrBufPtr, addrBufLen := GetSlicePointerAndLen(to.Bytes())
	code := sysSend(uintptr(unsafe.Pointer(send)), addrBufPtr, addrBufLen, method, params, fvmTokenAmount.Hi, fvmTokenAmount.Lo)
	if code != 0 {
		return nil, ferrors.NewFvmErrorNumber(ferrors.ErrorNumber(code), "failed to send")
	}

	return send, nil
}
