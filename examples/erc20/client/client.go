package client

import (
	"bytes"
	"context"
	contract "erc20/contract"
	"fmt"

	address "github.com/filecoin-project/go-address"
	abi "github.com/filecoin-project/go-state-types/abi"
	big "github.com/filecoin-project/go-state-types/big"
	builtin "github.com/filecoin-project/go-state-types/builtin"
	init8 "github.com/filecoin-project/specs-actors/v8/actors/builtin/init"
	actors "github.com/filecoin-project/venus/venus-shared/actors"
	types "github.com/filecoin-project/venus/venus-shared/types"
	sdkTypes "github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	cid "github.com/ipfs/go-cid"
	typegen "github.com/whyrusleeping/cbor-gen"

	v0 "github.com/filecoin-project/venus/venus-shared/api/chain/v0"
)

type FullNode interface {
	MpoolPushMessage(ctx context.Context, msg *types.Message, spec *types.MessageSendSpec) (*types.SignedMessage, error)
	StateWaitMsg(ctx context.Context, cid cid.Cid, confidence uint64) (*types.MsgLookup, error)
}

type IErc20TokenClient interface {
	Install(context.Context, []byte) (*init8.InstallReturn, error)
	CreateActor(context.Context, cid.Cid, []byte) (*init8.ExecReturn, error)

	GetName(context.Context) (sdkTypes.CborString, error)

	GetSymbol(context.Context) (sdkTypes.CborString, error)

	GetDecimal(context.Context) (typegen.CborInt, error)

	GetTotalSupply(context.Context) (*big.Int, error)

	GetBalanceOf(context.Context, *address.Address) (*big.Int, error)

	Transfer(context.Context, *contract.TransferReq) error

	TransferFrom(context.Context, *contract.TransferFromReq) error

	Approval(context.Context, *contract.ApprovalReq) error

	Allowance(context.Context, *contract.AllowanceReq) (*big.Int, error)

	FakeSetBalance(context.Context, *contract.FakeSetBalance) error
}

var _ IErc20TokenClient = (*Erc20TokenClient)(nil)

type Erc20TokenClient struct {
	node        v0.FullNode
	fromAddress address.Address
	actor       address.Address
	codeCid     cid.Cid
}

//Option option func
type Option func(opt ClientOption)

//ClientOption option for set client config
type ClientOption struct {
	fromAddress address.Address
	actor       address.Address
	codeCid     cid.Cid
}

//SetFromAddressOpt used to set from address who send actor messages
func SetFromAddressOpt(fromAddress address.Address) Option {
	return func(opt ClientOption) {
		opt.fromAddress = fromAddress
	}
}

//SetActorOpt used to set exit actoraddress
func SetActorOpt(actor address.Address) Option {
	return func(opt ClientOption) {
		opt.actor = actor
	}
}

//SetCodeCid used to set actor code cid
func SetCodeCid(codeCid cid.Cid) Option {
	return func(opt ClientOption) {
		opt.codeCid = codeCid
	}
}

func NewErc20TokenClient(fullNode v0.FullNode, opts ...Option) *Erc20TokenClient {
	cfg := ClientOption{}
	for _, opt := range opts {
		opt(cfg)
	}
	return &Erc20TokenClient{
		node:        fullNode,
		fromAddress: cfg.fromAddress,
		actor:       cfg.actor,
	}
}

func (c *Erc20TokenClient) CreateActor(ctx context.Context, codeCid cid.Cid, execParams []byte) (*init8.ExecReturn, error) {
	params, aErr := actors.SerializeParams(&init8.ExecParams{
		CodeCID:           codeCid,
		ConstructorParams: execParams,
	})
	if aErr != nil {
		return nil, fmt.Errorf("failed to serialize params: %w", aErr)
	}

	msg := &types.Message{
		To:     builtin.InitActorAddr,
		From:   c.fromAddress,
		Value:  big.Zero(),
		Method: 2,
		Params: params,
	}

	smsg, err := c.node.MpoolPushMessage(ctx, msg, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to push message: %w", err)
	}

	wait, err := c.node.StateWaitMsg(ctx, smsg.Cid(), 0)
	if err != nil {
		return nil, fmt.Errorf("error waiting for message: %w", err)
	}

	// check it executed successfully
	if wait.Receipt.ExitCode != 0 {
		return nil, fmt.Errorf("actor execution failed")
	}

	var result init8.ExecReturn
	r := bytes.NewReader(wait.Receipt.Return)
	if err := result.UnmarshalCBOR(r); err != nil {
		return nil, fmt.Errorf("error unmarshaling return value: %w", err)
	}
	c.actor = result.IDAddress
	return &result, nil
}

func (c *Erc20TokenClient) Install(ctx context.Context, code []byte) (*init8.InstallReturn, error) {
	params, aerr := actors.SerializeParams(&init8.InstallParams{
		Code: code,
	})
	if aerr != nil {
		return nil, fmt.Errorf("failed to serialize params: %w", aerr)
	}

	msg := &types.Message{
		To:     builtin.InitActorAddr,
		From:   c.fromAddress,
		Value:  big.Zero(),
		Method: 3,
		Params: params,
	}

	smsg, err := c.node.MpoolPushMessage(ctx, msg, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to push message: %w", err)
	}

	wait, err := c.node.StateWaitMsg(ctx, smsg.Cid(), 0)
	if err != nil {
		return nil, fmt.Errorf("error waiting for message: %w", err)
	}

	// check it executed successfully
	if wait.Receipt.ExitCode != 0 {
		return nil, fmt.Errorf("actor installation failed")
	}

	var result init8.InstallReturn
	r := bytes.NewReader(wait.Receipt.Return)
	if err := result.UnmarshalCBOR(r); err != nil {
		return nil, fmt.Errorf("error unmarshaling return value: %w", err)
	}
	c.codeCid = result.CodeCid
	return &result, nil
}

func (c *Erc20TokenClient) GetName(ctx context.Context) (sdkTypes.CborString, error) {
	if c.actor == address.Undef {
		return "", fmt.Errorf("unset actor address for call")
	}

	msg := &types.Message{
		To:     c.actor,
		From:   c.fromAddress,
		Value:  big.Zero(),
		Method: abi.MethodNum(2),
		Params: nil,
	}

	smsg, err := c.node.MpoolPushMessage(ctx, msg, nil)
	if err != nil {
		return "", fmt.Errorf("failed to push message: %w", err)
	}

	wait, err := c.node.StateWaitMsg(ctx, smsg.Cid(), 0)
	if err != nil {
		return "", fmt.Errorf("error waiting for message: %w", err)
	}

	// check it executed successfully
	if wait.Receipt.ExitCode != 0 {
		return "", fmt.Errorf("actor execution failed")
	}
	if len(wait.Receipt.Return) == 0 {
		return "", fmt.Errorf("expect get result for call")
	}

	result := new(sdkTypes.CborString)
	result.UnmarshalCBOR(bytes.NewReader(wait.Receipt.Return))

	return *result, nil

}

func (c *Erc20TokenClient) GetSymbol(ctx context.Context) (sdkTypes.CborString, error) {
	if c.actor == address.Undef {
		return "", fmt.Errorf("unset actor address for call")
	}

	msg := &types.Message{
		To:     c.actor,
		From:   c.fromAddress,
		Value:  big.Zero(),
		Method: abi.MethodNum(3),
		Params: nil,
	}

	smsg, err := c.node.MpoolPushMessage(ctx, msg, nil)
	if err != nil {
		return "", fmt.Errorf("failed to push message: %w", err)
	}

	wait, err := c.node.StateWaitMsg(ctx, smsg.Cid(), 0)
	if err != nil {
		return "", fmt.Errorf("error waiting for message: %w", err)
	}

	// check it executed successfully
	if wait.Receipt.ExitCode != 0 {
		return "", fmt.Errorf("actor execution failed")
	}
	if len(wait.Receipt.Return) == 0 {
		return "", fmt.Errorf("expect get result for call")
	}

	result := new(sdkTypes.CborString)
	result.UnmarshalCBOR(bytes.NewReader(wait.Receipt.Return))

	return *result, nil

}

func (c *Erc20TokenClient) GetDecimal(ctx context.Context) (typegen.CborInt, error) {
	if c.actor == address.Undef {
		return 0, fmt.Errorf("unset actor address for call")
	}

	msg := &types.Message{
		To:     c.actor,
		From:   c.fromAddress,
		Value:  big.Zero(),
		Method: abi.MethodNum(4),
		Params: nil,
	}

	smsg, err := c.node.MpoolPushMessage(ctx, msg, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to push message: %w", err)
	}

	wait, err := c.node.StateWaitMsg(ctx, smsg.Cid(), 0)
	if err != nil {
		return 0, fmt.Errorf("error waiting for message: %w", err)
	}

	// check it executed successfully
	if wait.Receipt.ExitCode != 0 {
		return 0, fmt.Errorf("actor execution failed")
	}
	if len(wait.Receipt.Return) == 0 {
		return 0, fmt.Errorf("expect get result for call")
	}

	result := new(typegen.CborInt)
	result.UnmarshalCBOR(bytes.NewReader(wait.Receipt.Return))

	return *result, nil

}

func (c *Erc20TokenClient) GetTotalSupply(ctx context.Context) (*big.Int, error) {
	if c.actor == address.Undef {
		return nil, fmt.Errorf("unset actor address for call")
	}

	msg := &types.Message{
		To:     c.actor,
		From:   c.fromAddress,
		Value:  big.Zero(),
		Method: abi.MethodNum(5),
		Params: nil,
	}

	smsg, err := c.node.MpoolPushMessage(ctx, msg, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to push message: %w", err)
	}

	wait, err := c.node.StateWaitMsg(ctx, smsg.Cid(), 0)
	if err != nil {
		return nil, fmt.Errorf("error waiting for message: %w", err)
	}

	// check it executed successfully
	if wait.Receipt.ExitCode != 0 {
		return nil, fmt.Errorf("actor execution failed")
	}
	if len(wait.Receipt.Return) == 0 {
		return nil, fmt.Errorf("expect get result for call")
	}

	result := new(big.Int)
	result.UnmarshalCBOR(bytes.NewReader(wait.Receipt.Return))

	return result, nil

}

func (c *Erc20TokenClient) GetBalanceOf(ctx context.Context, p0 *address.Address) (*big.Int, error) {
	if c.actor == address.Undef {
		return nil, fmt.Errorf("unset actor address for call")
	}

	buf := bytes.NewBufferString("")
	if err := p0.MarshalCBOR(buf); err != nil {
		return nil, err
	}
	msg := &types.Message{
		To:     c.actor,
		From:   c.fromAddress,
		Value:  big.Zero(),
		Method: abi.MethodNum(6),
		Params: buf.Bytes(),
	}

	smsg, err := c.node.MpoolPushMessage(ctx, msg, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to push message: %w", err)
	}

	wait, err := c.node.StateWaitMsg(ctx, smsg.Cid(), 0)
	if err != nil {
		return nil, fmt.Errorf("error waiting for message: %w", err)
	}

	// check it executed successfully
	if wait.Receipt.ExitCode != 0 {
		return nil, fmt.Errorf("actor execution failed")
	}
	if len(wait.Receipt.Return) == 0 {
		return nil, fmt.Errorf("expect get result for call")
	}

	result := new(big.Int)
	result.UnmarshalCBOR(bytes.NewReader(wait.Receipt.Return))

	return result, nil

}

func (c *Erc20TokenClient) Transfer(ctx context.Context, p0 *contract.TransferReq) error {
	if c.actor == address.Undef {
		return fmt.Errorf("unset actor address for call")
	}

	buf := bytes.NewBufferString("")
	if err := p0.MarshalCBOR(buf); err != nil {
		return err
	}
	msg := &types.Message{
		To:     c.actor,
		From:   c.fromAddress,
		Value:  big.Zero(),
		Method: abi.MethodNum(7),
		Params: buf.Bytes(),
	}

	smsg, err := c.node.MpoolPushMessage(ctx, msg, nil)
	if err != nil {
		return fmt.Errorf("failed to push message: %w", err)
	}

	wait, err := c.node.StateWaitMsg(ctx, smsg.Cid(), 0)
	if err != nil {
		return fmt.Errorf("error waiting for message: %w", err)
	}

	// check it executed successfully
	if wait.Receipt.ExitCode != 0 {
		return fmt.Errorf("actor execution failed")
	}
	return nil
}

func (c *Erc20TokenClient) TransferFrom(ctx context.Context, p0 *contract.TransferFromReq) error {
	if c.actor == address.Undef {
		return fmt.Errorf("unset actor address for call")
	}

	buf := bytes.NewBufferString("")
	if err := p0.MarshalCBOR(buf); err != nil {
		return err
	}
	msg := &types.Message{
		To:     c.actor,
		From:   c.fromAddress,
		Value:  big.Zero(),
		Method: abi.MethodNum(8),
		Params: buf.Bytes(),
	}

	smsg, err := c.node.MpoolPushMessage(ctx, msg, nil)
	if err != nil {
		return fmt.Errorf("failed to push message: %w", err)
	}

	wait, err := c.node.StateWaitMsg(ctx, smsg.Cid(), 0)
	if err != nil {
		return fmt.Errorf("error waiting for message: %w", err)
	}

	// check it executed successfully
	if wait.Receipt.ExitCode != 0 {
		return fmt.Errorf("actor execution failed")
	}
	return nil
}

func (c *Erc20TokenClient) Approval(ctx context.Context, p0 *contract.ApprovalReq) error {
	if c.actor == address.Undef {
		return fmt.Errorf("unset actor address for call")
	}

	buf := bytes.NewBufferString("")
	if err := p0.MarshalCBOR(buf); err != nil {
		return err
	}
	msg := &types.Message{
		To:     c.actor,
		From:   c.fromAddress,
		Value:  big.Zero(),
		Method: abi.MethodNum(9),
		Params: buf.Bytes(),
	}

	smsg, err := c.node.MpoolPushMessage(ctx, msg, nil)
	if err != nil {
		return fmt.Errorf("failed to push message: %w", err)
	}

	wait, err := c.node.StateWaitMsg(ctx, smsg.Cid(), 0)
	if err != nil {
		return fmt.Errorf("error waiting for message: %w", err)
	}

	// check it executed successfully
	if wait.Receipt.ExitCode != 0 {
		return fmt.Errorf("actor execution failed")
	}
	return nil
}

func (c *Erc20TokenClient) Allowance(ctx context.Context, p0 *contract.AllowanceReq) (*big.Int, error) {
	if c.actor == address.Undef {
		return nil, fmt.Errorf("unset actor address for call")
	}

	buf := bytes.NewBufferString("")
	if err := p0.MarshalCBOR(buf); err != nil {
		return nil, err
	}
	msg := &types.Message{
		To:     c.actor,
		From:   c.fromAddress,
		Value:  big.Zero(),
		Method: abi.MethodNum(10),
		Params: buf.Bytes(),
	}

	smsg, err := c.node.MpoolPushMessage(ctx, msg, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to push message: %w", err)
	}

	wait, err := c.node.StateWaitMsg(ctx, smsg.Cid(), 0)
	if err != nil {
		return nil, fmt.Errorf("error waiting for message: %w", err)
	}

	// check it executed successfully
	if wait.Receipt.ExitCode != 0 {
		return nil, fmt.Errorf("actor execution failed")
	}
	if len(wait.Receipt.Return) == 0 {
		return nil, fmt.Errorf("expect get result for call")
	}

	result := new(big.Int)
	result.UnmarshalCBOR(bytes.NewReader(wait.Receipt.Return))

	return result, nil

}

func (c *Erc20TokenClient) FakeSetBalance(ctx context.Context, p0 *contract.FakeSetBalance) error {
	if c.actor == address.Undef {
		return fmt.Errorf("unset actor address for call")
	}

	buf := bytes.NewBufferString("")
	if err := p0.MarshalCBOR(buf); err != nil {
		return err
	}
	msg := &types.Message{
		To:     c.actor,
		From:   c.fromAddress,
		Value:  big.Zero(),
		Method: abi.MethodNum(11),
		Params: buf.Bytes(),
	}

	smsg, err := c.node.MpoolPushMessage(ctx, msg, nil)
	if err != nil {
		return fmt.Errorf("failed to push message: %w", err)
	}

	wait, err := c.node.StateWaitMsg(ctx, smsg.Cid(), 0)
	if err != nil {
		return fmt.Errorf("error waiting for message: %w", err)
	}

	// check it executed successfully
	if wait.Receipt.ExitCode != 0 {
		return fmt.Errorf("actor execution failed")
	}
	return nil
}
