//go:build simulate

package sys

import (
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys/internal/simulated"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
)

func VMContext() (*types.InvocationContext, error) {
	return simulated.MockFvmInstance.VMContext()
}
