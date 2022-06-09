package main

import (
	"log"
	"path/filepath"

	"erc20/contract"

	gen "github.com/whyrusleeping/cbor-gen"
)

type genTarget struct {
	dir   string
	pkg   string
	types []interface{}
}

func main() {
	targets := []genTarget{
		{
			dir: "../contract",
			types: []interface{}{
				contract.Erc20Token{},
				contract.ConstructorReq{},
				contract.TransferReq{},
				contract.AllowanceReq{},
				contract.TransferFromReq{},
				contract.ApprovalReq{},
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
