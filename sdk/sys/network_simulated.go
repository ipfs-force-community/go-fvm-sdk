//go:build simulate

package sys

import (



	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"


	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys/internal/simulated"
)

func BaseFee() (*types.TokenAmount, error) {
	return simulated.MockFvmInstance.BaseFee()
}

func TotalFilCircSupply() (*types.TokenAmount, error) {
	return simulated.MockFvmInstance.TotalFilCircSupply()
}
