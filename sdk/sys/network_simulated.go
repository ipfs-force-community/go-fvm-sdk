//go:build simulate

package sys

import (
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/fvm"
)


func BaseFee() (*types.TokenAmount, error) {
	return fvm.MockFvmInstance.BaseFee()
}

func TotalFilCircSupply() (*types.TokenAmount, error) {
	return fvm.MockFvmInstance.TotalFilCircSupply()
}
