package main

import (
	"context"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"

	"github.com/filecoin-project/go-address"
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

	f1AddrStr := "f1dwyrbh74hr5nwqv2gjedjyvgphxxkffxug4rkkq"
	f1Addr, err := address.NewFromString(f1AddrStr)
	assert.Nil(t, err)
	assert.NotNil(t, f1Addr)

	f4AddrStr := "t410forsxg5c7l5pv6x27l5pwczdeojsxg4yczlame"
	f4Addr, err := address.NewFromString(f4AddrStr)
	assert.Nil(t, err)

	switch methodNum {
	case 1:
		actorID, err := sdk.ResolveAddress(ctx, f4Addr)
		assert.Nil(t, err)
		addr_, err := sdk.LookupDelegatedAddress(ctx, actorID)
		assert.Nil(t, err)
		assert.Equal(t, f4Addr, addr_.String())
	case 2:
		_, err = sdk.NextActorAddress(ctx) //todo how to check next actor address
		assert.Nil(t, err)
	case 3:
		palcehoderAddr, err := sdk.GetActorCodeCid(ctx, f4Addr)
		assert.Nil(t, err)
		assert.Equal(t, "bafk2bzacedfvut2myeleyq67fljcrw4kkmn5pb5dpyozovj7jpoez5irnc3ro", palcehoderAddr.String())
	case 4:
		actorID, err := sdk.ResolveAddress(ctx, f4Addr)
		assert.Nil(t, err)
		codeCid, err := sdk.GetCodeCidForType(ctx, types.PlaceHolder)
		assert.Nil(t, err)
		err = sdk.CreateActor(ctx, actorID, codeCid, f4Addr) //not allow create actor in non-builtin-actors
		assert.ErrorIs(t, err, ferrors.Forbidden)
	case 5:
		codeCid, err := sdk.GetCodeCidForType(ctx, types.Miner)
		assert.Nil(t, err)
		actorType, err := sdk.GetBuiltinActorType(ctx, codeCid)
		assert.Nil(t, err)
		assert.Equal(t, types.Miner, actorType)
	case 6:
		actorID, err := sdk.ResolveAddress(ctx, f4Addr)
		assert.Nil(t, err)
		balance, err := sdk.BalanceOf(ctx, actorID)
		assert.Nil(t, err)
		assert.Equal(t, "10000000000000000000000", balance.String())
	}
	return 0
}
