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
	var key crypto.PrivKey
	key.UnmarshalJSON([]byte(`
	  {"type": "ed25519",
		"data": "34480566DD7CF553E90732DA3B851AAE35A6719ABD12B83A5CC92C4056CFE048F1E7FA10D01D729BB5B72DEA6944113BCDDB5887EBEC79243CAE62E2D0D09C30"}`))
	expectedSignature := key.Sign(bytes)
	expectedJson, _ := expectedSignature.MarshalJSON()
	fmt.Printf("Expecting signature %s\n", expectedJson)
	fmt.Printf("Expected Signature verification result: %v\n", pr.PubKey.VerifyBytes(bytes, expectedSignature))
	json, _ := transaction.Signature.MarshalJSON()
	fmt.Printf("Signature: %s\n\n", json)
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
