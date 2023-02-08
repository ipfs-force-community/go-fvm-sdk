package simulated

import (
	"testing"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
)

func TestSetSend(t *testing.T) {
	data := []byte{1, 1, 1, 1, 1, 1, 1}
	sendmockcase := make([]SendMock, 0)
	sendmockcase = append(sendmockcase, SendMock{address.Undef, 1, data, big.NewInt(1), types.SendResult{}})
	defaultfsm := FvmSimulator{}
	blkId := defaultfsm.blockCreate(types.DAGCBOR, data)
	defaultfsm.ExpectSend(sendmockcase...)
	_, err := defaultfsm.sendMatch(address.Undef, 1, blkId, big.NewInt(1))
	if err != nil {
		t.Errorf("match is failed")
	}
}
