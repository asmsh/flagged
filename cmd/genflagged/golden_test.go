package main

import (
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// updateGolden regenerates the .golden files from current output.
// Run `go test -update` after an intentional template change, then
// review the diff before committing.
var updateGolden = flag.Bool("update", false, "update golden files")

// goldenFixtures are directories under testdata/ holding a single Go
// package plus a //go:generate genflagged directive describing how to
// generate its flags. The generated output is compared against
// <output>.golden files checked in alongside the inputs.
var goldenFixtures = []string{
	"options",
	"mix_options",
	"max_options",
	"multiple_types",
	"raw_options",
	"tested_options",
	"raw_tested_options",
}

func TestGolden(t *testing.T) {
	// Build the generator binary once, shared across fixtures.
	bin := filepath.Join(t.TempDir(), "genflagged")
	if out, err := exec.Command("go", "build", "-o", bin, ".").CombinedOutput(); err != nil {
		t.Fatalf("building genflagged: %v\n%s", err, out)
	}

	for _, fixture := range goldenFixtures {
		t.Run(fixture, func(t *testing.T) {
			srcDir := filepath.Join("testdata", fixture)

			// Copy the fixture's .go inputs into a temp module so the
			// generator can write its output without touching testdata.
			inputs := copyFixture(t, srcDir)

			// Reuse the fixture's own go:generate invocation.
			args := generateArgs(t, inputs)
			gen := exec.Command(bin, append(args, ".")...)
			gen.Dir = filepath.Dir(inputs[0])
			if out, err := gen.CombinedOutput(); err != nil {
				t.Fatalf("running genflagged %v: %v\n%s", args, err, out)
			}

			// Compare every produced file against its golden counterpart.
			for _, produced := range producedFiles(t, gen.Dir, inputs) {
				got, err := os.ReadFile(produced)
				if err != nil {
					t.Fatal(err)
				}
				golden := filepath.Join(srcDir, filepath.Base(produced)+".golden")
				if *updateGolden {
					if err := os.WriteFile(golden, got, 0o644); err != nil {
						t.Fatal(err)
					}
					continue
				}
				want, err := os.ReadFile(golden)
				if err != nil {
					t.Fatalf("reading golden (run `go test -update` to create it): %v", err)
				}
				if string(got) != string(want) {
					t.Errorf("%s does not match %s:\ngot:\n%s\nwant:\n%s",
						filepath.Base(produced), golden, got, want)
				}
			}
		})
	}
}

// copyFixture copies the .go files from srcDir into a fresh temp module
// and returns the paths of the copied files.
func copyFixture(t *testing.T, srcDir string) []string {
	t.Helper()
	tmp := t.TempDir()
	writeFile(t, filepath.Join(tmp, "go.mod"), "module fixture\n\ngo 1.23\n")

	entries, err := os.ReadDir(srcDir)
	if err != nil {
		t.Fatal(err)
	}
	var copied []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") {
			continue
		}
		content, err := os.ReadFile(filepath.Join(srcDir, e.Name()))
		if err != nil {
			t.Fatal(err)
		}
		dst := filepath.Join(tmp, e.Name())
		writeFile(t, dst, string(content))
		copied = append(copied, dst)
	}
	if len(copied) == 0 {
		t.Fatalf("no .go files in %s", srcDir)
	}
	return copied
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

// generateArgs extracts the genflagged flags from the //go:generate
// directive found in one of the given files.
func generateArgs(t *testing.T, files []string) []string {
	t.Helper()
	const marker = "genflagged"
	for _, f := range files {
		content, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}

		for line := range strings.SplitSeq(string(content), "\n") {
			if !strings.HasPrefix(strings.TrimSpace(line), "//go:generate") {
				continue
			}
			fields := strings.Fields(line)
			for i, fld := range fields {
				if fld == marker {
					return fields[i+1:]
				}
			}
		}
	}
	t.Fatalf("no //go:generate genflagged directive found in %v", files)
	return nil
}

// producedFiles returns the .go files in dir that were not part of the
// copied inputs (i.e. the generator's output).
func producedFiles(t *testing.T, dir string, inputs []string) []string {
	t.Helper()
	original := make(map[string]bool, len(inputs))
	for _, in := range inputs {
		original[filepath.Base(in)] = true
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	var produced []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") || original[e.Name()] {
			continue
		}
		produced = append(produced, filepath.Join(dir, e.Name()))
	}
	if len(produced) == 0 {
		t.Fatalf("generator produced no output files in %s", dir)
	}
	return produced
}
