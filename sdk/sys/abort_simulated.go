//go:build simulate

package sys

import (
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys/internal/simulated"
)

func Abort(code uint32, msg string) {
	simulated.MockFvmInstance.Abort(code, msg)
}
