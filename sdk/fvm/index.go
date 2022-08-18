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

// 执行 go  generate生成文件

//go:generate mockgen -destination ./mock_scheme.go -package=fvm -source ./index.go
type Fvm interface {
	Open(id cid.Cid) (*types.IpldOpen, error)
}

var MockFvmInstance *MockFvm
var mockFvmInstanceCtl *gomock.Controller

func EpochFinish() {
	mockFvmInstanceCtl.Finish()
}
func init() {

	t := FakeReporter{}
	mockFvmInstanceCtl = gomock.NewController(&t)
	// defer ctl.Finish()
	MockFvmInstance = NewMockFvm(mockFvmInstanceCtl)

}
