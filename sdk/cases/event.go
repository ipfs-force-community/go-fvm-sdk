package main

import (
	"context"

	"github.com/stretchr/testify/assert"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/testing"
)

func main() {} //nolint

//go:export invoke
func Invoke(_ uint32) uint32 { //nolint
	t := testing.NewTestingT()
	defer t.CheckResult()
	ctx := context.Background()

	err := sdk.EmitEvent(ctx, types.ActorEvent{Entries: []*types.Entry{
		{
			Flags: 0,
			Key:   "111",
			Codec: types.IPLDRAW,
			Value: nil,
		},
		{
			Flags: 0,
			Key:   "222",
			Codec: types.IPLDRAW,
			Value: nil,
		},
		{
			Flags: 0,
			Key:   "333",
			Codec: types.IPLDRAW,
			Value: nil,
		},
	}})
	assert.Nil(t, err)
	return 0
}
