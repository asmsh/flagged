package main

import (
	"flag"
	"fmt"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type input struct {
	sourceTypeNames []string
	outTypeNames    []string
	trimPrefix      string
	trimSuffix      string
	flagsSize       int

	outFile string
	outDir  string

	buildTags string

	patterns []string
}

func validateFlags() *input {
	// Validate that the type argument is passed and in correct format.
	if len(*typeFlag) == 0 {
		flag.Usage()
		os.Exit(2)
	}
	sourceTypeNames := strings.Split(*typeFlag, ",")
	if err := validateTypeNames(sourceTypeNames); err != nil {
		log.Fatalf("error: invalid type argument: %s", err)
	}

	// Validate that the type argument is passed and in correct format.
	// TODO: maybe add a validation to make sure sourceTypeNames and outTypeNames
	//  doesn't overlap, as it means the compilation will fail.
	var outTypeNames []string
	if len(*outTypeFlag) != 0 {
		outTypeNames = strings.Split(*outTypeFlag, ",")
		if err := validateTypeNames(outTypeNames); err != nil {
			log.Fatalf("error: invalid outType argument: %s", err)
		}
		if len(outTypeNames) != len(sourceTypeNames) {
			log.Fatalf("error: type argument doesn't match outType argument: %s", *outTypeFlag)
		}
	}

	// Validate the size argument, if passed.
	if *sizeFlag != 0 {
		switch *sizeFlag {
		case 8, 16, 32, 64:
		default:
			log.Fatalf("error: invalid size argument %d; supported values are 8,16,32,64", *sizeFlag)
		}
	}

	// We accept either one directory or a list of files. Which do we have?
	args := flag.Args()
	if len(args) == 0 {
		// Default: process whole package in current directory.
		args = []string{"."}
	}

	// Parse the package once.
	outputDir := getDirFromArgs(args, *buildTagsFlag)

	return &input{
		sourceTypeNames: sourceTypeNames,
		outTypeNames:    outTypeNames,
		trimPrefix:      *trimprefixFlag,
		trimSuffix:      *trimsuffixFlag,
		flagsSize:       *sizeFlag,
		outFile:         *outFileFlag,
		outDir:          outputDir,
		buildTags:       *buildTagsFlag,
		patterns:        args,
	}
}

func validateTypeNames(typeNames []string) error {
	for _, typeName := range typeNames {
		if !token.IsIdentifier(typeName) {
			return fmt.Errorf("invalid type identifier %q", typeName)
		}
	}
	return nil
}

func getDirFromArgs(args []string, tags string) string {
	var dir string
	if len(args) == 1 && isDirectory(args[0]) {
		dir = args[0]
	} else {
		if len(tags) != 0 {
			log.Fatal("error: -tags option applies only to directories, not when files are specified")
		}

		dir = filepath.Dir(args[0])
	}
	return dir
}

// isDirectory reports whether the named file is a directory.
func isDirectory(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		log.Fatal(err)
	}
	return info.IsDir()
}
