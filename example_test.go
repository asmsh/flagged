package flagged

import "fmt"

func ExampleBitFlags8() {
	const (
		permissionReadBitIndex BitIndex = iota
		permissionWriteBitIndex
		permissionExecBitIndex
	)

	var permFlags BitFlags8
	permFlags.Set(permissionReadBitIndex)
	permFlags.Set(permissionExecBitIndex)

	fmt.Println(permFlags.Is(permissionReadBitIndex))
	fmt.Println(permFlags.Is(permissionWriteBitIndex))

	fmt.Println(permFlags)

	permFlags.Toggle(permissionWriteBitIndex)
	fmt.Println(permFlags.Is(permissionWriteBitIndex))

	fmt.Println(permFlags)
	// Output:
	// true
	// false
	// 00000101
	// true
	// 00000111
}

func ExampleBitFlags() {
	var f BitFlags = New[BitFlags32](0)
	f.Set(1)
	f.Set(5)
	f.Toggle(1)
	fmt.Println(f)
	// Output:
	// 00000000000000000000000000100000
}
