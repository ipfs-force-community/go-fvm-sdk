package main

import (
	"erc20/contract"
	"fmt"
	"log"
	"reflect"

	"github.com/ipfs-force-community/go-fvm-sdk/gen"
)

func main() {
	if err := gen.GenCborType("../contract", "", contract.Erc20Token{},
		contract.ConstructorReq{},
		contract.TransferReq{},
		contract.AllowanceReq{},
		contract.TransferFromReq{},
		contract.ApprovalReq{},
		contract.FakeSetBalance{}); err != nil {
		log.Fatalf("gen for ../contract: %s", err)
	}
	stateT := reflect.TypeOf(contract.Erc20Token{})
	err := gen.GenEntry(stateT, "../entry_gen.go")
	if err != nil {
		log.Fatalf("gen for entry %s", err)
	}
	err = gen.GenContractClient(stateT, "../client/client_gen.go")
	if err != nil {
		log.Fatalf("gen for client %s", err)
	}
	fmt.Println("generate erc20 actor success")
}
