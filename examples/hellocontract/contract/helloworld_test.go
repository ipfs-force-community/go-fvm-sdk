//go:build simulate

package contract

import (
	"fmt"
	gomock "github.com/golang/mock/gomock"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"
	"testing"
	//mh "github.com/multiformats/go-multihash"
)

func TestSayHello(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("%v\n", err)
		}
	}()
	opresult := &types.IpldOpen{uint64(1), 1, 1}
	cidin, _ := cid.Decode("bafy2bzacecdjkk2tzogitpcybu3eszr4uptrjogstqmyt6u4q2p3hh4chmf3i")
	sys.SimulatedInstance.EXPECT().SelfRoot().Return(cidin, nil)
	sys.SimulatedInstance.EXPECT().Abort(gomock.Any(), gomock.Any())
	sys.SimulatedInstance.EXPECT().Open(cidin).Return(opresult, nil)
	sys.SimulatedInstance.EXPECT().Read(types.BlockID(1), uint32(0), uint32(1)).Return([]byte{}, uint32(0), nil)
}
