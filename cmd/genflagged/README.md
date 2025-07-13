### genflagged: a tool for generating compact, uint-backed bitflags types out of Go struct types containing bool fields.

[![PkgGoDev](https://pkg.go.dev/badge/github.com/asmsh/flagged/cmd/genflagged)](https://pkg.go.dev/github.com/asmsh/flagged/cmd/genflagged)
[![Go Report Card](https://goreportcard.com/badge/github.com/asmsh/flagged/cmd/genflagged)](https://goreportcard.com/report/github.com/asmsh/flagged/cmd/genflagged)

It scans the struct type, selects the minimal underlying `uint` type (`uint8`, `uint16`, `uint32`, or `uint64`) based on the number of `bool` fields, and generates a new type with accessor and mutator methods for each field.

The goal is to enable memory-efficient, type-safe, and readable flag representation without sacrificing performance or clarity.

It uses the `github.com/asmsh/flagged` as the implementation for the bitflags access and modification.

### Features:

* Generates compact types out of `struct` types that has `bool` fields, offering a way to replace big structs with a compact uint-based types.
* Compatible with `go:generate` for automated code generation.
* Generates strongly typed flag types, with named methods after each field.
* Auto-selects optimal `uint` size (`uint8`, `uint16`, `uint32`, `uint64`) to fit fields, with optional override.
* Creates 5 methods per field: `Is<Field>()`, `Set<Field>()`, `Reset<Field>()`, `Set<Field>To(bool)`, `Toggle<Field>()`.
* Also generates general methods: `BitFlags()`, `Clone()`, `TypedFlags()`, `SetTypedFlags()`.

### Installation:

```shell
go install github.com/asmsh/flagged/cmd/genflagged@latest
```

### Usage:

Basic usage:

```shell
genflagged -type=T
```

With options:

```shell
genflagged [flags] -type T [directory]
genflagged [flags] -type T files... # Must be a single package
```

| Flag          | Description                                                                                                                                                                        |
|---------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `-type`       | Comma-separated list of struct types to generate the bitflags types for. (required)                                                                                                |
| `-outType`    | Comma-separated list of names for generated types, matching the values in `-type`. (default: `<type>BitFlags`) <br/> Use `_` to fall back to default naming for the matching type. |
| `-outFile`    | Name of the output file. (default: `<type>_flagged.go`, or `<type>_flagged_test.go` for test types)                                                                                |
| `-size`       | Force bit size for generated types (one of `8`, `16`, `32`, or `64`). (default: auto, and depends on number of `bool` fields of each type in `-type`)                              |
| `-trimprefix` | Trim prefix from bool field names before generating methods.                                                                                                                       |
| `-trimsuffix` | Trim suffix from bool field names before generating methods.                                                                                                                       |
| `-tags`       | Build tags to be applied during processing.                                                                                                                                        |
| `-verbose`    | Enable extensive logging during processing.                                                                                                                                        |

### Example:

Given the struct:

```go
package permissions

type Permissions struct {
	Read  bool
	Write bool
	Exec  bool
}
```

Run:

```shell
genflagged -type=Permissions
```

This generates file `permissions_flagged.go` in the same directory with:

```go
package permissions

import "github.com/asmsh/flagged"

type PermissionsBitFlags flagged.BitFlags8

func (f *PermissionsBitFlags) IsRead() bool
func (f *PermissionsBitFlags) SetRead() bool
func (f *PermissionsBitFlags) ResetRead() bool
func (f *PermissionsBitFlags) SetReadTo(bool) bool
func (f *PermissionsBitFlags) ToggleRead() bool
// Same for Write and Exec...

func (f *PermissionsBitFlags) BitFlags() flagged.BitFlags
func (f *PermissionsBitFlags) Clone() PermissionsBitFlags
func (f *PermissionsBitFlags) TypedFlags() Permissions
func (f *PermissionsBitFlags) SetTypedFlags(Permissions)

// Plus other helper types and constants...
```

Or, instead you can have `go:generate` handle that, such as:

```go
package permissions

//go:generate genflagged -type=Permissions
type Permissions struct {
	Read  bool
	Write bool
	Exec  bool
}
```

Then run:

```shell
go generate
```

### Notes:

* It's based on the `golang.org/x/tools/cmd/stringer` source, but with a lot of changes to produce the wanted types.
* Only `struct` types that contain at least one `bool` field are supported.
