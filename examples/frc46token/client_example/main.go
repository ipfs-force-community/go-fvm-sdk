package main

import (
	"context"
	"flag"
	"fmt"
	"frc46token/client"
	"frc46token/contract"
	"log"
	"os"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"

	_ "github.com/filecoin-project/specs-actors/v8/actors/builtin/init"
	v0 "github.com/filecoin-project/venus/venus-shared/api/chain/v0"
)

func main() {
	var ip = flag.String("ip", "", "full node url")
	var token = flag.String("token", "", "full node token")
	var fromAddrStr = flag.String("from", "", "send message from, also frc46 owner address")
	var toAddrStr = flag.String("to", "", "used to recieve token")
	flag.Parse()

	ctx := context.Background()
	v0FullNode, closer, err := v0.DialFullNodeRPC(ctx, *ip, *token, nil)
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer closer()

	addr, err := address.NewFromString(*fromAddrStr)
	if err != nil {
		log.Fatalln(err)
		return
	}

	toAddr, err := address.NewFromString(*toAddrStr)
	if err != nil {
		log.Fatalln(err)
		return
	}

	actClient := client.NewFrc46TokenClient(v0FullNode, client.SetFromAddressOpt(addr))

	code, err := os.ReadFile("../erc20.wasm")
	if err != nil {
		log.Fatalln(err)
		return
	}

	installRet, err := actClient.Install(ctx, code)
	if err != nil {
		log.Fatalln(err)
		return
	}
	fmt.Printf("install code %s\n", installRet.CodeCid)

	execRet, err := actClient.CreateActor(ctx, installRet.CodeCid, &contract.ConstructorReq{
		Name:        "test_coin",
		Symbol:      "TC",
		Granularity: 1,
		Supply:      abi.NewTokenAmount(100),
	})
	if err != nil {
		log.Fatalln(err)
		return
	}
	fmt.Printf("actor id %s\n", execRet.IDAddress.String())
	balance, err := actClient.BalanceOf(ctx, &addr)
	if err != nil {
		log.Fatalln(err)
		return
	}
	fmt.Printf("address: %s, balance:%s\n", addr, balance)

	_, err = actClient.Transfer(ctx, &contract.TransferParams{
		To:     toAddr,
		Amount: abi.NewTokenAmount(10),
	})
	if err != nil {
		log.Fatalln(err)
		return
	}

	fmt.Printf("transfer from %s to %s amount 10\n", addr, toAddr)

	balance, err = actClient.BalanceOf(ctx, &toAddr)
	if err != nil {
		log.Fatalln(err)
		return
	}
	fmt.Printf("address: %s, balance:%s\n", toAddr, balance)
}
