package market

import (
	addr "github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/go-state-types/builtin/v8/util/adt"
	"github.com/filecoin-project/go-state-types/exitcode"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/ipfs/go-cid"
	xerrors "golang.org/x/xerrors"
)

const EpochUndefined = abi.ChainEpoch(-1)

// Bitwidth of AMTs determined empirically from mutation patterns and projections of mainnet data.
const ProposalsAmtBitwidth = 5
const StatesAmtBitwidth = 6

type State struct {
	// Proposals are deals that have been proposed and not yet cleaned up after expiry or termination.
	Proposals cid.Cid // AMT[DealID]DealProposal
	// States contains state for deals that have been activated and not yet cleaned up after expiry or termination.
	// After expiration, the state exists until the proposal is cleaned up too.
	// Invariant: keys(States) ⊆ keys(Proposals).
	States cid.Cid // AMT[DealID]DealState

	// PendingProposals tracks dealProposals that have not yet reached their deal start date.
	// We track them here to ensure that miners can't publish the same deal proposal twice
	PendingProposals cid.Cid // Set[DealCid]

	// Total amount held in escrow, indexed by actor address (including both locked and unlocked amounts).
	EscrowTable cid.Cid // BalanceTable

	// Amount locked, indexed by actor address.
	// Note: the amounts in this table do not affect the overall amount in escrow:
	// only the _portion_ of the total escrow amount that is locked.
	LockedTable cid.Cid // BalanceTable

	NextID abi.DealID

	// Metadata cached for efficient iteration over deals.
	DealOpsByEpoch cid.Cid // SetMultimap, HAMT[epoch]Set
	LastCron       abi.ChainEpoch

	// Total Client Collateral that is locked -> unlocked when deal is terminated
	TotalClientLockedCollateral abi.TokenAmount
	// Total Provider Collateral that is locked -> unlocked when deal is terminated
	TotalProviderLockedCollateral abi.TokenAmount
	// Total storage fee that is locked in escrow -> unlocked when payments are made
	TotalClientStorageFee abi.TokenAmount
}

func ConstructState(store adt.Store) (*State, error) {
	emptyProposalsArrayCid, err := adt.StoreEmptyArray(store, ProposalsAmtBitwidth)
	if err != nil {
		return nil, xerrors.Errorf("failed to create empty array: %w", err)
	}
	emptyStatesArrayCid, err := adt.StoreEmptyArray(store, StatesAmtBitwidth)
	if err != nil {
		return nil, xerrors.Errorf("failed to create empty states array: %w", err)
	}

	emptyPendingProposalsMapCid, err := adt.StoreEmptyMap(store, builtin.DefaultHamtBitwidth)
	if err != nil {
		return nil, xerrors.Errorf("failed to create empty map: %w", err)
	}
	emptyDealOpsHamtCid, err := StoreEmptySetMultimap(store, builtin.DefaultHamtBitwidth)
	if err != nil {
		return nil, xerrors.Errorf("failed to create empty multiset: %w", err)
	}
	emptyBalanceTableCid, err := adt.StoreEmptyMap(store, adt.BalanceTableBitwidth)
	if err != nil {
		return nil, xerrors.Errorf("failed to create empty balance table: %w", err)
	}

	return &State{
		Proposals:        emptyProposalsArrayCid,
		States:           emptyStatesArrayCid,
		PendingProposals: emptyPendingProposalsMapCid,
		EscrowTable:      emptyBalanceTableCid,
		LockedTable:      emptyBalanceTableCid,
		NextID:           abi.DealID(0),
		DealOpsByEpoch:   emptyDealOpsHamtCid,
		LastCron:         abi.ChainEpoch(-1),

		TotalClientLockedCollateral:   abi.NewTokenAmount(0),
		TotalProviderLockedCollateral: abi.NewTokenAmount(0),
		TotalClientStorageFee:         abi.NewTokenAmount(0),
	}, nil
}

// A specialization of a array to deals.
// It is an error to query for a key that doesn't exist.
type DealArray struct {
	*adt.Array
}

// Interprets a store as balance table with root `r`.
func AsDealProposalArray(s adt.Store, r cid.Cid) (*DealArray, error) {
	a, err := adt.AsArray(s, r, ProposalsAmtBitwidth)
	if err != nil {
		return nil, err
	}
	return &DealArray{a}, nil
}

func (d *DealArray) GetDealProposal(dealID abi.DealID) (*DealProposal, error) {
	proposal, found, err := d.Get(dealID)
	if err != nil {
		return nil, xerrors.Errorf("failed to load proposal: %w", err)
	}
	if !found {
		return nil, exitcode.ErrNotFound.Wrapf("no such deal %d", dealID)
	}

	return proposal, nil
}

// Returns the root cid of underlying AMT.
func (t *DealArray) Root() (cid.Cid, error) {
	return t.Array.Root()
}

// Gets the deal for a key. The entry must have been previously initialized.
func (t *DealArray) Get(id abi.DealID) (*DealProposal, bool, error) {
	var value DealProposal
	found, err := t.Array.Get(uint64(id), &value)
	return &value, found, err
}

func (t *DealArray) Set(k abi.DealID, value *DealProposal) error {
	return t.Array.Set(uint64(k), value)
}

func (t *DealArray) Delete(id abi.DealID) error {
	return t.Array.Delete(uint64(id))
}

// Validates a collection of deal dealProposals for activation, and returns their combined weight,
// split into regular deal weight and verified deal weight.
func ValidateDealsForActivation(
	st *State, store adt.Store, dealIDs []abi.DealID, minerAddr addr.Address, sectorExpiry, currEpoch abi.ChainEpoch,
) (big.Int, big.Int, uint64, error) {
	proposals, err := AsDealProposalArray(store, st.Proposals)
	if err != nil {
		return big.Int{}, big.Int{}, 0, xerrors.Errorf("failed to load dealProposals: %w", err)
	}

	return validateAndComputeDealWeight(proposals, dealIDs, minerAddr, sectorExpiry, currEpoch)
}

////////////////////////////////////////////////////////////////////////////////
// Checks
////////////////////////////////////////////////////////////////////////////////

func validateAndComputeDealWeight(proposals *DealArray, dealIDs []abi.DealID, minerAddr addr.Address,
	sectorExpiry abi.ChainEpoch, sectorActivation abi.ChainEpoch) (big.Int, big.Int, uint64, error) {

	seenDealIDs := make(map[abi.DealID]struct{}, len(dealIDs))
	totalDealSpace := uint64(0)
	totalDealSpaceTime := big.Zero()
	totalVerifiedSpaceTime := big.Zero()
	for _, dealID := range dealIDs {
		// Make sure we don't double-count deals.
		if _, seen := seenDealIDs[dealID]; seen {
			return big.Int{}, big.Int{}, 0, exitcode.ErrIllegalArgument.Wrapf("deal ID %d present multiple times", dealID)
		}
		seenDealIDs[dealID] = struct{}{}

		proposal, found, err := proposals.Get(dealID)
		if err != nil {
			return big.Int{}, big.Int{}, 0, xerrors.Errorf("failed to load deal %d: %w", dealID, err)
		}
		if !found {
			return big.Int{}, big.Int{}, 0, exitcode.ErrNotFound.Wrapf("no such deal %d", dealID)
		}
		if err = validateDealCanActivate(proposal, minerAddr, sectorExpiry, sectorActivation); err != nil {
			return big.Int{}, big.Int{}, 0, xerrors.Errorf("cannot activate deal %d: %w", dealID, err)
		}

		// Compute deal weight
		totalDealSpace += uint64(proposal.PieceSize)
		dealSpaceTime := DealWeight(proposal)
		if proposal.VerifiedDeal {
			totalVerifiedSpaceTime = big.Add(totalVerifiedSpaceTime, dealSpaceTime)
		} else {
			totalDealSpaceTime = big.Add(totalDealSpaceTime, dealSpaceTime)
		}
	}
	return totalDealSpaceTime, totalVerifiedSpaceTime, totalDealSpace, nil
}

func validateDealCanActivate(proposal *DealProposal, minerAddr addr.Address, sectorExpiration, sectorActivation abi.ChainEpoch) error {
	if proposal.Provider != minerAddr {
		return exitcode.ErrForbidden.Wrapf("proposal has provider %v, must be %v", proposal.Provider, minerAddr)
	}
	if sectorActivation > proposal.StartEpoch {
		return exitcode.ErrIllegalArgument.Wrapf("proposal start epoch %d has already elapsed at %d", proposal.StartEpoch, sectorActivation)
	}
	if proposal.EndEpoch > sectorExpiration {
		return exitcode.ErrIllegalArgument.Wrapf("proposal expiration %d exceeds sector expiration %d", proposal.EndEpoch, sectorExpiration)
	}
	return nil
}
