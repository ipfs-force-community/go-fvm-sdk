package power

import (
	addr "github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/go-state-types/builtin/v8/util/adt"
	"github.com/filecoin-project/go-state-types/builtin/v8/util/smoothing"
	"github.com/ipfs/go-cid"
	"golang.org/x/xerrors"
)

// genesis power in bytes = 750,000 GiB
var InitialQAPowerEstimatePosition = big.Mul(big.NewInt(750_000), big.NewInt(1<<30))

// max chain throughput in bytes per epoch = 120 ProveCommits / epoch = 3,840 GiB
var InitialQAPowerEstimateVelocity = big.Mul(big.NewInt(3_840), big.NewInt(1<<30))

// Bitwidth of CronEventQueue HAMT determined empirically from mutation
// patterns and projections of mainnet data.
const CronQueueHamtBitwidth = 6

// Bitwidth of CronEventQueue AMT determined empirically from mutation
// patterns and projections of mainnet data.
const CronQueueAmtBitwidth = 6

// Bitwidth of ProofValidationBatch AMT determined empirically from mutation
// pattersn and projections of mainnet data.
const ProofValidationBatchAmtBitwidth = 4

// The number of miners that must meet the consensus minimum miner power before that minimum power is enforced
// as a condition of leader election.
// This ensures a network still functions before any miners reach that threshold.
const ConsensusMinerMinMiners = 4 // PARAM_SPEC

type State struct {
	TotalRawBytePower abi.StoragePower
	// TotalBytesCommitted includes claims from miners below min power threshold
	TotalBytesCommitted  abi.StoragePower
	TotalQualityAdjPower abi.StoragePower
	// TotalQABytesCommitted includes claims from miners below min power threshold
	TotalQABytesCommitted abi.StoragePower
	TotalPledgeCollateral abi.TokenAmount

	// These fields are set once per epoch in the previous cron tick and used
	// for consistent values across a single epoch's state transition.
	ThisEpochRawBytePower     abi.StoragePower
	ThisEpochQualityAdjPower  abi.StoragePower
	ThisEpochPledgeCollateral abi.TokenAmount
	ThisEpochQAPowerSmoothed  smoothing.FilterEstimate

	MinerCount int64
	// Number of miners having proven the minimum consensus power.
	MinerAboveMinPowerCount int64

	// A queue of events to be triggered by cron, indexed by epoch.
	CronEventQueue cid.Cid // Multimap, (HAMT[ChainEpoch]AMT[CronEvent])

	// First epoch in which a cron task may be stored.
	// Cron will iterate every epoch between this and the current epoch inclusively to find tasks to execute.
	FirstCronEpoch abi.ChainEpoch

	// Claimed power for each miner.
	Claims cid.Cid // Map, HAMT[address]Claim

	ProofValidationBatch *cid.Cid // Multimap, (HAMT[Address]AMT[SealVerifyInfo])
}

func ConstructState(store adt.Store) (*State, error) {
	emptyClaimsMapCid, err := adt.StoreEmptyMap(store, builtin.DefaultHamtBitwidth)
	if err != nil {
		return nil, xerrors.Errorf("failed to create empty map: %w", err)
	}
	emptyCronQueueMMapCid, err := adt.StoreEmptyMultimap(store, CronQueueHamtBitwidth, CronQueueAmtBitwidth)
	if err != nil {
		return nil, xerrors.Errorf("failed to create empty multimap: %w", err)
	}

	return &State{
		TotalRawBytePower:         abi.NewStoragePower(0),
		TotalBytesCommitted:       abi.NewStoragePower(0),
		TotalQualityAdjPower:      abi.NewStoragePower(0),
		TotalQABytesCommitted:     abi.NewStoragePower(0),
		TotalPledgeCollateral:     abi.NewTokenAmount(0),
		ThisEpochRawBytePower:     abi.NewStoragePower(0),
		ThisEpochQualityAdjPower:  abi.NewStoragePower(0),
		ThisEpochPledgeCollateral: abi.NewTokenAmount(0),
		ThisEpochQAPowerSmoothed:  smoothing.NewEstimate(InitialQAPowerEstimatePosition, InitialQAPowerEstimateVelocity),
		FirstCronEpoch:            0,
		CronEventQueue:            emptyCronQueueMMapCid,
		Claims:                    emptyClaimsMapCid,
		MinerCount:                0,
		MinerAboveMinPowerCount:   0,
	}, nil
}

type Claim struct {
	// Miner's proof type used to determine minimum miner size
	WindowPoStProofType abi.RegisteredPoStProof

	// Sum of raw byte power for a miner's sectors.
	RawBytePower abi.StoragePower

	// Sum of quality adjusted power for a miner's sectors.
	QualityAdjPower abi.StoragePower
}

// MinerNominalPowerMeetsConsensusMinimum is used to validate Election PoSt
// winners outside the chain state. If the miner has over a threshold of power
// the miner meets the minimum.  If the network is a below a threshold of
// miners and has power > zero the miner meets the minimum.
func (st *State) MinerNominalPowerMeetsConsensusMinimum(s adt.Store, miner addr.Address) (bool, error) { //nolint:deadcode,unused
	claims, err := adt.AsMap(s, st.Claims, builtin.DefaultHamtBitwidth)
	if err != nil {
		return false, xerrors.Errorf("failed to load claims: %w", err)
	}

	claim, ok, err := getClaim(claims, miner)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, xerrors.Errorf("no claim for actor %w", miner)
	}

	minerNominalPower := claim.RawBytePower
	minerMinPower, err := builtin.ConsensusMinerMinPower(claim.WindowPoStProofType)
	if err != nil {
		return false, xerrors.Errorf("could not get miner min power from proof type: %w", err)
	}

	// if miner is larger than min power requirement, we're set
	if minerNominalPower.GreaterThanEqual(minerMinPower) {
		return true, nil
	}

	// otherwise, if ConsensusMinerMinMiners miners meet min power requirement, return false
	if st.MinerAboveMinPowerCount >= ConsensusMinerMinMiners {
		return false, nil
	}

	// If fewer than ConsensusMinerMinMiners over threshold miner can win a block with non-zero power
	return minerNominalPower.GreaterThan(abi.NewStoragePower(0)), nil
}

func (st *State) GetClaim(s adt.Store, a addr.Address) (*Claim, bool, error) {
	claims, err := adt.AsMap(s, st.Claims, builtin.DefaultHamtBitwidth)
	if err != nil {
		return nil, false, xerrors.Errorf("failed to load claims: %w", err)
	}
	return getClaim(claims, a)
}

func getClaim(claims *adt.Map, a addr.Address) (*Claim, bool, error) {
	var out Claim
	found, err := claims.Get(abi.AddrKey(a), &out)
	if err != nil {
		return nil, false, xerrors.Errorf("failed to get claim for address %v: %w", a, err)
	}
	if !found {
		return nil, false, nil
	}
	return &out, true, nil
}
