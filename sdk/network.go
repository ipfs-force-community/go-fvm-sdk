package sdk

import (
	"context"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
)

// BaseFee gets the base fee for the current epoch.
func BaseFee(ctx context.Context) (*types.TokenAmount, error) {
	return sys.BaseFee(ctx)
}

// TotalFilCircSupply gets the circulating supply.
func TotalFilCircSupply(ctx context.Context) (*types.TokenAmount, error) {
	return sys.TotalFilCircSupply(ctx)
}
