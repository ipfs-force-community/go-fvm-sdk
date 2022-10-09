package miner

import (
	"github.com/filecoin-project/go-state-types/builtin"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/ipfs/go-cid"
	mh "github.com/multiformats/go-multihash"
)

// The period over which a miner's active sectors are expected to be proven via WindowPoSt.
// This guarantees that (1) user data is proven daily, (2) user data is stored for 24h by a rational miner
// (due to Window PoSt cost assumption).
var WPoStProvingPeriod = abi.ChainEpoch(builtin.EpochsInDay) // 24 hours PARAM_SPEC

// The period between the opening and the closing of a WindowPoSt deadline in which the miner is expected to
// provide a Window PoSt proof.
// This provides a miner enough time to compute and propagate a Window PoSt proof.
var WPoStChallengeWindow = abi.ChainEpoch(30 * 60 / builtin.EpochDurationSeconds) // 30 minutes (48 per day) PARAM_SPEC

// WPoStDisputeWindow is the period after a challenge window ends during which
// PoSts submitted during that period may be disputed.
var WPoStDisputeWindow = 2 * ChainFinality // PARAM_SPEC

// The number of non-overlapping PoSt deadlines in a proving period.
// This spreads a miner's Window PoSt work across a proving period.
const WPoStPeriodDeadlines = uint64(48) // PARAM_SPEC

// MaxPartitionsPerDeadline is the maximum number of partitions that will be assigned to a deadline.
// For a minimum storage of upto 1Eib, we need 300 partitions per deadline.
// 48 * 32GiB * 2349 * 300 = 1.00808144 EiB
// So, to support upto 10Eib storage, we set this to 3000.
const MaxPartitionsPerDeadline = 3000

// The maximum number of partitions that can be loaded in a single invocation.
// This limits the number of simultaneous fault, recovery, or sector-extension declarations.
// We set this to same as MaxPartitionsPerDeadline so we can process that many partitions every deadline.
const AddressedPartitionsMax = MaxPartitionsPerDeadline

// Maximum number of unique "declarations" in batch operations.
const DeclarationsMax = AddressedPartitionsMax

// The maximum number of sector infos that can be loaded in a single invocation.
// This limits the amount of state to be read in a single message execution.
const AddressedSectorsMax = 25_000 // PARAM_SPEC

// Epochs after which chain state is final with overwhelming probability (hence the likelihood of two fork of this size is negligible)
// This is a conservative value that is chosen via simulations of all known attacks.
const ChainFinality = abi.ChainEpoch(900) // PARAM_SPEC

// Prefix for sealed sector CIDs (CommR).
var SealedCIDPrefix = cid.Prefix{
	Version:  1,
	Codec:    cid.FilCommitmentSealed,
	MhType:   mh.POSEIDON_BLS12_381_A1_FC1,
	MhLength: 32,
}

// List of proof types which may be used when creating a new miner actor.
// This is mutable to allow configuration of testing and development networks.
var WindowPoStProofTypes = map[abi.RegisteredPoStProof]struct{}{
	abi.RegisteredPoStProof_StackedDrgWindow32GiBV1: {},
	abi.RegisteredPoStProof_StackedDrgWindow64GiBV1: {},
}

// Maximum delay to allow between sector pre-commit and subsequent proof.
// The allowable delay depends on seal proof algorithm.
var MaxProveCommitDuration = map[abi.RegisteredSealProof]abi.ChainEpoch{
	abi.RegisteredSealProof_StackedDrg32GiBV1:  builtin.EpochsInDay + PreCommitChallengeDelay, // PARAM_SPEC
	abi.RegisteredSealProof_StackedDrg2KiBV1:   builtin.EpochsInDay + PreCommitChallengeDelay,
	abi.RegisteredSealProof_StackedDrg8MiBV1:   builtin.EpochsInDay + PreCommitChallengeDelay,
	abi.RegisteredSealProof_StackedDrg512MiBV1: builtin.EpochsInDay + PreCommitChallengeDelay,
	abi.RegisteredSealProof_StackedDrg64GiBV1:  builtin.EpochsInDay + PreCommitChallengeDelay,

	abi.RegisteredSealProof_StackedDrg32GiBV1_1:  30*builtin.EpochsInDay + PreCommitChallengeDelay, // PARAM_SPEC
	abi.RegisteredSealProof_StackedDrg2KiBV1_1:   30*builtin.EpochsInDay + PreCommitChallengeDelay,
	abi.RegisteredSealProof_StackedDrg8MiBV1_1:   30*builtin.EpochsInDay + PreCommitChallengeDelay,
	abi.RegisteredSealProof_StackedDrg512MiBV1_1: 30*builtin.EpochsInDay + PreCommitChallengeDelay,
	abi.RegisteredSealProof_StackedDrg64GiBV1_1:  30*builtin.EpochsInDay + PreCommitChallengeDelay,
}

// The maximum number of sector pre-commitments in a single batch.
// 32 sectors per epoch would support a single miner onboarding 1EiB of 32GiB sectors in 1 year.
const PreCommitSectorBatchMaxSize = 256

// The maximum number of sector replica updates in a single batch.
// Same as PreCommitSectorBatchMaxSize for consistency
const ProveReplicaUpdatesMaxSize = PreCommitSectorBatchMaxSize

// Maximum delay between challenge and pre-commitment.
// This prevents a miner sealing sectors far in advance of committing them to the chain, thus committing to a
// particular chain.
var MaxPreCommitRandomnessLookback = builtin.EpochsInDay + ChainFinality // PARAM_SPEC

// Number of epochs between publishing a sector pre-commitment and when the challenge for interactive PoRep is drawn.
// This (1) prevents a miner predicting a challenge before staking their pre-commit deposit, and
// (2) prevents a miner attempting a long fork in the past to insert a pre-commitment after seeing the challenge.
var PreCommitChallengeDelay = abi.ChainEpoch(150) // PARAM_SPEC

// Lookback from the deadline's challenge window opening from which to sample chain randomness for the WindowPoSt challenge seed.
// This means that deadline windows can be non-overlapping (which make the programming simpler) without requiring a
// miner to wait for chain stability during the challenge window.
// This value cannot be too large lest it compromise the rationality of honest storage (from Window PoSt cost assumptions).
const WPoStChallengeLookback = abi.ChainEpoch(20) // PARAM_SPEC

// Minimum period between fault declaration and the next deadline opening.
// If the number of epochs between fault declaration and deadline's challenge window opening is lower than FaultDeclarationCutoff,
// the fault declaration is considered invalid for that deadline.
// This guarantees that a miner is not likely to successfully fork the chain and declare a fault after seeing the challenges.
const FaultDeclarationCutoff = WPoStChallengeLookback + 50 // PARAM_SPEC

// The maximum age of a fault before the sector is terminated.
// This bounds the time a miner can lose client's data before sacrificing pledge and deal collateral.
var FaultMaxAge = WPoStProvingPeriod * 42 // PARAM_SPEC

// Staging period for a miner worker key change.
// This delay prevents a miner choosing a more favorable worker key that wins leader elections.
const WorkerKeyChangeDelay = ChainFinality // PARAM_SPEC

// Minimum number of epochs past the current epoch a sector may be set to expire.
const MinSectorExpiration = 180 * builtin.EpochsInDay // PARAM_SPEC

// The maximum number of epochs past the current epoch that sector lifetime may be extended.
// A sector may be extended multiple times, however, the total maximum lifetime is also bounded by
// the associated seal proof's maximum lifetime.
const MaxSectorExpirationExtension = 540 * builtin.EpochsInDay // PARAM_SPEC

// DealWeight and VerifiedDealWeight are spacetime occupied by regular deals and verified deals in a sector.
// Sum of DealWeight and VerifiedDealWeight should be less than or equal to total SpaceTime of a sector.
// Sectors full of VerifiedDeals will have a SectorQuality of VerifiedDealWeightMultiplier/QualityBaseMultiplier.
// Sectors full of Deals will have a SectorQuality of DealWeightMultiplier/QualityBaseMultiplier.
// Sectors with neither will have a SectorQuality of QualityBaseMultiplier/QualityBaseMultiplier.
// SectorQuality of a sector is a weighted average of multipliers based on their proportions.
func QualityForWeight(size abi.SectorSize, duration abi.ChainEpoch, dealWeight, verifiedWeight abi.DealWeight) abi.SectorQuality {
	// sectorSpaceTime = size * duration
	sectorSpaceTime := big.Mul(big.NewIntUnsigned(uint64(size)), big.NewInt(int64(duration)))
	// totalDealSpaceTime = dealWeight + verifiedWeight
	totalDealSpaceTime := big.Add(dealWeight, verifiedWeight)

	// Base - all size * duration of non-deals
	// weightedBaseSpaceTime = (sectorSpaceTime - totalDealSpaceTime) * QualityBaseMultiplier
	weightedBaseSpaceTime := big.Mul(big.Sub(sectorSpaceTime, totalDealSpaceTime), builtin.QualityBaseMultiplier)
	// Deal - all deal size * deal duration * 10
	// weightedDealSpaceTime = dealWeight * DealWeightMultiplier
	weightedDealSpaceTime := big.Mul(dealWeight, builtin.DealWeightMultiplier)
	// Verified - all verified deal size * verified deal duration * 100
	// weightedVerifiedSpaceTime = verifiedWeight * VerifiedDealWeightMultiplier
	weightedVerifiedSpaceTime := big.Mul(verifiedWeight, builtin.VerifiedDealWeightMultiplier)
	// Sum - sum of all spacetime
	// weightedSumSpaceTime = weightedBaseSpaceTime + weightedDealSpaceTime + weightedVerifiedSpaceTime
	weightedSumSpaceTime := big.Sum(weightedBaseSpaceTime, weightedDealSpaceTime, weightedVerifiedSpaceTime)
	// scaledUpWeightedSumSpaceTime = weightedSumSpaceTime * 2^20
	scaledUpWeightedSumSpaceTime := big.Lsh(weightedSumSpaceTime, builtin.SectorQualityPrecision)

	// Average of weighted space time: (scaledUpWeightedSumSpaceTime / sectorSpaceTime * 10)
	return big.Div(big.Div(scaledUpWeightedSumSpaceTime, sectorSpaceTime), builtin.QualityBaseMultiplier)
}

// The power for a sector size, committed duration, and weight.
func QAPowerForWeight(size abi.SectorSize, duration abi.ChainEpoch, dealWeight, verifiedWeight abi.DealWeight) abi.StoragePower {
	quality := QualityForWeight(size, duration, dealWeight, verifiedWeight)
	return big.Rsh(big.Mul(big.NewIntUnsigned(uint64(size)), quality), builtin.SectorQualityPrecision)
}

const MaxAggregatedSectors = 819
const MinAggregatedSectors = 4
const MaxAggregateProofSize = 81960

// Specification for a linear vesting schedule.
type VestSpec struct {
	InitialDelay abi.ChainEpoch // Delay before any amount starts vesting.
	VestPeriod   abi.ChainEpoch // Period over which the total should vest, after the initial delay.
	StepDuration abi.ChainEpoch // Duration between successive incremental vests (independent of vesting period).
	Quantization abi.ChainEpoch // Maximum precision of vesting table (limits cardinality of table).
}

// Returns maximum achievable QA power.
func QAPowerMax(size abi.SectorSize) abi.StoragePower {
	return big.Div(
		big.Mul(big.NewInt(int64(size)), builtin.VerifiedDealWeightMultiplier),
		builtin.QualityBaseMultiplier)
}
