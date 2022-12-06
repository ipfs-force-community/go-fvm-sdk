//go:build simulate
// +build simulate

package sys

import "github.com/ipfs-force-community/go-fvm-sdk/sdk/types"

func EmitEvent(evt types.ActorEvent) error {
	if env, ok := tryGetSimulator(ctx); ok {
		return env.AppendEvent(evt)
	}
	panic(ErrorEnvValid)
}
