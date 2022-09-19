//go:build simulated
// +build simulated

package sys

import (
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys/simulated"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
)

func VMContext() (*types.InvocationContext, error) {
	return simulated.DefaultFsm.VMContext()
}
