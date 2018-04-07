package types

import (
	"fmt"
)

type PixelRequest struct {
	X     int
	Y     int
	Color Color
}

func (pr *PixelRequest) String() string {
	if pr == nil {
		return "nil Pixel Request"
	}
	return fmt.Sprintf("Pixel{%d %d %d}",
		pr.X, pr.Y, pr.Color)
}

func (pr *PixelRequest) ToTransaction() *Transaction {
    return &Transaction{
        pr.X,
        pr.Y,
        pr.Color,
        "adf",
    }
}