package sdk

import (
	"fmt"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"
)

/// Get the IPLD root CID. Fails if the actor doesn't have state (before the first call to
/// `set_root` and after actor deletion).
func Root() (cid.Cid, error) {
	// I really hate this CID interface. Why can't I just have bytes?
	cidBuf := make([]byte, types.MAX_CID_LEN)
	cidBufLen, err := sys.SelfRoot(cidBuf)
	if int(cidBufLen) > len(cidBuf) {
		// TODO: re-try with a larger buffer?
		panic(fmt.Sprintf("CID too big: %d > %d", cidBufLen, len(cidBuf)))
	}
	_, cid, err := cid.CidFromBytes(cidBuf)
	return cid, err
}

/// Set the actor's state-tree root.
///
/// Fails if:
///
/// - The new root is not in the actor's "reachable" set.
/// - Fails if the actor has been deleted.
func SetRoot(c cid.Cid) error {
	return sys.SelfSetRoot(c)
}

/// Gets the current balance for the calling actor.
func CurrentBalance() *types.TokenAmount {
	tok, err := sys.SelfCurrentBalance()
	if err != nil {
		panic(err.Error())
	}
	return tok
}

func SelfDestruct(addr types.Address) error {
	return sys.SelfDestruct(addr)
}
