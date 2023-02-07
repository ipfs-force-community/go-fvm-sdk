package sdk

import (
	"context"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
)

// EmitEvent emit event to fvm
func EmitEvent(ctx context.Context, evt types.ActorEvent) error {
	return sys.EmitEvent(ctx, evt)
}
