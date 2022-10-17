package main

import (
	"context"

	"github.com/filecoin-project/go-address"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk"
	"github.com/stretchr/testify/assert"

	//"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/testing"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
)

func main() {} //nolint

//go:export invoke
func Invoke(_ uint32) uint32 { //nolint
	t := testing.NewTestingT()
	defer t.CheckResult()

	ctx := context.Background()
	originData := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}

	//set data
	stCid, err := sdk.Put(ctx, 0xb220, 32, types.DAGCbor, originData)
	assert.Nil(t, err)

	//set root
	err = sdk.SetRoot(ctx, stCid)
	assert.Nil(t, err)

	//get root
	root, err := sdk.Root(ctx)
	assert.Nil(t, err)

	//get data
	data, err := sdk.Get(ctx, root)
	assert.Nil(t, err)
	assert.Equal(t, originData, data)

	//check balance
	actorBalance := sdk.CurrentBalance(ctx)
	assert.Nil(t, err)
	assert.Equal(t, "20", actorBalance.Big().String())

	//destruct
	toAddr, err := address.NewFromString("f1dwyrbh74hr5nwqv2gjedjyvgphxxkffxug4rkkq")
	assert.Nil(t, err)
	err = sdk.SelfDestruct(ctx, toAddr)
	assert.Nil(t, err)
	_, err = sdk.Root(ctx)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "unable to create ipld")
	return 0
}
