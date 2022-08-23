//go:build simulate

package contract

import (
	"reflect"
	"testing"

	"github.com/filecoin-project/go-state-types/big"
	"github.com/ipfs-force-community/go-fvm-sdk/sdk/sys"
	"github.com/ipfs/go-cid"
)

func TestErc20Token_Approval(t *testing.T) {
	// defer func() {
	// 	if err := recover(); err != nil {
	// 		fmt.Printf("%v\n", err)
	// 	}
	// }()
	sys.Begin()
	_ = Erc20Token{}
	sys.End()

}

// func TestErc20Token_SelfRoot(t *testing.T) {

// 	defer func() {
// 		if err := recover(); err != nil {
// 			fmt.Printf("%v\n", err)
// 		}
// 	}()

// 	sys.Begin()
// 	cidout, _ := cid.Decode("bafy2bzacecdjkk2tzogitpcybu3eszr4uptrjogstqmyt6u4q2p3hh4chmf3i")
// 	sys.GetSimulated().EXPECT().SelfRoot().Return(cidout, nil)
// 	token := Erc20Token{}

// 	got, err := token.SelfRoot()
// 	fmt.Printf("%v-%v\n", got, err)
// 	sys.End()

// }

func TestErc20Token_SelfRoot(t *testing.T) {

	sys.Begin()
	cidout, _ := cid.Decode("bafy2bzacecdjkk2tzogitpcybu3eszr4uptrjogstqmyt6u4q2p3hh4chmf3i")
	sys.GetSimulated().EXPECT().SelfRoot().Return(cidout, nil)

	cidoutwant, _ := cid.Decode("bafy2bzacecdjkk2tzogitpcybu3eszr4uptrjogstqmyt6u4q2p3hh4chmf3i")
	type fields struct {
		Name        string
		Symbol      string
		Decimals    uint8
		TotalSupply *big.Int
		Balances    cid.Cid
		Allowed     cid.Cid
	}
	tests := []struct {
		name    string
		fields  fields
		want    cid.Cid
		wantErr bool
	}{
		{name: "test1", fields: fields{}, want: cidoutwant, wantErr: false},
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
			got, err := tr.SelfRoot()
			if (err != nil) != tt.wantErr {
				t.Errorf("Erc20Token.SelfRoot() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Erc20Token.SelfRoot() = %v, want %v", got, tt.want)
			}
		})
	}
	sys.End()
}
