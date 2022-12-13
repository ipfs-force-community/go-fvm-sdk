//go:build simulate
// +build simulate

package sys

import (
	"context"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
)

func EmitEvent(ctx context.Context, evt types.ActorEvent) error {
	if env, ok := tryGetSimulator(ctx); ok {
		env.AppendEvent(evt)
		return nil
	}
	panic(ErrorEnvValid)
}
