package main

import (
	"fmt"

	"github.com/asmsh/flagged"
)

func main() {
	var f flagged.BitFlags = flagged.New(flagged.BitFlags32(0))
	f.Set(1)
	f.Set(5)
	f.Toggle(1)
	fmt.Println(f) // 00000000000000000000000000100000
}
