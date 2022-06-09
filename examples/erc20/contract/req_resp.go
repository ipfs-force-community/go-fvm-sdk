package contract

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/big"
)

type ConstructorReq struct {
	Name        string
	Symbol      string
	Decimals    uint8
	TotalSupply big.Int
}

type TransferReq struct {
	ReceiverAddr   address.Address
	TransferAmount big.Int
}

type AllowanceReq struct {
	OwnerAddr   address.Address
	SpenderAddr address.Address
}

type TransferFromReq struct {
	OwnerAddr      address.Address
	SpenderAddr    address.Address
	TransferAmount big.Int
}

type ApprovalReq struct {
	SpenderAddr  address.Address
	NewAllowance big.Int
}
