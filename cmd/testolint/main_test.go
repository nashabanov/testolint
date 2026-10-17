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
		name, arg         string
		want              []string
		code, diagnostics int
	}{
		{"clean", "./clean", nil, 0, 0},
		{"diagnostic", ".", []string{
			"suite_test.go:8:14: TESTO001:",
			"runs_test.go:12:2: TESTO013:",
			"runs_test.go:19:2: TESTO013:",
			"runs_test.go:25:17: TESTO012:",
			"runs_test.go:27:9: TESTO012:",
			"runs_test.go:28:51: TESTO012:",
			"coverage_test.go:14:8: TESTO014:",
			"coverage_test.go:15:2: TESTO016:",
			"coverage_test.go:28:27: TESTO015:",
			"coverage_test.go:43:2: TESTO017:",
		}, 3, 20},
		{"help", "-h", []string{"Usage: testolint"}, 0, 0},
		{"long-help", "--help", []string{"Usage: testolint"}, 0, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cmd := exec.Command(binary, tc.arg)
			cmd.Dir = "../../plugin/golangci/testdata/integration"
			out, err := cmd.CombinedOutput()
			code := 0
			if exit, ok := errors.AsType[*exec.ExitError](err); ok {
				code = exit.ExitCode()
			} else if err != nil {
				t.Fatal(err)
			}
			if code != tc.code || len(tc.want) == 0 && len(out) != 0 {
				t.Fatalf("exit=%d, want %d; output:\n%s", code, tc.code, out)
			}
			for _, want := range tc.want {
				if !strings.Contains(string(out), want) {
					t.Errorf("missing %q in output:\n%s", want, out)
				}
			}
			if count := strings.Count(string(out), "TESTO"); count != tc.diagnostics {
				t.Errorf("diagnostics=%d, want %d; output:\n%s", count, tc.diagnostics, out)
			}
			if tc.name == "diagnostic" {
				repeat := exec.Command(binary, tc.arg)
				repeat.Dir = cmd.Dir
				repeated, repeatErr := repeat.CombinedOutput()
				var exit *exec.ExitError
				if !errors.As(repeatErr, &exit) || exit.ExitCode() != tc.code || string(repeated) != string(out) {
					t.Fatalf("diagnostics changed between runs: %v\n%s", repeatErr, repeated)
				}
			}
		})
	}
}
