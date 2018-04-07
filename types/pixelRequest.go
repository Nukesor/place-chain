package types

import (
    "fmt"
)

type PixelRequest struct {
    X       uint8
    Y       uint8
    Color  Color
}

func (pr *PixelRequest) String() string {
    if pr == nil {
        return "nil Pixel Request"
    }
    return fmt.Sprintf("Pixel{%d %d %d}",
        pr.X, pr.Y, pr.Color)
}