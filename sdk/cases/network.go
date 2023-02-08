package main

import (
	"context"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/testing"
	"github.com/stretchr/testify/assert"
)

func main() {} //nolint

//go:export invoke
func Invoke(_ uint32) uint32 { //nolint
	t := testing.NewTestingT()
	defer t.CheckResult()
	ctx := context.Background()
	epoch, err := sdk.CurrEpoch(ctx)
	assert.Nil(t, err)
	assert.Equal(t, 0, int(epoch), "epoch not match")

	ver, err := sdk.Version(ctx)
	assert.Nil(t, err)
	assert.Equal(t, 18, int(ver), "version not match")

	fee, err := sdk.BaseFee(ctx)
	assert.Nil(t, err)
	assert.Equal(t, "100", fee.String(), "base fee not match")

	value, err := sdk.TotalFilCircSupply(ctx)
	assert.Nil(t, err)
	assert.Equal(t, "2000000000000000000000000000", value.String(), "ful circsupply not match")
	return 0
}
