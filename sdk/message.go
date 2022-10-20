package sdk

import (
	"context"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/network"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys"
)

// Caller get caller, from address of message
// todo cache invocation in context
func Caller(ctx context.Context) (abi.ActorID, error) {
	invocationCtx, err := sys.VMContext(ctx)
	if err != nil {
		return 0, err
	}
	return invocationCtx.Caller, nil
}

// Receiver get recevier, to address of message
func Receiver(ctx context.Context) (abi.ActorID, error) {
	invocationCtx, err := sys.VMContext(ctx)
	if err != nil {
		return 0, err
	}
	return invocationCtx.Receiver, nil
}

// MethodNumber method number
func MethodNumber(ctx context.Context) (abi.MethodNum, error) {
	invocationCtx, err := sys.VMContext(ctx)
	if err != nil {
		return 0, err
	}
	return invocationCtx.MethodNumber, nil
}

// ValueReceived the amount was transferred in message
func ValueReceived(ctx context.Context) (abi.TokenAmount, error) {
	invocationCtx, err := sys.VMContext(ctx)
	if err != nil {
		return abi.TokenAmount{}, err
	}
	return invocationCtx.ValueReceived, nil
}

// CurrEpoch get network current epoch
func CurrEpoch(ctx context.Context) (abi.ChainEpoch, error) {
	invocationCtx, err := sys.VMContext(ctx)
	if err != nil {
		return 0, err
	}
	return invocationCtx.NetworkCurrEpoch, nil
}

// Version network version
func Version(ctx context.Context) (network.Version, error) {
	invocationCtx, err := sys.VMContext(ctx)
	if err != nil {
		return 0, err
	}
	return invocationCtx.NetworkVersion, nil
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
