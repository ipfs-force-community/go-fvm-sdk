package fvm

import (
	gomock "github.com/golang/mock/gomock"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"
)

type FakeReporter struct {
}

func (f *FakeReporter) Errorf(format string, args ...interface{}) {

}
func (f *FakeReporter) Fatalf(format string, args ...interface{}) {

}

type Fvm interface {
	Open(id cid.Cid) (*types.IpldOpen, error)
}

var mockFvmInstance *MockFvm
var mockFvmInstanceCtl *gomock.Controller

func EpochFinish() {
	mockFvmInstanceCtl.Finish()
}
func InitMockFvm(t gomock.TestReporter) {
	mockFvmInstanceCtl = gomock.NewController(t)
	// defer ctl.Finish()
	mockFvmInstance = NewMockFvm(mockFvmInstanceCtl)

}
