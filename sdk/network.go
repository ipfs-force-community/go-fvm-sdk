package sdk

import (
	"context"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/ipfs/go-cid"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys"
)

// BaseFee gets the base fee for the current epoch.
func BaseFee(ctx context.Context) (abi.TokenAmount, error) {
	return sys.BaseFee(ctx)
}

// TotalFilCircSupply gets the circulating supply.
func TotalFilCircSupply(ctx context.Context) (abi.TokenAmount, error) {
	return sys.TotalFilCircSupply(ctx)
}

func TipsetTimestamp(ctx context.Context) (uint64, error) {
	return sys.TipsetTimestamp(ctx)
}

func TipsetCid(ctx context.Context, epoch uint64) (*cid.Cid, error) {
	return sys.TipsetCid(ctx, epoch)
}
