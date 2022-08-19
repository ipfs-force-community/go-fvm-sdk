//go:build simulate

package contract

import (
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys"
	"testing"
)

func TestNewState(t *testing.T) {
	in := []byte{1, 2, 4}
	out1 := uint32(1)
	sys.MockFvmInstance.EXPECT().SelfRoot(in).Return(out1, nil)
	sys.MockFvmInstance.EXPECT().Abort(uint32(1), "1")
	//sys.SelfRoot(in1)
	NewState()
}
