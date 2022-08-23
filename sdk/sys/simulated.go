//go:build simulate

package sys

import (
	"github.com/filecoin-project/go-address"
	addr "github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/filecoin-project/go-state-types/proof"
	"github.com/golang/mock/gomock"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys/internal/simulated"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"
)

type ExpectOptions struct {
	// Do          func(f interface{})
	// DoAndReturn func(f interface{})
	MaxTimes *int
	MinTimes *int
	AnyTimes bool
	Times    *int
}

var (
	MatchAny                = gomock.Any
	MatchEq                 = gomock.Eq
	MatchNil                = gomock.Nil
	MatchLen                = gomock.Len
	MatchNot                = gomock.Not
	MatchAssignableToTypeOf = gomock.AssignableToTypeOf
	MatchInAnyOrder         = gomock.InAnyOrder
)

func initcall(call *gomock.Call, op *ExpectOptions) {
	if op == nil {
		return
	}
	// if op.Do != nil {
	// 	call = call.Do(op.Do)
	// }
	// if op.DoAndReturn != nil {
	// 	call = call.DoAndReturn(op.DoAndReturn)
	// }

	if op.MaxTimes != nil {
		call = call.MaxTimes(*op.MaxTimes)

	}

	if op.MinTimes != nil {
		call = call.MinTimes(*op.MinTimes)

	}

	if op.AnyTimes == true {
		call = call.AnyTimes()
	}

	if op.Times != nil {
		call = call.Times(*op.Times)
	}

}

// func CreateSimulated() (*simulated.MockSimulated, *gomock.Controller) {
// 	return simulated.CreateSimulated()
// }

func SetOpenExpect(id cid.Cid, out []interface{}, op *ExpectOptions) {
	call := simulated.SimulatedInstance.EXPECT().Open(id).Return(out...)
	initcall(call, op)
}
func SetSelfRootExpect(out []interface{}, op *ExpectOptions) {
	call := simulated.SimulatedInstance.EXPECT().SelfRoot().Return(out...)
	initcall(call, op)
}
func SetSelfSetRootExpect(id cid.Cid, out []interface{}, op *ExpectOptions) {
	call := simulated.SimulatedInstance.EXPECT().SelfSetRoot(id).Return(out...)
	initcall(call, op)
}
func SetSelfCurrentBalanceExpect(out []interface{}, op *ExpectOptions) {
	call := simulated.SimulatedInstance.EXPECT().SelfCurrentBalance().Return(out...)
	initcall(call, op)
}
func SetSelfDestructExpect(addr addr.Address, out []interface{}, op *ExpectOptions) {
	call := simulated.SimulatedInstance.EXPECT().SelfDestruct(addr).Return(out...)
	initcall(call, op)
}

func SetSendExpect(to address.Address, method uint64, params uint32, value types.TokenAmount, out []interface{}, op *ExpectOptions) {
	call := simulated.SimulatedInstance.EXPECT().Send(to, method, params, value).Return(out...)
	initcall(call, op)
}

func SetGetChainRandomnessExpect(dst int64, round int64, entropy []byte, out []interface{}, op *ExpectOptions) {
	call := simulated.SimulatedInstance.EXPECT().GetChainRandomness(dst, round, entropy).Return(out...)
	initcall(call, op)
}

func SetGetBeaconRandomnessExpect(dst int64, round int64, entropy []byte, out []interface{}, op *ExpectOptions) {
	call := simulated.SimulatedInstance.EXPECT().GetBeaconRandomness(dst, round, entropy).Return(out...)
	initcall(call, op)
}
func SetBaseFeeExpect(out []interface{}, op *ExpectOptions) {
	call := simulated.SimulatedInstance.EXPECT().BaseFee().Return(out...)
	initcall(call, op)
}
func SetTotalFilCircSupplyExpect(out []interface{}, op *ExpectOptions) {
	call := simulated.SimulatedInstance.EXPECT().TotalFilCircSupply().Return(out...)
	initcall(call, op)
}

func SetCreateExpect(codec uint64, data []byte, out []interface{}, op *ExpectOptions) {
	call := simulated.SimulatedInstance.EXPECT().Create(codec, data).Return(out...)
	initcall(call, op)
}
func SetReadExpect(id uint32, offset uint32, buf []byte, out []interface{}, op *ExpectOptions) {
	call := simulated.SimulatedInstance.EXPECT().Read(id, offset, buf).Return(out...)
	initcall(call, op)
}

func SetStatExpect(id uint32, out []interface{}, op *ExpectOptions) {
	call := simulated.SimulatedInstance.EXPECT().Stat(id).Return(out...)
	initcall(call, op)
}
func SetBlockLinkExpect(id uint32, hashFun uint64, hashLen uint32, cidBuf []byte, out []interface{}, op *ExpectOptions) {
	call := simulated.SimulatedInstance.EXPECT().BlockLink(id, hashFun, hashLen, cidBuf).Return(out...)
	initcall(call, op)
}

func SetChargeExpect(name string, compute uint64, out []interface{}, op *ExpectOptions) {
	call := simulated.SimulatedInstance.EXPECT().Charge(name, compute).Return(out...)
	initcall(call, op)
}

func SetEnabledExpect(out []interface{}, op *ExpectOptions) {
	call := simulated.SimulatedInstance.EXPECT().Enabled().Return(out...)
	initcall(call, op)
}
func SetLogExpect(msg string, out []interface{}, op *ExpectOptions) {
	call := simulated.SimulatedInstance.EXPECT().Log(msg).Return(out...)
	initcall(call, op)
}

func SetVerifySignatureExpect(
	signature *crypto.Signature,
	signer *address.Address,
	plaintext []byte, out []interface{}, op *ExpectOptions) {
	call := simulated.SimulatedInstance.EXPECT().VerifySignature(signature, signer, plaintext).Return(out...)
	initcall(call, op)
}

func SetHashBlake2bExpect(data []byte, out []interface{}, op *ExpectOptions) {

	call := simulated.SimulatedInstance.EXPECT().HashBlake2b(data).Return(out...)
	initcall(call, op)
}
func SetComputeUnsealedSectorCidExpect(
	proofType abi.RegisteredSealProof,
	pieces []abi.PieceInfo, out []interface{}, op *ExpectOptions) {
	call := simulated.SimulatedInstance.EXPECT().ComputeUnsealedSectorCid(proofType, pieces).Return(out...)
	initcall(call, op)
}

func SetVerifySealExpect(info *proof.SealVerifyInfo, out []interface{}, op *ExpectOptions) {
	call := simulated.SimulatedInstance.EXPECT().VerifySeal(info).Return(out...)
	initcall(call, op)
}

func SetVerifyPostExpect(info *proof.WindowPoStVerifyInfo, out []interface{}, op *ExpectOptions) {
	call := simulated.SimulatedInstance.EXPECT().VerifyPost(info).Return(out...)
	initcall(call, op)
}
func SetVerifyConsensusFaultExpect(
	h1 []byte,
	h2 []byte,
	extra []byte, out []interface{}, op *ExpectOptions) {
	call := simulated.SimulatedInstance.EXPECT().VerifyConsensusFault(h1, h2, extra).Return(out...)
	initcall(call, op)
}

func SetVerifyAggregateSealsExpect(info *types.AggregateSealVerifyProofAndInfos, out []interface{}, op *ExpectOptions) {
	call := simulated.SimulatedInstance.EXPECT().VerifyAggregateSeals(info).Return(out...)
	initcall(call, op)
}

func SetVerifyReplicaUpdateExpect(info *types.ReplicaUpdateInfo, out []interface{}, op *ExpectOptions) {
	call := simulated.SimulatedInstance.EXPECT().VerifyReplicaUpdate(info).Return(out...)
	initcall(call, op)
}

func SetBatchVerifySealsExpect(sealVerifyInfos []proof.SealVerifyInfo, out []interface{}, op *ExpectOptions) {
	call := simulated.SimulatedInstance.EXPECT().BatchVerifySeals(sealVerifyInfos).Return(out...)
	initcall(call, op)
}

func SetVMContextExpect(out []interface{}, op *ExpectOptions) {

	call := simulated.SimulatedInstance.EXPECT().VMContext().Return(out...)
	initcall(call, op)
}

func SetResolveAddressExpect(addr address.Address, out []interface{}, op *ExpectOptions) {
	call := simulated.SimulatedInstance.EXPECT().ResolveAddress(addr).Return(out...)
	initcall(call, op)
}

func SetGetActorCodeCidExpect(addr address.Address, out []interface{}, op *ExpectOptions) {
	call := simulated.SimulatedInstance.EXPECT().GetActorCodeCid(addr).Return(out...)
	initcall(call, op)
}

func SetResolveBuiltinActorTypeExpect(codeCid cid.Cid, out []interface{}, op *ExpectOptions) {
	call := simulated.SimulatedInstance.EXPECT().ResolveBuiltinActorType(codeCid).Return(out...)
	initcall(call, op)
}

func SetGetCodeCidForTypeExpect(actorT types.ActorType, out []interface{}, op *ExpectOptions) {
	call := simulated.SimulatedInstance.EXPECT().GetCodeCidForType(actorT).Return(out...)
	initcall(call, op)
}

func SetNewActorAddressExpect(out []interface{}, op *ExpectOptions) {
	call := simulated.SimulatedInstance.EXPECT().NewActorAddress().Return(out...)
	initcall(call, op)
}

func SetCreateActorExpect(actorID abi.ActorID, codeCid cid.Cid, out []interface{}, op *ExpectOptions) {
	call := simulated.SimulatedInstance.EXPECT().CreateActor(actorID, codeCid).Return(out...)
	initcall(call, op)
}

func SetAbortExpect(code uint32, msg string, op *ExpectOptions) {
	call := simulated.SimulatedInstance.EXPECT().Abort(code, msg)
	initcall(call, op)
}

func Begin() {
	simulated.Begin()
}

func End() {
	simulated.End()
}

func GetSimulated() *simulated.MockSimulated {
	return simulated.SimulatedInstance
}
