//go:build !simulate
// +build !simulate

package sys

import (
	"context"
	"unsafe"

	"github.com/filecoin-project/go-state-types/network"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
)

func VMContext(_ context.Context) (*types.InvocationContext, error) {
	var result invocationContext
	code := vmContext(uintptr(unsafe.Pointer(&result)))
	if code != 0 {
		return nil, ferrors.NewFvmError(ferrors.ExitCode(code), "unable to get invocation context")
	}
	return &types.InvocationContext{
		ValueReceived:    result.ValueReceived.TokenAmount(),
		Caller:           result.Caller,
		Receiver:         result.Receiver,
		MethodNumber:     result.MethodNumber,
		NetworkCurrEpoch: result.NetworkCurrEpoch,
		NetworkVersion:   network.Version(result.NetworkVersion),
	}, nil
}
