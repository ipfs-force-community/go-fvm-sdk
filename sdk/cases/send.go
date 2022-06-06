package main

import (
	"github.com/filecoin-project/go-address"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/testing"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/stretchr/testify/assert"
)

func main() {} //nolint

//go:export invoke
func Invoke(_ uint32) uint32 { //nolint
	t := testing.NewTestingT()
	defer t.CheckResult()

	methodsNum, err := sdk.MethodNumber()
	if err != nil {
		sdk.Abort(ferrors.USR_ILLEGAL_STATE, "unable to get method number")
	}
	switch methodsNum {
	case 1:
		// actor does not exist: 128788 (6: resource not found)
		addr, _ := address.NewFromString("f0128788")
		ret, err := sdk.Send(addr, 0, []byte{}, types.From64(1000))
		assert.Nil(t, err, "send %v", err)
		assert.Equal(t, 0, int(ret.ExitCode))
	case 2:
		addr, _ := address.NewFromString("f010000")
		ret, err := sdk.Send(addr, 0, []byte{}, types.From64(1))
		assert.Nil(t, err, "send %v", err)
		assert.Equal(t, 0, int(ret.ExitCode))
		assert.Equal(t, 0, int(ret.GasUsed))
		assert.Equal(t, "", string(ret.ReturnData))
	case 3:
		// sender does not have funds to transfer (balance 10, transfer 5000) (5: insufficient funds)
		addr, _ := address.NewFromString("f010000")
		ret, err := sdk.Send(addr, 0, []byte{}, types.From64(5000))
		assert.Nil(t, err, "send %v", err)
		assert.Equal(t, 0, int(ret.ExitCode))
	}

	return 0
}
