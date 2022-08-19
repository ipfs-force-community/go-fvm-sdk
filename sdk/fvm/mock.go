package fvm

import (
	"github.com/filecoin-project/go-address"
	addr "github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/filecoin-project/go-state-types/proof"
	"github.com/golang/mock/gomock"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"
)

type ExpectOptions struct {
	Do          func(f interface{})
	DoAndReturn func(f interface{})
	MaxTimes    *int
	MinTimes    *int
	AnyTimes    *int
	Times       *int
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
	if op.Do != nil {
		call = call.Do(op.Do)
	}
	if op.DoAndReturn != nil {
		call = call.DoAndReturn(op.DoAndReturn)
	}

	if op.MaxTimes != nil {
		call = call.MaxTimes(*op.MaxTimes)

	}

	if op.MinTimes != nil {
		call = call.MinTimes(*op.MinTimes)

	}

	if op.AnyTimes != nil {
		call = call.AnyTimes()
	}

	if op.Times != nil {
		call = call.Times(*op.Times)
	}

}

func OpenExpect(in interface{}, out interface{}, op *ExpectOptions) {
	call := MockFvmInstance.EXPECT().Open(in).Return(out, nil)
	initcall(call, op)
}
func SelfRootExpect(in interface{}, out interface{}, op *ExpectOptions) {
	call := MockFvmInstance.EXPECT().SelfRoot(in).Return(out, nil)
	initcall(call, op)
}
func SelfSetRootExpect(in interface{}, out interface{}, op *ExpectOptions) {
	call := MockFvmInstance.EXPECT().SelfSetRoot(in).Return(out, gomock.Nil())
	initcall(call, op)
}
func SelfCurrentBalanceExpect(out interface{}, op *ExpectOptions) {
	call := MockFvmInstance.EXPECT().SelfCurrentBalance().Return(out, nil)
	initcall(call, op)
}
func SelfDestructExpect(addr addr.Address, out []interface{}, op *ExpectOptions) {
	call := MockFvmInstance.EXPECT().SelfDestruct(addr).Return(out...)
	initcall(call, op)
}

func SendExpect(to address.Address, method uint64, params uint32, value types.TokenAmount, out []interface{}, op *ExpectOptions) {
	call := MockFvmInstance.EXPECT().Send(to, method, params, value).Return(out...)
	initcall(call, op)
}

func GetChainRandomnessExpect(dst int64, round int64, entropy []byte, out []interface{}, op *ExpectOptions) {
	call := MockFvmInstance.EXPECT().GetChainRandomness(dst, round, entropy).Return(out...)
	initcall(call, op)
}

func GetBeaconRandomnessExpect(dst int64, round int64, entropy []byte, out []interface{}, op *ExpectOptions) {
	call := MockFvmInstance.EXPECT().GetBeaconRandomness(dst, round, entropy).Return(out...)
	initcall(call, op)
}
func BaseFeeExpect(out []interface{}, op *ExpectOptions) {
	call := MockFvmInstance.EXPECT().BaseFee().Return(out...)
	initcall(call, op)
}
func TotalFilCircSupplyExpect(out []interface{}, op *ExpectOptions) {
	call := MockFvmInstance.EXPECT().TotalFilCircSupply().Return(out...)
	initcall(call, op)
}

func CreateExpect(codec uint64, data []byte, out []interface{}, op *ExpectOptions) {
	call := MockFvmInstance.EXPECT().Create(codec, data).Return(out...)
	initcall(call, op)
}
func ReadExpect(id uint32, offset uint32, buf []byte, out []interface{}, op *ExpectOptions) {
	call := MockFvmInstance.EXPECT().Read(id, offset, buf).Return(out...)
	initcall(call, op)
}

func StatExpect(id uint32, out []interface{}, op *ExpectOptions) {
	call := MockFvmInstance.EXPECT().Stat(id).Return(out...)
	initcall(call, op)
}
func BlockLinkExpect(id uint32, hashFun uint64, hashLen uint32, cidBuf []byte, out []interface{}, op *ExpectOptions) {
	call := MockFvmInstance.EXPECT().BlockLink(id, hashFun, hashLen, cidBuf).Return(out...)
	initcall(call, op)
}

func ChargeExpect(name string, compute uint64, out []interface{}, op *ExpectOptions) {
	call := MockFvmInstance.EXPECT().Charge(name, compute).Return(out...)
	initcall(call, op)
}

func EnabledExpect(out []interface{}, op *ExpectOptions) {
	call := MockFvmInstance.EXPECT().Enabled().Return(out...)
	initcall(call, op)
}
func LogExpect(msg string, out []interface{}, op *ExpectOptions) {
	call := MockFvmInstance.EXPECT().Log(msg).Return(out...)
	initcall(call, op)
}

func VerifySignatureExpect(
	signature *crypto.Signature,
	signer *address.Address,
	plaintext []byte, out []interface{}, op *ExpectOptions) {
	call := MockFvmInstance.EXPECT().VerifySignature(signature, signer, plaintext).Return(out...)
	initcall(call, op)
}

func HashBlake2bExpect(data []byte, out []interface{}, op *ExpectOptions) {

	call := MockFvmInstance.EXPECT().HashBlake2b(data).Return(out...)
	initcall(call, op)
}
func ComputeUnsealedSectorCidExpect(
	proofType abi.RegisteredSealProof,
	pieces []abi.PieceInfo, out []interface{}, op *ExpectOptions) {
	call := MockFvmInstance.EXPECT().ComputeUnsealedSectorCid(proofType, pieces).Return(out...)
	initcall(call, op)
}

func VerifySealExpect(info *proof.SealVerifyInfo, out []interface{}, op *ExpectOptions) {
	call := MockFvmInstance.EXPECT().VerifySeal(info).Return(out...)
	initcall(call, op)
}

func VerifyPostExpect(info *proof.WindowPoStVerifyInfo, out []interface{}, op *ExpectOptions) {
	call := MockFvmInstance.EXPECT().VerifyPost(info).Return(out...)
	initcall(call, op)
}
func VerifyConsensusFaultExpect(
	h1 []byte,
	h2 []byte,
	extra []byte, out []interface{}, op *ExpectOptions) {
	call := MockFvmInstance.EXPECT().VerifyConsensusFault(h1, h2, extra).Return(out...)
	initcall(call, op)
}

func VerifyAggregateSealsExpect(info *types.AggregateSealVerifyProofAndInfos, out []interface{}, op *ExpectOptions) {
	call := MockFvmInstance.EXPECT().VerifyAggregateSeals(info).Return(out...)
	initcall(call, op)
}

func VerifyReplicaUpdateExpect(info *types.ReplicaUpdateInfo, out []interface{}, op *ExpectOptions) {
	call := MockFvmInstance.EXPECT().VerifyReplicaUpdate(info).Return(out...)
	initcall(call, op)
}

func BatchVerifySealsExpect(sealVerifyInfos []proof.SealVerifyInfo, out []interface{}, op *ExpectOptions) {
	call := MockFvmInstance.EXPECT().BatchVerifySeals(sealVerifyInfos).Return(out...)
	initcall(call, op)
}

func VMContextExpect(out []interface{}, op *ExpectOptions) {

	call := MockFvmInstance.EXPECT().VMContext().Return(out...)
	initcall(call, op)
}

func ResolveAddressExpect(addr address.Address, out []interface{}, op *ExpectOptions) {
	call := MockFvmInstance.EXPECT().ResolveAddress(addr).Return(out...)
	initcall(call, op)
}

func GetActorCodeCidExpect(addr address.Address, out []interface{}, op *ExpectOptions) {
	call := MockFvmInstance.EXPECT().GetActorCodeCid(addr).Return(out...)
	initcall(call, op)
}

func ResolveBuiltinActorTypeExpect(codeCid cid.Cid, out []interface{}, op *ExpectOptions) {
	call := MockFvmInstance.EXPECT().ResolveBuiltinActorType(codeCid).Return(out...)
	initcall(call, op)
}

func GetCodeCidForTypeExpect(actorT types.ActorType, out []interface{}, op *ExpectOptions) {
	call := MockFvmInstance.EXPECT().GetCodeCidForType(actorT).Return(out...)
	initcall(call, op)
}

func NewActorAddressExpect(out []interface{}, op *ExpectOptions) {
	call := MockFvmInstance.EXPECT().NewActorAddress().Return(out...)
	initcall(call, op)
}

func CreateActorExpect(actorID abi.ActorID, codeCid cid.Cid, out []interface{}, op *ExpectOptions) {
	call := MockFvmInstance.EXPECT().CreateActor(actorID, codeCid).Return(out...)
	initcall(call, op)
}

func AbortExpect(code uint32, msg string, op *ExpectOptions) {
	call := MockFvmInstance.EXPECT().Abort(code, msg)
	initcall(call, op)
}
