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
	if rr.Profile.Name == "" || rr.PubKey.Empty() {
		return nil, errors.New("Account creation must specify `profile.name` and `pubkey`")
	}
	return &Account{
		rr.Profile,
		rr.PubKey,
	}, nil
	// TODO: account already exists? handle where?
}
