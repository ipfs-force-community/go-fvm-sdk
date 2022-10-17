package sdk

import (
	"context"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/network"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys"
)

// InvocationCtx the runtime context in actor method invocation.
var InvocationCtx *types.InvocationContext

// Caller get caller, from address of message
func Caller(ctx context.Context) (abi.ActorID, error) {
	if InvocationCtx == nil {
		var err error
		InvocationCtx, err = sys.VMContext(ctx)
		if err != nil {
			return 0, err
		}
	}
	return InvocationCtx.Caller, nil
}

// Receiver get recevier, to address of message
func Receiver(ctx context.Context) (abi.ActorID, error) {
	if InvocationCtx == nil {
		var err error
		InvocationCtx, err = sys.VMContext(ctx)
		if err != nil {
			return 0, err
		}
	}
	return InvocationCtx.Receiver, nil
}

// MethodNumber method number
func MethodNumber(ctx context.Context) (abi.MethodNum, error) {
	if InvocationCtx == nil {
		var err error
		InvocationCtx, err = sys.VMContext(ctx)
		if err != nil {
			return 0, err
		}
	}
	return InvocationCtx.MethodNumber, nil
}

// ValueReceived the amount was transferred in message
func ValueReceived(ctx context.Context) (*types.TokenAmount, error) {
	if InvocationCtx == nil {
		var err error
		InvocationCtx, err = sys.VMContext(ctx)
		if err != nil {
			return nil, err
		}
	}

	return &types.TokenAmount{ //avoud change
		Lo: InvocationCtx.ValueReceived.Lo,
		Hi: InvocationCtx.ValueReceived.Hi,
	}, nil
}

// ParamsRaw returns the message codec and parameters.
func ParamsRaw(ctx context.Context, id types.BlockID) (*types.ParamsRaw, error) {
	if id == types.NoDataBlockID {
		return &types.ParamsRaw{}, nil
	}
	state, err := sys.Stat(ctx, id)
	if err != nil {
		return nil, err
	}

	block, err := GetBlock(ctx, id, &state.Size)
	if err != nil {
		return nil, err
	}
	return &types.ParamsRaw{
		Codec: state.Codec,
		Raw:   block,
	}, nil
}

// CurrEpoch get network current epoch
func CurrEpoch(ctx context.Context) (abi.ChainEpoch, error) {
	if InvocationCtx == nil {
		var err error
		InvocationCtx, err = sys.VMContext(ctx)
		if err != nil {
			return 0, err
		}
	}
	return InvocationCtx.NetworkCurrEpoch, nil
}

// Version network version
func Version(ctx context.Context) (network.Version, error) {
	if InvocationCtx == nil {
		var err error
		InvocationCtx, err = sys.VMContext(ctx)
		if err != nil {
			return 0, err
		}
	}
	return network.Version(InvocationCtx.NetworkVersion), nil
}
