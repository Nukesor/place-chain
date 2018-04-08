package types

import (
	"github.com/tendermint/go-crypto"
)

type PixelRequest struct {
	X         int
	Y         int
	Color     Color
	Nonce     string
	PubKey    crypto.PubKey
	Signature crypto.Signature
}

func (pr PixelRequest) IsValid() bool {
	return pr.ToTransaction().IsValid()
}

func (pr PixelRequest) ToTransaction() *PixelTransaction {
	return &PixelTransaction{
		PIXEL_TRANSACTION,
		pr.X,
		pr.Y,
		pr.Color,
		pr.Nonce,
		pr.PubKey,
		pr.Signature,
	}
}
