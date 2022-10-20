//go:build simulate
// +build simulate

package sys

import (
	"context"

	"github.com/filecoin-project/go-state-types/abi"
)

// BaseFee gets the base fee for the current epoch.
func BaseFee(ctx context.Context) (abi.TokenAmount, error) {
	if env, ok := tryGetSimulator(ctx); ok {
		return env.BaseFee()
	}
	panic(ErrorEnvValid)
}

// TotalFilCircSupply gets the circulating supply.
func TotalFilCircSupply(ctx context.Context) (abi.TokenAmount, error) {
	if env, ok := tryGetSimulator(ctx); ok {
		return env.TotalFilCircSupply()
	}
	panic(ErrorEnvValid)
}
