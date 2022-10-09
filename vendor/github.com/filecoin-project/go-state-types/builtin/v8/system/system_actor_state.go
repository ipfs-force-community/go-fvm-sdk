package system

import (
	"context"

	"github.com/filecoin-project/go-state-types/manifest"

	"github.com/filecoin-project/go-state-types/builtin/v8/util/adt"

	"github.com/ipfs/go-cid"
	"golang.org/x/xerrors"
)

type State struct {
	BuiltinActors cid.Cid // ManifestData
}

func ConstructState(store adt.Store) (*State, error) {
	empty, err := store.Put(context.TODO(), &manifest.ManifestData{})
	if err != nil {
		return nil, xerrors.Errorf("failed to create empty manifest: %w", err)
	}

	return &State{BuiltinActors: empty}, nil
}
