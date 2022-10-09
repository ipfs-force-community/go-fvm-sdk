package miner

import (
	"github.com/filecoin-project/go-bitfield"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/ipfs/go-cid"
	"golang.org/x/xerrors"
)

type Partition struct {
	// Sector numbers in this partition, including faulty, unproven, and terminated sectors.
	Sectors bitfield.BitField
	// Unproven sectors in this partition. This bitfield will be cleared on
	// a successful window post (or at the end of the partition's next
	// deadline). At that time, any still unproven sectors will be added to
	// the faulty sector bitfield.
	Unproven bitfield.BitField
	// Subset of sectors detected/declared faulty and not yet recovered (excl. from PoSt).
	// Faults ∩ Terminated = ∅
	Faults bitfield.BitField
	// Subset of faulty sectors expected to recover on next PoSt
	// Recoveries ∩ Terminated = ∅
	Recoveries bitfield.BitField
	// Subset of sectors terminated but not yet removed from partition (excl. from PoSt)
	Terminated bitfield.BitField
	// Maps epochs sectors that expire in or before that epoch.
	// An expiration may be an "on-time" scheduled expiration, or early "faulty" expiration.
	// Keys are quantized to last-in-deadline epochs.
	ExpirationsEpochs cid.Cid // AMT[ChainEpoch]ExpirationSet
	// Subset of terminated that were before their committed expiration epoch, by termination epoch.
	// Termination fees have not yet been calculated or paid and associated deals have not yet been
	// canceled but effective power has already been adjusted.
	// Not quantized.
	EarlyTerminated cid.Cid // AMT[ChainEpoch]BitField

	// Power of not-yet-terminated sectors (incl faulty & unproven).
	LivePower PowerPair
	// Power of yet-to-be-proved sectors (never faulty).
	UnprovenPower PowerPair
	// Power of currently-faulty sectors. FaultyPower <= LivePower.
	FaultyPower PowerPair
	// Power of expected-to-recover sectors. RecoveringPower <= FaultyPower.
	RecoveringPower PowerPair
}

// Bitwidth of AMTs determined empirically from mutation patterns and projections of mainnet data.
const PartitionExpirationAmtBitwidth = 4
const PartitionEarlyTerminationArrayAmtBitwidth = 3

// Value type for a pair of raw and QA power.
type PowerPair struct {
	Raw abi.StoragePower
	QA  abi.StoragePower
}

// Live sectors are those that are not terminated (but may be faulty).
func (p *Partition) LiveSectors() (bitfield.BitField, error) {
	live, err := bitfield.SubtractBitField(p.Sectors, p.Terminated)
	if err != nil {
		return bitfield.BitField{}, xerrors.Errorf("failed to compute live sectors: %w", err)
	}
	return live, nil

}

// Active sectors are those that are neither terminated nor faulty nor unproven, i.e. actively contributing power.
func (p *Partition) ActiveSectors() (bitfield.BitField, error) {
	live, err := p.LiveSectors()
	if err != nil {
		return bitfield.BitField{}, err
	}
	nonFaulty, err := bitfield.SubtractBitField(live, p.Faults)
	if err != nil {
		return bitfield.BitField{}, xerrors.Errorf("failed to compute active sectors: %w", err)
	}
	active, err := bitfield.SubtractBitField(nonFaulty, p.Unproven)
	if err != nil {
		return bitfield.BitField{}, xerrors.Errorf("failed to compute active sectors: %w", err)
	}
	return active, err
}

// Activates unproven sectors, returning the activated power.
func (p *Partition) ActivateUnproven() PowerPair {
	newPower := p.UnprovenPower
	p.UnprovenPower = NewPowerPairZero()
	p.Unproven = bitfield.New()
	return newPower
}

//
// PowerPair
//

func NewPowerPairZero() PowerPair {
	return NewPowerPair(big.Zero(), big.Zero())
}

func NewPowerPair(raw, qa abi.StoragePower) PowerPair {
	return PowerPair{Raw: raw, QA: qa}
}

func (pp PowerPair) Add(other PowerPair) PowerPair {
	return PowerPair{
		Raw: big.Add(pp.Raw, other.Raw),
		QA:  big.Add(pp.QA, other.QA),
	}
}

func (pp PowerPair) Sub(other PowerPair) PowerPair {
	return PowerPair{
		Raw: big.Sub(pp.Raw, other.Raw),
		QA:  big.Sub(pp.QA, other.QA),
	}
}
