package contract

import (
	"context"
	"fmt"

	"github.com/filecoin-project/go-state-types/big"

	"github.com/filecoin-project/go-state-types/abi"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk"
)

type State struct {
	Owner abi.ActorID
}

func (e *State) Export() []interface{} {
	return []interface{}{
		Constructor,
		e.GetBalance,
		e.Withdraw,
	}
}

func Constructor(ctx context.Context) error {
	originer, err := sdk.Origin(ctx)
	if err != nil {
		return err
	}

	st := &State{
		Owner: originer,
	}
	_ = sdk.Constructor(ctx, st)
	return nil
}

func (st *State) GetBalance(ctx context.Context) (*big.Int, error) {
	receiver, err := sdk.Receiver(ctx)
	if err != nil {
		return nil, err
	}

	return sdk.BalanceOf(ctx, receiver)
}

func (st *State) Withdraw(ctx context.Context, amount *abi.TokenAmount) error {
	caller, err := sdk.Caller(ctx)
	if err != nil {
		return err
	}
	if caller != st.Owner {
		return fmt.Errorf("caller %d is not owner %d", caller, st.Owner)
	}

	receipt, err := sdk.Send(ctx, sdk.MustAddressFromActorId(caller), 0, nil, *amount)
	if err != nil {
		return err
	}
	if receipt.ExitCode != 0 {
		return fmt.Errorf("transfer fail exit code %d", receipt.ExitCode)
	}
	return nil
}
