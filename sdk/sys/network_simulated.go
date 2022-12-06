//go:build simulate
// +build simulate

package sys

import (
	"context"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"

	"github.com/filecoin-project/go-state-types/abi"
)

// TotalFilCircSupply gets the circulating supply.
func TotalFilCircSupply(ctx context.Context) (abi.TokenAmount, error) {
	if env, ok := tryGetSimulator(ctx); ok {
		return env.TotalFilCircSupply()
	}
	panic(ErrorEnvValid)
}

func TipsetCid(ctx context.Context, epoch abi.ChainEpoch) (*cid.Cid, error) {
	if env, ok := tryGetSimulator(ctx); ok {
		return env.TipsetCid(epoch)
	}
	panic(ErrorEnvValid)
}

func NetworkContext(_ context.Context) (*types.NetworkContext, error) {
	if env, ok := tryGetSimulator(ctx); ok {
		return env.NetworkContext(epoch)
	}
	panic(ErrorEnvValid)
}
