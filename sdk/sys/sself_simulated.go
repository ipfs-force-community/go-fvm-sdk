//go:build simulate

package sys

import (
	"github.com/ipfs/go-cid"
	addr "github.com/filecoin-project/go-address"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys/internal/simulated"
)

func SelfRoot(cidBuf []byte) (uint32, error) {
	return simulated.MockFvmInstance.SelfRoot(cidBuf)
}

func SelfSetRoot(id cid.Cid) error {
	return simulated.MockFvmInstance.SelfSetRoot(id)

}

func SelfCurrentBalance() (*types.TokenAmount, error) {
	return simulated.MockFvmInstance.SelfCurrentBalance()
}

func SelfDestruct(addr addr.Address) error {
	return simulated.MockFvmInstance.SelfDestruct(addr)
}
