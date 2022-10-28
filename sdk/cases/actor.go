package main

import (
	"context"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/testing"
	"github.com/stretchr/testify/assert"
)

func main() {} //nolint

//go:export invoke
func Invoke(_ uint32) uint32 { //nolint
	t := testing.NewTestingT()
	defer t.CheckResult()
	ctx := context.Background()
	methodNum, err := sdk.MethodNumber(ctx)
	if err != nil {
		sdk.Abort(ctx, ferrors.USR_ILLEGAL_STATE, "unable to get method number")
	}

	switch methodNum {
	case 1:
		addr, _ := address.NewFromString("f1dwyrbh74hr5nwqv2gjedjyvgphxxkffxug4rkkq")
		actorID, err := sdk.ResolveAddress(ctx, addr)
		addr_, err := sdk.LookupAddress(ctx, actorID)
		assert.Nil(t, err)
		sdk.Abort(ctx, ferrors.USR_NOT_FOUND, "address:"+addr_.String())
	case 2:
		_, err := sdk.BalanceOf(ctx, abi.ActorID(10))
		assert.Nil(t, err)
		sdk.Abort(ctx, ferrors.USR_NOT_FOUND, "test_actor USR_NOT_FOUND")
	case 3:
		addr, _ := address.NewFromString("f1dwyrbh74hr5nwqv2gjedjyvgphxxkffxug4rkkq")
		actorID, err := sdk.ResolveAddress(ctx, addr)
		assert.Nil(t, err)
		sdk.Abort(ctx, ferrors.USR_NOT_FOUND, actorID.String()+"----")
	}
	return 0
}
