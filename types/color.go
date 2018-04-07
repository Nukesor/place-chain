package types

import (
    "bytes"
    "encoding/json"
    "strconv"
    "errors"
)

type Color int

const (
	White Color = iota
	Blue
	Green
	Yellow
	Black
	Cyan
	Pink
	Orange
    Red
)

var dataTypesId = map[Color]int{
    White: 0,
    Blue: 1,
    Green: 2,
    Yellow: 3,
    Black: 4,
    Cyan: 5,
    Pink: 6,
    Orange: 7,
    Red: 8,
}

var dataTypesName = map[int]Color{
    0: White,
    1: Blue,
    2: Green,
    3: Yellow,
    4: Black,
    5: Cyan,
    6: Pink,
    7: Orange,
    8: Red,
}

func (c *Color) MarshalJSON() ([]byte, error) {
    buffer := bytes.NewBufferString(`"`)
    val := strconv.Itoa(dataTypesId[*c])
    buffer.WriteString(val)
    buffer.WriteString(`"`)
    return buffer.Bytes(), nil
}

func (c *Color) UnmarshalJSON(b []byte) error {
    var s int
    err := json.Unmarshal(b, &s)
    if err != nil {
        return err
    }
    col, ok := dataTypesName[s]
    if !ok {
        return errors.New("Unknown color")
    }
    *c = col
    return nil
}