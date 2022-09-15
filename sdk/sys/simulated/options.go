package simulated

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"
)

func SetActorAndAddress(actorId uint32, ActorState ActorState, addr address.Address) {
	DefaultFsm.actorMutex.Lock()
	defer DefaultFsm.actorMutex.Unlock()
	DefaultFsm.actorsMap.Store(actorId, ActorState)
	DefaultFsm.addressMap.Store(addr, actorId)
}

type SendMock struct {
	to     address.Address
	method uint64
	params uint32
	value  big.Int
	out    types.Send
}

func SetSend(mock ...SendMock) bool {
	temp := make([]SendMock, 0)
	for _, v := range mock {
		_, ok := DefaultFsm.sendMatch(v.to, v.method, v.params, v.value)
		if !ok {
			temp = append(temp, v)
		}
	}
	DefaultFsm.SendList = append(DefaultFsm.SendList, temp...)
	return true

}

func SetAccount(actorId uint32, addr address.Address) {
	DefaultFsm.actorMutex.Lock()
	defer DefaultFsm.actorMutex.Unlock()
	DefaultFsm.actorsMap.Store(actorId, ActorState{
		Code:     cid.Undef,
		State:    cid.Undef,
		Sequence: 0,
	})
	DefaultFsm.addressMap.Store(addr, actorId)
}

func SetBaseFee(ta types.TokenAmount) {
	DefaultFsm.baseFee = &ta
}

func SetTotalFilCircSupply(ta types.TokenAmount) {
	DefaultFsm.totalFilCircSupply = &ta
}

func SetCurrentBalance(ta types.TokenAmount) {
	DefaultFsm.totalFilCircSupply = &ta
}
