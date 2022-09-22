//go:build simulated
// +build simulated

package sys

import (
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys/simulated"
)

func BaseFee() (*types.TokenAmount, error) {
	return simulated.DefaultFsm.BaseFee()
}

func TotalFilCircSupply() (*types.TokenAmount, error) {
	return simulated.DefaultFsm.TotalFilCircSupply()
}
