package contract

import (
	"context"
	"errors"
	"fmt"

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

// Frc46Token token state IPLD structure
type Frc46Token struct {
	Name        string
	Symbol      string
	Granularity uint64
	Owner       abi.ActorID
	/// Total supply of token
	Supply abi.TokenAmount
	/// Map<ActorId, TokenAmount> of balances as a Hamt
	Balances cid.Cid
	/// Map<ActorId, Map<ActorId, TokenAmount>> as a Hamt. Allowances are stored balances[owner][operator]
	Allowances cid.Cid
}

var _ IFrc46Token = (*Frc46Token)(nil)
var _ IFrc46Unspecific = (*Frc46Token)(nil)

func (t *Frc46Token) Export() []interface{} {
	return []interface{}{
		Constructor,
		sdk.MethodInfo{
			Alias:    "Name",
			Func:     t.GetName,
			Readonly: true,
		},
		sdk.MethodInfo{
			Alias:    "Symbol",
			Func:     t.GetSymbol,
			Readonly: true,
		},
		sdk.MethodInfo{
			Alias:    "Granularity",
			Func:     t.GetGranularity,
			Readonly: true,
		},
		sdk.MethodInfo{
			Alias:    "TotalSupply",
			Func:     t.GetTotalSupply,
			Readonly: true,
		},
		t.Mint,
		sdk.MethodInfo{
			Func:     t.BalanceOf,
			Readonly: true,
		},
		sdk.MethodInfo{
			Func:     t.Allowance,
			Readonly: true,
		},
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

	originId, err := sdk.Origin(ctx)
	if err != nil {
		return err
	}

	err = emptyMap.Put(types.ActorKey(originId), &req.Supply)
	if err != nil {
		return err
	}

	balanceRoot, err := emptyMap.Root()
	if err != nil {
		return err
	}

	state := &Frc46Token{
		Name:        req.Name,
		Symbol:      req.Symbol,
		Granularity: req.Granularity,
		Supply:      req.Supply,
		Owner:       originId,
		Balances:    balanceRoot,
		Allowances:  emptyRoot,
	}

	logger.Logf(ctx, "create token %s, symbol %s owner %d supply %s", req.Name, state.Symbol, originId, &req.Supply)
	_ = sdk.Constructor(ctx, state)
	return nil
}

// GetName gets actor's name
func (t *Frc46Token) GetName(_ context.Context) types.CborString {
	return types.CborString(t.Name)
}

// GetSymbol gets actor's symbol
func (t *Frc46Token) GetSymbol(_ context.Context) types.CborString {
	return types.CborString(t.Symbol)
}

// GetGranularity gets the granularity of actor
func (t *Frc46Token) GetGranularity(_ context.Context) types.CborUint {
	return types.CborUint(t.Granularity)
}

// GetTotalSupply gets the total number of tokens in existence
// this equals the sum of `balance_of` called on all addresses. This equals sum of all
// successful `mint` calls minus the sum of all successful `burn`/`burn_from` calls
func (t *Frc46Token) GetTotalSupply(_ context.Context) *abi.TokenAmount {
	return &t.Supply
}

// Mint mints the specified value of tokens into an account
// the minter is implicitly defined as the caller of the actor, and must be an ID address.
// the mint amount must be non-negative or the method returns an error.
//
// returns a ReceiverHook to call the owner's token receiver hook,
// and the owner's new balance.
// receiverHook must be called or it will panic and abort the transaction.
//
// the hook call will return a MintIntermediate struct which must be passed to mint_return
// to get the final return data
func (t *Frc46Token) Mint(ctx context.Context, params *MintParams) (*MintReturn, error) {
	callerId, err := sdk.Caller(ctx)
	if err != nil {
		return nil, err
	}

	if callerId != t.Owner {
		return nil, fmt.Errorf("caller %d is not owner(%d) of actor %w", callerId, t.Owner, ferrors.USR_FORBIDDEN)
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

	sdk.SaveState(ctx, t)
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

// BalanceOf returns the balance associated with a particular address
// accounts that have never received transfers implicitly have a zero-balance
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
	balanceMap, err := adt.AsMap(adt.AdtStore(ctx), t.Balances, DEFAULTHAMTBITWIDTH)
	if err != nil {
		return nil, err
	}
	var balance = abi.NewTokenAmount(0)
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

// Allowance gets the allowance between owner and operator
// an allowance is the amount that the operator can transfer or burn out of the owner's account
// via the `transfer` and `burn` methods.
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

// Transfer transfers an amount from the caller to another address
//
//   - The requested value MUST be non-negative
//
//   - The requested value MUST NOT exceed the sender's balance
//
//   - The receiving actor MUST implement a method called `tokens_received`, corresponding to the
//     interface specified for FRC-0046 token receiver. If the receiving hook aborts, when called,
//     the transfer is discarded and this method returns an error
//
//     Upon successful transfer:
//
//   - The from balance decreases by the requested value
//
//   - The to balance increases by the requested value
//
//     Returns a ReceiverHook to call the recipient's token receiver hook,
//     and a TransferIntermediate struct
//     ReceiverHook must be called or it will panic and abort the transaction.
//
//     Return data from the hook should be passed to transfer_return which will generate
//     the Transfereturn struct
func (t *Frc46Token) Transfer(ctx context.Context, params *TransferParams) (*TransferReturn, error) {
	callerAddr, err := sdk.CallerAddress(ctx)
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
	sdk.SaveState(ctx, t)
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

// TransferFrom Transfers an amount from one address to another
//   - The requested value MUST be non-negative
//   - The requested value MUST NOT exceed the sender's balance
//   - The receiving actor MUST implement a method called `tokens_received`, corresponding to the
//     interface specified for FRC-0046 token receiver. If the receiving hook aborts, when called,
//     the transfer is discarded and this method returns an error
//   - The operator MUST be initialised AND have an allowance not less than the requested value
//
// Upon successful transfer:
//   - The from balance decreases by the requested value
//   - The to balance increases by the requested value
//   - The owner-operator allowance decreases by the requested value
//
// Returns a ReceiverHook to call the recipient's token receiver hook,
// and a TransferFromIntermediate struct.
// ReceiverHook must be called or it will panic and abort the transaction.
//
// Return data from the hook should be passed to transfer_from_return which will generate
// the TransferFromReturn struct
func (t *Frc46Token) TransferFrom(ctx context.Context, params *TransferFromParams) (*TransferFromReturn, error) {
	operatorAddr, err := sdk.CallerAddress(ctx)
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
	sdk.SaveState(ctx, t)
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

// IncreaseAllowance increase the allowance that an operator can control of an owner's balance by the requested delta
// returns an error if requested delta is negative or there are errors in (de)serialization of
// state.If either owner or operator addresses are not resolvable and cannot be initialised, this
// method returns MessagingError::AddressNotInitialized.
//
// else returns the new allowance
func (t *Frc46Token) IncreaseAllowance(ctx context.Context, params *IncreaseAllowanceParams) (*abi.TokenAmount, error) {
	addr, err := sdk.CallerAddress(ctx)
	if err != nil {
		return nil, err
	}

	amount, err := t.increaseAllowance(ctx, addr, params.Operator, params.Increase)
	if err != nil {
		return nil, err
	}
	sdk.SaveState(ctx, t)
	return amount, nil
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

// DecreaseAllowance decrease the allowance that an operator controls of the owner's balance by the requested delta
// returns an error if requested delta is negative or there are errors in (de)serialization of
// of state. If the resulting allowance would be negative, the allowance between owner and operator is
// set to zero.Returns an error if either the operator or owner addresses are not resolvable and
// cannot be initialized.
//
// else returns the new allowance
func (t *Frc46Token) DecreaseAllowance(ctx context.Context, params *DecreaseAllowanceParams) (*abi.TokenAmount, error) {
	addr, err := sdk.CallerAddress(ctx)
	if err != nil {
		return nil, err
	}

	amount, err := t.decreaseAllowance(ctx, addr, params.Operator, params.Decrease)
	if err != nil {
		return nil, err
	}
	sdk.SaveState(ctx, t)
	return amount, nil
}

func (t *Frc46Token) decreaseAllowance(ctx context.Context, owner address.Address, operator address.Address, delta abi.TokenAmount) (*abi.TokenAmount, error) {
	err := ValidateAmountWithGranularity(delta, "increase allowance delta", t.Granularity)
	if err != nil {
		return nil, err
	}

	err = ValidateAllowance(delta, "decrease allowance delta")
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

// RevokeAllowance sets the allowance between owner and operator to zero, returning the old allowance
func (t *Frc46Token) RevokeAllowance(ctx context.Context, params *RevokeAllowanceParams) (*abi.TokenAmount, error) {
	callerAddr, err := sdk.CallerAddress(ctx)
	if err != nil {
		return nil, err
	}

	amount, err := t.revokeAllowance(ctx, callerAddr, params.Operator)
	if err != nil {
		return nil, err
	}
	sdk.SaveState(ctx, t)
	return amount, nil

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

// Burn burns an amount of token from the specified address, decreasing total token supply
//   - The requested value MUST be non-negative
//   - The requested value MUST NOT exceed the target's balance
//   - If the burn operation would result in a negative balance for the owner, the burn is discarded and this method returns an error
//     Upon successful burn
//   - The target's balance decreases by the requested value
//   - The total_supply decreases by the requested value
func (t *Frc46Token) Burn(ctx context.Context, params *BurnParams) (*BurnReturn, error) {
	callerAddr, err := sdk.CallerAddress(ctx)
	if err != nil {
		return nil, err
	}

	burnRet, err := t.burn(ctx, callerAddr, params.Amount)
	if err != nil {
		return nil, err
	}
	sdk.SaveState(ctx, t)
	return burnRet, nil
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

// BurnFrom burns an amount of token from the specified address, decreasing total token supply
//
// If operator and owner are the same address, this method returns an InvalidOperator error.
//
//   - The requested value MUST be non-negative
//   - The requested value MUST NOT exceed the target's balance
//   - If the burn operation would result in a negative balance for the owner, the burn is discarded and this method returns an error
//   - The operator MUST have an allowance not less than the requested value
//
// Upon successful burn
//   - The target's balance decreases by the requested value
//   - The total_supply decreases by the requested value
//   - The operator's allowance is decreased by the requested value
func (t *Frc46Token) BurnFrom(ctx context.Context, params *BurnFromParams) (*BurnFromReturn, error) {
	callerAddr, err := sdk.CallerAddress(ctx)
	if err != nil {
		return nil, err
	}

	burnReturnRet, err := t.burnFrom(ctx, callerAddr, params.Owner, params.Amount)
	if err != nil {
		return nil, err
	}
	sdk.SaveState(ctx, t)
	return burnReturnRet, nil
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

////////////////////////Return ///////////////////////

// transferReturn generate TransferReturn from the intermediate data returned by a receiver hook call
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

// transferFromReturn converts a TransferFromIntermediate to a TransferFromReturn
// this function should be called on a freshly loaded or known-up-to-date state
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

// mintReturn finalise return data from MintIntermediate data returned by calling receiver hook after minting
// this is done to allow reloading the state if it changed as a result of the hook call
// so we can return an accurate balance even if the receiver transferred or burned tokens upon receipt
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
	balanceMap, err := adt.AsMap(adt.AdtStore(ctx), t.Balances, DEFAULTHAMTBITWIDTH)
	if err != nil {
		return nil, err
	}
	var balance = abi.NewTokenAmount(0)
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

// changeAllowanceBy change the allowance between owner and operator by the specified delta
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

// changeSupplyBy increase/decrease the total supply by the specified value
// returns the new total supply
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

// GetAllowanceBalance get allowance map of specific action id
func (allowanceMap *AllowanceMap) GetAllowanceBalance(actorId abi.ActorID) (*abi.TokenAmount, error) {
	var balance = abi.NewTokenAmount(0)
	_, err := allowanceMap.Get(types.ActorKey(actorId), &balance)
	if err != nil {
		return nil, err
	}
	return &balance, nil
}

// ValidateAmountWithGranularity validates that a token amount for burning/transfer/minting is non-negative, and an integer multiple of granularity.
// returns the argument, or an error.
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

// ValidateAllowance validates that an allowance is non-negative. Allowances do not need to be an integer multiple of granularity.
// returns the argument, or an error.
func ValidateAllowance(a abi.TokenAmount, name string) error {
	if a.Sign() < 0 {
		return fmt.Errorf("method %s allowance %s is negative %w", name, a, ferrors.USR_ILLEGAL_ARGUMENT)
	}
	return nil
}
