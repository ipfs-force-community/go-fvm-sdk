package simulated

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin/v9/migration"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
)

func SetActorAndAddress(actorId uint32, actorState migration.Actor, addr address.Address) {
	DefaultFsm.actorMutex.Lock()
	defer DefaultFsm.actorMutex.Unlock()
	DefaultFsm.actorsMap.Store(actorId, actorState)
	DefaultFsm.addressMap.Store(addr, actorId)
}

type SendMock struct {
	to     address.Address
	method uint64
	params uint32
	value  big.Int
	out    types.Send
}

func SetSend(mock ...SendMock) {
	temp := make([]SendMock, 0)
	for _, v := range mock {
		_, ok := DefaultFsm.sendMatch(v.to, v.method, v.params, v.value)
		if !ok {
			temp = append(temp, v)
		}
	}
	DefaultFsm.SendList = append(DefaultFsm.SendList, temp...)

}

func SetAccount(actorId uint32, addr address.Address, actor migration.Actor) {
	DefaultFsm.actorMutex.Lock()
	defer DefaultFsm.actorMutex.Unlock()

	DefaultFsm.actorsMap.Store(actorId, actor)
	DefaultFsm.addressMap.Store(addr, actorId)
}

func SetBaseFee(ta abi.TokenAmount) {
	DefaultFsm.baseFee = &ta
}

func SetTotalFilCircSupply(ta abi.TokenAmount) {
	DefaultFsm.totalFilCircSupply = &ta
}

func SetCurrentBalance(ta abi.TokenAmount) {
	DefaultFsm.totalFilCircSupply = &ta
}

func SetCallContext(callcontext *types.InvocationContext) {
	DefaultFsm.callContext = callcontext
}
