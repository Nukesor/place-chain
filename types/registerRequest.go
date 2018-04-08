package types

import (
	"errors"
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
	return fmt.Sprintf("Register{%s %s}", rr.Name, rr.PubKey)
}

func (rr *RegisterRequest) ToAccount() (*Account, error) {
	if rr.Name == "" || rr.PubKey.Empty() {
		return nil, errors.New("Account creation must specify `name` and `pubkey`")
	}
	return &Account{
		rr.Name,
		rr.PubKey,
	}, nil
	// TODO: account already exists? handle where?
}
