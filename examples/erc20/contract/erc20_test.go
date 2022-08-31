//go:build simulate
// +build simulate

package contract

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/filecoin-project/go-address"
	//"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys/simulated"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
	"github.com/ipfs/go-cid"
)

func TestErc20Token_Approval(t *testing.T) {
	simulated.Begin()
	_ = Erc20Token{}
	simulated.End()

}

func makeErc20Token() Erc20Token {
	cidtest, _ := cid.Decode("bafy2bzacecdjkk2tzogitpcybu3eszr4uptrjogstqmyt6u4q2p3hh4chmf3i")
	TotalSupplytest := big.NewInt(88)
	return Erc20Token{Name: "name", Symbol: "symbol", Decimals: 8, TotalSupply: &TotalSupplytest, Balances: cidtest, Allowed: cidtest}
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
		{name: "test1", fields: makeErc20Token(), args: args{req: makeFakeSetBalance()}, wantErr: true},
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

func TestErc20Token_GetName(t *testing.T) {
	simulated.Begin()
	tests := []struct {
		name   string
		fields Erc20Token
		want   types.CborString
	}{
		{name: "test1", fields: makeErc20Token(), want: types.CborString("name")},
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

func TestErc20Token_GetBalanceOf(t *testing.T) {
	simulated.Begin()
	addr, _ := address.NewIDAddress(uint64(rand.Int()))

	simulated.SetActorAndAddress( 33, simulated.ActorState{}, addr)

	//(actorId uint32, ActorState ActorState, addr address.Address)
	wantbig := big.NewInt(99)
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
		{name: "test1", fields: makeErc20Token(), args: args{addr: &addr}, want: &wantbig, wantErr: false},
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
			got, err := tr.GetBalanceOf(tt.args.addr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Erc20Token.GetBalanceOf() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Erc20Token.GetBalanceOf() = %v, want %v", got, tt.want)
			}
		})
	}
	simulated.End()
}
