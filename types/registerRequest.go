package types

import (
	"fmt"
	"github.com/tendermint/go-crypto"
)

type RegisterRequest struct {
	Name   string
	PubKey crypto.PubKey
}

func (rr *RegisterRequest) String() string {
	if rr == nil {
		return "nil Register Request"
	}
	return fmt.Sprintf("Register{%s %v}", rr.Name, rr.PubKey)
}

func (rr *RegisterRequest) ToAccount() *Account {
	return &Account{
		rr.Name,
		rr.PubKey,
	}
	// TODO: account already exists? handle where?
}
