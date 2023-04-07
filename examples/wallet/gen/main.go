package main

import (
	"fmt"
	"log"
	"reflect"
	"wallet/contract"

	"github.com/ipfs-force-community/go-fvm-sdk/gen"
)

func main() {
	if err := gen.GenCborType("../contract", "", contract.State{}); err != nil {
		log.Fatalf("gen for ../contract: %s", err)
	}

	stateT := reflect.TypeOf(contract.State{})
	err := gen.GenEntry(stateT, "../entry_gen.go")
	if err != nil {
		log.Fatalf("gen for entry %s", err)
	}
	err = gen.GenContractClient(stateT, "../client/client_gen.go")
	if err != nil {
		log.Fatalf("gen for client %s", err)
	}
	fmt.Println("generate wallet actor success")
}
