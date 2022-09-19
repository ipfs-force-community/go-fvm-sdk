package simulated

import (
	"testing"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
)

func TestSetSend(t *testing.T) {
	sendMockCase := make([]SendMock, 0)
	sendMockCase = append(sendMockCase, SendMock{address.Undef, 1, 1, big.NewInt(1), types.Send{}})
	SetSend(sendMockCase...)
	_, ok := DefaultFsm.sendMatch(address.Undef, 1, 1, big.NewInt(1))
	if ok != true {
		t.Errorf("match is failed")
	}
}
