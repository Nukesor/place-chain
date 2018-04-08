package types

import (
	"github.com/tendermint/go-crypto"
)

type Account struct {
	Profile Profile
	PubKey  crypto.PubKey
}

type Profile struct {
	Name      string
	Bio       string
	AvatarUrl string
}

func (acc *Account) ToTransaction() *RegisterTransaction {
	return &RegisterTransaction{
		REGISTER_TRANSACTION,
		acc.Profile,
		acc.PubKey,
	}
}
