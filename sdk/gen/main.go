package main

import (
	"log"
	"path/filepath"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	gen "github.com/whyrusleeping/cbor-gen"
)

// todo unable to generate State cbor if use state directly, sys call broke build
type State struct {
	Count uint64
}

type genTarget struct {
	dir   string
	pkg   string
	types []interface{}
}

func main() {
	targets := []genTarget{
		{
			dir: "./types",
			types: []interface{}{
				types.AggregateSealVerifyInfo{},
				types.AggregateSealVerifyProofAndInfos{},
				types.ReplicaUpdateInfo{},
			},
		},
	}

	for _, target := range targets {
		pkg := target.pkg
		if pkg == "" {
			pkg = filepath.Base(target.dir)
		}

		if err := gen.WriteTupleEncodersToFile(filepath.Join(target.dir, "cbor_gen.go"), pkg, target.types...); err != nil {
			log.Fatalf("gen for %s: %s", target.dir, err)
		}
	}
}
