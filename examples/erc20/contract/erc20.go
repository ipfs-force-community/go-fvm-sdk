package contract

import (
	"context"
	"errors"
	"fmt"
	stdbig "math/big"
	"strconv"

	"github.com/ipfs/go-cid"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/adt"
	typegen "github.com/whyrusleeping/cbor-gen"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"

	"github.com/filecoin-project/go-state-types/abi"

	"github.com/filecoin-project/go-state-types/big"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk"

	"github.com/filecoin-project/go-address"
)

/*
function Name() public view returns (string)
function Symbol() public view returns (string)
function Decimals() public view returns (uint8)
function TotalSupply() public view returns (uint256)
function balanceOf(address _owner) public view returns (uint256 balance)

function transfer(address _to, uint256 _value) public returns (bool success)

function transferFrom(address _from, address _to, uint256 _value) public returns (bool success)
function approve(address _spender, uint256 _value) public returns (bool success)
function allowance(address _owner, address _spender) public view returns (uint256 remaining)
*/

//keep unused for code generation

var logger sdk.Logger

func init() {
	logger, _ = sdk.NewLogger()
}

// Erc20Token basic Token20
type Erc20Token struct {
	Name        string
	Symbol      string
	Decimals    uint8
	TotalSupply abi.TokenAmount

	//todo cbor gen not support non-string key and map value
	Balances cid.Cid //map[string]*big.Int
	Allowed  cid.Cid // map[string]*big.Int //owner-spender
}

func (t *Erc20Token) Export() []interface{} {
	return []interface{}{
		Constructor,
		t.GetName,
		t.GetSymbol,
		t.GetDecimal,
		t.GetTotalSupply,
		t.GetBalanceOf,
		t.Transfer,
		t.TransferFrom,
		t.Approval,
		t.Allowance,
	}
}

type ConstructorReq struct {
	Name        string
	Symbol      string
	Decimals    uint8
	TotalSupply abi.TokenAmount
}

func Constructor(ctx context.Context, req *ConstructorReq) error {
	emptyMap, err := adt.MakeEmptyMap(adt.AdtStore(ctx), adt.BalanceTableBitwidth)
	if err != nil {
		return err
	}
	emptyRoot, err := emptyMap.Root()
	if err != nil {
		return err
	}
	caller, err := sdk.Caller(ctx)
	if err != nil {
		return err
	}

	//todo call is init actor but not message real sendor wait for ref fvm fix this issue
	err = emptyMap.Put(types.ActorKey(caller), &req.TotalSupply)
	if err != nil {
		return err
	}

	originId, err := sdk.Origin(ctx)
	if err != nil {
		return err
	}

	err = emptyMap.Put(types.ActorKey(originId), &req.TotalSupply)
	if err != nil {
		return err
	}

	balanceRoot, err := emptyMap.Root()
	if err != nil {
		return err
	}

	state := &Erc20Token{
		Name:        req.Name,
		Symbol:      req.Symbol,
		Decimals:    req.Decimals,
		TotalSupply: req.TotalSupply,
		Balances:    balanceRoot,
		Allowed:     emptyRoot,
	}

	logger.Logf(ctx, "construct token %s  issue %s token to %s", req.Name, req.TotalSupply.String(), actorToString(caller))
	_ = sdk.Constructor(ctx, state)
	return nil
}

type FakeSetBalance struct {
	Addr    address.Address
	Balance abi.TokenAmount
}

// GetName return token Name of erc20 token
func (t *Erc20Token) GetName() types.CborString {
	return types.CborString(t.Name)
}

// GetDecimal return token Symbol of erc20 token
func (t *Erc20Token) GetSymbol() types.CborString {
	return types.CborString(t.Symbol)
}

// GetDecimal return decimal of erc20 token
func (t *Erc20Token) GetDecimal() typegen.CborInt {
	return typegen.CborInt(t.Decimals)
}

// GetTotalSupply returns total number of tokens in existence
func (t *Erc20Token) GetTotalSupply() *abi.TokenAmount {
	return &t.TotalSupply
}

/*
GetBalanceOf sender by ID.

* `args[0]` - the ID of user.
*/
func (t *Erc20Token) GetBalanceOf(ctx context.Context, addr *address.Address) (*big.Int, error) {
	senderId, err := sdk.ResolveAddress(ctx, *addr)
	if err != nil {
		return nil, err
	}
	return t.getBalanceOf(ctx, senderId)
}

func (t *Erc20Token) getBalanceOf(ctx context.Context, act abi.ActorID) (*big.Int, error) {
	balanceMap, err := adt.AsMap(adt.AdtStore(ctx), t.Balances, adt.BalanceTableBitwidth)
	if err != nil {
		return nil, err
	}
	var balance = &big.Int{Int: stdbig.NewInt(0)}
	_, err = balanceMap.Get(types.ActorKey(act), balance)
	if err != nil {
		return nil, err
	}
	//return 0 if not exit
	return balance, nil
}

type TransferReq struct {
	ReceiverAddr   address.Address
	TransferAmount abi.TokenAmount
}

/*
Transfer token from current caller to a specified address.

* `receiverAddr` - the ID of receiver.

* `transferAmount` - the transfer amount.
*/
func (t *Erc20Token) Transfer(ctx context.Context, transferReq *TransferReq) error {
	senderID, err := sdk.Caller(ctx)

	if err != nil {
		return err
	}

	receiverID, err := sdk.ResolveAddress(ctx, transferReq.ReceiverAddr)
	if err != nil {
		return err
	}

	if transferReq.TransferAmount.LessThanEqual(big.Zero()) {
		return errors.New("transfer value must bigger than zero")
	}

	balanceOfSender, err := t.getBalanceOf(ctx, senderID)
	if err != nil {
		return err
	}
	balanceOfReceiver, err := t.getBalanceOf(ctx, receiverID)
	if err != nil {
		return err
	}

	if err := checkBalance(balanceOfSender, senderID); err != nil {
		return err
	}

	if err := isSmallerOrEqual(&transferReq.TransferAmount, balanceOfSender); err != nil {
		return fmt.Errorf("transfer amount should be less than balance of sender (%v): %v", senderID, err)
	}

	balanceMap, err := adt.AsMap(adt.AdtStore(ctx), t.Balances, adt.BalanceTableBitwidth)
	if err != nil {
		return err
	}

	if err = balanceMap.Put(types.ActorKey(senderID), sub(balanceOfSender, &transferReq.TransferAmount)); err != nil {
		return err
	}
	if err = balanceMap.Put(types.ActorKey(receiverID), add(balanceOfReceiver, &transferReq.TransferAmount)); err != nil {
		return err
	}
	newBalanceMapRoot, err := balanceMap.Root()
	if err != nil {
		return err
	}
	t.Balances = newBalanceMapRoot
	logger.Logf(ctx, "transfer from %d to %d amount %s", senderID, receiverID, transferReq.TransferAmount.String())
	_ = sdk.SaveState(ctx, t)
	return nil
}

type AllowanceReq struct {
	OwnerAddr   address.Address
	SpenderAddr address.Address
}

/*
Allowance checks the amount of tokens that an owner Allowed a spender to transfer in behalf of the owner to another receiver.

* `ownerAddr` - the ID of owner.

* `spenderAddr` - the ID of spender
*/
func (t *Erc20Token) Allowance(ctx context.Context, req *AllowanceReq) (*big.Int, error) {
	ownerID, err := sdk.ResolveAddress(ctx, req.OwnerAddr)
	if err != nil {
		return nil, err
	}

	spenderId, err := sdk.ResolveAddress(ctx, req.SpenderAddr)
	if err != nil {
		return nil, err
	}

	return t.getAllowance(ctx, ownerID, spenderId)
}

func (t *Erc20Token) getAllowance(ctx context.Context, ownerID, spenderId abi.ActorID) (*big.Int, error) {
	allowBalanceMap, err := adt.AsMap(adt.AdtStore(ctx), t.Allowed, adt.BalanceTableBitwidth)
	if err != nil {
		return nil, err
	}

	balance := &big.Int{Int: stdbig.NewInt(0)}
	if _, err = allowBalanceMap.Get(types.StringKey(getAllowKey(ownerID, spenderId)), balance); err != nil {
		return nil, err
	}

	return balance, nil
}

type TransferFromReq struct {
	OwnerAddr      address.Address
	ReceiverAddr   address.Address
	TransferAmount abi.TokenAmount
}

/*
TransferFrom transfer tokens from token owner to receiver.

* `ownerAddr` - the ID of token owner.

* `receiverAddr` - the ID of receiver.

* `transferAmount` - the transfer amount.
*/
func (t *Erc20Token) TransferFrom(ctx context.Context, req *TransferFromReq) error {
	tokenOwnerID, err := sdk.ResolveAddress(ctx, req.OwnerAddr)
	if err != nil {
		return err
	}

	receiverID, err := sdk.ResolveAddress(ctx, req.ReceiverAddr)
	if err != nil {
		return err
	}

	if req.TransferAmount.LessThanEqual(big.Zero()) {
		return errors.New("send value must bigger than zero")
	}

	spenderID, err := sdk.Caller(ctx)
	if err != nil {
		return err
	}
	balanceOfTokenOwner, err := t.getBalanceOf(ctx, tokenOwnerID)
	if err != nil {
		return err
	}
	balanceOfReceiver, err := t.getBalanceOf(ctx, receiverID)
	if err != nil {
		return err
	}
	approvedAmount, err := t.getAllowance(ctx, tokenOwnerID, spenderID)
	if err != nil {
		return err
	}

	if err := checkBalance(balanceOfTokenOwner, tokenOwnerID); err != nil {
		return err
	}
	if err := checkBalance(balanceOfReceiver, receiverID); err != nil {
		return err
	}

	if approvedAmount.LessThanEqual(big.Zero()) {
		return fmt.Errorf("approved amount for %v-%v less than zero", tokenOwnerID, spenderID)
	}

	if err := isSmallerOrEqual(&req.TransferAmount, balanceOfTokenOwner); err != nil {
		return fmt.Errorf("transfer amount should be less than balance of token owner (%v): %v", tokenOwnerID, err)
	}
	if err := isSmallerOrEqual(&req.TransferAmount, approvedAmount); err != nil {
		return fmt.Errorf("transfer amount should be less than approved spending amount of %v: %v", spenderID, err)
	}

	store := adt.AdtStore(ctx)
	balanceMap, err := adt.AsMap(store, t.Balances, adt.BalanceTableBitwidth)
	if err != nil {
		return err
	}

	allowBalanceMap, err := adt.AsMap(store, t.Allowed, adt.BalanceTableBitwidth)
	if err != nil {
		return err
	}

	if err = balanceMap.Put(types.ActorKey(tokenOwnerID), sub(balanceOfTokenOwner, &req.TransferAmount)); err != nil {
		return err
	}

	if err = balanceMap.Put(types.ActorKey(receiverID), add(balanceOfReceiver, &req.TransferAmount)); err != nil {
		return err
	}

	if err = allowBalanceMap.Put(types.StringKey(getAllowKey(tokenOwnerID, spenderID)), sub(approvedAmount, &req.TransferAmount)); err != nil {
		return err
	}

	if t.Balances, err = balanceMap.Root(); err != nil {
		return err
	}
	if t.Allowed, err = allowBalanceMap.Root(); err != nil {
		return err
	}
	_ = sdk.SaveState(ctx, t)
	return nil
}

type ApprovalReq struct {
	SpenderAddr  address.Address
	NewAllowance abi.TokenAmount
}

/*Approval approves the passed-in identity to spend/burn a maximum amount of tokens on behalf of the function caller.
* `spenderAddr` - the ID of approved user.
* `newAllowance` - the maximum approved amount.*/
func (t *Erc20Token) Approval(ctx context.Context, req *ApprovalReq) error {
	spenderID, err := sdk.ResolveAddress(ctx, req.SpenderAddr)
	if err != nil {
		return err
	}

	if req.NewAllowance.LessThanEqual(big.Zero()) {
		return errors.New("allow value must bigger than zero")
	}

	callerID, err := sdk.Caller(ctx)
	if err != nil {
		return err
	}

	allowance, err := t.getAllowance(ctx, callerID, spenderID)
	if err != nil {
		return err
	}

	allowBalanceMap, err := adt.AsMap(adt.AdtStore(ctx), t.Allowed, adt.BalanceTableBitwidth)
	if err != nil {
		return err
	}

	err = allowBalanceMap.Put(types.StringKey(getAllowKey(callerID, spenderID)), add(allowance, &req.NewAllowance))
	if err != nil {
		return err
	}
	t.Allowed, err = allowBalanceMap.Root()
	if err != nil {
		return err
	}
	_ = sdk.SaveState(ctx, t)
	logger.Logf(ctx, "approval %s for %s", getAllowKey(callerID, spenderID), req.NewAllowance.String())
	return nil
}

/*checkBalance checks if sender's balance is >= 0*/
func checkBalance(balance *big.Int, mspID abi.ActorID) error {
	if balance.LessThan(big.Zero()) {
		return fmt.Errorf("balance of sender %v is %v", mspID, balance)
	}
	return nil
}

/*isSmallerOrEqual returns `nil` if a is <= b*/
func isSmallerOrEqual(a *abi.TokenAmount, b *abi.TokenAmount) error {
	if a.GreaterThan(*b) {
		return fmt.Errorf("%v should be <= to %v", a, b)
	}
	return nil
}

func actorToString(act abi.ActorID) string {
	return strconv.FormatUint(uint64(act), 10)
}

func actorFromString(actStr string) abi.ActorID {
	val, _ := strconv.ParseUint(actStr, 10, 64)
	return abi.ActorID(val)
}

func getAllowKey(ownerID, spenderId abi.ActorID) string {
	return actorToString(ownerID) + actorToString(spenderId)
}

func sub(a, b *abi.TokenAmount) *big.Int {
	return &big.Int{Int: big.NewInt(0).Sub(a.Int, b.Int)}
}

func add(a, b *abi.TokenAmount) *big.Int {
	return &big.Int{Int: big.NewInt(0).Add(a.Int, b.Int)}
}
