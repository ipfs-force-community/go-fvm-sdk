package sdk

import (
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
)

// BaseFee gets the base fee for the current epoch.
func BaseFee() (*types.TokenAmount, error) {
	return sys.BaseFee()
}

// TotalFilCircSupply gets the circulating supply.
func TotalFilCircSupply() (*types.TokenAmount, error) {
	return sys.TotalFilCircSupply()
}
