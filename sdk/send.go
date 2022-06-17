package sdk

import (
	"fmt"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
)

func Send(to address.Address, method abi.MethodNum, params types.RawBytes, value types.TokenAmount) (*types.Receipt, error) {
	var (
		paramsID uint32
		err      error
	)
	if len(params) > 0 {
		paramsID, err = sys.Create(types.DAG_CBOR, params)
		if err != nil {
			return nil, fmt.Errorf("invalid params: %w", err)
		}
	} else {
		paramsID = types.NO_DATA_BLOCK_ID
	}

	send, err := sys.Send(to, uint64(method), paramsID, value)
	if err != nil {
		return nil, err
	}

	var returnData types.RawBytes
	var exitCode = ferrors.ExitCode(send.ExitCode)
	if exitCode == ferrors.OK && send.ReturnID != types.NO_DATA_BLOCK_ID {
		ipldStat, err := sys.Stat(send.ReturnID)
		if err != nil {
			return nil, fmt.Errorf("return id ipld-stat: %w", err)
		}

		// Now read the return data.
		readBuf := make([]byte, ipldStat.Size)
		read, err := sys.Read(send.ReturnID, 0, readBuf)
		if err != nil {
			return nil, fmt.Errorf("read return_data: %w", err)
		}

		if read != ipldStat.Size {
			return nil, fmt.Errorf("read size is not equal to stat-size %v-%v", read, ipldStat.Size)
		}

		returnData = readBuf
	}

	return &types.Receipt{
		ExitCode:   send.ExitCode,
		ReturnData: returnData,
		GasUsed:    0,
	}, nil
}
