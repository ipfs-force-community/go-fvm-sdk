package contract

import (
	"bytes"
	"context"
	"fmt"

	cbor2 "github.com/filecoin-project/go-state-types/cbor"

	"github.com/filecoin-project/go-state-types/big"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk"

	"github.com/filecoin-project/go-state-types/abi"

	"github.com/filecoin-project/go-address"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
)

var RECEIVERHOOKMETHODNUM abi.MethodNum = 0xde180de3 //3726118371

// UniversalReceiverParams parameters for universal receiver
// Actual payload varies with asset type
// eg: FRC46_TOKEN_TYPE will come with a payload of FRC46TokenReceived
type UniversalReceiverParams struct {
	// Type_ Asset type
	Type_ ReceiverType
	// Payload corresponding to asset type
	Payload types.RawBytes
}

// IUniversalReceiver standard interface for an actor that wishes to receive FRC-0046 tokens or other assets
type IUniversalReceiver interface {
	// Receive Invoked by a token actor during pending transfer or mint to the receiver's address
	// Within this hook, the token actor has optimistically persisted the new balance so
	// the receiving actor can immediately utilise the received funds. If the receiver wishes to
	// reject the incoming transfer, this function should abort which will cause the token actor
	// to rollback the transaction.
	Receive(UniversalReceiverParams)
}

type ReceiverType = uint64

type IRecipientData interface {
	SetRecipientData(bytes types.RawBytes)
}

type ReceiverHook struct {
	ToAddr      address.Address
	TokenType   ReceiverType
	TokenParams types.RawBytes
	ResultData  IRecipientData
	called      bool
}

func NewReceiverHook(toAddr address.Address, tokenType ReceiverType, tokenParams cbor2.Marshaler, resultData IRecipientData) (*ReceiverHook, error) {
	buf := bytes.NewBuffer(nil)
	if err := tokenParams.MarshalCBOR(buf); err != nil {
		return nil, err
	}

	return &ReceiverHook{
		ToAddr:      toAddr,
		TokenType:   tokenType,
		TokenParams: buf.Bytes(),
		ResultData:  resultData,
		called:      false,
	}, nil
}

func (hook *ReceiverHook) Call(ctx context.Context) error {
	if hook.called {
		return fmt.Errorf("receiver hook was already called %w", ferrors.USR_ASSERTION_FAILED)
	}

	hook.called = true

	params := UniversalReceiverParams{
		Type_:   hook.TokenType,
		Payload: hook.TokenParams,
	}

	buf := bytes.NewBuffer(nil)
	err := params.MarshalCBOR(buf)
	if err != nil {
		return fmt.Errorf("error encoding to ipld %w", ferrors.USR_SERIALIZATION)
	}

	receipt, err := sdk.Send(ctx, hook.ToAddr, RECEIVERHOOKMETHODNUM, buf.Bytes(), big.Zero())
	if err != nil {
		return err
	}

	if receipt.ExitCode != ferrors.OK {
		return fmt.Errorf("receiver hook error to %s: exit_code=%w, method_num=%d, return_data=%v", hook.ToAddr, ferrors.ExitCode(receipt.ExitCode), RECEIVERHOOKMETHODNUM, receipt.ReturnData)
	}
	hook.ResultData.SetRecipientData(receipt.ReturnData)
	return nil
}
