//nolint:unparam
package simulated

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
)

type SendMock struct {
	to     address.Address
	method uint64
	params uint32
	value  big.Int
	out    types.Send
}

func (s *FvmSimulator) Send(to address.Address, method uint64, params uint32, value types.TokenAmount) (*types.Send, error) {
	send, ok := s.sendMatch(to, method, params, *value.Big())
	if ok {
		return send, nil
	}
	return nil, ErrorKeyMatchFail
}

func (s *FvmSimulator) SetSend(mock ...SendMock) {
	temp := make([]SendMock, 0)
	for _, v := range mock {
		_, ok := s.sendMatch(v.to, v.method, v.params, v.value)
		if !ok {
			temp = append(temp, v)
		}
	}
	s.sendList = append(s.sendList, temp...)

}
