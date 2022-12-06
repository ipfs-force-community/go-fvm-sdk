package contract

import (
	"encoding/hex"
	"testing"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/adt"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys/simulated"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestX(t *testing.T) {
	println(hex.EncodeToString(sdk.MustCborMarshal(&ConstructorReq{
		Name:        "Venus",
		Symbol:      "V",
		Granularity: 1,
		Supply:      abi.NewTokenAmount(100000),
	})))
	println(hex.EncodeToString(sdk.MustCborMarshal(types.CborString("Venus"))))
	println(hex.EncodeToString(sdk.MustCborMarshal(types.CborString("V"))))
	println(hex.EncodeToString(sdk.MustCborMarshal(types.CborUint(1))))
	println(hex.EncodeToString(sdk.MustCborMarshal(newAmount(100000))))
	addr1, _ := address.NewFromString("f1m674sjwmga36qi3wkowt3wozwpahrkdlvd4tpci")
	println(hex.EncodeToString(sdk.MustCborMarshal(&MintParams{
		InitialOwner: addr1,
		Amount:       abi.NewTokenAmount(100),
	})))
	println(hex.EncodeToString(sdk.MustCborMarshal(newAmount(100100))))
}

func newAmount(vv int64) *abi.TokenAmount {
	v := abi.NewTokenAmount(vv)
	return &v
}

var supply int64 = 100000

func setup(t *testing.T, fromInitBalance abi.TokenAmount, granularity uint64) (*simulated.FvmSimulator, address.Address, address.Address, address.Address) {
	simulator, ctx := simulated.CreateSimulateEnv(&types.MessageContext{}, &types.NetworkContext{}, abi.NewTokenAmount(1), abi.NewTokenAmount(1))
	fromActor := abi.ActorID(1)
	fromAddr, err := simulated.NewF1Address()
	assert.NoError(t, err)
	simulator.SetActor(fromActor, fromAddr, builtin.Actor{Code: simulated.AccountCid})

	approvalActor := abi.ActorID(2)
	approvalAddr, err := simulated.NewF1Address()
	assert.NoError(t, err)
	simulator.SetActor(approvalActor, approvalAddr, builtin.Actor{Code: simulated.AccountCid})

	toActor := abi.ActorID(3)
	toAddr, err := simulated.NewF1Address()
	assert.NoError(t, err)
	simulator.SetActor(toActor, toAddr, builtin.Actor{Code: simulated.AccountCid})

	balanceMap, err := adt.MakeEmptyMap(adt.AdtStore(ctx), adt.BalanceTableBitwidth)
	assert.NoError(t, err)
	emptyRoot, err := balanceMap.Root()
	assert.NoError(t, err)
	assert.NoError(t, balanceMap.Put(types.ActorKey(fromActor), &fromInitBalance))
	balanceRoot, err := balanceMap.Root()
	assert.Nil(t, err)

	erc20State := &Frc46Token{Name: "Ep Coin", Symbol: "EP", Granularity: granularity, Supply: abi.NewTokenAmount(supply), Balances: balanceRoot, Allowances: emptyRoot}
	_ = sdk.SaveState(ctx, erc20State) //Save state

	// set info of context
	simulator.SetMessageContext(&types.MessageContext{
		Caller: fromActor,
	})
	return simulator, fromAddr, approvalAddr, toAddr
}

func TestFrc46TokenGetter(t *testing.T) {
	simulator, ctx := simulated.CreateSimulateEnv(&types.MessageContext{}, &types.NetworkContext{}, abi.NewTokenAmount(1), abi.NewTokenAmount(1))
	addr, err := simulated.NewF1Address()
	assert.NoError(t, err)
	simulator.SetActor(abi.ActorID(1), addr, builtin.Actor{})

	empMap, err := adt.MakeEmptyMap(adt.AdtStore(ctx), adt.BalanceTableBitwidth)
	assert.Nil(t, err)
	emptyBalance, err := empMap.Root()
	assert.Nil(t, err)

	frc46State := &Frc46Token{
		Name:        "EP Coin",
		Symbol:      "EP",
		Granularity: 1,
		Supply:      abi.NewTokenAmount(100000),
		Balances:    emptyBalance,
		Allowances:  emptyBalance,
	}

	_ = sdk.SaveState(ctx, frc46State)
	newFrc46State := &Frc46Token{}
	sdk.LoadState(ctx, newFrc46State)
	t.Run("get name", func(t *testing.T) {
		assert.Equal(t, "EP Coin", newFrc46State.GetName(ctx))
	})

	t.Run("get symbol", func(t *testing.T) {
		assert.Equal(t, "EP", newFrc46State.GetSymbol(ctx))
	})

	t.Run("get granularity", func(t *testing.T) {
		assert.Equal(t, uint64(1), newFrc46State.GetGranularity(ctx))
	})

	t.Run("get supply", func(t *testing.T) {
		assert.Equal(t, uint64(100000), newFrc46State.GetTotalSupply(ctx).Uint64())
	})
}

func TestErc20TokenGetBalanceOf(t *testing.T) {
	simulator, ctx := simulated.CreateSimulateEnv(&types.MessageContext{}, &types.NetworkContext{}, abi.NewTokenAmount(1), abi.NewTokenAmount(1))
	actor := abi.ActorID(1)
	addr, err := simulated.NewF1Address()
	assert.NoError(t, err)
	simulator.SetActor(actor, addr, builtin.Actor{})

	balanceMap, err := adt.MakeEmptyMap(adt.AdtStore(ctx), adt.BalanceTableBitwidth)
	assert.Nil(t, err)
	emptyRoot, err := balanceMap.Root()

	assert.Nil(t, balanceMap.Put(types.ActorKey(actor), simulated.NewPtrTokenAmount(100)))
	balanceRoot, err := balanceMap.Root()
	assert.Nil(t, err)

	erc20State := &Frc46Token{Name: "Ep Coin", Symbol: "EP", Granularity: 1, Supply: abi.NewTokenAmount(100000), Balances: balanceRoot, Allowances: emptyRoot}
	sdk.SaveState(ctx, erc20State) //Save state

	got, err := erc20State.BalanceOf(ctx, &addr)
	assert.Nil(t, err)
	assert.Equal(t, got.Uint64(), uint64(100))

	t.Run("got zero for not found address", func(t *testing.T) {
		addr2, err := simulated.NewF1Address()
		assert.NoError(t, err)
		got, err = erc20State.BalanceOf(ctx, &addr2)
		assert.Nil(t, err)
		assert.Equal(t, got.Uint64(), uint64(0))
	})
}

func TestFrc46TokenTransfer(t *testing.T) {
	t.Run("successful", func(t *testing.T) {
		simulator, fromAddr, toAddr, _ := setup(t, abi.NewTokenAmount(1000), 1)
		simulator.ExpectSend(simulated.SendMock{
			To:     toAddr,
			Method: RECEIVERHOOKMETHODNUM,
			Params: []byte{130, 26, 133, 34, 59, 223, 73, 134, 1, 2, 1, 66, 0, 100, 64, 64},
			Value:  big.Zero(),
			Out:    types.SendResult{},
		})

		var newState Frc46Token
		sdk.LoadState(simulator.Context, &newState)

		_, err := newState.Transfer(simulator.Context, &TransferParams{
			To:           toAddr,
			Amount:       abi.NewTokenAmount(100),
			OperatorData: nil,
		})
		assert.NoError(t, err)

		fromBalance, err := newState.BalanceOf(simulator.Context, &fromAddr)
		assert.NoError(t, err)
		assert.Equal(t, uint64(900), fromBalance.Uint64())

		toBalance, err := newState.BalanceOf(simulator.Context, &toAddr)
		assert.NoError(t, err)
		assert.Equal(t, uint64(100), toBalance.Uint64())
	})

	t.Run("successful with granularity", func(t *testing.T) {
		simulator, fromAddr, toAddr, _ := setup(t, abi.NewTokenAmount(1000), 5)
		simulator.ExpectSend(simulated.SendMock{
			To:     toAddr,
			Method: RECEIVERHOOKMETHODNUM,
			Params: []byte{130, 26, 133, 34, 59, 223, 73, 134, 1, 2, 1, 66, 0, 100, 64, 64},
			Value:  big.Zero(),
			Out:    types.SendResult{},
		})

		var newState Frc46Token
		sdk.LoadState(simulator.Context, &newState)

		_, err := newState.Transfer(simulator.Context, &TransferParams{
			To:           toAddr,
			Amount:       abi.NewTokenAmount(100),
			OperatorData: nil,
		})
		assert.NoError(t, err)

		fromBalance, err := newState.BalanceOf(simulator.Context, &fromAddr)
		assert.NoError(t, err)
		assert.Equal(t, uint64(900), fromBalance.Uint64())

		toBalance, err := newState.BalanceOf(simulator.Context, &toAddr)
		assert.NoError(t, err)
		assert.Equal(t, uint64(100), toBalance.Uint64())
	})

	t.Run("success transfer zero", func(t *testing.T) {
		simulator, _, toAddr, _ := setup(t, abi.NewTokenAmount(1000), 1)

		simulator.ExpectSend(simulated.SendMock{
			To:     toAddr,
			Method: RECEIVERHOOKMETHODNUM,
			Params: []byte{130, 26, 133, 34, 59, 223, 71, 134, 1, 2, 1, 64, 64, 64},
			Value:  big.Zero(),
			Out:    types.SendResult{},
		})
		var newState Frc46Token
		sdk.LoadState(simulator.Context, &newState)
		_, err := newState.Transfer(simulator.Context, &TransferParams{
			To:     toAddr,
			Amount: abi.NewTokenAmount(0),
		})
		assert.NoError(t, err)
	})

	t.Run("fail balance not enough", func(t *testing.T) {
		simulator, _, toAddr, _ := setup(t, abi.NewTokenAmount(1000), 1)

		var newState Frc46Token
		sdk.LoadState(simulator.Context, &newState)
		_, err := newState.Transfer(simulator.Context, &TransferParams{
			To:     toAddr,
			Amount: abi.NewTokenAmount(10000),
		})
		assert.EqualError(t, err, "negative balance caused by decreasing 1's balance of 1000 by -10000")
	})

	t.Run("fail with granularity", func(t *testing.T) {
		simulator, _, toAddr, _ := setup(t, abi.NewTokenAmount(1000), 9)
		var newState Frc46Token
		sdk.LoadState(simulator.Context, &newState)
		_, err := newState.Transfer(simulator.Context, &TransferParams{
			To:           toAddr,
			Amount:       abi.NewTokenAmount(100),
			OperatorData: nil,
		})
		assert.EqualError(t, err, "amount 100 for transfer must be a multiple of 9 16")
	})

}

func TestApprovalAndTransfer(t *testing.T) {
	t.Run("success approval and transfer", func(t *testing.T) {
		simulator, fromAddr, approvalAddr, toAddr := setup(t, abi.NewTokenAmount(1000), 1)
		fromId, err := simulator.ResolveAddress(fromAddr)
		assert.NoError(t, err)
		approvalId, err := simulator.ResolveAddress(approvalAddr)
		assert.NoError(t, err)
		newState := &Frc46Token{}
		sdk.LoadState(simulator.Context, newState)
		ctx := simulator.Context
		simulator.SetMessageContext(&types.MessageContext{
			ValueReceived: abi.NewTokenAmount(0),
			Caller:        fromId,
		})
		amount, err := newState.IncreaseAllowance(ctx, &IncreaseAllowanceParams{
			Operator: approvalAddr,
			Increase: abi.NewTokenAmount(100),
		})
		assert.NoError(t, err)
		assert.Equal(t, uint64(100), amount.Uint64())

		simulator.SetMessageContext(&types.MessageContext{
			ValueReceived: abi.NewTokenAmount(0),
			Caller:        approvalId,
		})
		simulator.ExpectSend(simulated.SendMock{
			To:     toAddr,
			Method: RECEIVERHOOKMETHODNUM,
			Params: []byte{130, 26, 133, 34, 59, 223, 73, 134, 1, 3, 1, 66, 0, 10, 64, 64},
			Value:  zero,
			Out: types.SendResult{
				ExitCode: ferrors.OK,
			},
		})
		transferFrom, err := newState.TransferFrom(ctx, &TransferFromParams{
			From:   fromAddr,
			To:     toAddr,
			Amount: abi.NewTokenAmount(10),
		})
		assert.NoError(t, err)
		assert.Equal(t, uint64(990), transferFrom.FromBalance.Uint64())
		assert.Equal(t, uint64(10), transferFrom.To_balance.Uint64())
		assert.Equal(t, uint64(90), transferFrom.Allowance.Uint64())
	})

	t.Run("success approval and transfer", func(t *testing.T) {
		simulator, fromAddr, approvalAddr, toAddr := setup(t, abi.NewTokenAmount(1000), 5)
		fromId, err := simulator.ResolveAddress(fromAddr)
		assert.NoError(t, err)
		approvalId, err := simulator.ResolveAddress(approvalAddr)
		assert.NoError(t, err)
		newState := &Frc46Token{}
		sdk.LoadState(simulator.Context, newState)
		ctx := simulator.Context
		simulator.SetMessageContext(&types.MessageContext{
			ValueReceived: abi.NewTokenAmount(0),
			Caller:        fromId,
		})
		amount, err := newState.IncreaseAllowance(ctx, &IncreaseAllowanceParams{
			Operator: approvalAddr,
			Increase: abi.NewTokenAmount(100),
		})
		assert.NoError(t, err)
		assert.Equal(t, uint64(100), amount.Uint64())

		simulator.SetMessageContext(&types.MessageContext{
			ValueReceived: abi.NewTokenAmount(0),
			Caller:        approvalId,
		})
		simulator.ExpectSend(simulated.SendMock{
			To:     toAddr,
			Method: RECEIVERHOOKMETHODNUM,
			Params: []byte{130, 26, 133, 34, 59, 223, 73, 134, 1, 3, 1, 66, 0, 10, 64, 64},
			Value:  zero,
			Out: types.SendResult{
				ExitCode: ferrors.OK,
			},
		})
		transferFrom, err := newState.TransferFrom(ctx, &TransferFromParams{
			From:   fromAddr,
			To:     toAddr,
			Amount: abi.NewTokenAmount(10),
		})
		assert.NoError(t, err)
		assert.Equal(t, uint64(990), transferFrom.FromBalance.Uint64())
		assert.Equal(t, uint64(10), transferFrom.To_balance.Uint64())
		assert.Equal(t, uint64(90), transferFrom.Allowance.Uint64())
	})

	t.Run("successful approval zero balance", func(t *testing.T) {
		simulator, fromAddr, approvalAddr, _ := setup(t, abi.NewTokenAmount(1000), 1)
		fromId, err := simulator.ResolveAddress(fromAddr)
		assert.NoError(t, err)
		var newState Frc46Token
		sdk.LoadState(simulator.Context, &newState)
		ctx := simulator.Context
		simulator.SetMessageContext(&types.MessageContext{
			ValueReceived: abi.NewTokenAmount(0),
			Caller:        fromId,
		})

		amount, err := newState.IncreaseAllowance(ctx, &IncreaseAllowanceParams{
			Operator: approvalAddr,
			Increase: abi.NewTokenAmount(0),
		})
		assert.NoError(t, err)
		assert.Equal(t, uint64(0), amount.Uint64())
	})

	t.Run("success transferfrom zero balance ", func(t *testing.T) {
		simulator, fromAddr, approvalAddr, toAddr := setup(t, abi.NewTokenAmount(1000), 1)
		fromId, err := simulator.ResolveAddress(fromAddr)
		assert.NoError(t, err)
		approvalId, err := simulator.ResolveAddress(approvalAddr)
		assert.NoError(t, err)
		var newState Frc46Token
		sdk.LoadState(simulator.Context, &newState)
		ctx := simulator.Context
		simulator.SetMessageContext(&types.MessageContext{
			ValueReceived: abi.NewTokenAmount(0),
			Caller:        fromId,
		})
		amount, err := newState.IncreaseAllowance(ctx, &IncreaseAllowanceParams{
			Operator: approvalAddr,
			Increase: abi.NewTokenAmount(0),
		})
		assert.NoError(t, err)
		assert.Equal(t, uint64(0), amount.Uint64())

		simulator.SetMessageContext(&types.MessageContext{
			ValueReceived: abi.NewTokenAmount(0),
			Caller:        approvalId,
		})
		simulator.ExpectSend(simulated.SendMock{
			To:     toAddr,
			Method: RECEIVERHOOKMETHODNUM,
			Params: []byte{130, 26, 133, 34, 59, 223, 71, 134, 1, 3, 1, 64, 64, 64},
			Value:  zero,
			Out: types.SendResult{
				ExitCode: ferrors.OK,
			},
		})
		transferRet, err := newState.TransferFrom(ctx, &TransferFromParams{
			From:   fromAddr,
			To:     toAddr,
			Amount: abi.NewTokenAmount(0),
		})
		assert.NoError(t, err)
		assert.Equal(t, uint64(0), transferRet.To_balance.Uint64())
	})
	t.Run("success transferfrom zero balance ", func(t *testing.T) {
		simulator, fromAddr, approvalAddr, toAddr := setup(t, abi.NewTokenAmount(1000), 1)
		approvalId, err := simulator.ResolveAddress(approvalAddr)
		assert.NoError(t, err)
		var newState Frc46Token
		sdk.LoadState(simulator.Context, &newState)
		ctx := simulator.Context

		simulator.SetMessageContext(&types.MessageContext{
			ValueReceived: abi.NewTokenAmount(0),
			Caller:        approvalId,
		})
		simulator.ExpectSend(simulated.SendMock{
			To:     toAddr,
			Method: RECEIVERHOOKMETHODNUM,
			Params: []byte{130, 26, 133, 34, 59, 223, 71, 134, 1, 3, 1, 64, 64, 64},
			Value:  zero,
			Out: types.SendResult{
				ExitCode: ferrors.OK,
			},
		})
		transferRet, err := newState.TransferFrom(ctx, &TransferFromParams{
			From:   fromAddr,
			To:     toAddr,
			Amount: abi.NewTokenAmount(0),
		})
		assert.NoError(t, err)
		assert.Equal(t, uint64(0), transferRet.To_balance.Uint64())
	})

	t.Run("fail transferfrom zero balance ", func(t *testing.T) {
		simulator, fromAddr, approvalAddr, toAddr := setup(t, abi.NewTokenAmount(1000), 1)
		fromId, err := simulator.ResolveAddress(fromAddr)
		assert.NoError(t, err)
		approvalId, err := simulator.ResolveAddress(approvalAddr)
		assert.NoError(t, err)
		var newState Frc46Token
		sdk.LoadState(simulator.Context, &newState)
		ctx := simulator.Context
		simulator.SetMessageContext(&types.MessageContext{
			ValueReceived: abi.NewTokenAmount(0),
			Caller:        fromId,
		})
		amount, err := newState.IncreaseAllowance(ctx, &IncreaseAllowanceParams{
			Operator: approvalAddr,
			Increase: abi.NewTokenAmount(100),
		})
		assert.NoError(t, err)
		assert.Equal(t, uint64(100), amount.Uint64())

		simulator.SetMessageContext(&types.MessageContext{
			ValueReceived: abi.NewTokenAmount(0),
			Caller:        approvalId,
		})

		_, err = newState.TransferFrom(ctx, &TransferFromParams{
			From:   fromAddr,
			To:     toAddr,
			Amount: abi.NewTokenAmount(200),
		})
		assert.EqualError(t, err, "t02 attempted to utilise 200 of allowance 100 set by t01 19")
	})

	t.Run("fail transferfrom zero balance ", func(t *testing.T) {
		simulator, fromAddr, approvalAddr, toAddr := setup(t, abi.NewTokenAmount(60), 1)
		fromId, err := simulator.ResolveAddress(fromAddr)
		assert.NoError(t, err)
		approvalId, err := simulator.ResolveAddress(approvalAddr)
		assert.NoError(t, err)
		var newState Frc46Token
		sdk.LoadState(simulator.Context, &newState)
		ctx := simulator.Context
		simulator.SetMessageContext(&types.MessageContext{
			ValueReceived: abi.NewTokenAmount(0),
			Caller:        fromId,
		})
		amount, err := newState.IncreaseAllowance(ctx, &IncreaseAllowanceParams{
			Operator: approvalAddr,
			Increase: abi.NewTokenAmount(100),
		})
		assert.NoError(t, err)
		assert.Equal(t, uint64(100), amount.Uint64())

		simulator.SetMessageContext(&types.MessageContext{
			ValueReceived: abi.NewTokenAmount(0),
			Caller:        approvalId,
		})
		_, err = newState.TransferFrom(ctx, &TransferFromParams{
			From:   fromAddr,
			To:     toAddr,
			Amount: abi.NewTokenAmount(80),
		})
		assert.EqualError(t, err, "negative balance caused by decreasing 1's balance of 60 by -80")
	})

	t.Run("fail granularity", func(t *testing.T) {
		simulator, fromAddr, approvalAddr, toAddr := setup(t, abi.NewTokenAmount(1000), 9)
		approvalId, err := simulator.ResolveAddress(approvalAddr)
		assert.NoError(t, err)
		newState := &Frc46Token{}
		sdk.LoadState(simulator.Context, newState)
		ctx := simulator.Context

		simulator.SetMessageContext(&types.MessageContext{
			ValueReceived: abi.NewTokenAmount(0),
			Caller:        approvalId,
		})
		_, err = newState.TransferFrom(ctx, &TransferFromParams{
			From:   fromAddr,
			To:     toAddr,
			Amount: abi.NewTokenAmount(10),
		})
		assert.EqualError(t, err, "amount 10 for transfer from must be a multiple of 9 16")
	})
}

func TestFrc46Token_In_DecreaseAllowance(t *testing.T) {
	t.Run("success increase/decrease allowance", func(t *testing.T) {
		simulator, fromAddr, approvalAddr, _ := setup(t, abi.NewTokenAmount(100), 1)
		fromId, err := simulator.ResolveAddress(fromAddr)
		assert.NoError(t, err)
		var newState Frc46Token
		sdk.LoadState(simulator.Context, &newState)
		ctx := simulator.Context
		simulator.SetMessageContext(&types.MessageContext{
			ValueReceived: abi.NewTokenAmount(0),
			Caller:        fromId,
		})

		{
			amount, err := newState.IncreaseAllowance(ctx, &IncreaseAllowanceParams{
				Operator: approvalAddr,
				Increase: abi.NewTokenAmount(100),
			})
			assert.NoError(t, err)
			assert.Equal(t, uint64(100), amount.Uint64())

			amount, err = newState.Allowance(ctx, &GetAllowanceParams{
				Owner:    fromAddr,
				Operator: approvalAddr,
			})
			assert.Equal(t, uint64(100), amount.Uint64())
		}

		{
			amount, err := newState.IncreaseAllowance(ctx, &IncreaseAllowanceParams{
				Operator: approvalAddr,
				Increase: abi.NewTokenAmount(50),
			})
			assert.NoError(t, err)
			assert.Equal(t, uint64(150), amount.Uint64())

			amount, err = newState.Allowance(ctx, &GetAllowanceParams{
				Owner:    fromAddr,
				Operator: approvalAddr,
			})
			assert.Equal(t, uint64(150), amount.Uint64())
		}

		{
			amount, err := newState.DecreaseAllowance(ctx, &DecreaseAllowanceParams{
				Operator: approvalAddr,
				Decrease: abi.NewTokenAmount(60),
			})
			assert.NoError(t, err)
			assert.Equal(t, uint64(90), amount.Uint64())

			amount, err = newState.Allowance(ctx, &GetAllowanceParams{
				Owner:    fromAddr,
				Operator: approvalAddr,
			})
			assert.Equal(t, uint64(90), amount.Uint64())
		}

		{
			amount, err := newState.DecreaseAllowance(ctx, &DecreaseAllowanceParams{
				Operator: approvalAddr,
				Decrease: abi.NewTokenAmount(100),
			})
			assert.NoError(t, err)
			assert.Equal(t, uint64(0), amount.Uint64())

			amount, err = newState.Allowance(ctx, &GetAllowanceParams{
				Owner:    fromAddr,
				Operator: approvalAddr,
			})
			assert.Equal(t, uint64(0), amount.Uint64())
		}
	})
}

func TestFrc46Token_RevokeAllowance(t *testing.T) {
	simulator, fromAddr, approvalAddr, _ := setup(t, abi.NewTokenAmount(100), 1)
	fromId, err := simulator.ResolveAddress(fromAddr)
	assert.NoError(t, err)
	var newState Frc46Token
	sdk.LoadState(simulator.Context, &newState)
	ctx := simulator.Context
	simulator.SetMessageContext(&types.MessageContext{
		ValueReceived: abi.NewTokenAmount(0),
		Caller:        fromId,
	})

	amount, err := newState.RevokeAllowance(ctx, &RevokeAllowanceParams{approvalAddr})
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), amount.Uint64())

	{
		amount, err := newState.IncreaseAllowance(ctx, &IncreaseAllowanceParams{
			Operator: approvalAddr,
			Increase: abi.NewTokenAmount(100),
		})
		assert.NoError(t, err)
		assert.Equal(t, uint64(100), amount.Uint64())

		amount, err = newState.Allowance(ctx, &GetAllowanceParams{
			Owner:    fromAddr,
			Operator: approvalAddr,
		})
		assert.Equal(t, uint64(100), amount.Uint64())
	}

	amount, err = newState.RevokeAllowance(ctx, &RevokeAllowanceParams{approvalAddr})
	assert.NoError(t, err)
	assert.Equal(t, uint64(100), amount.Uint64())
}

func TestFrc46Token_Burn(t *testing.T) {
	t.Run("success burn", func(t *testing.T) {
		simulator, fromAddr, _, _ := setup(t, abi.NewTokenAmount(100), 1)
		fromId, err := simulator.ResolveAddress(fromAddr)
		assert.NoError(t, err)
		var newState Frc46Token
		sdk.LoadState(simulator.Context, &newState)
		ctx := simulator.Context
		simulator.SetMessageContext(&types.MessageContext{
			ValueReceived: abi.NewTokenAmount(0),
			Caller:        fromId,
		})

		burnRet, err := newState.Burn(ctx, &BurnParams{Amount: abi.NewTokenAmount(1)})
		assert.NoError(t, err)
		assert.Equal(t, uint64(99), burnRet.Balance.Uint64())
		assert.Equal(t, uint64(supply-1), newState.Supply.Uint64())
	})

	t.Run("burn more than balance", func(t *testing.T) {
		simulator, fromAddr, _, _ := setup(t, abi.NewTokenAmount(100), 1)
		fromId, err := simulator.ResolveAddress(fromAddr)
		assert.NoError(t, err)
		var newState Frc46Token
		sdk.LoadState(simulator.Context, &newState)
		ctx := simulator.Context
		simulator.SetMessageContext(&types.MessageContext{
			ValueReceived: abi.NewTokenAmount(0),
			Caller:        fromId,
		})

		_, err = newState.Burn(ctx, &BurnParams{Amount: abi.NewTokenAmount(10000)})
		assert.EqualError(t, err, "negative balance caused by decreasing 1's balance of 100 by -10000")
	})

	t.Run("burn more than supply", func(t *testing.T) {
		simulator, fromAddr, _, _ := setup(t, abi.NewTokenAmount(100001), 1)
		fromId, err := simulator.ResolveAddress(fromAddr)
		assert.NoError(t, err)
		var newState Frc46Token
		sdk.LoadState(simulator.Context, &newState)
		ctx := simulator.Context
		simulator.SetMessageContext(&types.MessageContext{
			ValueReceived: abi.NewTokenAmount(0),
			Caller:        fromId,
		})

		_, err = newState.Burn(ctx, &BurnParams{Amount: abi.NewTokenAmount(100001)})
		assert.EqualError(t, err, "supply must big than 0 supply 100000, delta -100001 16")
	})
}

func TestFrc46Token_BurnFrom(t *testing.T) {
	t.Run("success burn from zero", func(t *testing.T) {
		simulator, fromAddr, _, approvalAddr := setup(t, abi.NewTokenAmount(100), 1)
		fromId, err := simulator.ResolveAddress(fromAddr)
		assert.NoError(t, err)
		approvalId, err := simulator.ResolveAddress(approvalAddr)
		assert.NoError(t, err)
		var newState Frc46Token
		sdk.LoadState(simulator.Context, &newState)
		ctx := simulator.Context
		simulator.SetMessageContext(&types.MessageContext{
			ValueReceived: abi.NewTokenAmount(0),
			Caller:        fromId,
		})

		_, err = newState.IncreaseAllowance(ctx, &IncreaseAllowanceParams{
			Operator: approvalAddr,
			Increase: abi.NewTokenAmount(100),
		})
		assert.NoError(t, err)

		simulator.SetMessageContext(&types.MessageContext{
			ValueReceived: abi.NewTokenAmount(0),
			Caller:        approvalId,
		})
		ret, err := newState.BurnFrom(ctx, &BurnFromParams{Owner: fromAddr, Amount: abi.NewTokenAmount(10)})
		assert.NoError(t, err)
		assert.Equal(t, uint64(90), ret.Balance.Uint64())
		assert.Equal(t, uint64(90), ret.Allowance.Uint64())
		assert.Equal(t, uint64(supply-10), newState.Supply.Uint64())
	})

	t.Run("success burn from zero", func(t *testing.T) {
		simulator, fromAddr, toAddr, _ := setup(t, abi.NewTokenAmount(100), 1)
		fromId, err := simulator.ResolveAddress(fromAddr)
		assert.NoError(t, err)
		var newState Frc46Token
		sdk.LoadState(simulator.Context, &newState)
		ctx := simulator.Context
		simulator.SetMessageContext(&types.MessageContext{
			ValueReceived: abi.NewTokenAmount(0),
			Caller:        fromId,
		})

		ret, err := newState.BurnFrom(ctx, &BurnFromParams{Owner: toAddr, Amount: abi.NewTokenAmount(0)})
		assert.NoError(t, err)
		assert.Equal(t, uint64(0), ret.Balance.Uint64())
		assert.Equal(t, uint64(0), ret.Allowance.Uint64())
	})

	t.Run("fail burn from more than allowance", func(t *testing.T) {
		simulator, fromAddr, _, approvalAddr := setup(t, abi.NewTokenAmount(100), 1)
		fromId, err := simulator.ResolveAddress(fromAddr)
		assert.NoError(t, err)
		approvalId, err := simulator.ResolveAddress(approvalAddr)
		assert.NoError(t, err)
		var newState Frc46Token
		sdk.LoadState(simulator.Context, &newState)
		ctx := simulator.Context
		simulator.SetMessageContext(&types.MessageContext{
			ValueReceived: abi.NewTokenAmount(0),
			Caller:        fromId,
		})

		_, err = newState.IncreaseAllowance(ctx, &IncreaseAllowanceParams{
			Operator: approvalAddr,
			Increase: abi.NewTokenAmount(100),
		})
		assert.NoError(t, err)

		simulator.SetMessageContext(&types.MessageContext{
			ValueReceived: abi.NewTokenAmount(0),
			Caller:        approvalId,
		})
		_, err = newState.BurnFrom(ctx, &BurnFromParams{Owner: fromAddr, Amount: abi.NewTokenAmount(101)})
		assert.EqualError(t, err, "t03 attempted to utilise 101 of allowance 100 set by t01 19")
	})

	t.Run("fail not enough balance", func(t *testing.T) {
		simulator, fromAddr, toAddr, _ := setup(t, abi.NewTokenAmount(100), 1)
		fromId, err := simulator.ResolveAddress(fromAddr)
		assert.NoError(t, err)
		var newState Frc46Token
		sdk.LoadState(simulator.Context, &newState)
		ctx := simulator.Context
		simulator.SetMessageContext(&types.MessageContext{
			ValueReceived: abi.NewTokenAmount(0),
			Caller:        fromId,
		})

		_, err = newState.BurnFrom(ctx, &BurnFromParams{Owner: toAddr, Amount: abi.NewTokenAmount(1)})
		assert.EqualError(t, err, "t01 attempted to utilise 1 of allowance 0 set by t02 19")
	})
}

func TestFrc46Token_Mint(t *testing.T) {
	t.Run("success mint", func(t *testing.T) {
		simulator, fromAddr, toAddr, _ := setup(t, abi.NewTokenAmount(100), 1)
		fromId, err := simulator.ResolveAddress(fromAddr)
		assert.NoError(t, err)
		var newState Frc46Token
		sdk.LoadState(simulator.Context, &newState)
		ctx := simulator.Context
		simulator.SetMessageContext(&types.MessageContext{
			ValueReceived: abi.NewTokenAmount(0),
			Caller:        fromId,
		})

		simulator.ExpectSend(simulated.SendMock{
			To:     toAddr,
			Method: RECEIVERHOOKMETHODNUM,
			Params: []byte{130, 26, 133, 34, 59, 223, 73, 134, 0, 2, 1, 66, 0, 100, 64, 64},
			Value:  abi.NewTokenAmount(0),
			Out: types.SendResult{
				ExitCode: ferrors.OK,
			},
		})
		ret, err := newState.Mint(ctx, &MintParams{
			InitialOwner: toAddr,
			Amount:       abi.NewTokenAmount(100),
			OperatorData: nil,
		})
		assert.NoError(t, err)
		assert.Equal(t, uint64(supply+100), newState.Supply.Uint64())
		assert.Equal(t, uint64(supply+100), ret.Supply.Uint64())
		assert.Equal(t, uint64(100), ret.Balance.Uint64())
	})

	t.Run("success mint zero", func(t *testing.T) {
		simulator, fromAddr, toAddr, _ := setup(t, abi.NewTokenAmount(100), 1)
		fromId, err := simulator.ResolveAddress(fromAddr)
		assert.NoError(t, err)
		var newState Frc46Token
		sdk.LoadState(simulator.Context, &newState)
		ctx := simulator.Context
		simulator.SetMessageContext(&types.MessageContext{
			ValueReceived: abi.NewTokenAmount(0),
			Caller:        fromId,
		})

		simulator.ExpectSend(simulated.SendMock{
			To:     toAddr,
			Method: RECEIVERHOOKMETHODNUM,
			Params: []byte{130, 26, 133, 34, 59, 223, 71, 134, 0, 2, 1, 64, 64, 64},
			Value:  abi.NewTokenAmount(0),
			Out: types.SendResult{
				ExitCode: ferrors.OK,
			},
		})
		ret, err := newState.Mint(ctx, &MintParams{
			InitialOwner: toAddr,
			Amount:       abi.NewTokenAmount(0),
			OperatorData: nil,
		})
		assert.NoError(t, err)
		assert.Equal(t, uint64(supply), newState.Supply.Uint64())
		assert.Equal(t, uint64(supply), ret.Supply.Uint64())
		assert.Equal(t, uint64(0), ret.Balance.Uint64())
	})
	t.Run("fail granularity", func(t *testing.T) {
		simulator, fromAddr, toAddr, _ := setup(t, abi.NewTokenAmount(100), 9)
		fromId, err := simulator.ResolveAddress(fromAddr)
		assert.NoError(t, err)
		var newState Frc46Token
		sdk.LoadState(simulator.Context, &newState)
		ctx := simulator.Context
		simulator.SetMessageContext(&types.MessageContext{
			ValueReceived: abi.NewTokenAmount(0),
			Caller:        fromId,
		})

		_, err = newState.Mint(ctx, &MintParams{
			InitialOwner: toAddr,
			Amount:       abi.NewTokenAmount(100),
			OperatorData: nil,
		})
		assert.EqualError(t, err, "amount 100 for mint must be a multiple of 9 16")

		simulator.ExpectSend(simulated.SendMock{
			To:     toAddr,
			Method: RECEIVERHOOKMETHODNUM,
			Params: []byte{130, 26, 133, 34, 59, 223, 73, 134, 0, 2, 1, 66, 0, 18, 64, 64},
			Value:  abi.NewTokenAmount(0),
			Out: types.SendResult{
				ExitCode: ferrors.OK,
			},
		})
		ret, err := newState.Mint(ctx, &MintParams{
			InitialOwner: toAddr,
			Amount:       abi.NewTokenAmount(18),
			OperatorData: nil,
		})
		assert.NoError(t, err)
		assert.Equal(t, uint64(supply+18), newState.Supply.Uint64())
		assert.Equal(t, uint64(supply+18), ret.Supply.Uint64())
		assert.Equal(t, uint64(18), ret.Balance.Uint64())
	})
}
