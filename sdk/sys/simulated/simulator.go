package simulated

import (
	"fmt"
	"time"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/ferrors"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"
)

const SimulateDebug = true

func (fvmSimulator *FvmSimulator) GetActor(addr address.Address) (builtin.Actor, error) {
	fvmSimulator.actorLk.Lock()
	defer fvmSimulator.actorLk.Unlock()
	actorId, err := fvmSimulator.ResolveAddress(addr)
	if err != nil {
		return builtin.Actor{}, err
	}
	actor, ok := fvmSimulator.actorsMap[actorId]
	if !ok {
		return builtin.Actor{}, ferrors.NotFound
	}
	return actor, nil
}

func (fvmSimulator *FvmSimulator) GetBuiltinActorType(codeCid cid.Cid) (types.ActorType, error) {
	for k, v := range EmbeddedBuiltinActors {
		if v == codeCid {
			av, err := stringToactorType(k)
			return av, err
		}
	}
	return types.ActorType(0), ferrors.NotFound
}

func (fvmSimulator *FvmSimulator) GetCodeCidForType(actorT types.ActorType) (cid.Cid, error) {
	actstr, err := actorTypeTostring(actorT)
	if err != nil {
		return cid.Undef, err
	}
	return EmbeddedBuiltinActors[actstr], nil
}

func (fvmSimulator *FvmSimulator) Exit(code ferrors.ExitCode, data []byte, msg string) {
	panic(fmt.Sprintf("%d:%v %s", code, data, msg))
}

func (fvmSimulator *FvmSimulator) ExitWithId(code ferrors.ExitCode, blkId types.BlockID, msg string) {
	data := fvmSimulator.blocks[int(blkId)].data
	fvmSimulator.Exit(code, data, msg)
}

func (fvmSimulator *FvmSimulator) Enabled() (bool, error) {
	return true, nil
}

func (fvmSimulator *FvmSimulator) Log(msg string) error {
	fmt.Println(msg)
	return nil
}
func (fvmSimulator *FvmSimulator) StoreArtifact(name string, data []byte) error {
	fmt.Printf("%s %v\n", name, data)
	return nil
}

func (fvmSimulator *FvmSimulator) SetMessageContext(messageCtx *types.MessageContext) {
	fvmSimulator.messageCtx = messageCtx
}

func (fvmSimulator *FvmSimulator) VMMessageContext() (*types.MessageContext, error) {
	return fvmSimulator.messageCtx, nil
}

func (fvmSimulator *FvmSimulator) SetNetworkContext(networkContext *types.NetworkContext) {
	fvmSimulator.networkCtx = networkContext
}

func (fvmSimulator *FvmSimulator) NetworkContext() (*types.NetworkContext, error) {
	return fvmSimulator.networkCtx, nil
}

func (fvmSimulator *FvmSimulator) SetTotalFilCircSupply(amount abi.TokenAmount) {
	fvmSimulator.totalFilCircSupply = amount
}

func (fvmSimulator *FvmSimulator) SetTipsetCid(epoch abi.ChainEpoch, cid *cid.Cid) {
	fvmSimulator.tipsetCidLk.Lock()
	defer fvmSimulator.tipsetCidLk.Unlock()
	fvmSimulator.tipsetCids[epoch] = cid
}

func (fvmSimulator *FvmSimulator) TotalFilCircSupply() (abi.TokenAmount, error) {
	return fvmSimulator.totalFilCircSupply, nil
}

func (fvmSimulator *FvmSimulator) TipsetTimestamp() (uint64, error) {
	return uint64(time.Now().Unix()), nil
}

func (fvmSimulator *FvmSimulator) TipsetCid(epoch abi.ChainEpoch) (*cid.Cid, error) {
	fvmSimulator.tipsetCidLk.Lock()
	defer fvmSimulator.tipsetCidLk.Unlock()
	if v, ok := fvmSimulator.tipsetCids[epoch]; ok {
		return v, nil
	}
	return nil, ferrors.NotFound
}
