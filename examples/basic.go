package main

import (
	"fmt"

	"github.com/asmsh/flagged"
)

const (
	permissionReadBitIndex flagged.BitIndex = iota
	permissionWriteBitIndex
	permissionExecBitIndex
)

func main() {
	var permFlags flagged.BitFlags8
	permFlags.Set(permissionReadBitIndex)
	permFlags.Set(permissionExecBitIndex)

	fmt.Println(permFlags.Is(permissionReadBitIndex))  // true
	fmt.Println(permFlags.Is(permissionWriteBitIndex)) // false

	permFlags.Toggle(permissionWriteBitIndex)
	fmt.Println(permFlags.Is(permissionWriteBitIndex)) // true

	fmt.Println(permFlags) // 00000111
}
