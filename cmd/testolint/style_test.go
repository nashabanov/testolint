package main

import (
	"go/format"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestStyleCLI(t *testing.T) {
	dir := t.TempDir()
	binary := filepath.Join(dir, "testolint")
	if out, err := exec.Command("go", "build", "-o", binary, ".").CombinedOutput(); err != nil {
		t.Fatalf("build: %v\n%s", err, out)
	}
	fixtures := filepath.Join(dir, "fixtures")
	if err := os.CopyFS(fixtures, os.DirFS("../../analyzer/testdata")); err != nil {
		t.Fatal(err)
	}
	for attempt := range 2 {
		cmd := exec.Command(binary, "-disable-all", "-enable=TESTO025,TESTO026", "-fix", "./style")
		cmd.Dir = filepath.Join(fixtures, "src")
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("fix attempt %d: %v\n%s", attempt, err, out)
		}
		for _, name := range []string{"test.go", "crlf.go"} {
			path := filepath.Join(fixtures, "src", "style", name)
			actual, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			expected, err := os.ReadFile(path + ".golden")
			if err != nil {
				t.Fatal(err)
			}
			actual, err = format.Source(actual)
			if err != nil {
				t.Fatal(err)
			}
			expected, err = format.Source(expected)
			if err != nil {
				t.Fatal(err)
			}
			if string(actual) != string(expected) {
				t.Errorf("attempt %d: %s differs from golden", attempt, name)
			}
		}
	}
}
