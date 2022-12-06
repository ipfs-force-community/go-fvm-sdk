package sdk

import (
	"context"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"

	"github.com/filecoin-project/go-address"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys"
)

// CallerAddress return caller address
func CallerAddress(ctx context.Context) (address.Address, error) {
	msgCtx, err := sys.VMMessageContext(ctx)
	if err != nil {
		return address.Undef, err
	}
	return address.NewIDAddress(uint64(msgCtx.Caller))
}

// ReceiverAddress return message to address
func ReceiverAddress(ctx context.Context) (address.Address, error) {
	msgCtx, err := sys.VMMessageContext(ctx)
	if err != nil {
		return address.Undef, err
	}

	return address.NewIDAddress(uint64(msgCtx.Receiver))
}

// OriginAddress return message origin caller address
func OriginAddress(ctx context.Context) (address.Address, error) {
	msgCtx, err := sys.VMMessageContext(ctx)
	if err != nil {
		return address.Undef, err
	}

	return address.NewIDAddress(uint64(msgCtx.Origin))
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
