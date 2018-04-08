package types

import (
	"github.com/tendermint/go-crypto"
)

type Account struct {
	Name   string
	PubKey crypto.PubKey
}
