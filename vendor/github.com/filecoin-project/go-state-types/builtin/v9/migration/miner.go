package migration

import (
	"context"

	"github.com/filecoin-project/go-state-types/builtin/v8/market"

	"golang.org/x/xerrors"

	commp "github.com/filecoin-project/go-commp-utils/nonffi"
	"github.com/filecoin-project/go-state-types/builtin"
	miner8 "github.com/filecoin-project/go-state-types/builtin/v8/miner"
	"github.com/filecoin-project/go-state-types/builtin/v8/util/adt"
	miner9 "github.com/filecoin-project/go-state-types/builtin/v9/miner"

	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"

	"github.com/filecoin-project/go-state-types/abi"
)

type minerMigrator struct {
	proposals  *market.DealArray
	OutCodeCID cid.Cid
}

func (m minerMigrator) migratedCodeCID() cid.Cid {
	return m.OutCodeCID
}

func (m minerMigrator) migrateState(ctx context.Context, store cbor.IpldStore, in actorMigrationInput) (*actorMigrationResult, error) {
	var inState miner8.State
	if err := store.Get(ctx, in.head, &inState); err != nil {
		return nil, err
	}
	var inInfo miner8.MinerInfo
	if err := store.Get(ctx, inState.Info, &inInfo); err != nil {
		return nil, err
	}
	wrappedStore := adt.WrapStore(ctx, store)

	oldPrecommitOnChainInfos, err := adt.AsMap(wrappedStore, inState.PreCommittedSectors, builtin.DefaultHamtBitwidth)
	if err != nil {
		return nil, xerrors.Errorf("failed to load old precommit onchain infos for miner %s: %w", in.address, err)
	}

	emptyMap, err := adt.StoreEmptyMap(wrappedStore, builtin.DefaultHamtBitwidth)
	if err != nil {
		return nil, xerrors.Errorf("failed to make empty map: %w", err)
	}

	newPrecommitOnChainInfos, err := adt.AsMap(wrappedStore, emptyMap, builtin.DefaultHamtBitwidth)
	if err != nil {
		return nil, xerrors.Errorf("failed to load empty map: %w", err)
	}

	var info miner8.SectorPreCommitOnChainInfo
	err = oldPrecommitOnChainInfos.ForEach(&info, func(key string) error {
		var unsealedCid *cid.Cid
		if len(info.Info.DealIDs) != 0 {
			pieces := make([]abi.PieceInfo, len(info.Info.DealIDs))
			for i, dealID := range info.Info.DealIDs {
				deal, err := m.proposals.GetDealProposal(dealID)
				if err != nil {
					return xerrors.Errorf("error getting deal proposal: %w", err)
				}

				pieces[i] = abi.PieceInfo{
					PieceCID: deal.PieceCID,
					Size:     deal.PieceSize,
				}
			}

			commd, err := commp.GenerateUnsealedCID(info.Info.SealProof, pieces)
			if err != nil {
				return xerrors.Errorf("failed to generate unsealed CID: %w", err)
			}

			unsealedCid = &commd
		}

		err = newPrecommitOnChainInfos.Put(miner9.SectorKey(info.Info.SectorNumber), &miner9.SectorPreCommitOnChainInfo{
			Info: miner9.SectorPreCommitInfo{
				SealProof:     info.Info.SealProof,
				SectorNumber:  info.Info.SectorNumber,
				SealedCID:     info.Info.SealedCID,
				SealRandEpoch: info.Info.SealRandEpoch,
				DealIDs:       info.Info.DealIDs,
				Expiration:    info.Info.Expiration,
				UnsealedCid:   unsealedCid,
			},
			PreCommitDeposit: info.PreCommitDeposit,
			PreCommitEpoch:   info.PreCommitEpoch,
		})

		if err != nil {
			return xerrors.Errorf("failed to write new precommitinfo: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, xerrors.Errorf("failed to iterate over precommitinfos: %w", err)
	}

	newPrecommits, err := newPrecommitOnChainInfos.Root()
	if err != nil {
		return nil, xerrors.Errorf("failed to flush new precommits: %w", err)
	}

	var newPendingWorkerKey *miner9.WorkerKeyChange
	if inInfo.PendingWorkerKey != nil {
		newPendingWorkerKey = &miner9.WorkerKeyChange{
			NewWorker:   inInfo.PendingWorkerKey.NewWorker,
			EffectiveAt: inInfo.PendingWorkerKey.EffectiveAt,
		}
	}

	outInfo := miner9.MinerInfo{
		Owner:       inInfo.Owner,
		Worker:      inInfo.Worker,
		Beneficiary: inInfo.Owner,
		BeneficiaryTerm: miner9.BeneficiaryTerm{
			Quota:      abi.NewTokenAmount(0),
			UsedQuota:  abi.NewTokenAmount(0),
			Expiration: 0,
		},
		PendingBeneficiaryTerm:     nil,
		ControlAddresses:           inInfo.ControlAddresses,
		PendingWorkerKey:           newPendingWorkerKey,
		PeerId:                     inInfo.PeerId,
		Multiaddrs:                 inInfo.Multiaddrs,
		WindowPoStProofType:        inInfo.WindowPoStProofType,
		SectorSize:                 inInfo.SectorSize,
		WindowPoStPartitionSectors: inInfo.WindowPoStPartitionSectors,
		ConsensusFaultElapsed:      inInfo.ConsensusFaultElapsed,
		PendingOwnerAddress:        inInfo.PendingOwnerAddress,
	}
	newInfoCid, err := store.Put(ctx, &outInfo)
	if err != nil {
		return nil, xerrors.Errorf("failed to flush new miner info: %w", err)
	}

	outState := miner9.State{
		Info:                       newInfoCid,
		PreCommitDeposits:          inState.PreCommitDeposits,
		LockedFunds:                inState.LockedFunds,
		VestingFunds:               inState.VestingFunds,
		FeeDebt:                    inState.FeeDebt,
		InitialPledge:              inState.InitialPledge,
		PreCommittedSectors:        newPrecommits,
		PreCommittedSectorsCleanUp: inState.PreCommittedSectorsCleanUp,
		AllocatedSectors:           inState.AllocatedSectors,
		Sectors:                    inState.Sectors,
		ProvingPeriodStart:         inState.ProvingPeriodStart,
		CurrentDeadline:            inState.CurrentDeadline,
		Deadlines:                  inState.Deadlines,
		EarlyTerminations:          inState.EarlyTerminations,
		DeadlineCronActive:         inState.DeadlineCronActive,
	}

	newHead, err := store.Put(ctx, &outState)
	return &actorMigrationResult{
		newCodeCID: m.migratedCodeCID(),
		newHead:    newHead,
	}, err
}
