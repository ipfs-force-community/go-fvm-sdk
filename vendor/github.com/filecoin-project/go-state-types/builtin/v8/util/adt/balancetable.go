package adt

import (
	addr "github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	cid "github.com/ipfs/go-cid"
)

// Bitwidth of balance table HAMTs, determined empirically from mutation
// patterns and projections of mainnet data
const BalanceTableBitwidth = 6

// A specialization of a map of addresses to (positive) token amounts.
// Absent keys implicitly have a balance of zero.
type BalanceTable Map

// Interprets a store as balance table with root `r`.
func AsBalanceTable(s Store, r cid.Cid) (*BalanceTable, error) {
	m, err := AsMap(s, r, BalanceTableBitwidth)
	if err != nil {
		return nil, err
	}

	return &BalanceTable{
		root:  m.root,
		store: s,
	}, nil
}

// Gets the balance for a key, which is zero if they key has never been added to.
func (t *BalanceTable) Get(key addr.Address) (abi.TokenAmount, error) {
	var value abi.TokenAmount
	found, err := (*Map)(t).Get(abi.AddrKey(key), &value)
	if !found || err != nil {
		value = big.Zero()
	}

	return value, err
}
