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
	h, _ := mh.Sum([]byte("TEST"), mh.SHA3, 4)

	rootcid := cid.NewCidV1(7, h)
	basefee_ := big.NewInt(1)
	basefee := types.FromBig(&basefee_)
	totalFilCircSupply_ := big.NewInt(1)
	totalFilCircSupply := types.FromBig(&totalFilCircSupply_)
	currentBalance_ := big.NewInt(999)
	currentBalance := types.FromBig(&currentBalance_)
	return simulated.CreateSimulateEnv(callcontext, rootcid, &basefee, &totalFilCircSupply, &currentBalance)
}

func TestSayHello(t *testing.T) {
	_, ctx := newSimulated()

	testState := State{}
	a := testState.SayHello(ctx)
	assert.Equal(t, string(a), "1", "The two words should be the same.")

	b := testState.SayHello(ctx)
	assert.Equal(t, string(b), "2", "The two words should be the same.")

}
