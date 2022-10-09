package sys

import (
	"context"
	"unsafe"

	"github.com/filecoin-project/go-address"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
)

func Send(ctx context.Context, to address.Address, method uint64, params uint32, value types.TokenAmount) (*types.Send, error) {
	if env, ok := isSimulatedEnv(ctx); ok {
		return env.Send(to, method, params, value)
	}
	send := new(types.Send)
	addrBufPtr, addrBufLen := GetSlicePointerAndLen(to.Bytes())
	code := sysSend(uintptr(unsafe.Pointer(send)), addrBufPtr, addrBufLen, method, params, value.Hi, value.Lo)
	if code != 0 {
		return nil, ferrors.NewFvmError(ferrors.ExitCode(code), "failed to send")
	}

	return send, nil
}
