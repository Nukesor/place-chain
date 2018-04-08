package types

import (
	"fmt"

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

func (pr *PixelRequest) String() string {
	if pr == nil {
		return "nil Pixel Request"
	}
	return fmt.Sprintf("Pixel{%d %d %d}",
		pr.X, pr.Y, pr.Color)
}

func (pr *PixelRequest) IsValid() bool {
	transaction := pr.ToTransaction()
	bytes, err := transaction.SignedBytes()
	fmt.Printf("=== Validating\n %v\n", transaction)
	fmt.Printf("Bytes to validate: %s\n", string(bytes))
	if err != nil {
		fmt.Println("Could not serialize transaction bytes for verifying signature")
		return false
	}

	return pr.PubKey.VerifyBytes(bytes, pr.Signature)
}

func (pr *PixelRequest) ToTransaction() *PixelTransaction {
	return &PixelTransaction{
		Tx{Type: 1},
		pr.X,
		pr.Y,
		pr.Color,
		pr.Nonce,
		pr.PubKey,
		pr.Signature,
	}
}
