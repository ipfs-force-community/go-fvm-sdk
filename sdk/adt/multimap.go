package adt

import (
	"fmt"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/cbor"
	"github.com/ipfs/go-cid"
	cbg "github.com/whyrusleeping/cbor-gen"
)

// Multimap stores multiple values per key in a HAMT of AMTs.
// The order of insertion of values for each key is retained.
type Multimap struct {
	mp            *Map
	innerBitwidth int
}

// AsMultimap interprets a store as a HAMT-based map of AMTs with root `r`.
// The outer map is interpreted with a branching factor of 2^bitwidth.
func AsMultimap(s Store, r cid.Cid, outerBitwidth, innerBitwidth int) (*Multimap, error) {
	m, err := AsMap(s, r, outerBitwidth)
	if err != nil {
		return nil, err
	}

	return &Multimap{m, innerBitwidth}, nil
}

// MakeEmptyMultimap creates a new map backed by an empty HAMT and flushes it to the store.
// The outer map has a branching factor of 2^bitwidth.
func MakeEmptyMultimap(s Store, outerBitwidth, innerBitwidth int) (*Multimap, error) {
	m, err := MakeEmptyMap(s, outerBitwidth)
	if err != nil {
		return nil, err
	}
	return &Multimap{m, innerBitwidth}, nil
}

// StoreEmptyMultimap creates and stores a new empty multimap, returning its CID.
func StoreEmptyMultimap(store Store, outerBitwidth, innerBitwidth int) (cid.Cid, error) {
	mmap, err := MakeEmptyMultimap(store, outerBitwidth, innerBitwidth)
	if err != nil {
		return cid.Undef, err
	}
	return mmap.Root()
}

// Root returns the root cid of the underlying HAMT.
func (mm *Multimap) Root() (cid.Cid, error) {
	return mm.mp.Root()
}

// Add adds a value for a key.
func (mm *Multimap) Add(key abi.Keyer, value cbor.Marshaler) error {
	// Load the array under key, or initialize a new empty one if not found.
	array, found, err := mm.Get(key)
	if err != nil {
		return err
	}
	if !found {
		array, err = MakeEmptyArray(mm.mp.store, mm.innerBitwidth)
		if err != nil {
			return err
		}
	}

	// Append to the array.
	if err = array.AppendContinuous(value); err != nil {
		return fmt.Errorf("failed to add multimap key %v value %v: %w", key, value, err)
	}

	c, err := array.Root()
	if err != nil {
		return fmt.Errorf("failed to flush child array: %w", err)
	}

	// Store the new array root under key.
	newArrayRoot := cbg.CborCid(c)
	err = mm.mp.Put(key, &newArrayRoot)
	if err != nil {
		return fmt.Errorf("failed to store multimap values: %w", err)
	}
	return nil
}

// RemoveAll removes all values for a key.
func (mm *Multimap) RemoveAll(key abi.Keyer) error {
	if _, err := mm.mp.TryDelete(key); err != nil {
		return fmt.Errorf("failed to delete multimap key %v root %v: %w", key, mm.mp.root, err)
	}
	return nil
}

// ForEach iterates all entries for a key in the order they were inserted, deserializing each value in turn into `out` and then
// calling a function.
// Iteration halts if the function returns an error.
// If the output parameter is nil, deserialization is skipped.
func (mm *Multimap) ForEach(key abi.Keyer, out cbor.Unmarshaler, fn func(i int64) error) error {
	array, found, err := mm.Get(key)
	if err != nil {
		return err
	}
	if found {
		return array.ForEach(out, fn)
	}
	return nil
}

// ForAll iterate all entries in map
func (mm *Multimap) ForAll(fn func(k string, arr *Array) error) error {
	var arrRoot cbg.CborCid
	if err := mm.mp.ForEach(&arrRoot, func(k string) error {
		arr, err := AsArray(mm.mp.store, cid.Cid(arrRoot), mm.innerBitwidth)
		if err != nil {
			return err
		}

		return fn(k, arr)
	}); err != nil {
		return err
	}

	return nil
}

// Get return entry for specify key
func (mm *Multimap) Get(key abi.Keyer) (*Array, bool, error) {
	var arrayRoot cbg.CborCid
	found, err := mm.mp.Get(key, &arrayRoot)
	if err != nil {
		return nil, false, fmt.Errorf("failed to load multimap key %v: %w", key, err)
	}
	var array *Array
	if found {
		array, err = AsArray(mm.mp.store, cid.Cid(arrayRoot), mm.innerBitwidth)
		if err != nil {
			return nil, false, fmt.Errorf("failed to load value %v as an array: %w", key, err)
		}
	}
	return array, found, nil
}
