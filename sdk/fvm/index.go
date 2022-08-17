package fvm

import (
	gomock "github.com/golang/mock/gomock"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"
)

type Fvm interface {
	Open(id cid.Cid) (*types.IpldOpen, error)
}

var mockFvmInstance *MockFvm

func InitMockFvm() {
	rep := FakeReporter{}
	ctl := gomock.NewController(&rep)
	defer ctl.Finish()
	mockFvm := NewMockFvm(ctl)
	mockFvmInstance = mockFvm
}
