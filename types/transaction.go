package types

import (
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
}
