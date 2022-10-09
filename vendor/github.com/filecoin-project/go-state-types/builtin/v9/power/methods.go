package power

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/proof"
)

var Methods = []interface{}{
	1: *new(func(interface{}, *abi.EmptyValue) *abi.EmptyValue),           // Constructor
	2: *new(func(interface{}, *CreateMinerParams) *CreateMinerReturn),     // CreateMiner
	3: *new(func(interface{}, *UpdateClaimedPowerParams) *abi.EmptyValue), // UpdateClaimedPower
	4: *new(func(interface{}, *EnrollCronEventParams) *abi.EmptyValue),    // EnrollCronEvent
	5: *new(func(interface{}, *abi.EmptyValue) *abi.EmptyValue),           // CronTick
	6: *new(func(interface{}, *abi.TokenAmount) *abi.EmptyValue),          // UpdatePledgeTotal
	7: nil,
	8: *new(func(interface{}, *proof.SealVerifyInfo) *abi.EmptyValue),    // SubmitPoRepForBulkVerify
	9: *new(func(interface{}, *abi.EmptyValue) *CurrentTotalPowerReturn), // CurrentTotalPower
}
