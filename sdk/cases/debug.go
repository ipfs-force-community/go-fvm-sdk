package main

import (
	"github.com/ipfs-force-community/go-fvm-sdk/sdk"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/testing"
	"github.com/stretchr/testify/assert"
)

func main() {} //nolint

//go:export invoke
func Invoke(_ uint32) uint32 { //nolint
	t := testing.NewTestingT()
	defer t.CheckResult()

	logger, err := sdk.NewLogger()
	assert.Nil(t, err, "create debug logger %v", err)

	enabled := logger.Enabled()
	assert.Equal(t, true, enabled)

	err = logger.Log("")
	assert.Nil(t, err)

	return 0
}
