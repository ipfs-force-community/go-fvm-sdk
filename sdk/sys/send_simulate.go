//go:build simulate
// +build simulate

package sys

import (
	"github.com/filecoin-project/go-address"


	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys/internal/simulated"
)

func Send(to address.Address, method uint64, params uint32, value types.TokenAmount) (*types.Send, error) {
	return simulated.SimulatedInstance.Send(to, method, params, value)
}
