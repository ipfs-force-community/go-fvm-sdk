package fvm

import (
	"github.com/golang/mock/gomock"
)

type FakeReporter struct {
}

func (f *FakeReporter) Errorf(format string, args ...interface{}) {

}
func (f *FakeReporter) Fatalf(format string, args ...interface{}) {

}

func OpenExpact(in interface{}, out interface{}) {
	if mockFvmInstance == nil {
		InitMockFvm()
	}
	gomock.InOrder(
		mockFvmInstance.EXPECT().Open(in).Return(out),
	)

}
