package sdk

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys"
)

// BaseFee gets the base fee for the current epoch.
func BaseFee() (*abi.TokenAmount, error) {
	return sys.BaseFee()
}

// TotalFilCircSupply gets the circulating supply.
func TotalFilCircSupply() (*abi.TokenAmount, error) {
	return sys.TotalFilCircSupply()
}
