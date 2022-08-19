//go:build simulate

package sys

import (
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/fvm"
)

func Send(to address.Address, method uint64, params uint32, value types.TokenAmount) (*types.Send, error) {
	return fvm.MockFvmInstance.send(to, method, params, value)
}
