### flagged is a lightweight Go library for working with typed, compact bitflags.

[![PkgGoDev](https://pkg.go.dev/badge/github.com/asmsh/flagged)](https://pkg.go.dev/github.com/asmsh/flagged)
[![Go Report Card](https://goreportcard.com/badge/github.com/asmsh/flagged)](https://goreportcard.com/report/github.com/asmsh/flagged)
[![Tests](https://github.com/asmsh/flagged/workflows/Tests/badge.svg)](https://github.com/asmsh/flagged/actions)
[![Go Coverage](https://github.com/asmsh/flagged/wiki/coverage.svg)](https://raw.githack.com/wiki/asmsh/flagged/coverage.html)

It provides a minimal, extensible API for manipulating and inspecting compact bitflags, while remaining dependency-free and allocation-free.

It’s ideal for scenarios where you need efficient and compact boolean state representation — whether for generated flags, boolean configurations, or packed state machines.

### Features:

* Exposes typed wrappers for `uint` types: `BitFlags8`, `BitFlags16`, `BitFlags32`, `BitFlags64` (matching `uint8`, `uint16`, `uint32`, `uint64`, respectively).
* Unified interface: all exposed types implement a common BitFlags interface
* Core bit operations, using only the bit index (normal integers, with no shifting required for inputs).
* Pure Go implementation, no reflection, no dependencies, suitable for any application, in any environment.
* Easy to use directly or as a backend for code generators (check `github.com/asmsh/flagged/cmd/genflagged`).

### Installation:

```shell
go get github.com/asmsh/flagged
```

### Usage:

Import the package:

```go
import "github.com/asmsh/flagged"
```

Choose a bit width depending on your needs:

```go
var flags flagged.BitFlags16
flags.Set(3)
if flags.Is(3) {
    fmt.Println("flags[3] is set")
}
flags.Reset(3)
if !flags.Is(3) {
    fmt.Println("flags[3] is not set")
}
```

Or work with the interface:

```go
var f flagged.BitFlags = flagged.New(flagged.BitFlags32(0))
f.Set(1)
f.Set(5)
f.Toggle(1)
fmt.Println(f) // 00000000000000000000000000100000
```

### Example:

```go
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
```
