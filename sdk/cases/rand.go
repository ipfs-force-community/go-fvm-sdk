package main

import (
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/testing"
	"github.com/stretchr/testify/assert"
)

func main() {}

//go:export invoke
func Invoke(_ uint32) uint32 {
	t := testing.NewTestingT()
	defer t.CheckResult()

	// todo  panicked at 'not yet implemented'
	_, err := sdk.GetChainRandomness(crypto.DomainSeparationTag_TicketProduction, 0, []byte{})
	assert.Nil(t, err, "get chain randomness %v", err)

	_, err = sdk.GetBeaconRandomness(crypto.DomainSeparationTag_SealRandomness, 0, []byte{})
	assert.Nil(t, err, "get beacon randomness %v", err)

	return 0
}
