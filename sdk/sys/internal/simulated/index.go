package simulated

import (
	"fmt"

	"github.com/filecoin-project/go-address"
	addr "github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/filecoin-project/go-state-types/proof"
	"github.com/filecoin-project/specs-actors/actors/runtime"
	gomock "github.com/golang/mock/gomock"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"
)

type FakeReporter struct {
}

func (f *FakeReporter) Errorf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}
func (f *FakeReporter) Fatalf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

// 执行 go  generate生成文件

//go:generate mockgen -destination ./mock_scheme.go -package=simulated -source ./index.go
type Fvm interface {
	Open(id cid.Cid) (*types.IpldOpen, error)
	SelfRoot() (cid.Cid, error)
	SelfSetRoot(id cid.Cid) error
	SelfCurrentBalance() (*types.TokenAmount, error)
	SelfDestruct(addr addr.Address) error
	Send(to address.Address, method uint64, params uint32, value types.TokenAmount) (*types.Send, error)
	GetChainRandomness(dst int64, round int64, entropy []byte) (abi.Randomness, error)
	GetBeaconRandomness(dst int64, round int64, entropy []byte) (abi.Randomness, error)
	BaseFee() (*types.TokenAmount, error)
	TotalFilCircSupply() (*types.TokenAmount, error)
	Create(codec uint64, data []byte) (uint32, error)
	Read(id uint32, offset, size uint32) ([]byte, uint32, error)
	Stat(id uint32) (*types.IpldStat, error)
	BlockLink(id uint32, hashFun uint64, hashLen uint32, cidBuf []byte) (uint32, error)
	Charge(name string, compute uint64) error
	Enabled() (bool, error)
	Log(msg string) error
	VerifySignature(
		signature *crypto.Signature,
		signer *address.Address,
		plaintext []byte,
	) (bool, error)
	HashBlake2b(data []byte) ([32]byte, error)
	ComputeUnsealedSectorCid(
		proofType abi.RegisteredSealProof,
		pieces []abi.PieceInfo,
	) (cid.Cid, error)
	VerifySeal(info *proof.SealVerifyInfo) (bool, error)
	VerifyPost(info *proof.WindowPoStVerifyInfo) (bool, error)
	VerifyConsensusFault(
		h1 []byte,
		h2 []byte,
		extra []byte,
	) (*runtime.ConsensusFault, error)
	VerifyAggregateSeals(info *types.AggregateSealVerifyProofAndInfos) (bool, error)
	VerifyReplicaUpdate(info *types.ReplicaUpdateInfo) (bool, error)
	BatchVerifySeals(sealVerifyInfos []proof.SealVerifyInfo) ([]bool, error)
	VMContext() (*types.InvocationContext, error)
	ResolveAddress(addr address.Address) (abi.ActorID, error)
	GetActorCodeCid(addr address.Address) (*cid.Cid, error)
	ResolveBuiltinActorType(codeCid cid.Cid) (types.ActorType, error)
	GetCodeCidForType(actorT types.ActorType) (cid.Cid, error)
	NewActorAddress() (address.Address, error)
	CreateActor(actorID abi.ActorID, codeCid cid.Cid) error
	Abort(code uint32, msg string)
}

var MockFvmInstance *MockFvm
var MockFvmInstanceCtl *gomock.Controller

func EpochFinish() {
	MockFvmInstanceCtl.Finish()
}
func init() {
	t := FakeReporter{}
	MockFvmInstanceCtl = gomock.NewController(&t)
	// defer ctl.Finish()
	MockFvmInstance = NewMockFvm(MockFvmInstanceCtl)

}
