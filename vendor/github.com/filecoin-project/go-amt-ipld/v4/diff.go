package amt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"golang.org/x/xerrors"

	"github.com/ipfs/go-cid"

	cbor "github.com/ipfs/go-ipld-cbor"
	cbg "github.com/whyrusleeping/cbor-gen"
)

// ChangeType denotes type of change in Change
type ChangeType int

// These constants define the changes that can be applied to a DAG.
const (
	Add ChangeType = iota
	Remove
	Modify
)

// Change represents a change to a DAG and contains a reference to the old and
// new CIDs.
type Change struct {
	Type   ChangeType
	Key    uint64
	Before *cbg.Deferred
	After  *cbg.Deferred
}

func (ch Change) String() string {
	b, _ := json.Marshal(ch)
	return string(b)
}

// Diff returns a set of changes that transform node 'a' into node 'b'. opts are applied to both prev and cur.
func Diff(ctx context.Context, prevBs, curBs cbor.IpldStore, prev, cur cid.Cid, opts ...Option) ([]*Change, error) {
	prevAmt, err := LoadAMT(ctx, prevBs, prev, opts...)
	if err != nil {
		return nil, xerrors.Errorf("loading previous root: %w", err)
	}

	prevCtx := &nodeContext{
		bs:       prevBs,
		bitWidth: prevAmt.bitWidth,
		height:   prevAmt.height,
	}

	curAmt, err := LoadAMT(ctx, curBs, cur, opts...)
	if err != nil {
		return nil, xerrors.Errorf("loading current root: %w", err)
	}

	// TODO: remove when https://github.com/filecoin-project/go-amt-ipld/issues/54 is closed.
	if curAmt.bitWidth != prevAmt.bitWidth {
		return nil, xerrors.Errorf("diffing AMTs with differing bitWidths not supported (prev=%d, cur=%d)", prevAmt.bitWidth, curAmt.bitWidth)
	}

	curCtx := &nodeContext{
		bs:       curBs,
		bitWidth: curAmt.bitWidth,
		height:   curAmt.height,
	}

	// edge case of diffing an empty AMT against non-empty
	if prevAmt.count == 0 && curAmt.count != 0 {
		return addAll(ctx, curCtx, curAmt.node, 0)
	}
	if prevAmt.count != 0 && curAmt.count == 0 {
		return removeAll(ctx, prevCtx, prevAmt.node, 0)
	}
	return diffNode(ctx, prevCtx, curCtx, prevAmt.node, curAmt.node, 0)
}

type nodeContext struct {
	bs       cbor.IpldStore // store containining AMT data
	bitWidth uint           // bit width of AMT
	height   int            // height of node
}

// nodesAtHeight returns the number of nodes that can be held at the context height
func (nc *nodeContext) nodesAtHeight() uint64 {
	return nodesForHeight(nc.bitWidth, nc.height)
}

func diffNode(ctx context.Context, prevCtx, curCtx *nodeContext, prev, cur *node, offset uint64) ([]*Change, error) {
	if prev == nil && cur == nil {
		return nil, nil
	}

	if prev == nil {
		return addAll(ctx, curCtx, cur, offset)
	}

	if cur == nil {
		return removeAll(ctx, prevCtx, prev, offset)
	}

	if prevCtx.height == 0 && curCtx.height == 0 {
		return diffLeaves(prev, cur, offset)
	}

	var changes []*Change

	if curCtx.height > prevCtx.height {
		subCount := curCtx.nodesAtHeight()
		for i, ln := range cur.links {
			if ln == nil || ln.cid == cid.Undef {
				continue
			}

			subCtx := &nodeContext{
				bs:       curCtx.bs,
				bitWidth: curCtx.bitWidth,
				height:   curCtx.height - 1,
			}

			subn, err := ln.load(ctx, subCtx.bs, subCtx.bitWidth, subCtx.height)
			if err != nil {
				return nil, err
			}

			offs := offset + (uint64(i) * subCount)
			if i == 0 {
				cs, err := diffNode(ctx, prevCtx, subCtx, prev, subn, offs)
				if err != nil {
					return nil, err
				}

				changes = append(changes, cs...)
			} else {
				cs, err := addAll(ctx, subCtx, subn, offs)
				if err != nil {
					return nil, err
				}

				changes = append(changes, cs...)
			}
		}

		return changes, nil
	}

	if prevCtx.height > curCtx.height {
		subCount := prevCtx.nodesAtHeight()
		for i, ln := range prev.links {
			if ln == nil || ln.cid == cid.Undef {
				continue
			}

			subCtx := &nodeContext{
				bs:       prevCtx.bs,
				bitWidth: prevCtx.bitWidth,
				height:   prevCtx.height - 1,
			}

			subn, err := ln.load(ctx, subCtx.bs, subCtx.bitWidth, subCtx.height)
			if err != nil {
				return nil, err
			}

			offs := offset + (uint64(i) * subCount)

			if i == 0 {
				cs, err := diffNode(ctx, subCtx, curCtx, subn, cur, offs)
				if err != nil {
					return nil, err
				}

				changes = append(changes, cs...)
			} else {
				cs, err := removeAll(ctx, subCtx, subn, offs)
				if err != nil {
					return nil, err
				}

				changes = append(changes, cs...)
			}
		}

		return changes, nil
	}

	// sanity check
	if prevCtx.height != curCtx.height {
		return nil, fmt.Errorf("comparing non-leaf nodes of unequal heights (%d, %d)", prevCtx.height, curCtx.height)
	}

	if len(prev.links) != len(cur.links) {
		return nil, fmt.Errorf("nodes have different numbers of links (prev=%d, cur=%d)", len(prev.links), len(cur.links))
	}

	if prev.links == nil || cur.links == nil {
		return nil, fmt.Errorf("nodes have no links")
	}

	subCount := prevCtx.nodesAtHeight()
	for i := range prev.links {
		// Neither previous or current links are in use
		if prev.links[i] == nil && cur.links[i] == nil {
			continue
		}

		// Previous had link, current did not
		if prev.links[i] != nil && cur.links[i] == nil {
			if prev.links[i].cid == cid.Undef {
				continue
			}

			subCtx := &nodeContext{
				bs:       prevCtx.bs,
				bitWidth: prevCtx.bitWidth,
				height:   prevCtx.height - 1,
			}

			subn, err := prev.links[i].load(ctx, subCtx.bs, subCtx.bitWidth, subCtx.height)
			if err != nil {
				return nil, err
			}

			offs := offset + (uint64(i) * subCount)
			cs, err := removeAll(ctx, subCtx, subn, offs)
			if err != nil {
				return nil, err
			}

			changes = append(changes, cs...)

			continue
		}

		// Current has link, previous did not
		if prev.links[i] == nil && cur.links[i] != nil {
			if cur.links[i].cid == cid.Undef {
				continue
			}

			subCtx := &nodeContext{
				bs:       curCtx.bs,
				bitWidth: curCtx.bitWidth,
				height:   curCtx.height - 1,
			}

			subn, err := cur.links[i].load(ctx, subCtx.bs, subCtx.bitWidth, subCtx.height)
			if err != nil {
				return nil, err
			}

			offs := offset + (uint64(i) * subCount)
			cs, err := addAll(ctx, subCtx, subn, offs)
			if err != nil {
				return nil, err
			}

			changes = append(changes, cs...)

			continue
		}

		// Both previous and current have links to diff

		prevSubCtx := &nodeContext{
			bs:       prevCtx.bs,
			bitWidth: prevCtx.bitWidth,
			height:   prevCtx.height - 1,
		}

		prevSubn, err := prev.links[i].load(ctx, prevSubCtx.bs, prevSubCtx.bitWidth, prevSubCtx.height)
		if err != nil {
			return nil, err
		}

		curSubCtx := &nodeContext{
			bs:       curCtx.bs,
			bitWidth: curCtx.bitWidth,
			height:   curCtx.height - 1,
		}

		curSubn, err := cur.links[i].load(ctx, curSubCtx.bs, curSubCtx.bitWidth, curSubCtx.height)
		if err != nil {
			return nil, err
		}

		offs := offset + (uint64(i) * subCount)

		cs, err := diffNode(ctx, prevSubCtx, curSubCtx, prevSubn, curSubn, offs)
		if err != nil {
			return nil, err
		}

		changes = append(changes, cs...)
	}

	return changes, nil
}

func addAll(ctx context.Context, nc *nodeContext, node *node, offset uint64) ([]*Change, error) {
	var changes []*Change
	err := node.forEachAt(ctx, nc.bs, nc.bitWidth, nc.height, 0, offset, func(index uint64, deferred *cbg.Deferred) error {
		changes = append(changes, &Change{
			Type:   Add,
			Key:    index,
			Before: nil,
			After:  deferred,
		})

		return nil
	})
	if err != nil {
		return nil, err
	}

	return changes, nil
}

func removeAll(ctx context.Context, nc *nodeContext, node *node, offset uint64) ([]*Change, error) {
	var changes []*Change

	err := node.forEachAt(ctx, nc.bs, nc.bitWidth, nc.height, 0, offset, func(index uint64, deferred *cbg.Deferred) error {
		changes = append(changes, &Change{
			Type:   Remove,
			Key:    index,
			Before: deferred,
			After:  nil,
		})

		return nil
	})
	if err != nil {
		return nil, err
	}

	return changes, nil
}

func diffLeaves(prev, cur *node, offset uint64) ([]*Change, error) {
	if len(prev.values) != len(cur.values) {
		return nil, fmt.Errorf("node leaves have different numbers of values (prev=%d, cur=%d)", len(prev.values), len(cur.values))
	}

	var changes []*Change
	for i, prevVal := range prev.values {
		index := offset + uint64(i)

		curVal := cur.values[i]
		if prevVal == nil && curVal == nil {
			continue
		}

		if prevVal == nil && curVal != nil {
			changes = append(changes, &Change{
				Type:   Add,
				Key:    index,
				Before: nil,
				After:  curVal,
			})

			continue
		}

		if prevVal != nil && curVal == nil {
			changes = append(changes, &Change{
				Type:   Remove,
				Key:    index,
				Before: prevVal,
				After:  nil,
			})

			continue
		}

		if !bytes.Equal(prevVal.Raw, curVal.Raw) {
			changes = append(changes, &Change{
				Type:   Modify,
				Key:    index,
				Before: prevVal,
				After:  curVal,
			})
		}

	}

	return changes, nil
}
