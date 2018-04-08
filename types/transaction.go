package types

import (
	"encoding/json"
	"fmt"
	"github.com/tendermint/go-crypto"
)

type TxType uint8

const (
	UNKNOWN TxType = iota
	PIXEL_TRANSACTION
	REGISTER_TRANSACTION
)

type Transaction interface {
	IsValid() bool
	SignedBytes() ([]byte, error)
}

// -------- TransactionWithType
// only used for unmarshalling a Transaction of unknown type. TransactionWithType only has a type, which can then be used to identify the specific transaction type (Pixel, Register)

type TransactionWithType struct {
	Type TxType
}

// -------- PixelTransaction

type PixelTransaction struct {
	Type      TxType
	X         int
	Y         int
	Color     Color
	Nonce     string
	PubKey    crypto.PubKey
	Signature crypto.Signature
}

func (pt PixelTransaction) GetTxType() TxType {
	return PIXEL_TRANSACTION
}

func (pt PixelTransaction) IsValid() bool {
	if pt.PubKey.Empty() {
		return false
	}
	bytes, err := pt.SignedBytes()
	if err != nil {
		fmt.Println("PixelTransaction: Could not serialize transaction bytes for verifying signature")
		return false
	}
	return pt.PubKey.VerifyBytes(bytes, pt.Signature)
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

// -------- RegisterTransaction

type RegisterTransaction struct {
	Type    TxType
	Profile Profile
	PubKey  crypto.PubKey
}

func (rt RegisterTransaction) GetTxType() TxType {
	return REGISTER_TRANSACTION
}

func (rt RegisterTransaction) IsValid() bool {
	// TODO: implement with signing
	return true
}

func (rt RegisterTransaction) SignedBytes() ([]byte, error) {
	return json.Marshal(rt.Profile)
}
