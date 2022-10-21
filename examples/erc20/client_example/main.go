package main

import (
	"context"
	"encoding/hex"
	"erc20/client"
	"erc20/contract"
	"flag"
	"io/ioutil"
	"log"

	"github.com/filecoin-project/go-state-types/big"

	"github.com/filecoin-project/go-address"

	_ "github.com/filecoin-project/specs-actors/v8/actors/builtin/init"
	v0 "github.com/filecoin-project/venus/venus-shared/api/chain/v0"
)

func main() {
	var ip = flag.String("ip", "", "full node url")
	var token = flag.String("token", "", "full node token")
	var fromAddrStr = flag.String("from", "", "send message from")
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
	createParams, err := hex.DecodeString("846556656e757361560647005af3107a3fff")
	if err != nil {
		log.Fatalln(err)
		return
	}

	execRet, err := ercClient.CreateActor(ctx, installRet.CodeCid, createParams)
	if err != nil {
		log.Fatalln(err)
		return
	}
	println("actor id", execRet.IDAddress.String())

	val := big.NewInt(1000)
	req := contract.FakeSetBalance{
		Addr:    addr,
		Balance: val,
	}
	err = ercClient.FakeSetBalance(ctx, &req)
	if err != nil {
		log.Fatalln(err)
		return
	}
}
