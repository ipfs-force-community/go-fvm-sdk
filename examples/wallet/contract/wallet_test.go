package contract

import (
	"testing"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"

	"github.com/stretchr/testify/assert"

	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/ipfs/go-cid"

	"github.com/filecoin-project/go-state-types/abi"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys/simulated"
)

func TestWallet(t *testing.T) {
	simu, ctx := simulated.CreateEmptySimulator()
	actorId := abi.ActorID(50)
	ownerId := abi.ActorID(100)
	callerId := abi.ActorID(200)

	{
		simu.SetActor(ownerId, sdk.MustAddressFromActorId(ownerId), builtin.Actor{
			Code:       cid.Undef,
			Head:       cid.Undef,
			CallSeqNum: 0,
			Balance:    big.NewInt(100),
		})

		simu.SetActor(ownerId, sdk.MustAddressFromActorId(callerId), builtin.Actor{
			Code:       cid.Undef,
			Head:       cid.Undef,
			CallSeqNum: 0,
			Balance:    big.NewInt(0),
		})
		//save state
		walletState := &State{
			Owner: ownerId,
		}

		headCid := sdk.SaveState(ctx, walletState)

		simu.SetActor(actorId, sdk.MustAddressFromActorId(actorId), builtin.Actor{
			Code:       cid.Undef,
			Head:       headCid,
			CallSeqNum: 0,
			Balance:    big.NewInt(100),
		})
	}

	{
		//get balance
		walletState := &State{}
		sdk.LoadState(ctx, walletState)

		simu.SetMessageContext(&types.MessageContext{
			Receiver: actorId,
		})
		balance, err := walletState.GetBalance(ctx)
		assert.Nil(t, err)
		assert.Equal(t, int64(100), balance.Int64())
	}

	{
		//withdraw
		walletState := &State{}
		sdk.LoadState(ctx, walletState)

		simu.SetMessageContext(&types.MessageContext{
			Caller: ownerId,
		})
		simu.ExpectSend(simulated.SendMock{
			To:     sdk.MustAddressFromActorId(ownerId),
			Method: 0,
			Params: nil,
			Value:  big.NewInt(50),
			Out:    types.SendResult{},
		})

		withdrawAmount := big.NewInt(50)
		err := walletState.Withdraw(ctx, &withdrawAmount)
		assert.Nil(t, err)
	}
}
