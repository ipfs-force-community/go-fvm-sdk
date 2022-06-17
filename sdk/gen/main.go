package main

import (
	"log"

	"github.com/ipfs-force-community/go-fvm-sdk/gen"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
)

/*type genTarget struct {
	dir   string
	pkg   string
	types []interface{}
}*/

func main() {
	if err := gen.GenCborType("../types", "", types.AggregateSealVerifyInfo{},
		types.AggregateSealVerifyProofAndInfos{},
		types.ReplicaUpdateInfo{}); err != nil {
		log.Fatalf("gen for ../types: %s", err)
	}
}
