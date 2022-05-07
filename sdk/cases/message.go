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

	method_num, err := sdk.MethodNumber()
	assert.Nil(t, err)
	assert.Equal(t, 1, int(method_num))

	valueRecieved, err := sdk.ValueReceived()
	assert.Nil(t, err)
	assert.Equal(t, "10", valueRecieved.Big().String())
	return 0
}
