//go:build simulate

package sys

import (
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/fvm"
)

func SelfRoot(cidBuf []byte) (uint32, error) {
	return fvm.MockFvmInstance.SelfRoot(cidBuf)
}

func SelfSetRoot(id cid.Cid) error {
	return fvm.MockFvmInstance.SelfRoot(id)

}

func SelfCurrentBalance() (*types.TokenAmount, error) {
	return fvm.MockFvmInstance.selfCurrentBalance()
}

func SelfDestruct(addr addr.Address) error {
	return fvm.MockFvmInstance.selfDestruct(addr)
}
