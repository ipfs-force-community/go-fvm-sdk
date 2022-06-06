package main

import (
	"github.com/ipfs-force-community/go-fvm-sdk/sdk"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys"

	//"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/testing"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/stretchr/testify/assert"
)

func main() {} //nolint

//go:export invoke
func Invoke(_ uint32) uint32 { //nolint
	t := testing.NewTestingT()
	defer t.CheckResult()

	data := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	//create
	stCid, err := sdk.Put(0xb220, 32, types.DAG_CBOR, data)
	assert.Nil(t, err, "unable to put block %v", err)
	//cid assert
	assert.Equal(t, stCid.String(), "bafy2bzacedpfdhph46exiifylwgpd5dwukzg763u5burfjpcesqhblyt4k5wg")

	//open
	block, err := sdk.Get(stCid)
	assert.Nil(t, err, "unable to get block %v", err)
	//state
	blockId, err := sdk.PutBlock(types.DAG_CBOR, data)
	assert.Nil(t, err, "unable to putblock %v", err)

	state, err := sys.Stat(blockId)
	assert.Nil(t, err, "unable to inspect state for block %d reason %v", blockId, err)
	assert.Equal(t, state.Size, uint32(len(data)))
	assert.Equal(t, state.Codec, types.DAG_CBOR)

	assert.Equal(t, block[:state.Size], data)
	return 0
}
