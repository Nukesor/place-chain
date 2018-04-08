package types

import (
<<<<<<< HEAD
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
=======
	"github.com/tendermint/go-crypto"
)

type TxType uint8

const (
	UNKNOWN TxType = iota
	PIXEL_TRANSACTION
	REGISTER_TRANSACTION
)

type Transaction interface{}

type Tx struct {
	Type TxType // pixel tx or register tx
}

type PixelTransaction struct {
	Tx
	X     int
	Y     int
	Color Color
	Nonce string
}

type RegisterTransaction struct {
	Tx
	Acc    *Account
	PubKey crypto.PubKey
}

func (pt PixelTransaction) GetTxType() TxType {
	return PIXEL_TRANSACTION
}

func (rt RegisterTransaction) GetTxType() TxType {
	return REGISTER_TRANSACTION
>>>>>>> 19c87c8bd869606fa3a77cd217b73ae5aeb37b77
}
