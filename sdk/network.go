package sdk

import (
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
)

func BaseFee() (*types.TokenAmount, error) {
	return sys.BaseFee()
}

func TotalFilCircSupply() (*types.TokenAmount, error) {
	return sys.TotalFilCircSupply()
}
