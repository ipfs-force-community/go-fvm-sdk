package types

import (
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
)

type RawBytes []byte

type Receipt struct {
	ExitCode   ferrors.ExitCode `json:"exit_code"`
	ReturnData RawBytes         `json:"return_data"`
	GasUsed    int64            `json:"gas_used"`
}
