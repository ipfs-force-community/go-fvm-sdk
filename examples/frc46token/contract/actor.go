package contract

import (
	"context"
	"errors"
	"fmt"
	stdbig "math/big"

	typegen "github.com/whyrusleeping/cbor-gen"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"

	"github.com/filecoin-project/go-state-types/big"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"

	"github.com/filecoin-project/go-address"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/adt"
	"github.com/ipfs/go-cid"
)

var logger sdk.Logger

func init() {
	logger, _ = sdk.NewLogger()
}

var zero = big.Zero()

// DEFAULTHAMTBITWIDTH  This value has been chosen to optimise to reduce gas-costs when accessing the balances map. Non-standard
// use cases of the token library might find a different value to be more efficient.
const DEFAULTHAMTBITWIDTH int = 3

// Token state IPLD structure
type Frc46Token struct {
	Name        string
	Symbol      string
	Granularity uint64
	/// Total supply of token
	Supply abi.TokenAmount
	/// Map<ActorId, TokenAmount> of balances as a Hamt
	Balances cid.Cid
	/// Map<ActorId, Map<ActorId, TokenAmount>> as a Hamt. Allowances are stored balances[owner][operator]
	Allowances cid.Cid
}

var _ IFrc46Token = (*Frc46Token)(nil)
var _ IFrc46Unspecific = (*Frc46Token)(nil)

type Alias struct {
	Name string
	Func interface{}
}

func (t *Frc46Token) Export() []interface{} {
	return []interface{}{
		Constructor,
		Alias{
			Name: "Name",
			Func: t.GetName,
		},
		Alias{
			Name: "Symbol",
			Func: t.GetSymbol,
		},
		Alias{
			Name: "Granularity",
			Func: t.GetGranularity,
		},
		Alias{
			Name: "TotalSupply",
			Func: t.GetTotalSupply,
		},
		t.Mint,
		t.BalanceOf,
		t.Allowance,
		t.Transfer,
		t.TransferFrom,
		t.IncreaseAllowance,
		t.DecreaseAllowance,
		t.RevokeAllowance,
		t.Burn,
		t.BurnFrom,
	}
}

func Constructor(ctx context.Context, req *ConstructorReq) error {
	emptyMap, err := adt.MakeEmptyMap(adt.AdtStore(ctx), DEFAULTHAMTBITWIDTH)
	if err != nil {
		return err
	}
	emptyRoot, err := emptyMap.Root()
	if err != nil {
		return err
	}
	state := &Frc46Token{
		Name:        req.Name,
		Symbol:      req.Symbol,
		Granularity: req.Granularity,
		Supply:      req.Supply,
		Balances:    emptyRoot,
		Allowances:  emptyRoot,
	}

	logger.Logf(ctx, "create token %s, symbol %s", req.Name, state.Symbol)
	_ = sdk.Constructor(ctx, state)
	return nil
}

func (t *Frc46Token) GetName(_ context.Context) types.CborString {
	return types.CborString(t.Name)
}

// 15026712600

func (t *Frc46Token) GetSymbol(_ context.Context) types.CborString {
	return types.CborString(t.Symbol)
}

func (t *Frc46Token) GetGranularity(_ context.Context) types.CborUint {
	return types.CborUint(t.Granularity)
}

func (t *Frc46Token) GetTotalSupply(_ context.Context) *abi.TokenAmount {
	return &t.Supply
}

func (t *Frc46Token) Mint(ctx context.Context, params *MintParams) (*MintReturn, error) {
	callerId, err := sdk.Caller(ctx)
	if err != nil {
		return nil, err
	}
	callerAddr, err := address.NewIDAddress(uint64(callerId))
	if err != nil {
		return nil, err
	}

	hook, err := t.mint(ctx, params.InitialOwner, callerAddr, params.Amount, params.OperatorData)
	if err != nil {
		return nil, err
	}

	rootCid := sdk.SaveState(ctx, t)
	err = sdk.SetRoot(ctx, rootCid)
	if err != nil {
		return nil, err
	}
	err = hook.Call(ctx)
	if err != nil {
		return nil, err
	}

	return t.mintReturn(ctx, hook.ResultData.(*MintIntermediate))
}

func (t *Frc46Token) mint(ctx context.Context, initAddr, operator address.Address, amount abi.TokenAmount, operatorData types.RawBytes) (*ReceiverHook, error) {
	err := ValidateAmountWithGranularity(amount, "mint", t.Granularity)
	if err != nil {
		return nil, err
	}

	initAddrId, err := sdk.ResolveOrInitAddress(ctx, initAddr)
	if err != nil {
		return nil, err
	}

	operatorId, err := sdk.ResolveOrInitAddress(ctx, operator)
	if err != nil {
		return nil, err
	}

	_, err = t.changeBalanceBy(ctx, initAddrId, amount)
	if err != nil {
		return nil, err
	}
	err = t.changeSupplyBy(amount)
	if err != nil {
		return nil, err
	}

	receiverId, err := sdk.Receiver(ctx)
	if err != nil {
		return nil, err
	}

	mintInterMedia := &MintIntermediate{Recipient: initAddrId}
	params := &FRC46TokenReceived{
		Operator:     operatorId,
		From:         receiverId,
		To:           initAddrId,
		Amount:       amount,
		OperatorData: operatorData,
	}

	return NewReceiverHook(initAddr, FRC46TOKENTYPE, params, mintInterMedia)
}

func (t *Frc46Token) BalanceOf(ctx context.Context, addr *address.Address) (*abi.TokenAmount, error) {
	actorId, err := sdk.ResolveAddress(ctx, *addr)
	if err != nil {
		if errors.Is(err, ferrors.NotFound) {
			return &zero, nil
		}
		return nil, err
	}

	return t.getBalance(ctx, actorId)
}

func (t *Frc46Token) getBalance(ctx context.Context, actorId abi.ActorID) (*abi.TokenAmount, error) {
	balanceMap, err := adt.AsMap(adt.AdtStore(ctx), t.Balances, adt.BalanceTableBitwidth)
	if err != nil {
		return nil, err
	}
	var balance = big.Int{Int: stdbig.NewInt(0)}
	found, err := balanceMap.Get(types.ActorKey(actorId), &balance)
	if err != nil {
		return nil, err
	}
	if !found {
		return &zero, nil
	}
	//return 0 if not exit
	return &balance, nil
}

func (t *Frc46Token) Allowance(ctx context.Context, params *GetAllowanceParams) (*abi.TokenAmount, error) {
	ownerId, err := sdk.ResolveAddress(ctx, params.Owner)
	if err != nil {
		if errors.Is(err, ferrors.NotFound) {
			return &zero, nil
		}
		return nil, err
	}

	operatorId, err := sdk.ResolveAddress(ctx, params.Operator)
	if err != nil {
		if errors.Is(err, ferrors.NotFound) {
			return &zero, nil
		}
		return nil, err
	}

	return t.getAllowanceBalance(ctx, ownerId, operatorId)
}

func (t *Frc46Token) Transfer(ctx context.Context, params *TransferParams) (*TransferReturn, error) {
	callerId, err := sdk.Caller(ctx)
	if err != nil {
		return nil, err
	}
	callerAddr, err := address.NewIDAddress(uint64(callerId))
	if err != nil {
		return nil, err
	}

	hook, err := t.transfer(ctx, callerAddr, params.To, params.Amount, params.OperatorData, nil)
	if err != nil {
		return nil, err
	}

	rootCid := sdk.SaveState(ctx, t)
	err = sdk.SetRoot(ctx, rootCid)
	if err != nil {
		return nil, err
	}
	err = hook.Call(ctx)
	if err != nil {
		return nil, err
	}
	return t.transferReturn(ctx, hook.ResultData.(*TransferIntermediate))
}

func (t *Frc46Token) transfer(ctx context.Context, from, to address.Address, amount abi.TokenAmount, operatorData, tokenData types.RawBytes) (*ReceiverHook, error) {
	err := ValidateAmountWithGranularity(amount, "transfer", t.Granularity)
	if err != nil {
		return nil, err
	}

	fromId, err := sdk.ResolveAddress(ctx, from)
	if err != nil {
		return nil, err
	}

	var transferIntermediate *TransferIntermediate
	toId, err := sdk.ResolveOrInitAddress(ctx, to)
	if fromId == toId {
		fromBalance, err := t.getBalance(ctx, fromId)
		if err != nil {
			return nil, err
		}
		if fromBalance.LessThan(amount) {
			return nil, fmt.Errorf("negative balance caused by decreasing %s's balance of %s by %s", from, fromBalance, amount)
		}
		transferIntermediate = &TransferIntermediate{
			From:          fromId,
			To:            toId,
			RecipientData: nil,
		}
	} else {
		_, err = t.changeBalanceBy(ctx, fromId, amount.Neg())
		if err != nil {
			return nil, err
		}
		_, err = t.changeBalanceBy(ctx, toId, amount)
		if err != nil {
			return nil, err
		}
		transferIntermediate = &TransferIntermediate{
			From:          fromId,
			To:            toId,
			RecipientData: nil,
		}
	}

	params := &FRC46TokenReceived{
		Operator:     fromId,
		From:         fromId,
		To:           toId,
		Amount:       amount,
		OperatorData: operatorData,
		TokenData:    tokenData,
	}

	return NewReceiverHook(to, FRC46TOKENTYPE, params, transferIntermediate)
}

func (t *Frc46Token) TransferFrom(ctx context.Context, params *TransferFromParams) (*TransferFromReturn, error) {
	operatorId, err := sdk.Caller(ctx)
	if err != nil {
		return nil, err
	}
	operatorAddr, err := address.NewIDAddress(uint64(operatorId))
	if err != nil {
		return nil, err
	}

	hook, err := t.transferFrom(ctx, operatorAddr, params.From, params.To, params.Amount, params.OperatorData, nil)
	if err != nil {
		return nil, err
	}

	rootCid := sdk.SaveState(ctx, t)
	err = sdk.SetRoot(ctx, rootCid)
	if err != nil {
		return nil, err
	}
	err = hook.Call(ctx)
	if err != nil {
		return nil, err
	}
	return t.transferFromReturn(ctx, hook.ResultData.(*TransferFromIntermediate))
}

func (t *Frc46Token) transferFrom(ctx context.Context, operator, from, to address.Address, amount abi.TokenAmount, operatorData, tokenData types.RawBytes) (*ReceiverHook, error) {
	err := ValidateAmountWithGranularity(amount, "transfer from", t.Granularity)
	if err != nil {
		return nil, err
	}

	if sdk.SameAddress(ctx, operator, from) {
		return nil, fmt.Errorf("operator cannot be the same as the debited address %s %w", operator, ferrors.USR_ILLEGAL_ARGUMENT)
	}

	operatorId, err := sdk.ResolveAddress(ctx, operator)
	if err != nil {
		if errors.Is(err, ferrors.NotFound) {
			return nil, fmt.Errorf("operator %s not found %w", operator, ferrors.USR_ILLEGAL_ARGUMENT)
		}
	}

	fromId, err := sdk.ResolveAddress(ctx, from)
	if err != nil {
		if errors.Is(err, ferrors.NotFound) {
			return nil, fmt.Errorf("from %s address not found %w", from, ferrors.USR_ILLEGAL_ARGUMENT)
		}
	}

	toId, err := sdk.ResolveOrInitAddress(ctx, to)
	if err != nil {
		return nil, err
	}

	var transferIntermediate *TransferFromIntermediate
	if amount.IsZero() {
		transferIntermediate = &TransferFromIntermediate{
			Operator:      operatorId,
			From:          fromId,
			To:            toId,
			RecipientData: nil,
		}
	} else {
		_, err = t.attemptUseAllowance(ctx, operatorId, fromId, amount)
		if err != nil {
			return nil, err
		}

		if fromId == toId {
			fromBalance, err := t.getBalance(ctx, fromId)
			if err != nil {
				return nil, err
			}
			if fromBalance.LessThan(amount) {
				return nil, fmt.Errorf("negative balance caused by decreasing %s's balance of %s by %s", from, fromBalance, amount)
			}
			transferIntermediate = &TransferFromIntermediate{
				Operator:      operatorId,
				From:          fromId,
				To:            toId,
				RecipientData: nil,
			}
		} else {
			_, err = t.changeBalanceBy(ctx, fromId, amount.Neg())
			if err != nil {
				return nil, err
			}
			_, err = t.changeBalanceBy(ctx, toId, amount)
			if err != nil {
				return nil, err
			}
			transferIntermediate = &TransferFromIntermediate{
				Operator:      operatorId,
				From:          fromId,
				To:            toId,
				RecipientData: nil,
			}
		}
	}

	params := &FRC46TokenReceived{
		Operator:     fromId,
		From:         fromId,
		To:           toId,
		Amount:       amount,
		OperatorData: operatorData,
		TokenData:    tokenData,
	}

	return NewReceiverHook(to, FRC46TOKENTYPE, params, transferIntermediate)
}

func (t *Frc46Token) IncreaseAllowance(ctx context.Context, params *IncreaseAllowanceParams) (*abi.TokenAmount, error) {
	addr, err := sdk.CallerAddress(ctx)
	if err != nil {
		return nil, err
	}
	return t.increaseAllowance(ctx, addr, params.Operator, params.Increase)
}

func (t *Frc46Token) increaseAllowance(ctx context.Context, owner address.Address, operator address.Address, delta abi.TokenAmount) (*abi.TokenAmount, error) {
	err := ValidateAmountWithGranularity(delta, "increase allowance delta", t.Granularity)
	if err != nil {
		return nil, err
	}

	err = ValidateAllowance(delta, "increase allowance delta")
	if err != nil {
		return nil, err
	}

	ownerId, err := sdk.ResolveOrInitAddress(ctx, owner)
	if err != nil {
		return nil, err
	}

	operatorId, err := sdk.ResolveOrInitAddress(ctx, operator)
	if err != nil {
		return nil, err
	}

	return t.changeAllowanceBy(ctx, ownerId, operatorId, delta)
}

func (t *Frc46Token) DecreaseAllowance(ctx context.Context, params *DecreaseAllowanceParams) (*abi.TokenAmount, error) {
	addr, err := sdk.CallerAddress(ctx)
	if err != nil {
		return nil, err
	}
	return t.decreaseAllowance(ctx, addr, params.Operator, params.Decrease)
}

func (t *Frc46Token) decreaseAllowance(ctx context.Context, owner address.Address, operator address.Address, delta abi.TokenAmount) (*abi.TokenAmount, error) {
	err := ValidateAmountWithGranularity(delta, "increase allowance delta", t.Granularity)
	if err != nil {
		return nil, err
	}

	err = ValidateAllowance(delta, "increase allowance delta")
	if err != nil {
		return nil, err
	}

	ownerId, err := sdk.ResolveOrInitAddress(ctx, owner)
	if err != nil {
		return nil, err
	}

	operatorId, err := sdk.ResolveOrInitAddress(ctx, operator)
	if err != nil {
		return nil, err
	}

	return t.changeAllowanceBy(ctx, ownerId, operatorId, delta.Neg())
}

func (t *Frc46Token) RevokeAllowance(ctx context.Context, params *RevokeAllowanceParams) (*abi.TokenAmount, error) {
	callerAddr, err := sdk.CallerAddress(ctx)
	if err != nil {
		return nil, err
	}

	return t.revokeAllowance(ctx, callerAddr, params.Operator)
}

func (t *Frc46Token) revokeAllowance(ctx context.Context, owner, operator address.Address) (*abi.TokenAmount, error) {
	ownerId, err := sdk.ResolveAddress(ctx, owner)
	if err != nil {
		if errors.Is(err, ferrors.NotFound) {
			return &zero, nil
		}
		return nil, err
	}

	operatorId, err := sdk.ResolveAddress(ctx, operator)
	if err != nil {
		if errors.Is(err, ferrors.NotFound) {
			return &zero, nil
		}
		return nil, err
	}

	ownerKey := types.ActorKey(ownerId)
	allowanceMap, found, err := t.getAllowanceMap(ctx, ownerId)
	if err != nil {
		return nil, err
	}
	if !found {
		return &zero, nil
	}

	operatorKey := types.ActorKey(operatorId)
	oldAmount, err := allowanceMap.GetAllowanceBalance(operatorId)
	if err != nil {
		return nil, err
	}

	if err = allowanceMap.Delete(operatorKey); err != nil {
		return nil, err
	}

	globalAllowanceMap, err := t.getGlobalAllowanceMap(ctx)
	if err != nil {
		return nil, err
	}

	if allowanceMap.IsEmpty() {
		err = globalAllowanceMap.Delete(ownerKey)
		if err != nil {
			return nil, err
		}
	} else {
		newAllowanceMapCid, err := allowanceMap.Root()
		if err != nil {
			return nil, err
		}
		globalAllowanceMap.Put(ownerKey, typegen.CborCid(newAllowanceMapCid))
	}

	t.Allowances, err = globalAllowanceMap.Root()
	if err != nil {
		return nil, err
	}
	return oldAmount, nil
}

func (t *Frc46Token) Burn(ctx context.Context, params *BurnParams) (*BurnReturn, error) {
	callerAddr, err := sdk.CallerAddress(ctx)
	if err != nil {
		return nil, err
	}

	return t.burn(ctx, callerAddr, params.Amount)
}

func (t *Frc46Token) burn(ctx context.Context, owner address.Address, amount abi.TokenAmount) (*BurnReturn, error) {
	err := ValidateAmountWithGranularity(amount, "burn", t.Granularity)
	if err != nil {
		return nil, err
	}
	ownerId, err := sdk.ResolveAddress(ctx, owner)
	if err != nil {
		return nil, err
	}

	// attempt to burn the requested amount
	newAmount, err := t.changeBalanceBy(ctx, ownerId, amount.Neg())
	if err != nil {
		return nil, err
	}
	// decrease total_supply
	err = t.changeSupplyBy(amount.Neg())
	if err != nil {
		return nil, err
	}
	return &BurnReturn{
		Balance: *newAmount,
	}, nil
}

func (t *Frc46Token) BurnFrom(ctx context.Context, params *BurnFromParams) (*BurnFromReturn, error) {
	callerAddr, err := sdk.CallerAddress(ctx)
	if err != nil {
		return nil, err
	}

	return t.burnFrom(ctx, callerAddr, params.Owner, params.Amount)
}

func (t *Frc46Token) burnFrom(ctx context.Context, operator, owner address.Address, amount abi.TokenAmount) (*BurnFromReturn, error) {
	err := ValidateAmountWithGranularity(amount, "transfer from", t.Granularity)
	if err != nil {
		return nil, err
	}

	if sdk.SameAddress(ctx, operator, owner) {
		return nil, fmt.Errorf("operator cannot be the same as the debited address %s %w", operator, ferrors.USR_ILLEGAL_ARGUMENT)
	}

	operatorId, err := sdk.ResolveAddress(ctx, operator)
	if err != nil {
		if errors.Is(err, ferrors.NotFound) {
			return nil, fmt.Errorf("operator %s not found %w", operator, ferrors.USR_ILLEGAL_ARGUMENT)
		}
	}

	ownerId, err := sdk.ResolveAddress(ctx, owner)
	if err != nil {
		if errors.Is(err, ferrors.NotFound) {
			return nil, fmt.Errorf("owner address %s not found %w", owner, ferrors.USR_ILLEGAL_ARGUMENT)
		}
	}

	newAllowance, err := t.attemptUseAllowance(ctx, operatorId, ownerId, amount)
	if err != nil {
		return nil, err
	}

	// attempt to burn the requested amount
	newBalance, err := t.changeBalanceBy(ctx, ownerId, amount.Neg())
	if err != nil {
		return nil, err
	}

	// decrease total_supply
	err = t.changeSupplyBy(amount.Neg())
	if err != nil {
		return nil, err
	}
	return &BurnFromReturn{
		Balance:   *newBalance,
		Allowance: *newAllowance,
	}, nil
}

////////////////////////Retrun ///////////////////////

// / Generate TransferReturn from the intermediate data returned by a receiver hook call
func (t *Frc46Token) transferReturn(ctx context.Context, intermediate *TransferIntermediate) (*TransferReturn, error) {
	fromBalance, err := t.getBalance(ctx, intermediate.From)
	if err != nil {
		return nil, err
	}
	toBalance, err := t.getBalance(ctx, intermediate.To)
	if err != nil {
		return nil, err
	}
	return &TransferReturn{
		FromBalance:   *fromBalance,
		ToBalance:     *toBalance,
		RecipientData: intermediate.RecipientData,
	}, nil
}

// / Converts a TransferFromIntermediate to a TransferFromReturn
// /
// / This function should be called on a freshly loaded or known-up-to-date state
func (t *Frc46Token) transferFromReturn(ctx context.Context, intermediate *TransferFromIntermediate) (*TransferFromReturn, error) {
	fromBalance, err := t.getBalance(ctx, intermediate.From)
	if err != nil {
		return nil, err
	}
	toBalance, err := t.getBalance(ctx, intermediate.To)
	if err != nil {
		return nil, err
	}
	allowanceBalance, err := t.getAllowanceBalance(ctx, intermediate.From, intermediate.Operator)
	if err != nil {
		return nil, err
	}
	return &TransferFromReturn{
		FromBalance:   *fromBalance,
		To_balance:    *toBalance,
		Allowance:     *allowanceBalance,
		RecipientData: intermediate.RecipientData,
	}, nil
}

// / Finalise return data from MintIntermediate data returned by calling receiver hook after minting
// / This is done to allow reloading the state if it changed as a result of the hook call
// / so we can return an accurate balance even if the receiver transferred or burned tokens upon receipt
func (t *Frc46Token) mintReturn(ctx context.Context, intermediate *MintIntermediate) (*MintReturn, error) {
	recipientBalance, err := t.getBalance(ctx, intermediate.Recipient)
	if err != nil {
		return nil, err
	}

	return &MintReturn{
		Balance:       *recipientBalance,
		Supply:        t.Supply,
		RecipientData: intermediate.RecipientData,
	}, nil
}

// //////////////////////Utils ///////////////////////

func (t *Frc46Token) getAllowanceBalance(ctx context.Context, owner, operator abi.ActorID) (*abi.TokenAmount, error) {
	allowanceMap, found, err := t.getAllowanceMap(ctx, owner)
	if err != nil {
		return nil, err
	}
	if !found {
		return &zero, nil
	}
	return allowanceMap.GetAllowanceBalance(operator)
}

func (t *Frc46Token) getGlobalAllowanceMap(ctx context.Context) (*adt.Map, error) {
	balanceMap, err := adt.AsMap(adt.AdtStore(ctx), t.Allowances, DEFAULTHAMTBITWIDTH)
	if err != nil {
		return nil, err
	}
	return balanceMap, nil
}

func (t *Frc46Token) getAllowanceMap(ctx context.Context, owner abi.ActorID) (*AllowanceMap, bool, error) {
	balanceMap, err := adt.AsMap(adt.AdtStore(ctx), t.Allowances, DEFAULTHAMTBITWIDTH)
	if err != nil {
		return nil, false, err
	}

	var allowanceMapCid typegen.CborCid
	found, err := balanceMap.Get(types.ActorKey(owner), &allowanceMapCid)
	if err != nil {
		return nil, false, err
	}

	if !found {
		return nil, false, nil
	}

	allowanceMap, err := adt.AsMap(adt.AdtStore(ctx), cid.Cid(allowanceMapCid), DEFAULTHAMTBITWIDTH)
	if err != nil {
		return nil, false, err
	}
	return &AllowanceMap{allowanceMap}, true, nil
}

func (t *Frc46Token) changeBalanceBy(ctx context.Context, owner abi.ActorID, delta abi.TokenAmount) (*abi.TokenAmount, error) {
	balanceMap, err := adt.AsMap(adt.AdtStore(ctx), t.Balances, adt.BalanceTableBitwidth)
	if err != nil {
		return nil, err
	}
	var balance = big.Int{Int: stdbig.NewInt(0)}
	_, err = balanceMap.Get(types.ActorKey(owner), &balance)
	if err != nil {
		return nil, err
	}

	if delta.IsZero() {
		// This is a no-op as far as mutating state
		return &balance, nil
	}

	newBalance := big.Add(balance, delta)
	if newBalance.Sign() < 0 {
		return nil, fmt.Errorf("negative balance caused by decreasing %s's balance of %s by %s", owner, balance, delta)
	}

	if newBalance.IsZero() {
		err = balanceMap.Delete(types.ActorKey(owner))
		if err != nil {
			return nil, err
		}
	} else {
		err = balanceMap.Put(types.ActorKey(owner), &newBalance)
		if err != nil {
			return nil, err
		}
	}

	t.Balances, err = balanceMap.Root()
	if err != nil {
		return nil, err
	}
	return &newBalance, nil
}

// / Change the allowance between owner and operator by the specified delta
func (t *Frc46Token) changeAllowanceBy(ctx context.Context, owner, operator abi.ActorID, delta abi.TokenAmount) (*abi.TokenAmount, error) {
	if delta.IsZero() {
		// This is a no-op as far as mutating state
		return t.getAllowanceBalance(ctx, owner, operator)
	}
	adtStore := adt.AdtStore(ctx)

	globalAllowancesMap, err := adt.AsMap(adtStore, t.Allowances, DEFAULTHAMTBITWIDTH)
	if err != nil {
		return nil, err
	}

	ownerKey := types.ActorKey(owner)
	var allowanceMapCid typegen.CborCid
	found, err := globalAllowancesMap.Get(ownerKey, &allowanceMapCid)
	if err != nil {
		return nil, err
	}

	var allowanceMap *adt.Map
	if found {
		allowanceMap, err = adt.AsMap(adtStore, cid.Cid(allowanceMapCid), DEFAULTHAMTBITWIDTH)
		if err != nil {
			return nil, err
		}
	} else {
		allowanceMap, err = adt.MakeEmptyMap(adtStore, DEFAULTHAMTBITWIDTH)
		if err != nil {
			return nil, err
		}
	}

	operatorKey := types.ActorKey(operator)
	oldAllowance := big.Zero()
	_, err = allowanceMap.Get(operatorKey, &oldAllowance)
	if err != nil {
		return nil, err
	}

	newAllowance := big.Max(big.Add(oldAllowance, delta), big.Zero())

	// if the new allowance is zero, we can remove the entry from the state tree
	if newAllowance.IsZero() {
		if err = allowanceMap.Delete(operatorKey); err != nil {
			return nil, err
		}
	} else {
		if err = allowanceMap.Put(operatorKey, &newAllowance); err != nil {
			return nil, err
		}
	}
	if allowanceMap.IsEmpty() {
		if err = globalAllowancesMap.Delete(ownerKey); err != nil {
			return nil, err
		}
	}
	newAllowanceMapCid, err := allowanceMap.Root()
	if err != nil {
		return nil, err
	}
	if err = globalAllowancesMap.Put(ownerKey, typegen.CborCid(newAllowanceMapCid)); err != nil {
		return nil, err
	}
	t.Allowances, err = globalAllowancesMap.Root()
	if err != nil {
		return nil, err
	}
	return &newAllowance, nil
}

// / Increase/decrease the total supply by the specified value
// /
// / Returns the new total supply
func (t *Frc46Token) changeSupplyBy(delta abi.TokenAmount) error {
	newSupply := big.Add(t.Supply, delta)
	if newSupply.Sign() < 0 {
		return fmt.Errorf("supply must big than 0 supply %s, delta %s %w", t.Supply, delta, ferrors.USR_ILLEGAL_ARGUMENT)
	}

	t.Supply = newSupply
	return nil
}

func (t *Frc46Token) attemptUseAllowance(ctx context.Context, operator, owner abi.ActorID, amount abi.TokenAmount) (*abi.TokenAmount, error) {
	curAllowBalance, err := t.getAllowanceBalance(ctx, owner, operator)
	if err != nil {
		return nil, err
	}

	if amount.IsZero() {
		return &zero, err
	}

	if curAllowBalance.IsZero() && operator != owner {
		return nil, InsufficientAllowanceError(sdk.MustAddressFromActorId(operator), sdk.MustAddressFromActorId(owner), amount, *curAllowBalance)
	}

	if curAllowBalance.LessThan(amount) {
		return nil, InsufficientAllowanceError(sdk.MustAddressFromActorId(operator), sdk.MustAddressFromActorId(owner), amount, *curAllowBalance)
	}
	// new_allowance = current_allowance - amount;
	return t.changeAllowanceBy(ctx, owner, operator, amount.Neg())
}

type AllowanceMap struct {
	*adt.Map
}

func (allowanceMap *AllowanceMap) GetAllowanceBalance(actorId abi.ActorID) (*abi.TokenAmount, error) {
	var balance = big.Int{Int: stdbig.NewInt(0)}
	_, err := allowanceMap.Get(types.ActorKey(actorId), &balance)
	if err != nil {
		return nil, err
	}
	return &balance, nil
}

// / Validates that a token amount for burning/transfer/minting is non-negative, and an integer
// / multiple of granularity.
// /
// / Returns the argument, or an error.
func ValidateAmountWithGranularity(a abi.TokenAmount, name string, granularity uint64) error {
	if a.Sign() < 0 {
		return fmt.Errorf("value %s for %s must be non-negative %w", a, name, ferrors.USR_ILLEGAL_ARGUMENT)
	}
	rem := big.NewInt(0).Rem(a.Int, big.NewIntUnsigned(granularity).Int)
	if rem.Sign() != 0 {
		return fmt.Errorf("amount %s for %s must be a multiple of %d %w", a, name, granularity, ferrors.USR_ILLEGAL_ARGUMENT)
	}
	return nil
}

// / Validates that an allowance is non-negative. Allowances do not need to be an integer multiple of
// / granularity.
// /
// / Returns the argument, or an error.
func ValidateAllowance(a abi.TokenAmount, name string) error {
	if a.Sign() < 0 {
		return fmt.Errorf("method %s allowance %s is negative %w", name, a, ferrors.USR_ILLEGAL_ARGUMENT)
	}
	return nil
}
