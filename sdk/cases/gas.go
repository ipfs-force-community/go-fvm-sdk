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
	err := sdk.Charge(ctx, "OnChainMessage", 38863)
	assert.Nil(t, err, "charge gas %v", err)
	_, err = sdk.AvailableGas(ctx)
	assert.Nil(t, err, "charge gas %v", err)
	return 0
}
