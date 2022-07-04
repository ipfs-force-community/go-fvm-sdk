package sdk

import (
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys"
)

// Abort abort execution
func Abort(code ferrors.ExitCode, msg string) {
	sys.Abort(uint32(code), msg)

}
