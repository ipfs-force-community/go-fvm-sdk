package sdk

import (
	"context"
	"fmt"
	"math"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
)

// sendCfg used to pass addition send params for send calls
type sendCfg struct {
	flags    types.SendFlags //default 0 means nothing, 1 means readonly
	gasLimit *uint64         //default 0 means no limit
}

// SendOption options for set send params
type SendOption func(cfg sendCfg)

// WithGasLimit used to set gas limit for send call
func WithGasLimit(gasLimit uint64) SendOption {
	return func(cfg sendCfg) {
		cfg.gasLimit = &gasLimit
	}
}

// WithReadonly used to set readonly mode for send call
func WithReadonly() SendOption {
	return func(cfg sendCfg) {
		cfg.flags = types.ReadonlyFlag
	}
}

// Send call another actor
func Send(ctx context.Context, to address.Address, method abi.MethodNum, params types.RawBytes, value abi.TokenAmount, opts ...SendOption) (*types.Receipt, error) {
	cfg := sendCfg{}
	for _, opt := range opts {
		opt(cfg)
	}

	var (
		paramsID uint32
		err      error
	)
	if len(params) > 0 {
		paramsID, err = sys.Create(ctx, types.DAGCBOR, params)
		if err != nil {
			return nil, fmt.Errorf("invalid params: %w", err)
		}
	} else {
		paramsID = types.NoDataBlockID
	}

	send, err := sys.Send(ctx, to, method, paramsID, value, toSysGasLimit(cfg.gasLimit), cfg.flags)
	if err != nil {
		return nil, err
	}

	var returnData types.RawBytes
	if send.ExitCode == ferrors.OK && send.ReturnID != types.NoDataBlockID {
		readBuf, read, err := sys.Read(ctx, send.ReturnID, 0, send.ReturnSize)
		if err != nil {
			return nil, fmt.Errorf("read return_data: %w", err)
		}

		if read != send.ReturnSize {
			return nil, fmt.Errorf("read size is not equal to stat-size %v-%v", read, send.ReturnSize)
		}

		returnData = readBuf
	}

	return &types.Receipt{
		ExitCode:   send.ExitCode,
		ReturnData: returnData,
		GasUsed:    0,
	}, nil
}

func toSysGasLimit(gas *uint64) uint64 {
	if gas == nil {
		return math.MaxUint64
	}
	return *gas
}
