package types

import (
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

func (rr *RegisterRequest) IsValid() bool {
	return rr != nil && rr.Profile.Name != "" && !rr.PubKey.Empty()
}

func (rr *RegisterRequest) ToTransaction() *RegisterTransaction {
	return &RegisterTransaction{
		REGISTER_TRANSACTION,
		rr.Profile,
		rr.PubKey,
	}
}
