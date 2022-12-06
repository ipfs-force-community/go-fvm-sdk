package sdk

import (
	"context"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
)

// Abort abort execution
func Abort(ctx context.Context, code ferrors.ExitCode, msg string) {
	if code == 0 {
		Exit(ctx, ferrors.USR_ASSERTION_FAILED, nil, msg)
	}
	Exit(ctx, code, nil, msg)

}

func Exit(ctx context.Context, code ferrors.ExitCode, buf []byte, msg string) {
	sys.Exit(ctx, code, buf, msg)
}
