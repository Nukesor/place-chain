package types

import (
	"errors"
	"fmt"
	"github.com/tendermint/go-crypto"
)

type RegisterRequest struct {
	PubKey  crypto.PubKey
	Profile Profile
}

func (rr *RegisterRequest) String() string {
	if rr == nil {
		return "nil Register Request"
	}
	return fmt.Sprintf("Register{%s %s}", rr.Profile, rr.PubKey)
}

func (rr *RegisterRequest) ToAccount() (*Account, error) {
	if rr.Profile == nil || rr.PubKey.Empty() {
		return nil, errors.New("Account creation failed")
	}
	return &Account{
		rr.Profile,
		rr.PubKey,
	}, nil
	// TODO: account already exists? handle where?
}
