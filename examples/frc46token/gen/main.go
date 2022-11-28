package main

import (
	"fmt"
	"frc46token/contract"
	"log"

	"github.com/ipfs-force-community/go-fvm-sdk/gen"
)

func main() {
	if err := gen.GenCborType("../contract", "", contract.Frc46Token{},
		contract.GetAllowanceParams{},
		contract.MintParams{},
		contract.MintReturn{},
		contract.TransferParams{},
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
	fmt.Println("generate erc20 actor success")
}
