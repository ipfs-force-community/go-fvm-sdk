package main

import (
	"erc20/contract"
	"fmt"
	"log"
	"path/filepath"
	"reflect"

	"github.com/filecoin-project/go-state-types/cbor"

	gen "github.com/whyrusleeping/cbor-gen"
)

var unMarshallerT = reflect.TypeOf((*cbor.Unmarshaler)(nil)).Elem()
var errorT = reflect.TypeOf((*error)(nil)).Elem()
var marshallerT = reflect.TypeOf((*cbor.Marshaler)(nil)).Elem()

type genTarget struct {
	dir   string
	pkg   string
	types []interface{}
}

func main() {
	gen_cbor_type()
	stateT := reflect.TypeOf(contract.Erc20Token{})
	err := gen_entry(stateT, "../entry.go")
	if err != nil {
		fmt.Println(err)
	}
	err = genContractClient(stateT, "../client/client.go")
	if err != nil {
		fmt.Println(err)
	}
}

func gen_cbor_type() {
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
				contract.FakeSetBalance{},
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
