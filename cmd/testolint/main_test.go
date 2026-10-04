package main

import (
	"errors"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestCLI(t *testing.T) {
	name := "testolint"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	binary := filepath.Join(t.TempDir(), name)
	if out, err := exec.Command("go", "build", "-o", binary, ".").CombinedOutput(); err != nil {
		t.Fatalf("build: %v\n%s", err, out)
	}
	for _, tc := range []struct {
		name, arg, want string
		code            int
	}{
		{"clean", "./clean", "", 0},
		{"diagnostic", ".", "suite_test.go:8:14: TESTO001:", 3},
		{"help", "-h", "Usage: testolint", 0},
		{"long-help", "--help", "Usage: testolint", 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cmd := exec.Command(binary, tc.arg)
			cmd.Dir = "../../plugin/golangci/testdata/integration"
			out, err := cmd.CombinedOutput()
			code := 0
			var exit *exec.ExitError
			if errors.As(err, &exit) {
				code = exit.ExitCode()
			} else if err != nil {
				t.Fatal(err)
			}
			if code != tc.code || tc.want == "" && len(out) != 0 || tc.want != "" && !strings.Contains(string(out), tc.want) {
				t.Fatalf("exit=%d, want %d; output:\n%s", code, tc.code, out)
			}
		})
	}
}
