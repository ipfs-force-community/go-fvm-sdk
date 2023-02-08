package main

import (
	"context"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"

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

	methodNum, err := sdk.MethodNumber(ctx)
	if err != nil {
		sdk.Abort(ctx, ferrors.USR_ILLEGAL_STATE, "unable to get method number")
	}

	switch methodNum {
	case 1:
		enabled := logger.Enabled(ctx)
		assert.Equal(t, true, enabled)

		logger.Log(ctx, "test")
	case 2:
		logger.StoreArtifact(ctx, "test_artifact", []byte{1, 2, 3, 4})
	}

	return 0
}
