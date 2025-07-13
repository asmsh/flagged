// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// genflagged generates typed bit flags representations for a struct type T,
// out of its bool fields.
// It generates a type with underlying uint type and bit width that's enough
// to represent all bool fields in type as a bitmask.
//
// The generated type will have methods to access and modify the individual
// fields based on the fields' names, among other useful methods.
//
// More specifically, for each bool field, 5 different methods are generated:
//   - Is<field name>: reports whether the field is set to true or not.
//   - Set<field name>: sets the field to true, and returns the old value.
//   - Reset<field name>: sets the field to false, and returns the old value.
//   - Set<field name>To: sets the field to the new value, and returns the old value.
//   - Toggle<field name>: toggles the field's value, and returns the new value.
//
// In addition to 4 other methods for the whole generated type:
//   - BitFlags: returns a [github.com/asmsh/flagged.BitFlags] value,
//     wrapping the receiver value, and exposing a wider range of methods.
//   - Clone: returns a copy of the receiver value.
//   - TypedFlags: returns a copy of the receiver value as a value of the
//     original type that was used to generate the new flags type.
//   - SetTypedFlags: takes a value of the original type and overrides the
//     receiver value based on its fields.
//
// For example, given this type:
//
//	package permissions
//
//	type Permissions struct {
//		Read  bool
//		Write bool
//		Exec  bool
//	}
//
// Running this command in the same directory:
//
//	go run genflagged -type=Permissions
//
// Will generate file 'permissions_flagged.go', in package 'permissions', containing:
//
//	type PermissionsFlags flagged.BitFlags8
//
//	func (f *PermissionsFlags) BitFlags() flagged.BitFlags
//	func (f *PermissionsFlags) Clone() PermissionsFlags
//	func (f *PermissionsFlags) TypedFlags() Permissions
//	func (f *PermissionsFlags) SetTypedFlags(Permissions)
//	func (f *PermissionsFlags) IsRead() bool
//	func (f *PermissionsFlags) SetRead() bool
//	func (f *PermissionsFlags) ResetRead() bool
//	func (f *PermissionsFlags) SetReadTo(bool) bool
//	func (f *PermissionsFlags) ToggleRead() bool
//	func (f *PermissionsFlags) IsWrite() bool
//	func (f *PermissionsFlags) SetWrite() bool
//	...
//	func (f *PermissionsFlags) IsExec() bool
//	func (f *PermissionsFlags) SetExec() bool
//	...
//
// This enables memory-efficient and type-safe storage of bool configurations,
// replacing big structs with a compact uint-based types, without sacrificing
// readability or performance.
//
// It is designed to work with go:generate:
//
//	package permissions
//
//	//go:generate genflagged -type=Permissions
//	type Permissions struct {
//		Read  bool
//		Write bool
//		Exec  bool
//	}
//
// With no arguments, it processes the package in the current directory.
// Otherwise, the arguments must name a single directory holding a Go package
// or a set of Go source files that represent a single Go package.
//
// The -type flag accepts a comma-separated list of types, so a single run
// can generate multiple types.
// The default output file is 't_flagged.go', where 't' is the lower-cased
// name of the first type listed.
// The output file can be overridden with the -outFile flag.
//
// Types can also be declared in tests, in which case type declarations in
// the non-test package or its test variant are preferred over types defined
// in the package with suffix "_test".
// The default output file for type declarations in tests is 't_flagged_test.go'
// with t picked as above.
//
// The -outType flag accepts a comma-separated list of type names, specifying
// the names of the generated types.
// The default name for each provided type T in the -type flag is 'TBitFlags'.
// The out type names has to match the types provided in the -type flag, in length,
// with each type name in -outType being the generated type for the source type
// in the -type flag at the same index.
// If the '_' is provided as an out type name, the default name is used for its
// matching source type.
//
// The -size flag accepts one of 8, 16, 32 or 64, specifying the underlying
// uint type's bit width.
// The default underlying type's bit width depends on the number of bool
// fields in each of the source types provided in the -type flag, with:
//   - 1 to 8 bool fields: the underlying type is uint8.
//   - 9 to 16 bool fields: the underlying type is uint16.
//   - 17 to 32 bool fields: the underlying type is uint32.
//   - 33 to 64 bool fields: the underlying type is uint64.
//
// The -trimprefix and -trimsuffix flags specifies a prefix and suffix to be
// removed from each bool field's name, in each source type in the -type flag,
// before it's used to generated the different methods.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/types"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"golang.org/x/tools/go/packages"
)

var (
	typeFlag       = flag.String("type", "", "comma-separated list of type names to generate flags for; must be set")
	outTypeFlag    = flag.String("outType", "", "comma-separated list of generated type names; default <type>BitFlags")
	outFileFlag    = flag.String("outFile", "", "output file name; default srcdir/<type>_flagged.go")
	sizeFlag       = flag.Int("size", 0, "generated type size; one of 8,16,32,64; default depends on number of flags in <type>")
	trimprefixFlag = flag.String("trimprefix", "", "trim the `prefix` from each field in <type> before using it")
	trimsuffixFlag = flag.String("trimsuffix", "", "trim the `suffix` from each field in <type> before using it")

	buildTagsFlag = flag.String("tags", "", "comma-separated list of build tags to apply")

	verboseFlag = flag.Bool("verbose", false, "enable detailed logging during execution, including while loading packages")

	// TODO: add a flag to generate tests for the generated types (maybe only if outfile is a test file)
	// TODO: adda flag to generate benchmarks for the generated types.
)

// Usage is a replacement usage function for the flags package.
func Usage() {
	_, _ = fmt.Fprintf(os.Stderr, "Usage of genflagged:\n")
	_, _ = fmt.Fprintf(os.Stderr, "\tgenflagged [flags] -type T [directory]\n")
	_, _ = fmt.Fprintf(os.Stderr, "\tgenflagged [flags] -type T files... # Must be a single package\n")
	_, _ = fmt.Fprintf(os.Stderr, "For more information, see:\n")
	_, _ = fmt.Fprintf(os.Stderr, "\thttps://pkg.go.dev/github.com/asmsh/flagged/cmd/genflagged\n")
	_, _ = fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("genflagged: ")

	flag.Usage = Usage
	flag.Parse()

	// Init the verbose logger, if verbose mode is enabled.
	if *verboseFlag {
		verbose.logf = log.Printf
	}

	// Validate the provided flags and get ready-to-use input values.
	in := validateFlags()

	// Load the needed templates.
	headerTmpl, err := template.New("header").Parse(flaggedHeaderTemplate)
	if err != nil {
		log.Fatalf("error: internal: failed to load header template: %s", err)
	}
	bodyTmpl, err := template.New("body").Parse(flaggedTypeTemplate)
	if err != nil {
		log.Fatalf("error: internal: failed to load type template: %s", err)
	}

	// For each type, generate code in the first package where the type is declared.
	// The order of packages is as follows:
	// package x
	// package x compiled for tests
	// package x_test
	//
	// Each package pass could result in a separate generated file.
	// These files must have the same package and test/not-test nature as the types
	// from which they were generated.
	//
	// Types will be excluded when generated, to avoid repetitions.
	pkgs := loadPackages(in)
	sort.Slice(pkgs, func(i, j int) bool {
		// Put x_test packages last.
		iTest := strings.HasSuffix(pkgs[i].name, "_test")
		jTest := strings.HasSuffix(pkgs[j].name, "_test")
		if iTest != jTest {
			return !iTest
		}

		return len(pkgs[i].files) < len(pkgs[j].files)
	})
	for _, pkg := range pkgs {
		g := Generator{
			pkg: pkg,
		}

		verbose.Printf(
			"info: processing pacakge %s with %d remaining types\n",
			pkg.name,
			len(in.sourceTypeNames),
		)

		g.generateHeader(headerTmpl)

		// Run generate for types that can be found. Keep the rest for the remainingTypes iteration.
		var foundTypes, remainingTypes []string
		for idx, sourceTypeName := range in.sourceTypeNames {
			outTypeName := ""
			if len(in.outTypeNames) > 0 {
				outTypeName = in.outTypeNames[idx]

				if outTypeName == "_" {
					outTypeName = ""
					verbose.Printf(
						"info: skip specified out type name %s for source type %s while processing pacakge %s\n",
						outTypeName,
						sourceTypeName,
						pkg.name,
					)
				} else {
					verbose.Printf(
						"info: using specified out type name %s for source type %s while processing pacakge %s\n",
						outTypeName,
						sourceTypeName,
						pkg.name,
					)
				}
			}

			if len(outTypeName) == 0 {
				outTypeName = defaultOutTypeName(sourceTypeName)

				verbose.Printf(
					"info: using generated out type name %s for source type %s while processing pacakge %s\n",
					outTypeName,
					sourceTypeName,
					pkg.name,
				)
			}

			file := pkg.findStructTypeFile(sourceTypeName)
			if file != nil {
				if !file.isValidStructFile() {
					log.Fatalf(
						"error: found unsupported type %s (%s) for name %s in package %s."+
							"\n\tsupported types are struct types with bool fields.",
						file.foundSourceType.Name(),
						file.foundSourceType.Type().Underlying(),
						sourceTypeName,
						pkg.name,
					)
				}

				g.generateForStruct(sourceTypeName, outTypeName, bodyTmpl, file)
				foundTypes = append(foundTypes, sourceTypeName)
			} else {
				remainingTypes = append(remainingTypes, sourceTypeName)
			}
		}

		// Skip writing the file if not matching types are found in the current package.
		if n := len(foundTypes); n == 0 {
			verbose.Printf("info: no matching types found in pacakge %s\n", pkg.name)

			continue
		} else {
			verbose.Printf(
				"info: %d matching types found in pacakge %s\n",
				n,
				pkg.name,
			)
		}

		if n := len(remainingTypes); n > 0 {
			verbose.Printf(
				"info: %d remaining types after processing pacakge %s\n",
				n,
				pkg.name,
			)

			if *outFileFlag != "" {
				log.Fatalf(
					"error: cannot write to single file (-outFile=%q) when matching types are found in multiple packages",
					*outFileFlag,
				)
			}
		}

		// Update the source types to the remaining types, to try to find
		// them in the rest of the loaded packages.
		in.sourceTypeNames = remainingTypes

		// Format the output.
		src := g.format()

		// Write to file.
		outFileName := in.outFile
		if outFileName == "" {
			// Type names will be unique across packages since only the first
			// match is picked.
			// So there won't be collisions between a package compiled for tests
			// and the separate package of tests (package foo_test).
			outFileName = filepath.Join(in.outDir, defaultFileName(pkg, foundTypes[0]))
		}
		verbose.Printf(
			"info: writing output to file %s after processing pacakge %s\n",
			outFileName,
			pkg.name,
		)
		if err := os.WriteFile(outFileName, src, 0644); err != nil {
			log.Fatalf("error: failed to write to out file: %s", err)
		}
	}

	if len(in.sourceTypeNames) > 0 {
		log.Fatalf(
			"error: no matching types found for names: %s",
			strings.Join(in.sourceTypeNames, ","),
		)
	}
}

// Generator holds the state of the analysis.
// Primarily used to buffer the output for format.Source.
type Generator struct {
	buf bytes.Buffer // Accumulated output.
	pkg *Package     // Package we are scanning.
}

type Package struct {
	name         string
	defs         map[*ast.Ident]types.Object
	files        []*File
	hasTestFiles bool

	// options that apply to all files.
	trimPrefix string
	trimSuffix string
	flagsSize  int
}

// File holds a single parsed file and associated data.
type File struct {
	pkg  *Package  // Package to which this file belongs.
	file *ast.File // Parsed AST.

	// These fields are reset for each type being generated.
	sourceTypeName  string // Name of the source flag type.
	foundSourceType types.Object
	flagValues      []flagValue // Accumulator for flag values of that type.
	flagsSize       int         // Actual value based on number of flagValues
}

// loadPackages analyzes the single package constructed from the patterns and tags.
// loadPackages exits if there is an error.
//
// Returns all variants (such as tests) of the package.
func loadPackages(in *input) []*Package {
	cfg := &packages.Config{
		Mode: packages.LoadSyntax,
		// Tests are included, let the caller decide how to fold them in.
		Tests:      true,
		BuildFlags: []string{fmt.Sprintf("-tags=%s", in.buildTags)},
		Logf:       verbose.logf,
	}
	pkgs, err := packages.Load(cfg, in.patterns...)
	if err != nil {
		log.Fatalf("error: failed to load packages: %s", err)
	}
	if len(pkgs) == 0 {
		log.Fatalf(
			"error: no packages matching %v",
			strings.Join(in.patterns, " "),
		)
	}

	out := make([]*Package, len(pkgs))
	for i, pkg := range pkgs {
		p := &Package{
			name:       pkg.Name,
			defs:       pkg.TypesInfo.Defs,
			files:      make([]*File, len(pkg.Syntax)),
			trimPrefix: in.trimPrefix,
			trimSuffix: in.trimSuffix,
			flagsSize:  in.flagsSize,
		}

		for j, file := range pkg.Syntax {
			p.files[j] = &File{
				pkg:  p,
				file: file,
			}
		}

		// Keep track of test files, since we might want to generated
		// code that ends up in that kind of package.
		// Can be replaced once https://go.dev/issue/38445 lands.
		for _, f := range pkg.GoFiles {
			if strings.HasSuffix(f, "_test.go") {
				p.hasTestFiles = true
				break
			}
		}

		out[i] = p
	}
	return out
}

func (pkg *Package) findStructTypeFile(sourceTypeName string) *File {
	for _, file := range pkg.files {
		// Set the state for this run of the walker.
		file.sourceTypeName = sourceTypeName
		file.foundSourceType = nil
		file.flagValues = nil
		file.flagsSize = 0

		// Return the first file we find the matching sourceTypeName in.
		ast.Inspect(file.file, file.genStructDecl)
		if file.foundSourceType != nil {
			return file
		}
	}
	return nil
}

func (g *Generator) generateHeader(headerTmpl *template.Template) {
	// Print the header and package clause.
	headerInput := templateHeaderInput{
		CmdArgs:     strings.Join(os.Args[1:], " "),
		PackageName: g.pkg.name,
	}
	if err := headerTmpl.Execute(&g.buf, headerInput); err != nil {
		log.Fatalf("error: failed to generate header: %s", err)
	}
}

func (g *Generator) generateForStruct(
	sourceTypeName string,
	outTypeName string,
	bodyTmpl *template.Template,
	structFile *File,
) {
	// Make sure the flags size is within allowed limit.
	size := structFile.flagsSize
	if size > 64 {
		log.Fatalf(
			"error: type %s contains %d bool fields which is more than supported; maximum supported is 64",
			sourceTypeName,
			size,
		)
	}

	// Make sure the size is valid, if it's provided.
	if g.pkg.flagsSize != 0 {
		// If the want size is less than the required for the current file,
		// return with an error.
		if g.pkg.flagsSize < size {
			log.Fatalf(
				"error: type %s flags size is too small; required at least %d, requested %d",
				sourceTypeName,
				size,
				g.pkg.flagsSize,
			)
		}
		size = g.pkg.flagsSize
	}

	tmplInput := templateTypeInput{
		SourceTypeName:   sourceTypeName,
		OutTypeName:      outTypeName,
		OutTypeSize:      size,
		OutInterfaceName: outTypeName + "Interface",
		FlagValues:       structFile.flagValues,
	}
	if err := bodyTmpl.Execute(&g.buf, tmplInput); err != nil {
		log.Fatalf(
			"error: failed to generated implementation for type %s: %s",
			sourceTypeName,
			err,
		)
	}
}

// format returns the gofmt-ed contents of the Generator's buffer.
func (g *Generator) format() []byte {
	src, err := format.Source(g.buf.Bytes())
	if err != nil {
		// Should never happen, but can arise when developing this code.
		// The user can compile the output to see the error.
		log.Printf("warning: internal error: invalid Go generated: %s", err)
		log.Printf("warning: compile the package to analyze the error")
		return g.buf.Bytes()
	}
	return src
}
