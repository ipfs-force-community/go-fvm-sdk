package amt

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"sort"

	cid "github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
	cbg "github.com/whyrusleeping/cbor-gen"

	"github.com/filecoin-project/go-amt-ipld/v4/internal"
)

// MaxIndex is the maximum index for elements in the AMT. This MaxUint64-1 so we
// don't overflow MaxUint64 when computing the length.
const MaxIndex = math.MaxUint64 - 1

// Root is described in more detail in its internal serialized form,
// internal.Root
type Root struct {
	bitWidth uint
	height   int
	count    uint64

	node *node

	store cbor.IpldStore
}

// NewAMT creates a new, empty AMT root with the given IpldStore and options.
func NewAMT(bs cbor.IpldStore, opts ...Option) (*Root, error) {
	cfg := defaultConfig()
	for _, opt := range opts {
		if err := opt(cfg); err != nil {
			return nil, err
		}
	}

	return &Root{
		bitWidth: cfg.bitWidth,
		store:    bs,
		node:     new(node),
	}, nil
}

// LoadAMT loads an existing AMT from the given IpldStore using the given
// root CID. An error will be returned where the AMT identified by the CID
// does not exist within the IpldStore. If the given options, or their defaults,
// do not match the AMT found at the given CID, an error will be returned.
func LoadAMT(ctx context.Context, bs cbor.IpldStore, c cid.Cid, opts ...Option) (*Root, error) {
	cfg := defaultConfig()
	for _, opt := range opts {
		if err := opt(cfg); err != nil {
			return nil, err
		}
	}

	var r internal.Root
	if err := bs.Get(ctx, c, &r); err != nil {
		return nil, err
	}

	// Check the bitwidth but don't rely on it. We may add an option in the
	// future to just discover the bitwidth from the AMT, but we need to be
	// careful to not just trust the value.
	if r.BitWidth != uint64(cfg.bitWidth) {
		return nil, fmt.Errorf("expected bitwidth %d but AMT has bitwidth %d", cfg.bitWidth, r.BitWidth)
	}

	// Make sure the height is sane to prevent any integer overflows later
	// (e.g., height+1). While MaxUint64-1 would solve the "+1" issue, we
	// might as well use 64 because the height cannot be greater than 62
	// (min width = 2, 2**64 == max elements).
	if r.Height > 64 {
		return nil, fmt.Errorf("height greater than 64: %d", r.Height)
	}

	maxNodes := nodesForHeight(cfg.bitWidth, int(r.Height+1))
	// nodesForHeight saturates. If "max nodes" is max uint64, the maximum
	// number of nodes at the previous level muss be less. This is the
	// simplest way to check to see if the height is sane.
	if maxNodes == math.MaxUint64 && nodesForHeight(cfg.bitWidth, int(r.Height)) == math.MaxUint64 {
		return nil, fmt.Errorf("failed to load AMT: height %d out of bounds", r.Height)
	}

	// If max nodes is less than the count, something is wrong.
	if maxNodes < r.Count {
		return nil, fmt.Errorf(
			"failed to load AMT: not tall enough (%d) for count (%d)", r.Height, r.Count,
		)
	}

	nd, err := newNode(r.Node, cfg.bitWidth, r.Height == 0, r.Height == 0)
	if err != nil {
		return nil, err
	}

	return &Root{
		bitWidth: cfg.bitWidth,
		height:   int(r.Height),
		count:    r.Count,
		node:     nd,
		store:    bs,
	}, nil
}

// FromArray creates a new AMT and performs a BatchSet on it using the vals and
// options provided. Indexes from the array are used as the indexes for the same
// values in the AMT.
func FromArray(ctx context.Context, bs cbor.IpldStore, vals []cbg.CBORMarshaler, opts ...Option) (cid.Cid, error) {
	r, err := NewAMT(bs, opts...)
	if err != nil {
		return cid.Undef, err
	}
	if err := r.BatchSet(ctx, vals); err != nil {
		return cid.Undef, err
	}

	return r.Flush(ctx)
}

// Set will add or update entry at index i with value val. The index must be
// within lower than MaxIndex.
//
// Where val has a compatible CBORMarshaler() it will be used to serialize the
// object into CBOR. Otherwise the generic go-ipld-cbor DumbObject() will be
// used.
//
// Setting a new index that is greater than the current capacity of the
// existing AMT structure will result in the creation of additional nodes to
// form a structure of enough height to contain the new index.
//
// The height required to store any given index can be calculated by finding
// the lowest (width^(height+1) - 1) that is higher than the index. For example,
// a height of 1 on an AMT with a width of 8 (bitWidth of 3) can fit up to
// indexes of 8^2 - 1, or 63. At height 2, indexes up to 511 can be stored. So a
// Set operation for an index between 64 and 511 will require that the AMT have
// a height of at least 3. Where an AMT has a height less than 3, additional
// nodes will be added until the height is 3.
func (r *Root) Set(ctx context.Context, i uint64, val cbg.CBORMarshaler) error {
	if i > MaxIndex {
		return fmt.Errorf("index %d is out of range for the amt", i)
	}

	var d cbg.Deferred
	if val == nil {
		d.Raw = cbg.CborNull
	} else {
		valueBuf := new(bytes.Buffer)
		if err := val.MarshalCBOR(valueBuf); err != nil {
			return err
		}
		d.Raw = valueBuf.Bytes()
	}

	// where the index is greater than the number of elements we can fit into the
	// current AMT, grow it until it will fit.
	for i >= nodesForHeight(r.bitWidth, r.height+1) {
		// if we have existing data, perform the re-height here by pushing down
		// the existing tree into the left-most portion of a new root
		if !r.node.empty() {
			nd := r.node
			// since all our current elements fit in the old height, we _know_ that
			// they will all sit under element [0] of this new node.
			r.node = &node{links: make([]*link, 1<<r.bitWidth)}
			r.node.links[0] = &link{
				dirty:  true,
				cached: nd,
			}
		}
		// else we still need to add new nodes to form the right height, but we can
		// defer that to our set() call below which will lazily create new nodes
		// where it expects there to be some
		r.height++
	}

	addVal, err := r.node.set(ctx, r.store, r.bitWidth, r.height, i, &d)
	if err != nil {
		return err
	}

	if addVal {
		// Something is wrong, so we'll just do our best to not overflow.
		if r.count >= (MaxIndex - 1) {
			return errInvalidCount
		}
		r.count++
	}

	return nil
}

// BatchSet takes an array of vals and performs a Set on each of them on an
// existing AMT. Indexes from the array are used as indexes for the same values
// in the AMT.
//
// This is currently a convenience method and does not perform optimizations
// above iterative Set calls for each entry.
func (r *Root) BatchSet(ctx context.Context, vals []cbg.CBORMarshaler) error {
	// TODO: there are more optimized ways of doing this method
	for i, v := range vals {
		if err := r.Set(ctx, uint64(i), v); err != nil {
			return err
		}
	}
	return nil
}

// Get retrieves a value from index i.
// If the index is set, returns true and, if the `out` parameter is not nil,
// deserializes the value into that interface. Returns false if the index is not set.
func (r *Root) Get(ctx context.Context, i uint64, out cbg.CBORUnmarshaler) (bool, error) {
	if i > MaxIndex {
		return false, fmt.Errorf("index %d is out of range for the amt", i)
	}

	// easy shortcut case, index is too large for our height, don't bother looking
	// further
	if i >= nodesForHeight(r.bitWidth, r.height+1) {
		return false, nil
	}
	return r.node.get(ctx, r.store, r.bitWidth, r.height, i, out)
}

// BatchDelete performs a bulk Delete operation on an array of indices. Each
// index in the given indices array will be removed from the AMT, if it is present.
// If `strict` is true, all indices are expected to be present, and this will return an error
// if one is not found.
//
// Returns true if the AMT was modified as a result of this operation.
//
// There is no special optimization applied to this method, it is simply a
// convenience wrapper around Delete for an array of indices.
func (r *Root) BatchDelete(ctx context.Context, indices []uint64, strict bool) (modified bool, err error) {
	// TODO: theres a faster way of doing this, but this works for now

	// Sort by index so we can safely implement these optimizations in the future.
	less := func(i, j int) bool { return indices[i] < indices[j] }
	if !sort.SliceIsSorted(indices, less) {
		// Copy first so we don't modify our inputs.
		indices = append(indices[0:0:0], indices...)
		sort.Slice(indices, less)
	}

	for _, i := range indices {
		found, err := r.Delete(ctx, i)
		if err != nil {
			return false, err
		} else if strict && !found {
			return false, fmt.Errorf("no such index %d", i)
		}
		modified = modified || found
	}
	return modified, nil
}

// Delete removes an index from the AMT.
// Returns true if the index was present and removed, or false if the index
// was not set.
//
// If this delete operation leaves nodes with no remaining elements, the height
// will be reduced to fit the maximum remaining index, leaving the AMT in
// canonical form for the given set of data that it contains.
func (r *Root) Delete(ctx context.Context, i uint64) (bool, error) {
	if i > MaxIndex {
		return false, fmt.Errorf("index %d is out of range for the amt", i)
	}

	// shortcut, index is greater than what we hold so we know it's not there
	if i >= nodesForHeight(r.bitWidth, r.height+1) {
		return false, nil
	}

	found, err := r.node.delete(ctx, r.store, r.bitWidth, r.height, i)
	if err != nil {
		return false, err
	} else if !found {
		return false, nil
	}

	// The AMT invariant dictates that for any non-empty AMT, the root node must
	// not address only its left-most child node. Where a deletion has created a
	// state where the current root node only consists of a link to the left-most
	// child and no others, that child node must become the new root node (i.e.
	// the height is reduced by 1). We perform the same check on the new root node
	// such that we reduce the AMT to canonical form for this data set.
	// In the extreme case, it is possible to perform a collapse from a large
	// `height` to height=0 where the index being removed is very large and there
	// remains no other indexes or the remaining indexes are in the range of 0 to
	// bitWidth^8.
	// See node.collapse() for more notes.
	newHeight, err := r.node.collapse(ctx, r.store, r.bitWidth, r.height)
	if err != nil {
		return false, err
	}
	r.height = newHeight

	// Something is very wrong but there's not much we can do. So we perform
	// the operation and then tell the user that something is wrong.
	if r.count == 0 {
		return false, errInvalidCount
	}

	r.count--
	return true, nil
}

// ForEach iterates over the entire AMT and calls the cb function for each
// entry found in the leaf nodes. The callback will receive the index and the
// value of each element.
func (r *Root) ForEach(ctx context.Context, cb func(uint64, *cbg.Deferred) error) error {
	return r.node.forEachAt(ctx, r.store, r.bitWidth, r.height, 0, 0, cb)
}

// ForEachAt iterates over the AMT beginning from the given start index. See
// ForEach for more details.
func (r *Root) ForEachAt(ctx context.Context, start uint64, cb func(uint64, *cbg.Deferred) error) error {
	return r.node.forEachAt(ctx, r.store, r.bitWidth, r.height, start, 0, cb)
}

// FirstSetIndex finds the lowest index in this AMT that has a value set for
// it. If this operation is called on an empty AMT, an ErrNoValues will be
// returned.
func (r *Root) FirstSetIndex(ctx context.Context) (uint64, error) {
	return r.node.firstSetIndex(ctx, r.store, r.bitWidth, r.height)
}

// Flush saves any unsaved node data and recompacts the in-memory forms of each
// node where they have been expanded for operational use.
func (r *Root) Flush(ctx context.Context) (cid.Cid, error) {
	nd, err := r.node.flush(ctx, r.store, r.bitWidth, r.height)
	if err != nil {
		return cid.Undef, err
	}
	root := internal.Root{
		BitWidth: uint64(r.bitWidth),
		Height:   uint64(r.height),
		Count:    r.count,
		Node:     *nd,
	}
	return r.store.Put(ctx, &root)
}

// Len returns the "Count" property that is stored in the root of this AMT.
// It's correctness is only guaranteed by the consistency of the build of the
// AMT (i.e. this code). A "secure" count would require iterating the entire
// tree, but if all nodes are part of a trusted structure (e.g. one where we
// control the entire build, or verify all incoming blocks from untrusted
// sources) then we ought to be able to say "count" is correct.
func (r *Root) Len() uint64 {
	return r.count
}
