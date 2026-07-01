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

// testFileName derives the companion test file name from the generated
// output file name, e.g. "options_flagged.go" -> "options_flagged_test.go".
// When the output is itself a test file (source declared in tests), it uses
// a "_gen_test.go" suffix to avoid colliding with the generated code file.
func testFileName(outFileName string) string {
	base := strings.TrimSuffix(outFileName, ".go")
	if trimmed, ok := strings.CutSuffix(base, "_test"); ok {
		return trimmed + "_gen_test.go"
	}
	return base + "_test.go"
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
