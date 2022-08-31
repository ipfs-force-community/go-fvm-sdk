//go:build simulate
// +build simulate

package sys

import (
	addr "github.com/filecoin-project/go-address"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys/simulated"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"
)

func SelfRoot() (cid.Cid, error) {
	return simulated.DefaultFsm.SelfRoot()
}

func SelfSetRoot(id cid.Cid) error {
	return simulated.DefaultFsm.SelfSetRoot(id)
}

func SelfCurrentBalance() (*types.TokenAmount, error) {
	return simulated.DefaultFsm.SelfCurrentBalance()
}

func SelfDestruct(addr addr.Address) error {
	return simulated.DefaultFsm.SelfDestruct(addr)
}
