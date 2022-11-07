package sdk

import (
	"context"

	"github.com/filecoin-project/go-state-types/abi"

	addr "github.com/filecoin-project/go-address"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys"
	"github.com/ipfs/go-cid"
)

// Root Get the IPLD root CID. Fails if the actor doesn't have state (before the first call to
// `set_root` and after actor deletion).
func Root(ctx context.Context) (cid.Cid, error) {
	return sys.SelfRoot(ctx)
}

// SetRoot set the actor's state-tree root.
//
// Fails if:
//
// - The new root is not in the actor's "reachable" set.
// - Fails if the actor has been deleted.
func SetRoot(ctx context.Context, c cid.Cid) error {
	return sys.SelfSetRoot(ctx, c)
}

// CurrentBalance gets the current balance for the calling actor.
func CurrentBalance(ctx context.Context) abi.TokenAmount {
	tok, err := sys.SelfCurrentBalance(ctx)
	if err != nil {
		panic(err.Error())
	}
	return *tok
}

// SelfDestruct destroys the calling actor, sending its current balance
// to the supplied address, which cannot be itself.
//
// Fails if the beneficiary doesn't exist or is the actor being deleted.
func SelfDestruct(ctx context.Context, addr addr.Address) error {
	return sys.SelfDestruct(ctx, addr)
}
