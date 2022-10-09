package amt

import (
	"context"

	"github.com/filecoin-project/go-amt-ipld/v4/internal"
	cid "github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
)

type link struct {
	cid cid.Cid

	cached *node
	dirty  bool
}

func (l *link) load(ctx context.Context, bs cbor.IpldStore, bitWidth uint, height int) (*node, error) {
	if l.cached == nil {
		var nd internal.Node
		if err := bs.Get(ctx, l.cid, &nd); err != nil {
			return nil, err
		}

		n, err := newNode(nd, bitWidth, false, height == 0)
		if err != nil {
			return nil, err
		}
		l.cached = n
	}
	return l.cached, nil
}
