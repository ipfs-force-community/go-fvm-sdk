//go:build simulate

package sys

import (
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/fvm"
)

func Open(id cid.Cid) (*types.IpldOpen, error) {
	return fvm.MockFvmInstance.Open(id)
}

func Create(codec uint64, data []byte) (uint32, error) {
	return fvm.MockFvmInstance.Create(codec, data)
}

func Read(id uint32, offset uint32, buf []byte) (uint32, error) {
	return fvm.MockFvmInstance.Read(id, offset, buf)
}

func Stat(id uint32) (*types.IpldStat, error) {
	return fvm.MockFvmInstance.Stat(id)
}

func BlockLink(id uint32, hashFun uint64, hashLen uint32, cidBuf []byte) (uint32, error) {
	return fvm.MockFvmInstance.BlockLink(id, hashFun, hashLen, cidBuf)
}
