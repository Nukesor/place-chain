package types

import (
	"encoding/json"
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
	X         int
	Y         int
	Color     Color
	Nonce     string
	PubKey    crypto.PubKey
	Signature crypto.Signature
}

func (tx *PixelTransaction) SignedBytes() (result []byte, err error) {
	data := struct {
		X     int
		Y     int
		Color Color
		Nonce string
	}{
		tx.X, tx.Y, tx.Color, tx.Nonce,
	}

	return json.Marshal(data)
}

func (tx *PixelTransaction) String() string {
	if tx == nil {
		return "nil Transaction"
	}
	res, err := json.Marshal(tx)
	if err != nil {
		return "Transaction that could not be json encoded"
	}
	return string(res)
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
}
