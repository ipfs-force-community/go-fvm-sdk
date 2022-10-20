package contract

import (
	"fmt"
	"testing"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk"
	"github.com/stretchr/testify/assert"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys/simulated"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
)

func TestSayHello(t *testing.T) {
	_, ctx := simulated.CreateSimulateEnv(&types.InvocationContext{}, abi.NewTokenAmount(1), abi.NewTokenAmount(1))
	{
		//save state
		helloState := &State{
			Count: 0,
		}
		sdk.SaveState(ctx, helloState)
	}

	for i := 0; i < 10; i++ {
		helloState := &State{}
		sdk.LoadState(ctx, helloState)
		bytes := helloState.SayHello(ctx)
		assert.Equal(t, string(bytes), fmt.Sprintf("Hello World %d", i+1))
	}
}
