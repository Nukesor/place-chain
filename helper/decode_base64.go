package main

import (
	"encoding/base64"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ./decode_base64 <base64 string>")
		return
	}
	input := os.Args[1]
	decoded, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		fmt.Printf("Error: cannot decode '%s': %s\n", input, err)
		return
	}
	fmt.Println(string(decoded))
}
