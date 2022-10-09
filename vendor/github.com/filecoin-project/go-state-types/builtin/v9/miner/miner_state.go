package miner

import (
	addr "github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-bitfield"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/go-state-types/builtin/v9/util/adt"
	"github.com/filecoin-project/go-state-types/dline"
	xc "github.com/filecoin-project/go-state-types/exitcode"
	cid "github.com/ipfs/go-cid"
	"golang.org/x/xerrors"
)

// Balance of Miner Actor should be greater than or equal to
// the sum of PreCommitDeposits and LockedFunds.
// It is possible for balance to fall below the sum of
// PCD, LF and InitialPledgeRequirements, and this is a bad
// state (IP Debt) that limits a miner actor's behavior (i.e. no balance withdrawals)
// Excess balance as computed by st.GetAvailableBalance will be
// withdrawable or usable for pre-commit deposit or pledge lock-up.
type State struct {
	// Information not related to sectors.
	Info cid.Cid

	PreCommitDeposits abi.TokenAmount // Total funds locked as PreCommitDeposits
	LockedFunds       abi.TokenAmount // Total rewards and added funds locked in vesting table

	VestingFunds cid.Cid // VestingFunds (Vesting Funds schedule for the miner).

	FeeDebt abi.TokenAmount // Absolute value of debt this miner owes from unpaid fees

	InitialPledge abi.TokenAmount // Sum of initial pledge requirements of all active sectors

	// Sectors that have been pre-committed but not yet proven.
	PreCommittedSectors cid.Cid // Map, HAMT[SectorNumber]SectorPreCommitOnChainInfo

	// PreCommittedSectorsCleanUp maintains the state required to cleanup expired PreCommittedSectors.
	PreCommittedSectorsCleanUp cid.Cid // BitFieldQueue (AMT[Epoch]*BitField)

	// Allocated sector IDs. Sector IDs can never be reused once allocated.
	AllocatedSectors cid.Cid // BitField

	// Information for all proven and not-yet-garbage-collected sectors.
	//
	// Sectors are removed from this AMT when the partition to which the
	// sector belongs is compacted.
	Sectors cid.Cid // Array, AMT[SectorNumber]SectorOnChainInfo (sparse)

	// DEPRECATED. This field will change names and no longer be updated every proving period in a future upgrade
	// The first epoch in this miner's current proving period. This is the first epoch in which a PoSt for a
	// partition at the miner's first deadline may arrive. Alternatively, it is after the last epoch at which
	// a PoSt for the previous window is valid.
	// Always greater than zero, this may be greater than the current epoch for genesis miners in the first
	// WPoStProvingPeriod epochs of the chain; the epochs before the first proving period starts are exempt from Window
	// PoSt requirements.
	// Updated at the end of every period by a cron callback.
	ProvingPeriodStart abi.ChainEpoch

	// DEPRECATED. This field will be removed from state in a future upgrade.
	// Index of the deadline within the proving period beginning at ProvingPeriodStart that has not yet been
	// finalized.
	// Updated at the end of each deadline window by a cron callback.
	CurrentDeadline uint64

	// The sector numbers due for PoSt at each deadline in the current proving period, frozen at period start.
	// New sectors are added and expired ones removed at proving period boundary.
	// Faults are not subtracted from this in state, but on the fly.
	Deadlines cid.Cid

	// Deadlines with outstanding fees for early sector termination.
	EarlyTerminations bitfield.BitField

	// True when miner cron is active, false otherwise
	DeadlineCronActive bool
}

// Bitwidth of AMTs determined empirically from mutation patterns and projections of mainnet data.
const PrecommitCleanUpAmtBitwidth = 6
const SectorsAmtBitwidth = 5

type MinerInfo struct {
	// Account that owns this miner.
	// - Income and returned collateral are paid to this address.
	// - This address is also allowed to change the worker address for the miner.
	Owner addr.Address // Must be an ID-address.

	// Worker account for this miner.
	// The associated pubkey-type address is used to sign blocks and messages on behalf of this miner.
	Worker addr.Address // Must be an ID-address.

	// Additional addresses that are permitted to submit messages controlling this actor (optional).
	ControlAddresses []addr.Address // Must all be ID addresses.

	PendingWorkerKey *WorkerKeyChange

	// Byte array representing a Libp2p identity that should be used when connecting to this miner.
	PeerId abi.PeerID

	// Slice of byte arrays representing Libp2p multi-addresses used for establishing a connection with this miner.
	Multiaddrs []abi.Multiaddrs

	// The proof type used for Window PoSt for this miner.
	// A miner may commit sectors with different seal proof types (but compatible sector size and
	// corresponding PoSt proof types).
	WindowPoStProofType abi.RegisteredPoStProof

	// Amount of space in each sector committed by this miner.
	// This is computed from the proof type and represented here redundantly.
	SectorSize abi.SectorSize

	// The number of sectors in each Window PoSt partition (proof).
	// This is computed from the proof type and represented here redundantly.
	WindowPoStPartitionSectors uint64

	// The next epoch this miner is eligible for certain permissioned actor methods
	// and winning block elections as a result of being reported for a consensus fault.
	ConsensusFaultElapsed abi.ChainEpoch

	// A proposed new owner account for this miner.
	// Must be confirmed by a message from the pending address itself.
	PendingOwnerAddress *addr.Address

	// Beneficiary address for this miner.
	// This is the address that tokens will be withdrawn to
	Beneficiary addr.Address

	// Beneficiary's withdrawal quota, how much of the quota has been withdrawn,
	// and when the Beneficiary expires.
	BeneficiaryTerm BeneficiaryTerm

	// A proposed change to `BenificiaryTerm`
	PendingBeneficiaryTerm *PendingBeneficiaryChange
}

type WorkerKeyChange struct {
	NewWorker   addr.Address // Must be an ID address
	EffectiveAt abi.ChainEpoch
}

// Information provided by a miner when pre-committing a sector.
type SectorPreCommitInfo struct {
	SealProof     abi.RegisteredSealProof
	SectorNumber  abi.SectorNumber
	SealedCID     cid.Cid `checked:"true"` // CommR
	SealRandEpoch abi.ChainEpoch
	DealIDs       []abi.DealID
	Expiration    abi.ChainEpoch
	UnsealedCid   *cid.Cid
}

// Information stored on-chain for a pre-committed sector.
type SectorPreCommitOnChainInfo struct {
	Info             SectorPreCommitInfo
	PreCommitDeposit abi.TokenAmount
	PreCommitEpoch   abi.ChainEpoch
}

// Information stored on-chain for a proven sector.
type SectorOnChainInfo struct {
	SectorNumber          abi.SectorNumber
	SealProof             abi.RegisteredSealProof // The seal proof type implies the PoSt proof/s
	SealedCID             cid.Cid                 // CommR
	DealIDs               []abi.DealID
	Activation            abi.ChainEpoch  // Epoch during which the sector proof was accepted
	Expiration            abi.ChainEpoch  // Epoch during which the sector expires
	DealWeight            abi.DealWeight  // Integral of active deals over sector lifetime
	VerifiedDealWeight    abi.DealWeight  // Integral of active verified deals over sector lifetime
	InitialPledge         abi.TokenAmount // Pledge collected to commit this sector
	ExpectedDayReward     abi.TokenAmount // Expected one day projection of reward for sector computed at activation time
	ExpectedStoragePledge abi.TokenAmount // Expected twenty day projection of reward for sector computed at activation time
	ReplacedSectorAge     abi.ChainEpoch  // Age of sector this sector replaced or zero
	ReplacedDayReward     abi.TokenAmount // Day reward of sector this sector replace or zero
	SectorKeyCID          *cid.Cid        // The original SealedSectorCID, only gets set on the first ReplicaUpdate
}

func (st *State) GetInfo(store adt.Store) (*MinerInfo, error) {
	var info MinerInfo
	if err := store.Get(store.Context(), st.Info, &info); err != nil {
		return nil, xerrors.Errorf("failed to get miner info %w", err)
	}
	return &info, nil
}

type BeneficiaryTerm struct {
	Quota      abi.TokenAmount
	UsedQuota  abi.TokenAmount
	Expiration abi.ChainEpoch
}

type PendingBeneficiaryChange struct {
	NewBeneficiary        addr.Address
	NewQuota              abi.TokenAmount
	NewExpiration         abi.ChainEpoch
	ApprovedByBeneficiary bool
	ApprovedByNominee     bool
}

// Returns deadline calculations for the state recorded proving period and deadline. This is out of date if the a
// miner does not have an active miner cron
func (st *State) RecordedDeadlineInfo(currEpoch abi.ChainEpoch) *dline.Info {
	return NewDeadlineInfo(st.ProvingPeriodStart, st.CurrentDeadline, currEpoch)
}

// Returns deadline calculations for the current (according to state) proving period
func (st *State) QuantSpecForDeadline(dlIdx uint64) builtin.QuantSpec {
	return QuantSpecForDeadline(NewDeadlineInfo(st.ProvingPeriodStart, dlIdx, 0))
}

func (st *State) GetPrecommittedSector(store adt.Store, sectorNo abi.SectorNumber) (*SectorPreCommitOnChainInfo, bool, error) {
	precommitted, err := adt.AsMap(store, st.PreCommittedSectors, builtin.DefaultHamtBitwidth)
	if err != nil {
		return nil, false, err
	}

	var info SectorPreCommitOnChainInfo
	found, err := precommitted.Get(SectorKey(sectorNo), &info)
	if err != nil {
		return nil, false, xerrors.Errorf("failed to load precommitment for %v: %w", sectorNo, err)
	}
	return &info, found, nil
}

func (st *State) GetSector(store adt.Store, sectorNo abi.SectorNumber) (*SectorOnChainInfo, bool, error) {
	sectors, err := LoadSectors(store, st.Sectors)
	if err != nil {
		return nil, false, err
	}

	return sectors.Get(sectorNo)
}

func (st *State) FindSector(store adt.Store, sno abi.SectorNumber) (uint64, uint64, error) {
	deadlines, err := st.LoadDeadlines(store)
	if err != nil {
		return 0, 0, err
	}
	return FindSector(store, deadlines, sno)
}

func (st *State) LoadDeadlines(store adt.Store) (*Deadlines, error) {
	var deadlines Deadlines
	if err := store.Get(store.Context(), st.Deadlines, &deadlines); err != nil {
		return nil, xc.ErrIllegalState.Wrapf("failed to load deadlines (%s): %w", st.Deadlines, err)
	}

	return &deadlines, nil
}

func (st *State) SaveDeadlines(store adt.Store, deadlines *Deadlines) error {
	c, err := store.Put(store.Context(), deadlines)
	if err != nil {
		return err
	}
	st.Deadlines = c
	return nil
}

// LoadVestingFunds loads the vesting funds table from the store
func (st *State) LoadVestingFunds(store adt.Store) (*VestingFunds, error) {
	var funds VestingFunds
	if err := store.Get(store.Context(), st.VestingFunds, &funds); err != nil {
		return nil, xerrors.Errorf("failed to load vesting funds (%s): %w", st.VestingFunds, err)
	}

	return &funds, nil
}

// CheckVestedFunds returns the amount of vested funds that have vested before the provided epoch.
func (st *State) CheckVestedFunds(store adt.Store, currEpoch abi.ChainEpoch) (abi.TokenAmount, error) {
	vestingFunds, err := st.LoadVestingFunds(store)
	if err != nil {
		return big.Zero(), xerrors.Errorf("failed to load vesting funds: %w", err)
	}

	amountVested := abi.NewTokenAmount(0)

	for i := range vestingFunds.Funds {
		vf := vestingFunds.Funds[i]
		epoch := vf.Epoch
		amount := vf.Amount

		if epoch >= currEpoch {
			break
		}

		amountVested = big.Add(amountVested, amount)
	}

	return amountVested, nil
}

// Unclaimed funds that are not locked -- includes free funds and does not
// account for fee debt.  Always greater than or equal to zero
func (st *State) GetUnlockedBalance(actorBalance abi.TokenAmount) (abi.TokenAmount, error) {
	unlockedBalance := big.Subtract(actorBalance, st.LockedFunds, st.PreCommitDeposits, st.InitialPledge)
	if unlockedBalance.LessThan(big.Zero()) {
		return big.Zero(), xerrors.Errorf("negative unlocked balance %v", unlockedBalance)
	}
	return unlockedBalance, nil
}

// Unclaimed funds.  Actor balance - (locked funds, precommit deposit, initial pledge, fee debt)
// Can go negative if the miner is in IP debt
func (st *State) GetAvailableBalance(actorBalance abi.TokenAmount) (abi.TokenAmount, error) {
	unlockedBalance, err := st.GetUnlockedBalance(actorBalance)
	if err != nil {
		return big.Zero(), err
	}
	return big.Subtract(unlockedBalance, st.FeeDebt), nil
}

//
// Misc helpers
//

func SectorKey(e abi.SectorNumber) abi.Keyer {
	return abi.UIntKey(uint64(e))
}
