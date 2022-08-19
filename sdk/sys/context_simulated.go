//go:build simulate

package sys

import (
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/fvm"
)

func VMContext() (*types.InvocationContext, error) {
	return fvm.MockFvmInstance.VMContext()
}
