package amt

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
	cbg "github.com/whyrusleeping/cbor-gen"

	"github.com/filecoin-project/go-amt-ipld/v4/internal"
)

// node is described in more detail in its internal serialized form,
// internal.Node. This form contains a fully expanded form of internal.Node
// where the Bmap is used to expand the contracted form of either Values (leaf)
// or Links (non-leaf) for ease of addressing.
// Both properties may be nil if the node is empty (a root node).
type node struct {
	links  []*link
	values []*cbg.Deferred
}

var (
	errEmptyNode      = errors.New("unexpected empty amt node")
	errUndefinedCID   = errors.New("amt node has undefined CID")
	errLinksAndValues = errors.New("amt node has both links and values")
	errLeafUnexpected = errors.New("amt leaf not expected at height")
	errLeafExpected   = errors.New("amt expected at height")
	errInvalidCount   = errors.New("amt count does not match number of elements")
)

// the number of bytes required such that there is a single bit for each element
// in the links or value array. This is (bitWidth^2)/8.
func bmapBytes(bitWidth uint) int {
	if bitWidth <= 3 {
		return 1
	}
	return 1 << (bitWidth - 3)
}

// Create a new from a serialized form. This operation takes an internal.Node
// and returns a node. internal.Node uses bitmap compaction of links or values
// arrays, while node uses the expanded form. This method performs the expansion
// such that we can use simple addressing of this node's child elements.
func newNode(nd internal.Node, bitWidth uint, allowEmpty, expectLeaf bool) (*node, error) {
	if len(nd.Links) > 0 && len(nd.Values) > 0 {
		// malformed AMT, a node cannot be both leaf and non-leaf
		return nil, errLinksAndValues
	}

	// strictly require the bitmap to be the correct size for the given bitWidth
	if expWidth := bmapBytes(bitWidth); expWidth != len(nd.Bmap) {
		return nil, fmt.Errorf(
			"expected bitfield to be %d bytes long, found bitfield with %d bytes",
			expWidth, len(nd.Bmap),
		)
	}

	width := uint(1 << bitWidth)
	i := 0
	n := new(node)
	if len(nd.Values) > 0 { // leaf node, height=0
		if !expectLeaf {
			return nil, errLeafUnexpected
		}
		n.values = make([]*cbg.Deferred, width)
		for x := uint(0); x < width; x++ {
			// check if this value exists in the bitmap, pull it out of the compacted
			// list if it does
			if nd.Bmap[x/8]&(1<<(x%8)) > 0 {
				if i >= len(nd.Values) {
					// too many bits were set in the bitmap for the number of values
					// available
					return nil, fmt.Errorf("expected at least %d values, found %d", i+1, len(nd.Values))
				}
				n.values[x] = nd.Values[i]
				i++
			}
		}
		if i != len(nd.Values) {
			// the number of bits set in the bitmap was not the same as the number of
			// values in the array
			return nil, fmt.Errorf("expected %d values, got %d", i, len(nd.Values))
		}
	} else if len(nd.Links) > 0 {
		// non-leaf node, height>0
		if expectLeaf {
			return nil, errLeafExpected
		}

		n.links = make([]*link, width)
		for x := uint(0); x < width; x++ {
			// check if this child link exists in the bitmap, pull it out of the
			// compacted list if it does
			if nd.Bmap[x/8]&(1<<(x%8)) > 0 {
				if i >= len(nd.Links) {
					// too many bits were set in the bitmap for the number of values
					// available
					return nil, fmt.Errorf("expected at least %d links, found %d", i+1, len(nd.Links))
				}
				c := nd.Links[i]
				if !c.Defined() {
					return nil, errUndefinedCID
				}
				// TODO: check link hash function.
				prefix := c.Prefix()
				if prefix.Codec != cid.DagCBOR {
					return nil, fmt.Errorf("internal amt nodes must be cbor, found %d", prefix.Codec)
				}
				n.links[x] = &link{cid: c}
				i++
			}
		}
		if i != len(nd.Links) {
			// the number of bits set in the bitmap was not the same as the number of
			// values in the array
			return nil, fmt.Errorf("expected %d links, got %d", i, len(nd.Links))
		}
	} else if !allowEmpty { // only THE empty AMT case can allow this
		return nil, errEmptyNode
	}
	return n, nil
}

// collapse occurs when we only have the single child node. If this is the case
// we need to reduce height by one. Continue down the tree, reducing height
// until we're either at a single height=0 node or we have something other than
// a single child node.
func (nd *node) collapse(ctx context.Context, bs cbor.IpldStore, bitWidth uint, height int) (int, error) {
	// No links at all?
	if nd.links == nil {
		return 0, nil
	}

	// If we have any links going "to the right", we can't collapse any
	// more.
	for _, l := range nd.links[1:] {
		if l != nil {
			return height, nil
		}
	}

	// If we have _no_ links, we've collapsed everything.
	if nd.links[0] == nil {
		return 0, nil
	}

	// only one child, collapse it.

	subn, err := nd.links[0].load(ctx, bs, bitWidth, height-1)
	if err != nil {
		return 0, err
	}

	// Collapse recursively.
	newHeight, err := subn.collapse(ctx, bs, bitWidth, height-1)
	if err != nil {
		return 0, err
	}

	*nd = *subn

	return newHeight, nil
}

// does this node contain any child nodes or values?
func (nd *node) empty() bool {
	for _, l := range nd.links {
		if l != nil {
			return false
		}
	}
	for _, v := range nd.values {
		if v != nil {
			return false
		}
	}
	return true
}

// Recursive get() called through the tree in order to retrieve values from
// leaf nodes. We start at the root and navigate until height=0 where the
// entries themselves should exist. At any point in the navigation we can
// assert that a value does not exist in this AMT if an expected intermediate
// doesn't exist, so we don't need to do full height traversal for many cases
// where we don't have that index.
func (n *node) get(ctx context.Context, bs cbor.IpldStore, bitWidth uint, height int, i uint64, out cbg.CBORUnmarshaler) (bool, error) {
	// height=0 means we're operating on a leaf node where the entries themselves
	// are stores, we have a `set` so it must exist if the node is correctly
	// formed
	if height == 0 {
		d := n.getValue(i)
		found := d != nil
		var err error
		if found && out != nil {
			err = out.UnmarshalCBOR(bytes.NewReader(d.Raw))
		}
		return found, err
	}

	// Non-leaf case where we need to navigate further down toward the correct
	// leaf by consuming some of the provided index to form the index at this
	// height and passing the remainder down.
	// The calculation performed is to divide the addressible indexes of each
	// child node such that each child has the ability to contain that range of
	// indexes somewhere in its graph. e.g. at height=1 for bitWidth=3, the total
	// addressible index space we can contain is in the range of 0 to
	// `(bitWidth^2) ^ (height+1) = 8^2 = 64`. Where each child node can contain
	// 64/8 of indexes. This is true regardless of the position in the overall
	// AMT and original index from the Get() operation because we modify the index
	// before passing it to lower nodes to remove the bits relevant to higher
	// addressing. e.g. at height=1, a call to any child's get() will receive an
	// index in the range of 0 to bitWidth^2.
	nfh := nodesForHeight(bitWidth, height)
	ln := n.getLink(i / nfh)
	if ln == nil {
		// This can occur at any point in the traversal, not just height=0, it just
		// means that the higher up it occurs that a larger range of indexes in this
		// region don't exist.
		return false, nil
	}
	subn, err := ln.load(ctx, bs, bitWidth, height-1)
	if err != nil {
		return false, err
	}

	// `i%nfh` discards index information for this height so the child only gets
	// the part of the index that is relevant for it.
	// e.g. get(50) at height=1 for width=8 would be 50%8=2, i.e. the child will
	// be asked to get(2) and it will have leaf nodes (because it's height=0) so
	// the actual value will be at index=2 of its values array.
	return subn.get(ctx, bs, bitWidth, height-1, i%nfh, out)
}

// Recursively handle a delete through the tree, navigating down in the same
// way as is documented in get().
func (n *node) delete(ctx context.Context, bs cbor.IpldStore, bitWidth uint, height int, i uint64) (bool, error) {
	// at the leaf node where the value is, expand out the values array and
	// zero out the value and bit in the bitmap to indicate its deletion
	if height == 0 {
		if n.getValue(i) == nil {
			return false, nil
		}

		n.setValue(bitWidth, i, nil)
		return true, nil
	}

	// see get() documentation on how nfh and subi describes the index at this
	// height
	nfh := nodesForHeight(bitWidth, height)
	subi := i / nfh

	ln := n.getLink(subi)
	if ln == nil {
		return false, nil
	}

	// we're at a non-leaf node, so navigate down to the appropriate child and
	// continue
	subn, err := ln.load(ctx, bs, bitWidth, height-1)
	if err != nil {
		return false, err
	}

	// see get() documentation for how the i%... calculation trims the index down
	// to only that which is applicable for the height below
	if deleted, err := subn.delete(ctx, bs, bitWidth, height-1, i%nfh); err != nil {
		return false, err
	} else if !deleted {
		return false, nil
	}

	// if the child node we just deleted from now has no children or elements of
	// its own, we need to zero it out in this node. This compaction process may
	// recursively chain back up through the calling nodes, removing more than
	// one node in total for this delete operation (i.e. where an index contains
	// the only entry on a particular branch of the tree).
	if subn.empty() {
		n.setLink(bitWidth, subi, nil)
	} else {
		ln.dirty = true
	}

	return true, nil
}

// Recursive implementation backing ForEach and ForEachAt. Performs a
// depth-first walk of the tree, beginning at the 'start' index. The 'offset'
// argument helps us locate the lateral position of the current node so we can
// figure out the appropriate 'index', since indexes are not stored with values
// and can only be determined by knowing how far a leaf node is removed from
// the left-most leaf node.
func (n *node) forEachAt(ctx context.Context, bs cbor.IpldStore, bitWidth uint, height int, start, offset uint64, cb func(uint64, *cbg.Deferred) error) error {
	if height == 0 {
		// height=0 means we're at leaf nodes and get to use our callback
		for i, v := range n.values {
			if v != nil {
				ix := offset + uint64(i)
				if ix < start {
					// if we're here, 'start' is probably somewhere in the
					// middle of this node's elements
					continue
				}

				// use 'offset' to determine the actual index for this element, it
				// tells us how distant we are from the left-most leaf node
				if err := cb(offset+uint64(i), v); err != nil {
					return err
				}
			}
		}

		return nil
	}

	subCount := nodesForHeight(bitWidth, height)
	for i, ln := range n.links {
		if ln == nil {
			continue
		}

		// 'offs' tells us the index of the left-most element of the subtree defined
		// by 'sub'
		offs := offset + (uint64(i) * subCount)
		nextOffs := offs + subCount
		// nextOffs > offs checks for overflow at MaxIndex (where the next offset wraps back
		// to 0).
		if nextOffs >= offs && start >= nextOffs {
			// if we're here, 'start' lets us skip this entire sub-tree
			continue
		}

		subn, err := ln.load(ctx, bs, bitWidth, height-1)
		if err != nil {
			return err
		}

		// recurse into the child node, providing 'offs' to tell it where it's
		// located in the tree
		if err := subn.forEachAt(ctx, bs, bitWidth, height-1, start, offs, cb); err != nil {
			return err
		}
	}
	return nil
}

var errNoVals = fmt.Errorf("no values")

// Recursive implementation of FirstSetIndex that's performed on the left-most
// nodes of the tree down to the leaf. In order to return a correct index, we
// need to accumulate the appropriate number of spaces to the left of the
// left-most that exist at each level, taking into account the number of
// blank leaf-entry positions that exist.
func (n *node) firstSetIndex(ctx context.Context, bs cbor.IpldStore, bitWidth uint, height int) (uint64, error) {
	if height == 0 {
		for i, v := range n.values {
			if v != nil {
				// returning 'i' here is a local index (0<=i<width), which isn't the
				// actual index unless this is a single-node, height=0 AMT.
				return uint64(i), nil
			}
		}
		// if we're here, we're either dealing with a malformed AMT or an empty AMT
		return 0, errNoVals
	}

	// we're dealing with a non-leaf node

	for i, ln := range n.links {
		if ln == nil {
			// nothing here.
			continue
		}
		subn, err := ln.load(ctx, bs, bitWidth, height-1)
		if err != nil {
			return 0, err
		}
		ix, err := subn.firstSetIndex(ctx, bs, bitWidth, height-1)
		if err != nil {
			return 0, err
		}

		// 'ix' is the child's understanding of it's left-most set index, we have
		// to add to it the number of _gaps_ that are present on the left of
		// the child node's position. So if the child node is index (i) 0 then
		// it's the left-most and i*subCount will be 0. But if it's 1, subCount
		// will account for an entire missing branch to the left in position 0.
		// This operation continues as we reverse back through the call stack
		// building up to the correct index.
		subCount := nodesForHeight(bitWidth, height)
		return ix + (uint64(i) * subCount), nil
	}

	return 0, errNoVals
}

// Recursive implementation of the set operation that calls through child nodes
// down into the appropriate leaf node for the given index. The index 'i' is
// relative to this current node, so must be adjusted as we recurse down
// through the tree. The same operation is used for get, see the documentation
// there for how the index is calculated for each height and adjusted as we
// move down.
// Returns a bool that indicates whether a new value was added or an existing
// one was overwritten. This is useful for adjusting the Count in the root node
// when we reverse back out of the calls.
func (n *node) set(ctx context.Context, bs cbor.IpldStore, bitWidth uint, height int, i uint64, val *cbg.Deferred) (bool, error) {
	if height == 0 {
		// we're at the leaf, we can either overwrite the value that already exists
		// or set a new one if there is none
		alreadySet := n.getValue(i) != nil
		n.setValue(bitWidth, i, val)
		return !alreadySet, nil
	}

	// see get() documentation on how nfh and subi describes the index at this
	// height
	nfh := nodesForHeight(bitWidth, height)

	// Load but don't mark dirty or actually link in any _new_ intermediate
	// nodes. We'll do that on return if nothing goes wrong.
	ln := n.getLink(i / nfh)
	if ln == nil {
		ln = &link{cached: new(node)}
	}
	subn, err := ln.load(ctx, bs, bitWidth, height-1)
	if err != nil {
		return false, err
	}

	// see get() documentation for how the i%... calculation trims the index down
	// to only that which is applicable for the height below
	nodeAdded, err := subn.set(ctx, bs, bitWidth, height-1, i%nfh, val)
	if err != nil {
		return false, err
	}

	// Make all modifications on the way back up if there was no error.
	ln.dirty = true // only mark dirty on success.
	n.setLink(bitWidth, i/nfh, ln)

	return nodeAdded, nil
}

// flush is the per-node form of Flush() that operates on each node, and calls
// flush() on each child node. It generates the serialized form of this node,
// which includes the bitmap and compacted links or values array.
func (n *node) flush(ctx context.Context, bs cbor.IpldStore, bitWidth uint, height int) (*internal.Node, error) {
	nd := new(internal.Node)
	nd.Bmap = make([]byte, bmapBytes(bitWidth))

	if height == 0 {
		// leaf node, we're storing values in this node
		for i, val := range n.values {
			if val == nil {
				continue
			}
			nd.Values = append(nd.Values, val)
			// set the bit in the bitmap for this position to indicate its presence
			nd.Bmap[i/8] |= 1 << (uint(i) % 8)
		}
		return nd, nil
	}

	// non-leaf node, we're only storing Links in this node
	for i, ln := range n.links {
		if ln == nil {
			continue
		}
		if ln.dirty {
			if ln.cached == nil {
				return nil, fmt.Errorf("expected dirty node to be cached")
			}
			subn, err := ln.cached.flush(ctx, bs, bitWidth, height-1)
			if err != nil {
				return nil, err
			}
			c, err := bs.Put(ctx, subn)
			if err != nil {
				return nil, err
			}

			ln.cid = c
			ln.dirty = false
		}
		nd.Links = append(nd.Links, ln.cid)
		// set the bit in the bitmap for this position to indicate its presence
		nd.Bmap[i/8] |= 1 << (uint(i) % 8)
	}

	return nd, nil
}

func (n *node) setLink(bitWidth uint, i uint64, l *link) {
	if n.links == nil {
		if l == nil {
			return
		}
		n.links = make([]*link, 1<<bitWidth)
	}
	n.links[i] = l
}

func (n *node) getLink(i uint64) *link {
	if n.links == nil {
		return nil
	}
	return n.links[i]
}

func (n *node) setValue(bitWidth uint, i uint64, v *cbg.Deferred) {
	if n.values == nil {
		if v == nil {
			return
		}
		n.values = make([]*cbg.Deferred, 1<<bitWidth)
	}
	n.values[i] = v
}

func (n *node) getValue(i uint64) *cbg.Deferred {
	if n.values == nil {
		return nil
	}
	return n.values[i]
}
