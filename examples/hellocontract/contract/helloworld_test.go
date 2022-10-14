package contract

import (
	"context"

	"github.com/ipfs/go-cid"

	"testing"

	"github.com/filecoin-project/go-state-types/big"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys/simulated"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/stretchr/testify/assert"

	mh "github.com/multiformats/go-multihash"
)

func newSimulated() (*simulated.FvmSimulator, context.Context) {
	callcontext := &types.InvocationContext{}
	return simulated.CreateSimulateEnv(callcontext, big.NewInt(1), big.NewInt(1), big.NewInt(1))
}

func TestSayHello(t *testing.T) {
	_, ctx := newSimulated()

	testState := State{}
	a := testState.SayHello(ctx)
	assert.Equal(t, string(a), "1", "The two words should be the same.")

	b := testState.SayHello(ctx)
	assert.Equal(t, string(b), "2", "The two words should be the same.")

}
