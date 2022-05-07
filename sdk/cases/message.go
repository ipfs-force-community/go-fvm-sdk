package main

import (
	"github.com/ipfs-force-community/go-fvm-sdk/sdk"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/testing"
	"github.com/stretchr/testify/assert"
)


func main() {
}

//go:export invoke
func Invoke(_ uint32) uint32 {
	t := testing.NewTestingT()
	defer t.CheckResult()

	_, err := sdk.Caller()
	assert.Nil(t, err)
//	assert.Equal(t, caller, 1) todo unable to verify caller, its random value in tester

	receiver, err := sdk.Receiver()
	assert.Nil(t, err)
	assert.Equal(t, 10000, int(receiver))

	valueRecieved, err := sdk.ValueReceived()
	assert.Nil(t, err)
	assert.Equal(t, "0", valueRecieved.Big().String())

	method_num, err := sdk.MethodNumber()
	assert.Nil(t, err)
	assert.Equal(t, 1, method_num)

	return 0
}
