//go:build simulate

package sys

import (
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/fvm"
)

func Abort(code uint32, msg string) {
	fvm.MockFvmInstance.Abort(code, msg)
}
