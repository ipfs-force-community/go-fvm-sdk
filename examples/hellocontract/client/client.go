package client

import (
	bytes "bytes"
	context "context"
	fmt "fmt"

	address "github.com/filecoin-project/go-address"
	abi "github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
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

type IStateClient interface {
	Install(context.Context, []byte) (*init8.InstallReturn, error)
	CreateActor(context.Context, cid.Cid, []byte) (*init8.ExecReturn, error)

	Constructor(context.Context) error

	SayHello(context.Context) (types2.CBORBytes, error)
}

var _ IStateClient = (*StateClient)(nil)

type StateClient struct {
	Node        FullNode
	FromAddress address.Address
	Actor       address.Address
}

func (c *StateClient) CreateActor(ctx context.Context, codeCid cid.Cid, execParams []byte) (*init8.ExecReturn, error) {
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

func (c *StateClient) Install(ctx context.Context, code []byte) (*init8.InstallReturn, error) {
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

func (c *StateClient) Constructor(ctx context.Context) error {
	if c.Actor == address.Undef {
		return fmt.Errorf("unset actor address for call")
	}

	msg := &types.Message{
		To:     c.Actor,
		From:   c.FromAddress,
		Value:  big.Zero(),
		Method: abi.MethodNum(1),
		Params: nil,
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

func (c *StateClient) SayHello(ctx context.Context) (types2.CBORBytes, error) {
	if c.Actor == address.Undef {
		return nil, fmt.Errorf("unset actor address for call")
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

	result := new(types2.CBORBytes)
	result.UnmarshalCBOR(bytes.NewReader(wait.Receipt.Return))

	return *result, nil

}
