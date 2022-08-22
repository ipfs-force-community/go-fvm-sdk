package simulated

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"
)

func TestOpenExpect(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockMale := NewMockSimulated(ctl)
	result := types.IpldOpen{1, 2, 3}
	in, _ := cid.Cast([]byte("bafy2bzacecdjkk2tzogitpcybu3eszr4uptrjogstqmyt6u4q2p3hh4chmf3i"))
	gomock.InOrder(
		mockMale.EXPECT().Open(in).Return(&result, nil))
	in2, _ := cid.Cast([]byte("bafy2bzacecdjkk2tzogitpcybu3eszr4uptrjogstqmyt6u4q2p3hh4chmf3i"))
	var err error
	mockMale.Open(in2)
	if err != nil {
		t.Errorf("user.GetUserInfo err: %v", err)
	}
}
