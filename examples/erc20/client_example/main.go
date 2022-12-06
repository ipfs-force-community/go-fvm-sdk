package main

import (
	"context"
	"erc20/client"
	"erc20/contract"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/filecoin-project/go-state-types/abi"

	"github.com/filecoin-project/go-address"

	_ "github.com/filecoin-project/specs-actors/v8/actors/builtin/init"
	v0 "github.com/filecoin-project/venus/venus-shared/api/chain/v0"
)

func main() {
	var ip = flag.String("ip", "", "full node url")
	var token = flag.String("token", "", "full node token")
	var fromAddrStr = flag.String("from", "", "send message from, also erc20 mint address")
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

	ercClient := client.NewErc20TokenClient(v0FullNode, client.SetFromAddressOpt(addr))

	code, err := os.ReadFile("../erc20.wasm")
	if err != nil {
		log.Fatalln(err)
		return
	}

	installRet, err := ercClient.Install(ctx, code)
	if err != nil {
		log.Fatalln(err)
		return
	}
	fmt.Printf("install code %s\n", installRet.CodeCid)

	execRet, err := ercClient.CreateActor(ctx, installRet.CodeCid, &contract.ConstructorReq{
		Name:        "test_coin",
		Symbol:      "TC",
		Decimals:    8,
		TotalSupply: abi.NewTokenAmount(100),
	})
	if err != nil {
		log.Fatalln(err)
		return
	}
	fmt.Printf("actor id %s\n", execRet.IDAddress.String())
	balance, err := ercClient.GetBalanceOf(ctx, &addr)
	if err != nil {
		log.Fatalln(err)
		return
	}
	fmt.Printf("address: %s, balance:%s\n", addr, balance)

	err = ercClient.Transfer(ctx, &contract.TransferReq{
		ReceiverAddr:   toAddr,
		TransferAmount: abi.NewTokenAmount(10),
	})
	if err != nil {
		log.Fatalln(err)
		return
	}

	fmt.Printf("transfer from %s to %s amount 10\n", addr, toAddr)

	balance, err = ercClient.GetBalanceOf(ctx, &toAddr)
	if err != nil {
		log.Fatalln(err)
		return
	}
	fmt.Printf("address: %s, balance:%s\n", toAddr, balance)
}
