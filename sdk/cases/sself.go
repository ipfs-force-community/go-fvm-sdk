package main

import (
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

	originData := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}

	//set data
	stCid, err := sdk.Put(0xb220, 32, types.DAG_CBOR, originData)
	assert.Nil(t, err)

	//set root
	err = sdk.SetRoot(stCid)
	assert.Nil(t, err)

	//get root
	root, err := sdk.Root()
	assert.Nil(t, err)

	//get data
	data, err := sdk.Get(root)
	assert.Nil(t, err)
	assert.Equal(t, originData, data)

	//check balance
	actorBalance := sdk.CurrentBalance()
	assert.Nil(t, err)
	assert.Equal(t, "20", actorBalance.Big().String())

	//destruct
	toAddr, err := address.NewFromString("f1dwyrbh74hr5nwqv2gjedjyvgphxxkffxug4rkkq")
	assert.Nil(t, err)
	err = sdk.SelfDestruct(toAddr)
	assert.Nil(t, err)
	_, err = sdk.Root()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "unable to create ipld")
	return 0
}
