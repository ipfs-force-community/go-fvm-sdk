package main

import (
	"context"
	"erc20/client"
	"erc20/contract"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/filecoin-project/go-state-types/abi"

	"github.com/filecoin-project/go-address"

	_ "github.com/filecoin-project/specs-actors/v8/actors/builtin/init"
	v0 "github.com/filecoin-project/venus/venus-shared/api/chain/v0"
)

func main() {
	var ip = flag.String("ip", "", "full node url")
	var token = flag.String("token", "", "full node token")
	var fromAddrStr = flag.String("from", "", "send message from, also erc20 mint address")
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
	ercClient := client.NewErc20TokenClient(v0FullNode, client.SetFromAddressOpt(addr))

	code, err := ioutil.ReadFile("../erc20.wasm")
	if err != nil {
		log.Fatalln(err)
		return
	}

	installRet, err := ercClient.Install(ctx, code)
	if err != nil {
		log.Fatalln(err)
		return
	}

	execRet, err := ercClient.CreateActor(ctx, installRet.CodeCid, &contract.ConstructorReq{
		Name:        "test_coin",
		Symbol:      "TC",
		Decimals:    8,
		TotalSupply: abi.NewTokenAmount(100),
		MintAddr:    addr,
	})
	if err != nil {
		log.Fatalln(err)
		return
	}
	fmt.Println("code cid %s actor id", installRet.CodeCid, execRet.IDAddress.String())
}
