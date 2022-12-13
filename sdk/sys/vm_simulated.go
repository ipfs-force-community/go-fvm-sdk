//go:build simulate
// +build simulate

package sys

import (
	"context"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
)

func VMMessageContext(ctx context.Context) (*types.MessageContext, error) {
	if env, ok := tryGetSimulator(ctx); ok {
		return env.VMMessageContext()
	}
	panic(ErrorEnvValid)
}

func Exit(ctx context.Context, code ferrors.ExitCode, data []byte, msg string) {
	if env, ok := tryGetSimulator(ctx); ok {
		env.Exit(code, data, msg)
		return
	}
	panic(ErrorEnvValid)
}

// Exit exit actor, panic to stop actor instead of return error
func ExitWithBlkId(ctx context.Context, code ferrors.ExitCode, blkId types.BlockID, msg string) {
	if env, ok := tryGetSimulator(ctx); ok {
		env.ExitWithId(code, blkId, msg)
		return
	}
	panic(ErrorEnvValid)
}
