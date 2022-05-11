package sdk

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/network"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
)

func CurrEpoch() (abi.ChainEpoch, error) {
	return sys.CurrEpoch()
}

func Version() (network.Version, error) {
	return sys.Version()
}

func BaseFee() (*types.TokenAmount, error) {
	return sys.BaseFee()
}

func TotalFilCircSupply() (*types.TokenAmount, error) {
	return sys.TotalFilCircSupply()
}
