package sdk

import (
	"context"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys"
)

// Charge charges the gas
func Charge(ctx context.Context, name string, compute uint64) error {
	return sys.Charge(ctx, name, compute)
}
