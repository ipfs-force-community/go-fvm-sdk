package contract

import (
	"context"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin/v9/migration"
	"github.com/ipfs/go-cid"
	"math/rand"
	"reflect"
	"testing"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/adt"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys/simulated"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/stretchr/testify/assert"
)

func makeErc20Token() Erc20Token {
	map_, err := adt.MakeEmptyMap(adt.AdtStore(context.Background()), adt.BalanceTableBitwidth)
	if err != nil {
		panic(err)
	}
	cidtest, err := map_.Root()
	if err != nil {
		panic(err)
	}
	totalsupplytest := big.NewInt(888888)

	return Erc20Token{Name: "name", Symbol: "symbol", Decimals: 8, TotalSupply: &totalsupplytest, Balances: cidtest, Allowed: cidtest}
}

func makeFakeSetBalance() *FakeSetBalance {
	balance := big.NewInt(0)
	addr, err := address.NewIDAddress(uint64(rand.Int()))
	if err != nil {
		panic(err)
	}
	FakeSetBalance := FakeSetBalance{Addr: addr, Balance: &balance}
	return &FakeSetBalance
}

func TestErc20TokenFakeSetBalance(t *testing.T) {
	simulated.Begin()

	type args struct {
		req *FakeSetBalance
	}
	tests := []struct {
		name    string
		fields  Erc20Token
		args    args
		wantErr bool
	}{
		{name: "case1", fields: makeErc20Token(), args: args{req: makeFakeSetBalance()}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Erc20Token{
				Name:        tt.fields.Name,
				Symbol:      tt.fields.Symbol,
				Decimals:    tt.fields.Decimals,
				TotalSupply: tt.fields.TotalSupply,
				Balances:    tt.fields.Balances,
				Allowed:     tt.fields.Allowed,
			}
			if err := tr.FakeSetBalance(tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("Erc20Token.FakeSetBalance() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	simulated.End()
}

func TestErc20TokenGetName(t *testing.T) {
	simulated.Begin()
	tests := []struct {
		name   string
		fields Erc20Token
		want   types.CborString
	}{
		{name: "pass", fields: makeErc20Token(), want: types.CborString("name")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Erc20Token{
				Name:        tt.fields.Name,
				Symbol:      tt.fields.Symbol,
				Decimals:    tt.fields.Decimals,
				TotalSupply: tt.fields.TotalSupply,
				Balances:    tt.fields.Balances,
				Allowed:     tt.fields.Allowed,
			}
			if got := tr.GetName(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Erc20Token.GetName() = %v, want %v", got, tt.want)
			}
		})
	}
	simulated.End()
}

func TestErc20TokenSaveState(t *testing.T) {
	simulated.Begin()

	erc20 := makeErc20Token()
	sdk.SaveState(&erc20)

	newSt := new(Erc20Token)
	sdk.LoadState(newSt)
	assert.Equal(t, *newSt, erc20)
	simulated.End()
}

func TestErc20TokenGetBalanceOf(t1 *testing.T) {

	simulated.Begin()
	erc20 := makeErc20Token()
	balanceMap, _ := adt.AsMap(adt.AdtStore(context.Background()), erc20.Balances, adt.BalanceTableBitwidth)
	addr, _ := address.NewIDAddress(uint64(rand.Int()))
	simulated.SetAccount(8899, addr, migration.Actor{})
	balance := big.NewInt(100)
	if err := balanceMap.Put(types.ActorKey(8899), &balance); err != nil {
		panic(err)
	}
	newRoot, _ := balanceMap.Root()
	erc20.Balances = newRoot
	sdk.SaveState(&erc20)

	type args struct {
		addr *address.Address
	}

	tests := []struct {
		name    string
		fields  Erc20Token
		args    args
		want    *big.Int
		wantErr bool
	}{
		{name: "pass", fields: erc20, args: args{addr: &addr}, want: &balance, wantErr: false},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Erc20Token{
				Name:        tt.fields.Name,
				Symbol:      tt.fields.Symbol,
				Decimals:    tt.fields.Decimals,
				TotalSupply: tt.fields.TotalSupply,
				Balances:    tt.fields.Balances,
				Allowed:     tt.fields.Allowed,
			}
			got, err := t.GetBalanceOf(tt.args.addr)
			if (err != nil) != tt.wantErr {
				t1.Errorf("GetBalanceOf() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("GetBalanceOf() got = %v, want %v", got, tt.want)
			}
		})
	}
	simulated.End()
}

func TestErc20TokenTransfer(t *testing.T) {
	simulated.Begin()

	erc20 := makeErc20Token()

	// set info of caller
	callactorid := uint32(8899)
	calladdr, _ := address.NewIDAddress(uint64(rand.Int()))
	simulated.SetAccount(callactorid, calladdr, migration.Actor{Code: cid.Undef, Head: cid.Undef, CallSeqNum: 0, Balance: big.NewInt(99)})

	//  push  balance of caller
	balanceMap, _ := adt.AsMap(adt.AdtStore(context.Background()), erc20.Balances, adt.BalanceTableBitwidth)
	balance := big.NewInt(100000)
	if err := balanceMap.Put(types.ActorKey(callactorid), &balance); err != nil {
		panic(err)
	}

	newRoot, _ := balanceMap.Root()
	erc20.Balances = newRoot
	sdk.SaveState(&erc20)

	// set info of receiver
	receiactorid := uint32(7788)
	receiveaddr, _ := address.NewIDAddress(uint64(rand.Int()))
	simulated.SetAccount(receiactorid, receiveaddr, migration.Actor{Code: cid.Undef, Head: cid.Undef, CallSeqNum: 0, Balance: big.NewInt(99)})

	// set info of context
	callcontext := types.InvocationContext{Caller: abi.ActorID(callactorid)}
	simulated.SetCallContext(&callcontext)

	toamount := big.NewInt(9)
	type args struct {
		transferReq *TransferReq
	}
	tests := []struct {
		name    string
		fields  Erc20Token
		args    args
		wantErr bool
	}{
		{name: "pass", fields: makeErc20Token(), args: args{transferReq: &TransferReq{ReceiverAddr: receiveaddr, TransferAmount: &toamount}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := erc20.Transfer(tt.args.transferReq); (err != nil) != tt.wantErr {
				t.Errorf("Erc20Token.Transfer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	simulated.End()
}
