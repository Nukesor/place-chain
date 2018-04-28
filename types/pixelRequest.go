package types

import (
	"github.com/tendermint/go-crypto"
)

type PixelRequest struct {
	X             int
	Y             int
	Color         Color
	Nonce         string
	TwitterHandle string
	Signature     crypto.Signature
}

func (pr PixelRequest) ToTransaction() *PixelTransaction {
	return &PixelTransaction{
		PIXEL_TRANSACTION,
		pr.X,
		pr.Y,
		pr.Color,
		pr.Nonce,
		pr.TwitterHandle,
		pr.Signature,
	}
}
