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

type Transaction interface {
	SignedBytes() ([]byte, error)
	GetTwitterHandle() string
	GetTxType() TxType
}

// -------- TransactionWithType
// only used for unmarshalling a Transaction of unknown type. TransactionWithType only has a type, which can then be used to identify the specific transaction type (Pixel, Register)

type TransactionWithType struct {
	Type TxType
}

// -------- PixelTransaction

type PixelTransaction struct {
	Type          TxType
	X             int
	Y             int
	Color         Color
	Nonce         string
	TwitterHandle string
	Signature     crypto.Signature
}

func (pt PixelTransaction) GetTxType() TxType {
	return PIXEL_TRANSACTION
}

func (pt PixelTransaction) SignedBytes() ([]byte, error) {
	data := struct {
		X     int
		Y     int
		Color Color
		Nonce string
	}{
		pt.X, pt.Y, pt.Color, pt.Nonce,
	}

	return json.Marshal(data)
}

func (pt PixelTransaction) GetTwitterHandle() string {
	return pt.TwitterHandle
}

// -------- RegisterTransaction

type RegisterTransaction struct {
	Type            TxType
	TwitterHandle   string
	UserPubKey      crypto.PubKey
	ValidatorPubKey crypto.PubKey
	Signature       crypto.Signature
}

func (rt RegisterTransaction) GetTxType() TxType {
	return REGISTER_TRANSACTION
}

func (rt RegisterTransaction) SignedBytes() ([]byte, error) {
	data := struct {
		TwitterHandle string
		UserPubKey    crypto.PubKey
	}{
		rt.TwitterHandle, rt.UserPubKey,
	}

	return json.Marshal(data)
}

func (rt RegisterTransaction) GetTwitterHandle() string {
	return rt.TwitterHandle
}
