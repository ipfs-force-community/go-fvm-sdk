//nolint:param
package main

import (
	"github.com/ipfs-force-community/go-fvm-sdk/sdk"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/testing"
	"github.com/stretchr/testify/assert"
)

func main() {} //nolint

//go:export invoke
func Invoke(_ uint32) uint32 { //nolint
	t := testing.NewTestingT()
	defer t.CheckResult()

	epoch, err := sdk.CurrEpoch()
	assert.Nil(t, err)
	assert.Equal(t, 0, int(epoch))

	ver, err := sdk.Version()
	assert.Nil(t, err)
	assert.Equal(t, 15, int(ver))

	fee, err := sdk.BaseFee()
	assert.Nil(t, err)
	assert.Equal(t, "100", fee.Big().String())

	value, err := sdk.TotalFilCircSupply()
	assert.Nil(t, err)
	assert.Equal(t, "2000000000000000000000000000", value.Big().String())
	return 0
}
