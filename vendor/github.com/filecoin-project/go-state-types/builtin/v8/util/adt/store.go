package adt

import (
	"context"

	ipldcbor "github.com/ipfs/go-ipld-cbor"
)

// Store defines an interface required to back the ADTs in this package.
type Store interface {
	Context() context.Context
	ipldcbor.IpldStore
}

// Adapts a vanilla IPLD store as an ADT store.
func WrapStore(ctx context.Context, store ipldcbor.IpldStore) Store {
	return &wstore{
		ctx:       ctx,
		IpldStore: store,
	}
}

type wstore struct {
	ctx context.Context
	ipldcbor.IpldStore
}

var _ Store = &wstore{}

func (s *wstore) Context() context.Context {
	return s.ctx
}
