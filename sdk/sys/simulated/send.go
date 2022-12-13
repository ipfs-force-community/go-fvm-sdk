//nolint:unparam
package simulated

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
)

type SendMock struct {
	To     address.Address
	Method abi.MethodNum
	Params []byte
	Value  big.Int
	Out    types.SendResult
}

func (fvmSimulator *FvmSimulator) Send(to address.Address, method abi.MethodNum, params uint32, value abi.TokenAmount) (*types.SendResult, error) {
	return fvmSimulator.sendMatch(to, method, params, value)
}

func (fvmSimulator *FvmSimulator) ExpectSend(mock ...SendMock) {
	fvmSimulator.sendList = append(fvmSimulator.sendList, mock...)

}
