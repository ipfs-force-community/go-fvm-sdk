package market

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
)

var Methods = []interface{}{
	1: *new(func(interface{}, *abi.EmptyValue) *abi.EmptyValue),                                 // Constructor
	2: *new(func(interface{}, *address.Address) *abi.EmptyValue),                                // AddBalance
	3: *new(func(interface{}, *WithdrawBalanceParams) *abi.TokenAmount),                         // WithdrawBalance
	4: *new(func(interface{}, *PublishStorageDealsParams) *PublishStorageDealsReturn),           // PublishStorageDeals
	5: *new(func(interface{}, *VerifyDealsForActivationParams) *VerifyDealsForActivationReturn), // VerifyDealsForActivation
	6: *new(func(interface{}, *ActivateDealsParams) *abi.EmptyValue),                            // ActivateDeals
	7: *new(func(interface{}, *OnMinerSectorsTerminateParams) *abi.EmptyValue),                  // OnMinerSectorsTerminate
	8: *new(func(interface{}, *ComputeDataCommitmentParams) *ComputeDataCommitmentReturn),       // ComputeDataCommitment
	9: *new(func(interface{}, *abi.EmptyValue) *abi.EmptyValue),                                 // CronTick
}
