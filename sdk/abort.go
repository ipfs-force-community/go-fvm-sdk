package sdk

import (
	"context"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys"
)

// Abort abort execution
func Abort(ctx context.Context, code ferrors.ExitCode, msg string) {
	sys.Abort(ctx, uint32(code), msg)

}
