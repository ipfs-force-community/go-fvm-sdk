package adt

import (
	"bytes"
	"context"
	"fmt"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"

	"github.com/filecoin-project/go-state-types/cbor"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk"
	"github.com/ipfs/go-cid"
)

// IpldStore wraps a Blockstore and provides an interface for storing and retrieving CBOR encoded data.
// type definition from github.com/ipfs/go-ipld-cbor  link https://github.com/ipfs/go-ipld-cbor/blob/0aba2b8dbb9aa1f9c335b50ea1c1ff1646239e7b/store.go#L17
// go-ipld-cbor/github.com/polydawn/refmt library use reflect heavily, and tinygo unable to build it.
// so redefine the same interface here
type IpldStore interface {
	Get(ctx context.Context, c cid.Cid, out interface{}) error
	Put(ctx context.Context, v interface{}) (cid.Cid, error)
}

// Store define store with context
type Store interface {
	Context() context.Context
	IpldStore
}

// AdtStore Adapts a vanilla IPLD store as an ADT store.
func AdtStore(ctx context.Context) Store { //nolint
	return &fvmStore{
		ctx:       ctx,
		IpldStore: &fvmStore{},
	}
}

type fvmStore struct {
	ctx       context.Context
	IpldStore IpldStore
}

var _ Store = &fvmStore{}

func (r fvmStore) Context() context.Context {
	return r.ctx
}
func (r fvmStore) Get(ctx context.Context, c cid.Cid, out interface{}) error {
	data, err := sdk.Get(ctx, c)
	if err != nil {
		return err
	}
	unmarshalableObj, ok := out.(cbor.Unmarshaler)
	if !ok {
		return fmt.Errorf("ipld store get method must be Unmarshalable")
	}
	return unmarshalableObj.UnmarshalCBOR(bytes.NewBuffer(data))
}

func (r fvmStore) Put(ctx context.Context, in interface{}) (cid.Cid, error) {
	marshalableObj, ok := in.(cbor.Marshaler)
	if !ok {
		return cid.Undef, fmt.Errorf("ipld store put method must be marshalable")
	}
	buf := bytes.NewBuffer(nil)
	err := marshalableObj.MarshalCBOR(buf)
	if err != nil {
		return cid.Undef, fmt.Errorf("marshal object fail")
	}
	return sdk.Put(ctx, types.BLAKE2B256, types.BLAKE2BLEN, types.DAGCBOR, buf.Bytes())
}
