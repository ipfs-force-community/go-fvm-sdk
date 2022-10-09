package market

import (
	"github.com/filecoin-project/go-state-types/builtin/v8/util/adt"

	cid "github.com/ipfs/go-cid"
)

type SetMultimap struct {
	mp            *adt.Map
	store         adt.Store
	innerBitwidth int
}

// Creates a new map backed by an empty HAMT and flushes it to the store.
// Both inner and outer HAMTs have branching factor 2^bitwidth.
func MakeEmptySetMultimap(s adt.Store, bitwidth int) (*SetMultimap, error) {
	m, err := adt.MakeEmptyMap(s, bitwidth)
	if err != nil {
		return nil, err
	}
	return &SetMultimap{mp: m, store: s, innerBitwidth: bitwidth}, nil
}

// Writes a new empty map to the store and returns its CID.
func StoreEmptySetMultimap(s adt.Store, bitwidth int) (cid.Cid, error) {
	mm, err := MakeEmptySetMultimap(s, bitwidth)
	if err != nil {
		return cid.Undef, err
	}
	return mm.Root()
}

// Returns the root cid of the underlying HAMT.
func (mm *SetMultimap) Root() (cid.Cid, error) {
	return mm.mp.Root()
}
