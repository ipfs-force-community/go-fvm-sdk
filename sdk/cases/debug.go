package main

import (
	"context"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/testing"
	"github.com/stretchr/testify/assert"
)

func main() {} //nolint

//go:export invoke
func Invoke(_ uint32) uint32 { //nolint
	t := testing.NewTestingT()
	defer t.CheckResult()

	ctx := context.Background()
	logger, err := sdk.NewLogger()
	assert.Nil(t, err, "create debug logger %v", err)

	enabled := logger.Enabled(ctx)
	assert.Equal(t, true, enabled)

	err = logger.Log(ctx, "")
	assert.Nil(t, err)

	return 0
}
