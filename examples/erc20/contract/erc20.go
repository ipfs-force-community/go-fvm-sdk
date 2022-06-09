package contract

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"

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

var zero = big.NewInt(0)

/*basic Token20*/
type Erc20Token struct {
	Name        string
	Symbol      string
	Decimals    uint8
	TotalSupply big.Int

	//todo cbor gen not support non-string key and map value
	Balances map[string]big.Int
	Allowed  map[string]big.Int //owner-spender
}

func Constructor(name string, symbol string, decimals uint8, totalSupply big.Int) *Erc20Token {

	return &Erc20Token{
		Name:        name,
		Symbol:      symbol,
		Decimals:    decimals,
		TotalSupply: totalSupply,
		Balances:    make(map[string]big.Int),
		Allowed:     make(map[string]big.Int),
	}
}

func LoadToken() *Erc20Token {
	root, err := sdk.Root()
	if err != nil {
		sdk.Abort(ferrors.USR_ILLEGAL_STATE, fmt.Sprintf("failed to get root: %v", err))
	}

	data, err := sdk.Get(root)
	if err != nil {
		sdk.Abort(ferrors.USR_ILLEGAL_STATE, fmt.Sprintf("failed to get data: %v", err))
	}
	st := new(Erc20Token)
	err = st.UnmarshalCBOR(bytes.NewReader(data))
	if err != nil {
		sdk.Abort(ferrors.USR_ILLEGAL_STATE, fmt.Sprintf("failed to get data: %v", err))
	}
	return st
}

/*GetDecimal return token Name of erc20 token*/
func (t *Erc20Token) GetName() string {
	return t.Name
}

/*GetDecimal return token Symbol of erc20 token*/
func (t *Erc20Token) GetSymbol() string {
	return t.Symbol
}

/*GetDecimal return decimal of erc20 token*/
func (t *Erc20Token) GetDecimal() uint8 {
	return t.Decimals
}

/*GetTotalSupply returns total number of tokens in existence*/
func (t *Erc20Token) GetTotalSupply() big.Int {
	return t.TotalSupply
}

/*GetBalanceOf sender by ID.

* `args[0]` - the ID of user.*/
func (t *Erc20Token) GetBalanceOf(addr address.Address) (big.Int, error) {
	senderId, err := sdk.ResolveAddress(addr)
	if err != nil {
		return big.Int{}, err
	}
	return t.getBalanceOf(senderId)
}

func (t *Erc20Token) getBalanceOf(act abi.ActorID) (big.Int, error) {
	if balance, ok := t.Balances[actorToString(act)]; ok {
		return balance, nil
	}
	return big.Int{}, fmt.Errorf("actor %s not exit", act)
}

/*Transfer token from current caller to a specified address.

* `receiverAddr` - the ID of receiver.

* `transferAmount` - the transfer amount.
 */
func (t *Erc20Token) Transfer(receiverAddr address.Address, transferAmount big.Int) error {
	senderID, err := sdk.Caller()
	if err != nil {
		return err
	}
	receiverID, err := sdk.ResolveAddress(receiverAddr)
	if err != nil {
		return err
	}

	if transferAmount.LessThanEqual(big.Zero()) {
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

	if err := isSmallerOrEqual(transferAmount, balanceOfSender); err != nil {
		return fmt.Errorf("transfer amount should be less than balance of sender (%v): %v", senderID, err)
	}

	t.Balances[actorToString(senderID)] = big.Sub(balanceOfSender, transferAmount)
	t.Balances[actorToString(receiverID)] = big.Add(balanceOfReceiver, transferAmount)
	return nil
}

/*GetAllowance checks the amount of tokens that an owner Allowed a spender to transfer in behalf of the owner to another receiver.

* `ownerAddr` - the ID of owner.

* `spenderAddr` - the ID of spender*/
func (t *Erc20Token) Allowance(ownerAddr, spenderAddr address.Address) (big.Int, error) {
	ownerID, err := sdk.ResolveAddress(ownerAddr)
	if err != nil {
		return big.Int{}, err
	}

	spenderId, err := sdk.ResolveAddress(spenderAddr)
	if err != nil {
		return big.Int{}, err
	}

	return t.getAllowance(ownerID, spenderId)
}

func (t *Erc20Token) getAllowance(ownerID, spenderId abi.ActorID) (big.Int, error) {
	if val, ok := t.Allowed[getAllowKey(ownerID, spenderId)]; ok {
		return val, nil
	}
	return big.Zero(), nil
}

/*TransferFrom transfer tokens from token owner to receiver.

* `ownerAddr` - the ID of token owner.

* `spenderAddr` - the ID of receiver.

* `transferAmount` - the transfer amount.
 */
func (t *Erc20Token) TransferFrom(ownerAddr, spenderAddr address.Address, transferAmount big.Int) error {
	tokenOwnerID, err := sdk.ResolveAddress(ownerAddr)
	if err != nil {
		return err
	}

	receiverID, err := sdk.ResolveAddress(spenderAddr)
	if err != nil {
		return err
	}

	if transferAmount.LessThanEqual(big.Zero()) {
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
		return fmt.Errorf("approved amount for %v-%v less than zero", ownerAddr, spenderAddr)
	}

	if err := isSmallerOrEqual(transferAmount, balanceOfTokenOwner); err != nil {
		return fmt.Errorf("transfer amount should be less than balance of token owner (%v): %v", tokenOwnerID, err)
	}
	if err := isSmallerOrEqual(transferAmount, approvedAmount); err != nil {
		return fmt.Errorf("transfer amount should be less than approved spending amount of %v: %v", spenderID, err)
	}

	t.Balances[actorToString(tokenOwnerID)] = big.Sub(balanceOfTokenOwner, transferAmount)
	t.Balances[actorToString(receiverID)] = big.Add(balanceOfReceiver, transferAmount)
	t.Allowed[getAllowKey(tokenOwnerID, spenderID)] = big.Sub(approvedAmount, transferAmount)
	return nil
}

/*Approval approves the passed-in identity to spend/burn a maximum amount of tokens on behalf of the function caller.
* `spenderAddr` - the ID of approved user.
* `newAllowance` - the maximum approved amount.*/
func (t *Erc20Token) Approval(spenderAddr address.Address, newAllowance big.Int) error {
	spenderID, err := sdk.ResolveAddress(spenderAddr)
	if err != nil {
		return err
	}

	if newAllowance.LessThanEqual(big.Zero()) {
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

	t.Allowed[getAllowKey(callerID, spenderID)] = big.Add(allowance, newAllowance)
	return nil
}

/*checkBalance checks if sender's balance is >= 0*/
func checkBalance(balance big.Int, mspID abi.ActorID) error {
	if balance.LessThan(zero) {
		return fmt.Errorf("Balance of sender %v is %v", mspID, balance)
	}
	return nil
}

/*isSmallerOrEqual returns `nil` if a is <= b*/
func isSmallerOrEqual(a big.Int, b big.Int) error {
	if a.GreaterThan(b) {
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
