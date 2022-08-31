//go:build simulate
// +build simulate

package sys

import (
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys/simulated"
)

func Abort(code uint32, msg string) {
	simulated.DefaultFsm.Abort(code, msg)
}
