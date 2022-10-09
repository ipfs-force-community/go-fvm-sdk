package adt

import "github.com/ipfs/go-cid"

// Multimap stores multiple values per key in a HAMT of AMTs.
// The order of insertion of values for each key is retained.
type Multimap struct {
	mp            *Map
	innerBitwidth int
}

// Creates a new map backed by an empty HAMT and flushes it to the store.
// The outer map has a branching factor of 2^bitwidth.
func MakeEmptyMultimap(s Store, outerBitwidth, innerBitwidth int) (*Multimap, error) {
	m, err := MakeEmptyMap(s, outerBitwidth)
	if err != nil {
		return nil, err
	}
	return &Multimap{m, innerBitwidth}, nil
}

// Creates and stores a new empty multimap, returning its CID.
func StoreEmptyMultimap(store Store, outerBitwidth, innerBitwidth int) (cid.Cid, error) {
	mmap, err := MakeEmptyMultimap(store, outerBitwidth, innerBitwidth)
	if err != nil {
		return cid.Undef, err
	}
	return mmap.Root()
}

// Returns the root cid of the underlying HAMT.
func (mm *Multimap) Root() (cid.Cid, error) {
	return mm.mp.Root()
}
