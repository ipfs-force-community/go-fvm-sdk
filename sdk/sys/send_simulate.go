//go:build simulate
// +build simulate

package sys

import (
	"github.com/filecoin-project/go-address"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys/simulated"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
)

func Send(to address.Address, method uint64, params uint32, value types.TokenAmount) (*types.Send, error) {
	return simulated.DefaultFsm.Send(to, method, params, value)
}
