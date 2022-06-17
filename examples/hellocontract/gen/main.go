package main

import (
	"fmt"
	"hellocontract/contract"
	"log"
	"reflect"

	"github.com/ipfs-force-community/go-fvm-sdk/gen"
)

func main() {
	if err := gen.GenCborType("../contract", "", contract.State{}); err != nil {
		log.Fatalf("gen for ../contract: %s", err)
	}
	stateT := reflect.TypeOf(contract.State{})
	err := gen.GenEntry(stateT, "../entry.go")
	if err != nil {
		log.Fatalf("gen for entry %s", err)
	}
	err = gen.GenContractClient(stateT, "../client/client.go")
	if err != nil {
		log.Fatalf("gen for client %s", err)
	}
	fmt.Println("generate hello actor success")
}
