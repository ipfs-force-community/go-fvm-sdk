package contract

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"

	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/ipfs/go-cid"

	"github.com/filecoin-project/go-state-types/abi"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys/simulated"
)

func TestQueueCall(t *testing.T) {
	simu, ctx := simulated.CreateEmptySimulator()
	ownerId := abi.ActorID(10)

	{
		simu.SetActor(ownerId, sdk.MustAddressFromActorId(ownerId), builtin.Actor{
			Code:       cid.Undef,
			Head:       cid.Undef,
			CallSeqNum: 0,
			Balance:    big.NewInt(100),
		})

		//save state
		simu.SetMessageContext(&types.MessageContext{Origin: ownerId, Caller: abi.ActorID(1)})
		Constructor(ctx)
	}

	{
		//queue call
		callState := &State{}
		sdk.LoadState(ctx, callState)
		//call
		toAddrID := abi.ActorID(200)
		simu.SetActor(ownerId, sdk.MustAddressFromActorId(toAddrID), builtin.Actor{
			Code:       cid.Undef,
			Head:       cid.Undef,
			CallSeqNum: 0,
			Balance:    big.NewInt(0),
		})

		call := &Call{
			To:        sdk.MustAddressFromActorId(toAddrID),
			Method:    0,
			Value:     abi.NewTokenAmount(0),
			Params:    nil,
			TimeStamp: 5000,
		}

		//expect 100+10,100+1000

		//below min delay
		simu.SetNetworkContext(&types.NetworkContext{
			Timestamp: 3000,
		})
		simu.SetMessageContext(&types.MessageContext{
			Caller: ownerId,
		})
		_, err := callState.Queue(ctx, call)
		assert.NotNil(t, err)

		//exceed max delay
		simu.SetNetworkContext(&types.NetworkContext{
			Timestamp: 8000,
		})
		simu.SetMessageContext(&types.MessageContext{
			Caller: ownerId,
		})
		_, err = callState.Queue(ctx, call)
		assert.NotNil(t, err)

		//success
		simu.SetMessageContext(&types.MessageContext{
			Caller: ownerId,
		})
		simu.SetNetworkContext(&types.NetworkContext{
			Timestamp: 4500,
		})
		id, err := callState.Queue(ctx, call)
		assert.Nil(t, err)
		key, _, err := call.Key()
		assert.Nil(t, err)
		assert.Equal(t, key.Key(), string(id))
	}
}

func TestExecuteCall(t *testing.T) {
	simu, ctx := simulated.CreateEmptySimulator()
	ownerId := abi.ActorID(10)

	{
		simu.SetActor(ownerId, sdk.MustAddressFromActorId(ownerId), builtin.Actor{
			Code:       cid.Undef,
			Head:       cid.Undef,
			CallSeqNum: 0,
			Balance:    big.NewInt(100),
		})

		//save state
		simu.SetMessageContext(&types.MessageContext{Origin: ownerId, Caller: abi.ActorID(1)})
		Constructor(ctx)
	}

	{
		//queue call
		callState := &State{}
		sdk.LoadState(ctx, callState)
		//call
		toAddrID := abi.ActorID(200)
		simu.SetActor(ownerId, sdk.MustAddressFromActorId(toAddrID), builtin.Actor{
			Code:       cid.Undef,
			Head:       cid.Undef,
			CallSeqNum: 0,
			Balance:    big.NewInt(0),
		})

		call := &Call{
			To:        sdk.MustAddressFromActorId(toAddrID),
			Method:    0,
			Value:     abi.NewTokenAmount(0),
			Params:    nil,
			TimeStamp: 5000,
		}
		//success
		simu.SetMessageContext(&types.MessageContext{
			Caller: ownerId,
		})
		simu.SetNetworkContext(&types.NetworkContext{
			Timestamp: 4500,
		})
		id, err := callState.Queue(ctx, call)
		assert.Nil(t, err)
		key, _, err := call.Key()
		assert.Nil(t, err)
		assert.Equal(t, key.Key(), string(id))

		//execute
		simu.SetNetworkContext(&types.NetworkContext{
			Timestamp: 5100,
		})
		simu.ExpectSend(simulated.SendMock{
			To:     sdk.MustAddressFromActorId(toAddrID),
			Method: 0,
			Params: nil,
			Value:  abi.NewTokenAmount(0),
			Out:    types.SendResult{},
		})
		_, err = callState.Execute(ctx, &id)
		assert.Nil(t, err)
	}
}
