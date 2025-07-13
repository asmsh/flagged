package main

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

// defaultFileName that will put the generated code together with pkg.
func defaultFileName(pkg *Package, sourceTypeName string) string {
	suffix := "flagged.go"
	if pkg.hasTestFiles {
		suffix = "flagged_test.go"
	}
	return fmt.Sprintf("%s_%s", strings.ToLower(sourceTypeName), suffix)
}

func defaultOutTypeName(sourceTypeName string) string {
	return sourceTypeName + "BitFlags"
}

// TODO: what happens if there's nothing left after applying both trimmings?
func flagName(fieldName, trimPrefix, trimSuffix string) string {
	fn := []byte(fieldName)
	if len(trimPrefix) > 0 {
		fn = bytes.TrimPrefix(fn, []byte(trimPrefix))
	}
	if len(trimSuffix) > 0 {
		fn = bytes.TrimSuffix(fn, []byte(trimSuffix))
	}
	fn[0] = uint8(unicode.ToUpper(rune(fn[0])))
	return string(fn)
}

func flagSize(numFields int) int {
	switch {
	case 0 < numFields && numFields <= 8:
		return 8
	case 8 < numFields && numFields <= 16:
		return 16
	case 16 < numFields && numFields <= 32:
		return 32
	case 32 < numFields && numFields <= 64:
		return 64
	default:
		return numFields
	}
}
