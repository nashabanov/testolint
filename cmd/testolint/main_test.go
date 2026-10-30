package main

import (
	"errors"
	"os/exec"
	"path/filepath"
	"regexp"
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
		flags             []string
	}{
		{"clean", "./clean", nil, 0, 0, nil},
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
			"execution_test.go:9:29: TESTO019:",
			"execution_test.go:10:36: TESTO019:",
			"execution_test.go:11:23: TESTO018:",
			"teardown_test.go:16:51: TESTO020:",
			"construction_test.go:11:2: TESTO021:",
			"construction_test.go:15:2: TESTO022:",
			"plugin_parent_test.go:11:5: TESTO024:",
			"empty_make_test.go:8:23: TESTO010:",
			"provider_test.go:13:15: TESTO023:",
		}, 3, 29, nil},
		{"selection-default", "./testdata/selection", []string{"TESTO006:", "TESTO012:", "TESTO018:", "TESTO019:", "TESTO021:", "TESTO022:"}, 3, 6, nil},
		{"selection-only", "./testdata/selection", []string{"selection_test.go:28:22: TESTO019:"}, 3, 1, []string{"-disable-all", "-enable=TESTO019"}},
		{"selection-policy", "./testdata/selection", []string{"selection_test.go:12:10: TESTO006:"}, 3, 1, []string{"-disable-all", "-enable=TESTO006"}},
		{"selection-empty", "./testdata/selection", nil, 0, 0, []string{"-disable-all"}},
		{"selection-disable", "./testdata/selection", []string{"selection_test.go:29:23: TESTO018:"}, 3, 1, []string{"-disable=TESTO006,TESTO012,TESTO019,TESTO021,TESTO022"}},
		{"help", "-h", []string{"Usage: testolint"}, 0, 0, nil},
		{"long-help", "--help", []string{"Usage: testolint"}, 0, 0, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cmd := exec.Command(binary, append(append([]string{}, tc.flags...), tc.arg)...)
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
			if count := len(regexp.MustCompile(`TESTO[0-9]{3}:`).FindAll(out, -1)); count != tc.diagnostics {
				t.Errorf("diagnostics=%d, want %d; output:\n%s", count, tc.diagnostics, out)
			}
			if tc.name == "diagnostic" {
				repeat := exec.Command(binary, append(append([]string{}, tc.flags...), tc.arg)...)
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

func TestCLIInvalidConfiguration(t *testing.T) {
	binary := filepath.Join(t.TempDir(), "testolint")
	if out, err := exec.Command("go", "build", "-o", binary, ".").CombinedOutput(); err != nil {
		t.Fatalf("build: %v\n%s", err, out)
	}
	for _, tc := range []struct {
		name  string
		flags []string
		want  string
	}{
		{"unknown", []string{"-enable=TESTO999"}, "unknown rule"},
		{"malformed", []string{"-enable=TESTO019,"}, "empty rule ID"},
		{"duplicate", []string{"-enable=TESTO019,TESTO019"}, "duplicate rule"},
		{"conflict", []string{"-enable=TESTO019", "-disable=TESTO019"}, "both enable and disable"},
		{"reverse conflict", []string{"-disable=TESTO019", "-enable=TESTO019"}, "both enable and disable"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cmd := exec.Command(binary, append(tc.flags, "./does-not-exist")...)
			out, err := cmd.CombinedOutput()
			if err == nil || !strings.Contains(string(out), tc.want) || strings.Contains(string(out), "matched no packages") {
				t.Fatalf("expected configuration error %q before loading: %v\n%s", tc.want, err, out)
			}
		})
	}
}
