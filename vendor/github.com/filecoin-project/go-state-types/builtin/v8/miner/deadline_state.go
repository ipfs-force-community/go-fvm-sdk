package miner

import (
	"github.com/filecoin-project/go-bitfield"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin/v8/util/adt"
	xc "github.com/filecoin-project/go-state-types/exitcode"
	"github.com/filecoin-project/go-state-types/proof"
	"github.com/ipfs/go-cid"
	"golang.org/x/xerrors"
)

// Deadlines contains Deadline objects, describing the sectors due at the given
// deadline and their state (faulty, terminated, recovering, etc.).
type Deadlines struct {
	// Note: we could inline part of the deadline struct (e.g., active/assigned sectors)
	// to make new sector assignment cheaper. At the moment, assigning a sector requires
	// loading all deadlines to figure out where best to assign new sectors.
	Due [WPoStPeriodDeadlines]cid.Cid // []Deadline
}

// Deadline holds the state for all sectors due at a specific deadline.
type Deadline struct {
	// Partitions in this deadline, in order.
	// The keys of this AMT are always sequential integers beginning with zero.
	Partitions cid.Cid // AMT[PartitionNumber]Partition

	// Maps epochs to partitions that _may_ have sectors that expire in or
	// before that epoch, either on-time or early as faults.
	// Keys are quantized to final epochs in each proving deadline.
	//
	// NOTE: Partitions MUST NOT be removed from this queue (until the
	// associated epoch has passed) even if they no longer have sectors
	// expiring at that epoch. Sectors expiring at this epoch may later be
	// recovered, and this queue will not be updated at that time.
	ExpirationsEpochs cid.Cid // AMT[ChainEpoch]BitField

	// Partitions that have been proved by window PoSts so far during the
	// current challenge window.
	// NOTE: This bitfield includes both partitions whose proofs
	// were optimistically accepted and stored in
	// OptimisticPoStSubmissions, and those whose proofs were
	// verified on-chain.
	PartitionsPoSted bitfield.BitField

	// Partitions with sectors that terminated early.
	EarlyTerminations bitfield.BitField

	// The number of non-terminated sectors in this deadline (incl faulty).
	LiveSectors uint64

	// The total number of sectors in this deadline (incl dead).
	TotalSectors uint64

	// Memoized sum of faulty power in partitions.
	FaultyPower PowerPair

	// AMT of optimistically accepted WindowPoSt proofs, submitted during
	// the current challenge window. At the end of the challenge window,
	// this AMT will be moved to OptimisticPoStSubmissionsSnapshot. WindowPoSt proofs
	// verified on-chain do not appear in this AMT.
	OptimisticPoStSubmissions cid.Cid // AMT[]WindowedPoSt

	// Snapshot of the miner's sectors AMT at the end of the previous challenge
	// window for this deadline.
	SectorsSnapshot cid.Cid

	// Snapshot of partition state at the end of the previous challenge
	// window for this deadline.
	PartitionsSnapshot cid.Cid

	// Snapshot of the proofs submitted by the end of the previous challenge
	// window for this deadline.
	//
	// These proofs may be disputed via DisputeWindowedPoSt. Successfully
	// disputed window PoSts are removed from the snapshot.
	OptimisticPoStSubmissionsSnapshot cid.Cid
}

type WindowedPoSt struct {
	// Partitions proved by this WindowedPoSt.
	Partitions bitfield.BitField
	// Array of proofs, one per distinct registered proof type present in
	// the sectors being proven. In the usual case of a single proof type,
	// this array will always have a single element (independent of number
	// of partitions).
	Proofs []proof.PoStProof
}

// Bitwidth of AMTs determined empirically from mutation patterns and projections of mainnet data.
const DeadlinePartitionsAmtBitwidth = 3 // Usually a small array
const DeadlineExpirationAmtBitwidth = 5

// Given that 4 partitions can be proven in one post, this AMT's height will
// only exceed the partition AMT's height at ~0.75EiB of storage.
const DeadlineOptimisticPoStSubmissionsAmtBitwidth = 2

//
// Deadlines (plural)
//

func (d *Deadlines) LoadDeadline(store adt.Store, dlIdx uint64) (*Deadline, error) {
	if dlIdx >= uint64(len(d.Due)) {
		return nil, xc.ErrIllegalArgument.Wrapf("invalid deadline %d", dlIdx)
	}
	deadline := new(Deadline)
	err := store.Get(store.Context(), d.Due[dlIdx], deadline)
	if err != nil {
		return nil, xc.ErrIllegalState.Wrapf("failed to lookup deadline %d: %w", dlIdx, err)
	}
	return deadline, nil
}

func (d *Deadlines) ForEach(store adt.Store, cb func(dlIdx uint64, dl *Deadline) error) error {
	for dlIdx := range d.Due {
		dl, err := d.LoadDeadline(store, uint64(dlIdx))
		if err != nil {
			return err
		}
		err = cb(uint64(dlIdx), dl)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Deadlines) UpdateDeadline(store adt.Store, dlIdx uint64, deadline *Deadline) error {
	if dlIdx >= uint64(len(d.Due)) {
		return xerrors.Errorf("invalid deadline %d", dlIdx)
	}

	if err := deadline.ValidateState(); err != nil {
		return err
	}

	dlCid, err := store.Put(store.Context(), deadline)
	if err != nil {
		return err
	}
	d.Due[dlIdx] = dlCid

	return nil
}

//
// Deadline (singular)
//

func (d *Deadline) PartitionsArray(store adt.Store) (*adt.Array, error) {
	arr, err := adt.AsArray(store, d.Partitions, DeadlinePartitionsAmtBitwidth)
	if err != nil {
		return nil, xc.ErrIllegalState.Wrapf("failed to load partitions: %w", err)
	}
	return arr, nil
}

func (d *Deadline) OptimisticProofsSnapshotArray(store adt.Store) (*adt.Array, error) {
	arr, err := adt.AsArray(store, d.OptimisticPoStSubmissionsSnapshot, DeadlineOptimisticPoStSubmissionsAmtBitwidth)
	if err != nil {
		return nil, xerrors.Errorf("failed to load proofs snapshot: %w", err)
	}
	return arr, nil
}

func (d *Deadline) LoadPartition(store adt.Store, partIdx uint64) (*Partition, error) {
	partitions, err := d.PartitionsArray(store)
	if err != nil {
		return nil, err
	}
	var partition Partition
	found, err := partitions.Get(partIdx, &partition)
	if err != nil {
		return nil, xc.ErrIllegalState.Wrapf("failed to lookup partition %d: %w", partIdx, err)
	}
	if !found {
		return nil, xc.ErrNotFound.Wrapf("no partition %d", partIdx)
	}
	return &partition, nil
}

func (d *Deadline) ValidateState() error {
	if d.LiveSectors > d.TotalSectors {
		return xerrors.Errorf("Deadline left with more live sectors than total: %v", d)
	}

	if d.FaultyPower.Raw.LessThan(big.Zero()) || d.FaultyPower.QA.LessThan(big.Zero()) {
		return xerrors.Errorf("Deadline left with negative faulty power: %v", d)
	}

	return nil
}
