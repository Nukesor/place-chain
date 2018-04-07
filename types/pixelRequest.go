package types

import (
	"fmt"

	"github.com/satori/go.uuid"
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
	uuid4 := uuid.Must(uuid.NewV4())
	uuid4String := fmt.Sprintf("%s", uuid4)

	return &Transaction{
		pr.X,
		pr.Y,
		pr.Color,
		uuid4String,
	}
}
