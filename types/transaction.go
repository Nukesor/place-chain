package types

import (
	"fmt"
)

type Transaction struct {
	X     int
	Y     int
	Color Color
	nonce string
}

func (tx *Transaction) String() string {
	if tx == nil {
		return "nil Pixel Request"
	}
	return fmt.Sprintf("Pixel{%d %d %d}",
		tx.X, tx.Y, tx.Color)
}
