package sdk

import (
	"context"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
)

// Abort abort execution, code must be non-okay value
func Abort(ctx context.Context, code ferrors.ExitCode, msg string) {
	if code == 0 {
		Exit(ctx, ferrors.USR_ASSERTION_FAILED, nil, msg)
	}
	Exit(ctx, code, nil, msg)

}

// Exit abort contract with data and exit message
func Exit(ctx context.Context, code ferrors.ExitCode, data []byte, msg string) {
	sys.Exit(ctx, code, data, msg)
}

// ExitWithBlkId abort contract with specify block id and exit message avoid to create block again sometimes
func ExitWithBlkId(ctx context.Context, code ferrors.ExitCode, blkId types.BlockID, msg string) {
	sys.ExitWithBlkId(ctx, code, blkId, msg)
}
