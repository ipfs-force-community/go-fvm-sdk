//go:build simulate
// +build simulate

package sys

import (
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys/simulated"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
)

func BaseFee() (*types.TokenAmount, error) {
	return simulated.DefaultFsm.BaseFee()
}

func TotalFilCircSupply() (*types.TokenAmount, error) {
	return simulated.DefaultFsm.TotalFilCircSupply()
}
