//go:build simulated
// +build simulated

package sys

import (
	"context"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
)

// BaseFee gets the base fee for the current epoch.
func BaseFee(ctx context.Context) (*types.TokenAmount, error) {
	if env, ok := isSimulatedEnv(ctx); ok {
		return env.BaseFee()
	}
	return &types.TokenAmount{}, nil

}

// TotalFilCircSupply gets the circulating supply.
func TotalFilCircSupply(ctx context.Context) (*types.TokenAmount, error) {
	if env, ok := isSimulatedEnv(ctx); ok {
		return env.TotalFilCircSupply()
	}
	return &types.TokenAmount{}, nil

}
