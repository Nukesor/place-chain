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

<<<<<<< HEAD
func (pr *PixelRequest) IsValid() bool {
	transaction := pr.ToTransaction()
	bytes, err := transaction.SignedBytes()
	if err != nil {
		fmt.Println("Could not serialize transaction bytes for verifying signature")
		return false
	}
	return pr.PubKey.VerifyBytes(bytes, pr.Signature)
}

func (pr *PixelRequest) ToTransaction() *Transaction {
	return &Transaction{
=======
func (pr *PixelRequest) ToTransaction() *PixelTransaction {
	uuid4 := uuid.Must(uuid.NewV4())
	uuid4String := fmt.Sprintf("%s", uuid4)

	return &PixelTransaction{
		Tx{Type: 1},
>>>>>>> 19c87c8bd869606fa3a77cd217b73ae5aeb37b77
		pr.X,
		pr.Y,
		pr.Color,
		pr.Nonce,
		pr.PubKey,
		pr.Signature,
	}
}
