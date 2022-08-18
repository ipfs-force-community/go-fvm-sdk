//go:build simulate

package sys

import (
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/fvm"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"
)

// Abort abort execution
func Open(id cid.Cid) (*types.IpldOpen, error) {
	return fvm.MockFvmInstance.Open(id)
}
