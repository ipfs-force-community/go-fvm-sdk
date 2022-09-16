//go:build simulate
// +build simulate

package contract

import (
	"context"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/adt"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys/simulated"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"reflect"
	"testing"
)

func TestErc20TokenApproval(t *testing.T) {
	simulated.Begin()
	_ = Erc20Token{}
	simulated.End()
}

func makeErc20Token() Erc20Token {
	map_, _ := adt.MakeEmptyMap(adt.AdtStore(context.Background()), adt.BalanceTableBitwidth)
	cidtest, err := map_.Root()
	if err != nil {
		panic(err)
	}
	TotalSupplytest := big.NewInt(0)

	return Erc20Token{Name: "pass", Symbol: "symbol", Decimals: 8, TotalSupply: &TotalSupplytest, Balances: cidtest, Allowed: cidtest}
}

func makeFakeSetBalance() *FakeSetBalance {
	Balance := big.NewInt(0)
	addr, _ := address.NewIDAddress(uint64(rand.Int()))
	FakeSetBalance := FakeSetBalance{Addr: addr, Balance: &Balance}
	return &FakeSetBalance
}

func TestErc20Token_FakeSetBalance(t *testing.T) {
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
	simulated.SetAccount(8899, addr)
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
