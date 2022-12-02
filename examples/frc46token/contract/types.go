package contract

import (
	"context"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
)

const FRC46TOKENTYPE ReceiverType = 0x85223bdf

type IFrc46Token interface {
	GetName(context.Context) types.CborString
	GetSymbol(context.Context) types.CborString
	GetGranularity(context.Context) types.CborUint
	GetTotalSupply(context.Context) *abi.TokenAmount
	BalanceOf(context.Context, *address.Address) (*abi.TokenAmount, error)
	Allowance(context.Context, *GetAllowanceParams) (*abi.TokenAmount, error)
	Transfer(context.Context, *TransferParams) (*TransferReturn, error)
	TransferFrom(context.Context, *TransferFromParams) (*TransferFromReturn, error)
	IncreaseAllowance(context.Context, *IncreaseAllowanceParams) (*abi.TokenAmount, error)
	DecreaseAllowance(context.Context, *DecreaseAllowanceParams) (*abi.TokenAmount, error)
	RevokeAllowance(context.Context, *RevokeAllowanceParams) (*abi.TokenAmount, error)
	Burn(context.Context, *BurnParams) (*BurnReturn, error)
	BurnFrom(context.Context, *BurnFromParams) (*BurnFromReturn, error)
}

type IFrc46Unspecific interface {
	Mint(context.Context, *MintParams) (*MintReturn, error)
}

type GetAllowanceParams struct {
	Owner    address.Address
	Operator address.Address
}

// / Return value after a successful mint.
// / The mint method is not standardised so this is merely a useful library-level type
// / and recommendation for token implementations.
type MintReturn struct {
	/// The new balance of the owner address
	Balance abi.TokenAmount
	/// The new total supply.
	Supply abi.TokenAmount
	/// (Optional) data returned from receiver hook
	RecipientData types.RawBytes
}

// / Intermediate data used by mint_return to construct the return data
type MintIntermediate struct {
	/// Recipient address to use for querying balance
	Recipient abi.ActorID
	/// (Optional) data returned from receiver hook
	RecipientData types.RawBytes
}

var _ IRecipientData = (*MintIntermediate)(nil)

func (m *MintIntermediate) SetRecipientData(bytes types.RawBytes) {
	m.RecipientData = bytes
}

// / Instruction to transfer tokens to another address
type TransferParams struct {
	To address.Address
	/// A non-negative amount to transfer
	Amount abi.TokenAmount
	/// Arbitrary data to pass on via the receiver hook
	OperatorData types.RawBytes
}

// / Return value after a successful transfer
type TransferReturn struct {
	/// The new balance of the `from` address
	FromBalance abi.TokenAmount
	/// The new balance of the `to` address
	ToBalance abi.TokenAmount
	/// (Optional) data returned from receiver hook
	RecipientData types.RawBytes
}

// / Intermediate data used by transfer_return to construct the return data
type TransferIntermediate struct {
	From abi.ActorID
	To   abi.ActorID
	/// (Optional) data returned from receiver hook
	RecipientData types.RawBytes
}

var _ IRecipientData = (*TransferIntermediate)(nil)

func (m *TransferIntermediate) SetRecipientData(bytes types.RawBytes) {
	m.RecipientData = bytes
}

// / Instruction to transfer tokens between two addresses as an operator
type TransferFromParams struct {
	From address.Address
	To   address.Address
	/// A non-negative amount to transfer
	Amount abi.TokenAmount
	/// Arbitrary data to pass on via the receiver hook
	OperatorData types.RawBytes
}

// / Return value after a successful delegated transfer
type TransferFromReturn struct {
	/// The new balance of the `from` address
	FromBalance abi.TokenAmount
	/// The new balance of the `to` address
	To_balance abi.TokenAmount
	/// The new remaining allowance between `owner` and `operator` (caller)
	Allowance abi.TokenAmount
	/// (Optional) data returned from receiver hook
	RecipientData types.RawBytes
}

// / Intermediate data used by transfer_from_return to construct the return data
type TransferFromIntermediate struct {
	Operator abi.ActorID
	From     abi.ActorID
	To       abi.ActorID
	/// (Optional) data returned from receiver hook
	RecipientData types.RawBytes
}

var _ IRecipientData = (*TransferFromIntermediate)(nil)

func (m *TransferFromIntermediate) SetRecipientData(bytes types.RawBytes) {
	m.RecipientData = bytes
}

// / Instruction to increase an allowance between two addresses
type IncreaseAllowanceParams struct {
	Operator address.Address
	/// A non-negative amount to increase the allowance by
	Increase abi.TokenAmount
}

// / Instruction to decrease an allowance between two addresses
type DecreaseAllowanceParams struct {
	Operator address.Address
	/// A non-negative amount to decrease the allowance by
	Decrease abi.TokenAmount
}

// / Instruction to revoke (set to 0) an allowance
type RevokeAllowanceParams struct {
	Operator address.Address
}

// / Instruction to burn an amount of tokens
type BurnParams struct {
	/// A non-negative amount to burn
	Amount abi.TokenAmount
}

// / The updated value after burning
type BurnReturn struct {
	/// New balance in the account after the successful burn
	Balance abi.TokenAmount
}

// / Instruction to burn an amount of tokens from another address
type BurnFromParams struct {
	Owner address.Address
	/// A non-negative amount to burn
	Amount abi.TokenAmount
}

// BurnFromReturn the updated value after a delegated burn
type BurnFromReturn struct {
	/// Balance new balance in the account after the successful burn
	Balance abi.TokenAmount
	/// Balance new remaining allowance between the owner and operator (caller)
	Allowance abi.TokenAmount
}

type ConstructorReq struct {
	Name        string
	Symbol      string
	Granularity uint64
	Supply      abi.TokenAmount
}

// / Receive parameters for an FRC46 token
type FRC46TokenReceived struct {
	/// The account that the tokens are being pulled from (the token actor address itself for mint)
	From abi.ActorID
	/// The account that the tokens are being sent to (the receiver address)
	To abi.ActorID
	/// Address of the operator that initiated the transfer/mint
	Operator abi.ActorID
	/// Amount of tokens being transferred/minted
	Amount abi.TokenAmount
	/// Data specified by the operator during transfer/mint
	OperatorData types.RawBytes
	/// Additional data specified by the token-actor during transfer/mint
	TokenData types.RawBytes
}

type MintParams struct {
	InitialOwner address.Address
	Amount       abi.TokenAmount
	OperatorData types.RawBytes
}
