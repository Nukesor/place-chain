package types

import (
	"encoding/json"
	"fmt"
	"github.com/tendermint/go-crypto"
)

type Transaction struct {
	X         int
	Y         int
	Color     Color
	Nonce     string
	PubKey    crypto.PubKey
	Signature crypto.Signature
}

func (tx *Transaction) SignedBytes() (result []byte, err error) {
	data := struct {
		X      int
		Y      int
		Color  Color
		Nonce  string
		PubKey crypto.PubKey
	}{
		tx.X, tx.Y, tx.Color, tx.Nonce, tx.PubKey,
	}

	result, err = json.Marshal(data)
	return
}

func (tx *Transaction) String() string {
	if tx == nil {
		return "nil Transaction"
	}
	return fmt.Sprintf("Transaction{%d %d %d %s}",
		tx.X, tx.Y, tx.Color, tx.PubKey)
}
