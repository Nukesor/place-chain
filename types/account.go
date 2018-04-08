package types

import (
	"github.com/tendermint/go-crypto"
)

type Account struct {
	Name   string
	PubKey crypto.PubKey
}

func (acc *Account) ToTransaction() *RegisterTransaction {

	return &RegisterTransaction{
		Tx{Type: 2},
		acc,
		acc.PubKey,
	}
}
