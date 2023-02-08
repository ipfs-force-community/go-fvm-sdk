package main

import (
	"context"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
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
	addr, _ := address.NewFromString("f010000")
	actorId, err := address.IDFromAddress(addr)
	assert.Nil(t, err)
	preBalance, err := sdk.BalanceOf(ctx, abi.ActorID(actorId))
	assert.Nil(t, err)

	ret, err := sdk.Send(ctx, addr, 0, []byte{}, abi.NewTokenAmount(1), sdk.WithReadonly())
	assert.Nil(t, err, "send %v", err)
	assert.Equal(t, 0, int(ret.ExitCode))
	assert.Equal(t, 0, int(ret.GasUsed))
	assert.Equal(t, "", string(ret.ReturnData))

	balance, err := sdk.BalanceOf(ctx, abi.ActorID(actorId))
	assert.Nil(t, err)
	assert.Equal(t, preBalance.String(), balance.String())
	return 0
}
