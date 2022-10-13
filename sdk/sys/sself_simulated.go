//go:build simulate
// +build simulate

package sys

import (
	"context"

	addr "github.com/filecoin-project/go-address"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"
)

func SelfRoot(ctx context.Context) (cid.Cid, error) {
	if env, ok := tryGetSimulator(ctx); ok {
		return env.SelfRoot()
	}
	return cid.Undef, ErrorEnvValid
}

func SelfSetRoot(ctx context.Context, id cid.Cid) error {
	if env, ok := tryGetSimulator(ctx); ok {
		return env.SelfSetRoot(id)
	}
	return ErrorEnvValid

}

func SelfCurrentBalance(ctx context.Context) (*types.TokenAmount, error) {
	if env, ok := tryGetSimulator(ctx); ok {
		return env.SelfCurrentBalance()
	}
	return &types.TokenAmount{}, ErrorEnvValid

}

func SelfDestruct(ctx context.Context, addr addr.Address) error {
	if env, ok := tryGetSimulator(ctx); ok {
		return env.SelfDestruct(addr)
	}
	return ErrorEnvValid

}
