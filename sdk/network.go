package sdk

import (
	"context"

	"github.com/filecoin-project/go-state-types/big"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/network"
	"github.com/ipfs/go-cid"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys"
)

// TotalFilCircSupply gets the circulating supply.
func TotalFilCircSupply(ctx context.Context) (abi.TokenAmount, error) {
	return sys.TotalFilCircSupply(ctx)
}

// TipsetTimestamp gets Timestamp
func TipsetTimestamp(ctx context.Context) (uint64, error) {
	networkCtx, err := sys.NetworkContext(ctx)
	if err != nil {
		return 0, err
	}
	return networkCtx.Timestamp, nil
}

// TipsetCid gets cid
func TipsetCid(ctx context.Context, epoch abi.ChainEpoch) (*cid.Cid, error) {
	return sys.TipsetCid(ctx, epoch)
}

// CurrEpoch get network current epoch
func CurrEpoch(ctx context.Context) (abi.ChainEpoch, error) {
	networkCtx, err := sys.NetworkContext(ctx)
	if err != nil {
		return 0, err
	}
	return networkCtx.Epoch, nil
}

// Version network version
func Version(ctx context.Context) (network.Version, error) {
	networkCtx, err := sys.NetworkContext(ctx)
	if err != nil {
		return 0, err
	}
	return network.Version(networkCtx.NetworkVersion), nil
}

// BaseFee gets the base fee for the current epoch.
func BaseFee(ctx context.Context) (abi.TokenAmount, error) {
	networkCtx, err := sys.NetworkContext(ctx)
	if err != nil {
		return big.Zero(), err
	}
	return networkCtx.BaseFee, nil
}
