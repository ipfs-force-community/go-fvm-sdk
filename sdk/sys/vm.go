//go:build !simulate
// +build !simulate

package sys

import (
	"context"
	"strconv"
	"unsafe"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
)

func VMMessageContext(_ context.Context) (*types.MessageContext, error) {
	var result messageContext
	code := vmMessageContext(uintptr(unsafe.Pointer(&result)))
	if code != 0 {
		return nil, ferrors.NewSysCallError(ferrors.ErrorNumber(code), "unable to get invocation context")
	}
	return &types.MessageContext{
		Origin:        result.Origin,
		Caller:        result.Caller,
		Receiver:      result.Receiver,
		MethodNumber:  result.MethodNumber,
		ValueReceived: *result.ValueReceived.TokenAmount(),
		GasPremium:    *result.ValueReceived.TokenAmount(),
		Flags:         0,
	}, nil
}

// Exit abort actor, panic to stop actor instead of return error
func Exit(ctx context.Context, code ferrors.ExitCode, data []byte, msg string) {
	blkId := types.NoDataBlockID
	if len(data) == 0 {
		var err error
		blkId, err = Create(ctx, types.DAGCbor, data)
		if err != nil {
			panic("failed create block when exit " + err.Error())
		}
	}

	strPtr, strLen := GetStringPointerAndLen(msg)
	exitCode := vmExit(uint32(code), blkId, strPtr, strLen)
	if exitCode != 0 {
		panic("fail to exit " + strconv.Itoa(int(exitCode)))
	}
}

// ExitWithBlkId exit actor, panic to stop actor instead of return error
func ExitWithBlkId(ctx context.Context, code ferrors.ExitCode, blkId types.BlockID, msg string) {
	strPtr, strLen := GetStringPointerAndLen(msg)
	exitCode := vmExit(uint32(code), blkId, strPtr, strLen)
	if exitCode != 0 {
		panic("fail to exit " + strconv.Itoa(int(exitCode)))
	}
}
