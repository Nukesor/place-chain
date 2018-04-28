package types

import (
	"fmt"
	"github.com/tendermint/go-crypto"
)

type RegisterRequest struct {
	PubKey               crypto.PubKey
	TwitterHandle        string
	VerificationTweetURL string
}

func (rr *RegisterRequest) String() string {
	if rr == nil {
		return "nil Register Request"
	}
	return fmt.Sprintf("Register{%s %s}", rr.TwitterHandle, rr.PubKey)
}

func (rr *RegisterRequest) IsValid() bool {
	return rr != nil && rr.TwitterHandle != "" && !rr.PubKey.Empty() && rr.VerificationTweetURL != ""
}

func (rr *RegisterRequest) ToTransaction(validatorKey crypto.PubKey, signature crypto.Signature) *RegisterTransaction {
	return &RegisterTransaction{
		REGISTER_TRANSACTION,
		rr.TwitterHandle,
		rr.PubKey,
		validatorKey,
		signature,
	}
}
