package simulated

import (
	"github.com/filecoin-project/go-address"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"
)

func SetActorAndAddress(actorId uint32, ActorState ActorState, addr address.Address) {
	DefaultFsm.actorMutex.Lock()
	defer DefaultFsm.actorMutex.Unlock()
	DefaultFsm.actors.Store(actorId, ActorState)
	DefaultFsm.address.Store(addr, actorId)
}

func SetAccount(actorId uint32, addr address.Address) {
	DefaultFsm.actorMutex.Lock()
	defer DefaultFsm.actorMutex.Unlock()
	DefaultFsm.actors.Store(actorId, ActorState{
		Code:     cid.Undef,
		State:    cid.Undef,
		Sequence: 0,
	})
	DefaultFsm.address.Store(addr, actorId)
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
