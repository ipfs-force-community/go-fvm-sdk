package sdk

import "github.com/ipfs-force-community/go-fvm-sdk/sdk/sys"

func Charge(name string, compute uint64) error {
	return sys.Charge(name, compute)
}
