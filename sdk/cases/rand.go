package main

import (
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/testing"
	"github.com/stretchr/testify/assert"
)

func main() {} //nolint

//go:export invoke
func Invoke(_ uint32) uint32 { //nolint
	t := testing.NewTestingT()
	defer t.CheckResult()

	randValue, err := sdk.GetChainRandomness(crypto.DomainSeparationTag_TicketProduction, 0, []byte{})
	assert.Nil(t, err, "get chain randomness %v", err)
	t.Infof("got chain randomness %v", randValue)
	assert.NotEmpty(t, randValue)

	randValue, err = sdk.GetBeaconRandomness(crypto.DomainSeparationTag_SealRandomness, 0, []byte{})
	assert.Nil(t, err, "get beacon randomness %v", err)
	t.Infof("got beacon randomness %v", randValue)
	assert.NotEmpty(t, randValue)
	return 0
}
