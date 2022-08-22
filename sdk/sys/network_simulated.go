//go:build simulate
// +build simulate

package sys

import (



	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"


	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys/internal/simulated"
)

func BaseFee() (*types.TokenAmount, error) {
	return simulated.SimulatedInstance.BaseFee()
}

func TotalFilCircSupply() (*types.TokenAmount, error) {
	return simulated.SimulatedInstance.TotalFilCircSupply()
}
