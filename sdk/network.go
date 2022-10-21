package sdk

import (
	"context"

	"github.com/filecoin-project/go-state-types/abi"

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
