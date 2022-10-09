package migration

import (
	"context"

	"github.com/filecoin-project/go-state-types/builtin/v8/system"

	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
)

// System Actor migrator
type systemActorMigrator struct {
	OutCodeCID   cid.Cid
	ManifestData cid.Cid
}

func (m systemActorMigrator) migrateState(ctx context.Context, store cbor.IpldStore, in actorMigrationInput) (*actorMigrationResult, error) {
	// The ManifestData itself is already in the blockstore
	state := system.State{BuiltinActors: m.ManifestData}
	stateHead, err := store.Put(ctx, &state)
	if err != nil {
		return nil, err
	}

	return &actorMigrationResult{
		newCodeCID: m.OutCodeCID,
		newHead:    stateHead,
	}, nil
}
