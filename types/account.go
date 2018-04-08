package types

import (
	"fmt"
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

func (profile *Profile) String() string {
	if profile == nil {
		return "nil Profile"
	}
	return fmt.Sprintf("Profile{%s %s %s}",
		profile.Name, profile.AvatarUrl, profile.Bio)
}
