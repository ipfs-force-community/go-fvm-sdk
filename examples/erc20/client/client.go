package client

import (
	bytes "bytes"
	context "context"
	contract "erc20/contract"
	fmt "fmt"

	address "github.com/filecoin-project/go-address"
	abi "github.com/filecoin-project/go-state-types/abi"
	big "github.com/filecoin-project/go-state-types/big"
	builtin "github.com/filecoin-project/go-state-types/builtin"
	init8 "github.com/filecoin-project/specs-actors/v8/actors/builtin/init"
	actors "github.com/filecoin-project/venus/venus-shared/actors"
	types "github.com/filecoin-project/venus/venus-shared/types"
	types2 "github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	cid "github.com/ipfs/go-cid"
)

type FullNode interface {
	MpoolPushMessage(ctx context.Context, msg *types.Message, spec *types.MessageSendSpec) (*types.SignedMessage, error)
	StateWaitMsg(ctx context.Context, cid cid.Cid, confidence uint64) (*types.MsgLookup, error)
}

type IErc20TokenClient interface {
	Install(context.Context, []byte) (*init8.InstallReturn, error)
	CreateActor(context.Context, cid.Cid, []byte) (*init8.ExecReturn, error)

	Constructor(context.Context, *contract.ConstructorReq) error

	GetName(context.Context) (types2.CborString, error)

	GetSymbol(context.Context) (types2.CborString, error)

	GetDecimal(context.Context) (typegen.CborInt, error)

	GetTotalSupply(context.Context) (*big.Int, error)

	GetBalanceOf(context.Context, *address.Address) (*big.Int, error)

	Transfer(context.Context, *contract.TransferReq) error

	TransferFrom(context.Context, *contract.TransferFromReq) error

	Approval(context.Context, *contract.ApprovalReq) error
}

var _ IErc20TokenClient = (*Erc20TokenClient)(nil)

type Erc20TokenClient struct {
	Node        FullNode
	FromAddress address.Address
	Actor       address.Address
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
		From:   c.FromAddress,
		Value:  big.Zero(),
		Method: 2,
		Params: params,
	}

	smsg, err := c.Node.MpoolPushMessage(ctx, msg, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to push message: %w", err)
	}

	wait, err := c.Node.StateWaitMsg(ctx, smsg.Cid(), 0)
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
		From:   c.FromAddress,
		Value:  big.Zero(),
		Method: 3,
		Params: params,
	}

	smsg, err := c.Node.MpoolPushMessage(ctx, msg, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to push message: %w", err)
	}

	wait, err := c.Node.StateWaitMsg(ctx, smsg.Cid(), 0)
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
	return &result, nil
}

func (c *Erc20TokenClient) Constructor(ctx context.Context, p0 *contract.ConstructorReq) error {
	if c.Actor == address.Undef {
		return fmt.Errorf("unset actor address for call")
	}

	buf := bytes.NewBufferString("")
	if err := p0.MarshalCBOR(buf); err != nil {
		return err
	}
	msg := &types.Message{
		To:     c.Actor,
		From:   c.FromAddress,
		Value:  big.Zero(),
		Method: abi.MethodNum(1),
		Params: buf.Bytes(),
	}

	smsg, err := c.Node.MpoolPushMessage(ctx, msg, nil)
	if err != nil {
		return fmt.Errorf("failed to push message: %w", err)
	}

	wait, err := c.Node.StateWaitMsg(ctx, smsg.Cid(), 0)
	if err != nil {
		return fmt.Errorf("error waiting for message: %w", err)
	}

	// check it executed successfully
	if wait.Receipt.ExitCode != 0 {
		return fmt.Errorf("actor execution failed")
	}
	return nil
}

func (c *Erc20TokenClient) GetName(ctx context.Context) (types2.CborString, error) {
	if c.Actor == address.Undef {
		return "", fmt.Errorf("unset actor address for call")
	}

	msg := &types.Message{
		To:     c.Actor,
		From:   c.FromAddress,
		Value:  big.Zero(),
		Method: abi.MethodNum(2),
		Params: nil,
	}

	smsg, err := c.Node.MpoolPushMessage(ctx, msg, nil)
	if err != nil {
		return "", fmt.Errorf("failed to push message: %w", err)
	}

	wait, err := c.Node.StateWaitMsg(ctx, smsg.Cid(), 0)
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

	result := new(types2.CborString)
	result.UnmarshalCBOR(bytes.NewReader(wait.Receipt.Return))

	return *result, nil

}

func (c *Erc20TokenClient) GetSymbol(ctx context.Context) (types2.CborString, error) {
	if c.Actor == address.Undef {
		return "", fmt.Errorf("unset actor address for call")
	}

	msg := &types.Message{
		To:     c.Actor,
		From:   c.FromAddress,
		Value:  big.Zero(),
		Method: abi.MethodNum(3),
		Params: nil,
	}

	smsg, err := c.Node.MpoolPushMessage(ctx, msg, nil)
	if err != nil {
		return "", fmt.Errorf("failed to push message: %w", err)
	}

	wait, err := c.Node.StateWaitMsg(ctx, smsg.Cid(), 0)
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

	result := new(types2.CborString)
	result.UnmarshalCBOR(bytes.NewReader(wait.Receipt.Return))

	return *result, nil

}

func (c *Erc20TokenClient) GetDecimal(ctx context.Context) (typegen.CborInt, error) {
	if c.Actor == address.Undef {
		return 0, fmt.Errorf("unset actor address for call")
	}

	msg := &types.Message{
		To:     c.Actor,
		From:   c.FromAddress,
		Value:  big.Zero(),
		Method: abi.MethodNum(4),
		Params: nil,
	}

	smsg, err := c.Node.MpoolPushMessage(ctx, msg, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to push message: %w", err)
	}

	wait, err := c.Node.StateWaitMsg(ctx, smsg.Cid(), 0)
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
	if c.Actor == address.Undef {
		return nil, fmt.Errorf("unset actor address for call")
	}

	msg := &types.Message{
		To:     c.Actor,
		From:   c.FromAddress,
		Value:  big.Zero(),
		Method: abi.MethodNum(5),
		Params: nil,
	}

	smsg, err := c.Node.MpoolPushMessage(ctx, msg, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to push message: %w", err)
	}

	wait, err := c.Node.StateWaitMsg(ctx, smsg.Cid(), 0)
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
	if c.Actor == address.Undef {
		return nil, fmt.Errorf("unset actor address for call")
	}

	buf := bytes.NewBufferString("")
	if err := p0.MarshalCBOR(buf); err != nil {
		return nil, err
	}
	msg := &types.Message{
		To:     c.Actor,
		From:   c.FromAddress,
		Value:  big.Zero(),
		Method: abi.MethodNum(6),
		Params: buf.Bytes(),
	}

	smsg, err := c.Node.MpoolPushMessage(ctx, msg, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to push message: %w", err)
	}

	wait, err := c.Node.StateWaitMsg(ctx, smsg.Cid(), 0)
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
	if c.Actor == address.Undef {
		return fmt.Errorf("unset actor address for call")
	}

	buf := bytes.NewBufferString("")
	if err := p0.MarshalCBOR(buf); err != nil {
		return err
	}
	msg := &types.Message{
		To:     c.Actor,
		From:   c.FromAddress,
		Value:  big.Zero(),
		Method: abi.MethodNum(7),
		Params: buf.Bytes(),
	}

	smsg, err := c.Node.MpoolPushMessage(ctx, msg, nil)
	if err != nil {
		return fmt.Errorf("failed to push message: %w", err)
	}

	wait, err := c.Node.StateWaitMsg(ctx, smsg.Cid(), 0)
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
	if c.Actor == address.Undef {
		return fmt.Errorf("unset actor address for call")
	}

	buf := bytes.NewBufferString("")
	if err := p0.MarshalCBOR(buf); err != nil {
		return err
	}
	msg := &types.Message{
		To:     c.Actor,
		From:   c.FromAddress,
		Value:  big.Zero(),
		Method: abi.MethodNum(8),
		Params: buf.Bytes(),
	}

	smsg, err := c.Node.MpoolPushMessage(ctx, msg, nil)
	if err != nil {
		return fmt.Errorf("failed to push message: %w", err)
	}

	wait, err := c.Node.StateWaitMsg(ctx, smsg.Cid(), 0)
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
	if c.Actor == address.Undef {
		return fmt.Errorf("unset actor address for call")
	}

	buf := bytes.NewBufferString("")
	if err := p0.MarshalCBOR(buf); err != nil {
		return err
	}
	msg := &types.Message{
		To:     c.Actor,
		From:   c.FromAddress,
		Value:  big.Zero(),
		Method: abi.MethodNum(9),
		Params: buf.Bytes(),
	}

	smsg, err := c.Node.MpoolPushMessage(ctx, msg, nil)
	if err != nil {
		return fmt.Errorf("failed to push message: %w", err)
	}

	wait, err := c.Node.StateWaitMsg(ctx, smsg.Cid(), 0)
	if err != nil {
		return fmt.Errorf("error waiting for message: %w", err)
	}

	// check it executed successfully
	if wait.Receipt.ExitCode != 0 {
		return fmt.Errorf("actor execution failed")
	}
	return nil
}
