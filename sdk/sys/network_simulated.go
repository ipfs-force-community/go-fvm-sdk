//go:build simulate
// +build simulate

package sys

import (
	"github.com/filecoin-project/go-state-types/big"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys/simulated"
)

func BaseFee() (*big.Int, error) {
	return simulated.DefaultFsm.BaseFee()
}

func TotalFilCircSupply() (*big.Int, error) {
	return simulated.DefaultFsm.TotalFilCircSupply()
}
