package miner

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-bitfield"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin/v8/power"
)

var Methods = []interface{}{
	1:  *new(func(interface{}, *power.MinerConstructorParams) *abi.EmptyValue),   // Constructor
	2:  *new(func(interface{}, *abi.EmptyValue) *GetControlAddressesReturn),      // ControlAddresses
	3:  *new(func(interface{}, *ChangeWorkerAddressParams) *abi.EmptyValue),      // ChangeWorkerAddress
	4:  *new(func(interface{}, *ChangePeerIDParams) *abi.EmptyValue),             // ChangePeerID
	5:  *new(func(interface{}, *SubmitWindowedPoStParams) *abi.EmptyValue),       // SubmitWindowedPoSt
	6:  *new(func(interface{}, *PreCommitSectorParams) *abi.EmptyValue),          // PreCommitSector
	7:  *new(func(interface{}, *ProveCommitSectorParams) *abi.EmptyValue),        // ProveCommitSector
	8:  *new(func(interface{}, *ExtendSectorExpirationParams) *abi.EmptyValue),   // ExtendSectorExpiration
	9:  *new(func(interface{}, *TerminateSectorsParams) *TerminateSectorsReturn), // TerminateSectors
	10: *new(func(interface{}, *DeclareFaultsParams) *abi.EmptyValue),            // DeclareFaults
	11: *new(func(interface{}, *DeclareFaultsRecoveredParams) *abi.EmptyValue),   // DeclareFaultsRecovered
	12: *new(func(interface{}, *DeferredCronEventParams) *abi.EmptyValue),        // OnDeferredCronEvent
	13: *new(func(interface{}, *CheckSectorProvenParams) *abi.EmptyValue),        // CheckSectorProven
	14: *new(func(interface{}, *ApplyRewardParams) *abi.EmptyValue),              // ApplyRewards
	15: *new(func(interface{}, *ReportConsensusFaultParams) *abi.EmptyValue),     // ReportConsensusFault
	16: *new(func(interface{}, *WithdrawBalanceParams) *abi.TokenAmount),         // WithdrawBalance
	17: *new(func(interface{}, *ConfirmSectorProofsParams) *abi.EmptyValue),      // ConfirmSectorProofsValid
	18: *new(func(interface{}, *ChangeMultiaddrsParams) *abi.EmptyValue),         // ChangeMultiaddrs
	19: *new(func(interface{}, *CompactPartitionsParams) *abi.EmptyValue),        // CompactPartitions
	20: *new(func(interface{}, *CompactSectorNumbersParams) *abi.EmptyValue),     // CompactSectorNumbers
	21: *new(func(interface{}, *abi.EmptyValue) *abi.EmptyValue),                 // ConfirmUpdateWorkerKey
	22: *new(func(interface{}, *abi.EmptyValue) *abi.EmptyValue),                 // RepayDebt
	23: *new(func(interface{}, *address.Address) *abi.EmptyValue),                // ChangeOwnerAddress
	24: *new(func(interface{}, *DisputeWindowedPoStParams) *abi.EmptyValue),      // DisputeWindowedPoSt
	25: *new(func(interface{}, *PreCommitSectorBatchParams) *abi.EmptyValue),     // PreCommitSectorBatch
	26: *new(func(interface{}, *ProveCommitAggregateParams) *abi.EmptyValue),     // ProveCommitAggregate
	27: *new(func(interface{}, *ProveReplicaUpdatesParams) *bitfield.BitField),   // ProveReplicaUpdates
}
