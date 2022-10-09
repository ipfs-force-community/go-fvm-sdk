package internal

import (
	cid "github.com/ipfs/go-cid"
	cbg "github.com/whyrusleeping/cbor-gen"
)

// Node represents any node within the AMT, including the root, intermediate
// and leaf nodes. For the minimal case of AMT, there may be a single Node
// containing all data. As the highest-index grows, more intermediate nodes
// are added.
//
// A Node will strictly either be a leaf node (height 0) or a non-leaf (root
// or intermediate, height 1). Leaf nodes contain an array of one or more
// Values, where non-leaf nodes contain an array of one or more Links to child
// nodes.
//
// The Bmap (bitmap) has the same number of bits as the "width" of the AMT
// (bitWidth^2), where each bit in the bitmap indicates the presence (1) or
// absence (0) of a value or link to a child node. In this way, the serialized
// form, and in-memory form of a Node contains only the value or links present.
// There must be at least one value for height=0 nodes and at least one link for
// height>0 nodes. Nodes with no links or values are invalid and the AMT will
// not have canonical form.
//
// Each node is serialized in the following form, described as an IPLD Schema:
//
//	type Node struct {
//		bmap Bytes
//		links [&Node]
//		values [Any]
//	} representation tuple
//
// Where bmap is strictly a byte array of length (bitWidth^2)/8 and the links
// and values arrays are between zero and the width of this AMT (bitWidth^2).
// One of links or values arrays must contain zero elements and one must contain
// at least one element since a node is strictly either a leaf or a non-leaf.
type Node struct {
	Bmap   []byte
	Links  []cid.Cid
	Values []*cbg.Deferred
}

// Root is the single entry point for this AMT. It is serialized with an inner
// root Node element.
//
// The bitWidth property dictates the number of bits used to generate an index
// at each level from the addressible index supplied by the user.
//
// The height property is essential for understanding how deep to navigate to
// value-holding leaf nodes and therefore how many bits of an index will be
// required for navigation.
//
// The count property is maintained during ongoing mutation of the AMT and can
// be used as a fast indicator of the size of the structure. It is assumed to
// be correct if the nodes of the AMT were part of a trusted construction or
// have been verified. It is not essential to the construction or navigation of
// the AMT but is helpful for fast Len() calls.
// Performing a secure count would require navigating through all leaf nodes
// and adding up the number of occupied slots.
//
// The root is serialized in the following form, described as an IPLD Schema:
//
// 	type Root struct {
//    bitWidth Int
//		height Int
//		count Int
//		node Node
//	} representation tuple
//
// Where bitWidth, height and count are unsigned integers and Node is the
// initial root node, see below.
type Root struct {
	BitWidth uint64
	Height   uint64
	Count    uint64
	Node     Node
}
