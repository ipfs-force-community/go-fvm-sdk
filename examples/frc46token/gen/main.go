package main

import (
	"fmt"
	"frc46token/contract"
	"log"
	"reflect"

	"github.com/ipfs-force-community/go-fvm-sdk/gen"
)

func main() {
	if err := gen.GenCborType("../contract", "", contract.Frc46Token{},
		contract.GetAllowanceParams{},
		contract.MintParams{},
		contract.MintReturn{},
		contract.TransferParams{},
		contract.TransferReturn{},
		contract.TransferFromParams{},
		contract.TransferFromReturn{},
		contract.IncreaseAllowanceParams{},
		contract.DecreaseAllowanceParams{},
		contract.RevokeAllowanceParams{},
		contract.BurnParams{},
		contract.BurnReturn{},
		contract.BurnFromParams{},
		contract.BurnFromReturn{},
		contract.ConstructorReq{},
		contract.FRC46TokenReceived{},
		contract.UniversalReceiverParams{},
	); err != nil {
		log.Fatalf("gen for ../contract: %s", err)
	}

	stateT := reflect.TypeOf(contract.Frc46Token{})
	err := gen.GenEntry(stateT, "../entry_gen.go")
	if err != nil {
		log.Fatalf("gen for entry %s", err)
	}
	err = gen.GenContractClient(stateT, "../client/client_gen.go")
	if err != nil {
		log.Fatalf("gen for client %s", err)
	}
	fmt.Println("generate frc46 actor success")
}
