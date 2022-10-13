//nolint:unparam
package simulated

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin/v9/migration"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
)

func (s *Fsm) SetActorAndAddress(actorID uint32, actorState migration.Actor, addr address.Address) {
	s.actorMutex.Lock()
	defer s.actorMutex.Unlock()
	s.actorsMap.Store(actorID, actorState)
	s.addressMap.Store(addr, actorID)
}

type SendMock struct {
	to     address.Address
	method uint64
	params uint32
	value  big.Int
	out    types.Send
}

func (s *Fsm) SetSend(mock ...SendMock) {
	temp := make([]SendMock, 0)
	for _, v := range mock {
		_, ok := s.sendMatch(v.to, v.method, v.params, v.value)
		if !ok {
			temp = append(temp, v)
		}
	}
	s.sendList = append(s.sendList, temp...)

}

func (s *Fsm) SetAccount(actorID uint32, addr address.Address, actor migration.Actor) {
	s.actorMutex.Lock()
	defer s.actorMutex.Unlock()

	s.actorsMap.Store(actorID, actor)
	s.addressMap.Store(addr, actorID)
}

func (s *Fsm) SetBaseFee(ta big.Int) {
	amount, _ := types.FromString(ta.String())
	s.baseFee = &amount
}

func (s *Fsm) SetTotalFilCircSupply(ta big.Int) {
	amount, _ := types.FromString(ta.String())
	s.totalFilCircSupply = &amount
}

func (s *Fsm) SetCurrentBalance(ta big.Int) {
	amount, _ := types.FromString(ta.String())
	s.currentBalance = &amount
}

func (s *Fsm) SetCallContext(callcontext *types.InvocationContext) {
	s.callContext = callcontext
}
