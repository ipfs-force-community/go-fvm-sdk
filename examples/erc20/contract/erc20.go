package contract

import (
	"errors"
	"fmt"
	"strconv"

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

var zero = big.Zero()

var logger sdk.Logger

func init() {
	logger, _ = sdk.NewLogger()
}

/*basic Token20*/
type Erc20Token struct {
	Name        string
	Symbol      string
	Decimals    uint8
	TotalSupply *big.Int

	//todo cbor gen not support non-string key and map value
	Balances map[string]*big.Int
	Allowed  map[string]*big.Int //owner-spender
}

func (e *Erc20Token) Export() map[int]interface{} {
	return map[int]interface{}{
		1:  e.Constructor,
		2:  e.GetName,
		3:  e.GetSymbol,
		4:  e.GetDecimal,
		5:  e.GetTotalSupply,
		6:  e.GetBalanceOf,
		7:  e.Transfer,
		8:  e.TransferFrom,
		9:  e.Approval,
		10: e.Allowance,
		11: e.FakeSetBalance,
	}
}

type ConstructorReq struct {
	Name        string
	Symbol      string
	Decimals    uint8
	TotalSupply *big.Int
}

func (t *Erc20Token) Constructor(req *ConstructorReq) error {
	state := &Erc20Token{
		Name:        req.Name,
		Symbol:      req.Symbol,
		Decimals:    req.Decimals,
		TotalSupply: req.TotalSupply,
		Balances:    make(map[string]*big.Int),
		Allowed:     make(map[string]*big.Int),
	}
	caller, err := sdk.Caller()
	if err != nil {
		return err
	}
	state.Balances[actorToString(caller)] = req.TotalSupply //todo call is init actor but not message real sendor wait for ref fvm fix this issue
	logger.Logf("construct %s  set %s token to %s", req.Name, actorToString(caller), req.TotalSupply.String())
	_ = sdk.Constructor(state)
	return nil
}

type FakeSetBalance struct {
	Addr    address.Address
	Balance *big.Int
}

// FakeSetBalance NOTICE **beacause unable to set init balance in constructor, so just add this function for test contract, and must be remove after fvm support get really caller**
func (t *Erc20Token) FakeSetBalance(req *FakeSetBalance) error {
	addrId, err := sdk.ResolveAddress(req.Addr)
	if err != nil {
		return err
	}
	t.Balances[actorToString(addrId)] = req.Balance
	_ = sdk.SaveState(t)
	return nil
}

/*GetDecimal return token Name of erc20 token*/
func (t *Erc20Token) GetName() types.CborString {
	return types.CborString(t.Name)
}

/*GetDecimal return token Symbol of erc20 token*/
func (t *Erc20Token) GetSymbol() types.CborString {
	return types.CborString(t.Symbol)
}

/*GetDecimal return decimal of erc20 token*/
func (t *Erc20Token) GetDecimal() typegen.CborInt {
	return typegen.CborInt(t.Decimals)
}

/*GetTotalSupply returns total number of tokens in existence*/
func (t *Erc20Token) GetTotalSupply() *big.Int {
	return t.TotalSupply
}

/*GetBalanceOf sender by ID.

* `args[0]` - the ID of user.*/
func (t *Erc20Token) GetBalanceOf(addr *address.Address) (*big.Int, error) {
	senderId, err := sdk.ResolveAddress(*addr)
	if err != nil {
		return nil, err
	}
	return t.getBalanceOf(senderId)
}

func (t *Erc20Token) getBalanceOf(act abi.ActorID) (*big.Int, error) {
	if balance, ok := t.Balances[actorToString(act)]; ok {
		return balance, nil
	}
	return &zero, nil
}

type TransferReq struct {
	ReceiverAddr   address.Address
	TransferAmount *big.Int
}

/*Transfer token from current caller to a specified address.

* `receiverAddr` - the ID of receiver.

* `transferAmount` - the transfer amount.
 */
func (t *Erc20Token) Transfer(transferReq *TransferReq) error {
	senderID, err := sdk.Caller()
	if err != nil {
		return err
	}
	receiverID, err := sdk.ResolveAddress(transferReq.ReceiverAddr)
	if err != nil {
		return err
	}

	if transferReq.TransferAmount.LessThanEqual(big.Zero()) {
		return errors.New("trasfer value must bigger than zero")
	}

	balanceOfSender, err := t.getBalanceOf(senderID)
	if err != nil {
		return err
	}
	balanceOfReceiver, err := t.getBalanceOf(receiverID)
	if err != nil {
		return err
	}

	if err := checkBalance(balanceOfSender, senderID); err != nil {
		return err
	}

	if err := isSmallerOrEqual(transferReq.TransferAmount, balanceOfSender); err != nil {
		return fmt.Errorf("transfer amount should be less than balance of sender (%v): %v", senderID, err)
	}

	t.Balances[actorToString(senderID)] = sub(balanceOfSender, transferReq.TransferAmount)
	t.Balances[actorToString(receiverID)] = add(balanceOfReceiver, transferReq.TransferAmount)
	logger.Logf("from %d to %d amount %s", senderID, receiverID, transferReq.TransferAmount.String())
	_ = sdk.SaveState(t)
	return nil
}

type AllowanceReq struct {
	OwnerAddr   address.Address
	SpenderAddr address.Address
}

/*GetAllowance checks the amount of tokens that an owner Allowed a spender to transfer in behalf of the owner to another receiver.

* `ownerAddr` - the ID of owner.

* `spenderAddr` - the ID of spender*/
func (t *Erc20Token) Allowance(req *AllowanceReq) (*big.Int, error) {
	ownerID, err := sdk.ResolveAddress(req.OwnerAddr)
	if err != nil {
		return nil, err
	}

	spenderId, err := sdk.ResolveAddress(req.SpenderAddr)
	if err != nil {
		return nil, err
	}

	return t.getAllowance(ownerID, spenderId)
}

func (t *Erc20Token) getAllowance(ownerID, spenderId abi.ActorID) (*big.Int, error) {
	if val, ok := t.Allowed[getAllowKey(ownerID, spenderId)]; ok {
		return val, nil
	}
	return &zero, nil
}

type TransferFromReq struct {
	OwnerAddr      address.Address
	ReceiverAddr   address.Address
	TransferAmount *big.Int
}

/*TransferFrom transfer tokens from token owner to receiver.

* `ownerAddr` - the ID of token owner.

* `receiverAddr` - the ID of receiver.

* `transferAmount` - the transfer amount.
 */
func (t *Erc20Token) TransferFrom(req *TransferFromReq) error {
	tokenOwnerID, err := sdk.ResolveAddress(req.OwnerAddr)
	if err != nil {
		return err
	}

	receiverID, err := sdk.ResolveAddress(req.ReceiverAddr)
	if err != nil {
		return err
	}

	if req.TransferAmount.LessThanEqual(big.Zero()) {
		return errors.New("send value must bigger than zero")
	}

	spenderID, err := sdk.Caller()
	if err != nil {
		return err
	}
	balanceOfTokenOwner, err := t.getBalanceOf(tokenOwnerID)
	if err != nil {
		return err
	}
	balanceOfReceiver, err := t.getBalanceOf(receiverID)
	if err != nil {
		return err
	}
	approvedAmount, err := t.getAllowance(tokenOwnerID, spenderID)
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

	if err := isSmallerOrEqual(req.TransferAmount, balanceOfTokenOwner); err != nil {
		return fmt.Errorf("transfer amount should be less than balance of token owner (%v): %v", tokenOwnerID, err)
	}
	if err := isSmallerOrEqual(req.TransferAmount, approvedAmount); err != nil {
		return fmt.Errorf("transfer amount should be less than approved spending amount of %v: %v", spenderID, err)
	}

	t.Balances[actorToString(tokenOwnerID)] = sub(balanceOfTokenOwner, req.TransferAmount)
	t.Balances[actorToString(receiverID)] = add(balanceOfReceiver, req.TransferAmount)
	t.Allowed[getAllowKey(tokenOwnerID, spenderID)] = sub(approvedAmount, req.TransferAmount)
	_ = sdk.SaveState(t)
	return nil
}

type ApprovalReq struct {
	SpenderAddr  address.Address
	NewAllowance *big.Int
}

/*Approval approves the passed-in identity to spend/burn a maximum amount of tokens on behalf of the function caller.
* `spenderAddr` - the ID of approved user.
* `newAllowance` - the maximum approved amount.*/
func (t *Erc20Token) Approval(req *ApprovalReq) error {
	spenderID, err := sdk.ResolveAddress(req.SpenderAddr)
	if err != nil {
		return err
	}

	if req.NewAllowance.LessThanEqual(big.Zero()) {
		return errors.New("allow value must bigger than zero")
	}

	callerID, err := sdk.Caller()
	if err != nil {
		return err
	}

	allowance, err := t.getAllowance(callerID, spenderID)
	if err != nil {
		return err
	}

	t.Allowed[getAllowKey(callerID, spenderID)] = add(allowance, req.NewAllowance)
	_ = sdk.SaveState(t)
	logger.Logf("approval %s for %s", getAllowKey(callerID, spenderID), req.NewAllowance.String())
	return nil
}

/*checkBalance checks if sender's balance is >= 0*/
func checkBalance(balance *big.Int, mspID abi.ActorID) error {
	if balance.LessThan(zero) {
		return fmt.Errorf("Balance of sender %v is %v", mspID, balance)
	}
	return nil
}

/*isSmallerOrEqual returns `nil` if a is <= b*/
func isSmallerOrEqual(a *big.Int, b *big.Int) error {
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

func sub(a, b *big.Int) *big.Int {
	return &big.Int{big.NewInt(0).Sub(a.Int, b.Int)}
}

func add(a, b *big.Int) *big.Int {
	return &big.Int{big.NewInt(0).Add(a.Int, b.Int)}
}
