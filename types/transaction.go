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
	MarshalJSON() ([]byte, error)
}

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

func (pt PixelTransaction) IsValid() bool {
	if pt.PubKey.Empty() {
		return false
	}
	bytes, err := pt.SignedBytes()
	fmt.Printf("=== Validating\n %v\n", pt)
	fmt.Printf("Bytes to validate: %s\n", string(bytes))
	if err != nil {
		fmt.Println("Could not serialize transaction bytes for verifying signature")
		return false
	}
	return pt.PubKey.VerifyBytes(bytes, pt.Signature)
}

func (rt RegisterTransaction) IsValid() bool {
	return true
}

func (tx *PixelTransaction) SignedBytes() ([]byte, error) {
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

func (tx *PixelTransaction) MarshalJSON() ([]byte, error) {
	data := struct {
		Type TxType
		X      int
		Y      int
		Color  Color
		Nonce  string
		PubKey crypto.PubKey
	}{
		tx.GetTxType(), tx.X, tx.Y, tx.Color, tx.Nonce, tx.PubKey,
	}
	return json.Marshal(data)
}

func (rt RegisterTransaction) SignedBytes() ([]byte, error) {
	return json.Marshal(rt.Acc.Profile)
}

func (tx *RegisterTransaction) MarshalJSON() ([]byte, error) {
	data := struct {
		Type TxType
		Acc    Account
		PubKey crypto.PubKey
	}{
		tx.GetTxType(), *tx.Acc, tx.PubKey,
	}
	return json.Marshal(data)
}
