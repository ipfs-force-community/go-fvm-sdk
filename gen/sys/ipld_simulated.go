//go:build simulatedd
// +build simulatedd

package sys

import (
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys/simulated"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"
)

func Open(id cid.Cid) (*types.IpldOpen, error) {
	return simulated.DefaultFsm.Open(id)
}

func Create(codec uint64, data []byte) (uint32, error) {
	return simulated.DefaultFsm.Create(codec, data)
}

func Read(id uint32, offset, size uint32) ([]byte, uint32, error) {
	return simulated.DefaultFsm.Read(id, offset, size)
}

func Stat(id uint32) (*types.IpldStat, error) {
	return simulated.DefaultFsm.Stat(id)
}

func BlockLink(id uint32, hashFun uint64, hashLen uint32, cidBuf []byte) (cided cid.Cid, err error) {
	return simulated.DefaultFsm.BlockLink(id, hashFun, hashLen, cidBuf)
}
