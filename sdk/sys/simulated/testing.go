package simulated

import (
	"github.com/filecoin-project/go-state-types/abi"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-crypto"
)

// NewF1Address create f1 address, f3 address not support for now
func NewF1Address() (address.Address, error) {
	priv, err := crypto.GenerateKey()
	if err != nil {
		return address.Undef, err
	}
	return address.NewSecp256k1Address(crypto.PublicKey(priv))
}

func NewPtrTokenAmount(t int64) *abi.TokenAmount {
	v := abi.NewTokenAmount(t)
	return &v
}
