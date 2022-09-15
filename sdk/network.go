package sdk

import (
	"github.com/filecoin-project/go-state-types/big"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys"
)

// BaseFee gets the base fee for the current epoch.
func BaseFee() (*big.Int, error) {
	return sys.BaseFee()
}

// TotalFilCircSupply gets the circulating supply.
func TotalFilCircSupply() (*big.Int, error) {
	return sys.TotalFilCircSupply()
}
