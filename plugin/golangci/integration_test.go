package golangci_test

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Run with TESTOLINT_GOLANGCI_BINARY pointing to a binary built by golangci-lint custom.
func TestIntegration(t *testing.T) {
	binary := os.Getenv("TESTOLINT_GOLANGCI_BINARY")
	if binary == "" {
		t.Skip("set TESTOLINT_GOLANGCI_BINARY to run the custom golangci-lint boundary test")
	}
	binary, err := filepath.Abs(binary)
	if err != nil {
		t.Fatal(err)
	}
	// Assert every diagnostic: golangci-lint otherwise caps identical messages
	// at three, which can hide a correct fourth TESTO017 from this fixture.
	cmd := exec.Command(binary, "run", "--max-same-issues=0", "--config", ".golangci.yml", "./...")
	cmd.Dir = "testdata/integration"
	output, err := cmd.CombinedOutput()
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) || exitErr.ExitCode() != 1 {
		t.Fatalf("expected diagnostic exit code 1, got %v\n%s", err, output)
	}
	for _, diagnostic := range []struct {
		id    string
		count int
	}{
		{"TESTO001:", 1},
		{"TESTO012:", 3},
		{"TESTO013:", 2},
		{"TESTO014:", 4},
		{"TESTO015:", 3},
		{"TESTO016:", 2},
		{"TESTO017:", 5},
		{"TESTO018:", 1},
		{"TESTO019:", 2},
		{"TESTO020:", 1},
		{"TESTO021:", 1},
		{"TESTO022:", 1},
		{"TESTO024:", 1},
		{"TESTO010:", 1},
		{"TESTO023:", 1},
		{"(testolint)", 29},
	} {
		if count := strings.Count(string(output), diagnostic.id); count != diagnostic.count {
			t.Errorf("%s count=%d, want %d\n%s", diagnostic.id, count, diagnostic.count, output)
		}
	}
	for _, position := range []string{
		"plugin_parent_test.go:11:5: TESTO024:",
		"empty_make_test.go:8:23: TESTO010:",
		"provider_test.go:13:15: TESTO023:",
		"construction_test.go:11:2: TESTO021:",
		"construction_test.go:15:2: TESTO022:",
		"teardown_test.go:16:51: TESTO020:",
		"execution_test.go:9:29: TESTO019:",
		"execution_test.go:10:36: TESTO019:",
		"execution_test.go:11:23: TESTO018:",
		"coverage_test.go:14:8: TESTO014:",
		"coverage_test.go:15:2: TESTO016:",
		"coverage_test.go:28:27: TESTO015:",
		"coverage_test.go:43:2: TESTO017:",
	} {
		if !strings.Contains(string(output), position) {
			t.Errorf("missing position %q\n%s", position, output)
		}
	}
	clean := exec.Command(binary, "run", "--config", ".golangci.yml", "./clean")
	clean.Dir = "testdata/integration"
	if output, err := clean.CombinedOutput(); err != nil || strings.Contains(string(output), "TESTO") {
		t.Fatalf("clean suite: %v\n%s", err, output)
	}
}

func TestIntegrationConfiguration(t *testing.T) {
	binary := os.Getenv("TESTOLINT_GOLANGCI_BINARY")
	if binary == "" {
		t.Skip("set TESTOLINT_GOLANGCI_BINARY to run configuration boundary tests")
	}
	binary, err := filepath.Abs(binary)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name, settings string
		ids            []string
		errorText      string
	}{
		{name: "defaults", ids: []string{"TESTO006", "TESTO012", "TESTO018", "TESTO019", "TESTO021", "TESTO022"}},
		{name: "only callback", settings: "disable-all: true\n            enable: [TESTO019]", ids: []string{"TESTO019"}},
		{name: "policy opt in", settings: "disable-all: true\n            enable: [TESTO006]", ids: []string{"TESTO006"}},
		{name: "none", settings: "disable-all: true"},
		{name: "empty lists", settings: "enable: []\n            disable: []", ids: []string{"TESTO006", "TESTO012", "TESTO018", "TESTO019", "TESTO021", "TESTO022"}},
		{name: "only cleanup", settings: "disable: [TESTO006, TESTO012, TESTO019, TESTO021, TESTO022]", ids: []string{"TESTO018"}},
		{name: "only fields", settings: "disable-all: true\n            enable: [TESTO021]", ids: []string{"TESTO021"}},
		{name: "only cycles", settings: "disable-all: true\n            enable: [TESTO022]", ids: []string{"TESTO022"}},
		{name: "unknown ID", settings: "enable: [TESTO999]", errorText: "unknown rule"},
		{name: "conflict", settings: "enable: [TESTO019]\n            disable: [TESTO019]", errorText: "both enable and disable"},
		{name: "duplicate", settings: "enable: [TESTO019, TESTO019]", errorText: "duplicate rule"},
		{name: "wrong type", settings: "enable: TESTO019", errorText: "decoding settings"},
		{name: "unknown setting", settings: "enabel: [TESTO019]", errorText: "unknown field"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			config := "version: \"2\"\nlinters:\n  default: none\n  enable: [testolint]\n  settings:\n    custom:\n      testolint:\n        type: module\n"
			if tc.settings != "" {
				config += "        settings:\n            " + tc.settings + "\n"
			}
			path := filepath.Join(t.TempDir(), "config.yml")
			if err := os.WriteFile(path, []byte(config), 0600); err != nil {
				t.Fatal(err)
			}
			target := "./testdata/selection"
			if tc.errorText != "" {
				target = "./does-not-exist"
			}
			cmd := exec.Command(binary, "run", "--max-same-issues=0", "--config", path, target)
			cmd.Dir = "testdata/integration"
			output, err := cmd.CombinedOutput()
			if tc.errorText != "" {
				if err == nil || !strings.Contains(string(output), tc.errorText) {
					t.Fatalf("expected %q: %v\n%s", tc.errorText, err, output)
				}
				return
			}
			if len(tc.ids) == 0 {
				if err != nil || strings.Contains(string(output), "TESTO") {
					t.Fatalf("expected no diagnostics: %v\n%s", err, output)
				}
			} else {
				var exit *exec.ExitError
				if !errors.As(err, &exit) || exit.ExitCode() != 1 {
					t.Fatalf("expected diagnostic exit: %v\n%s", err, output)
				}
				if strings.Count(string(output), "(testolint)") != len(tc.ids) {
					t.Fatalf("unexpected diagnostic count\n%s", output)
				}
				for _, id := range tc.ids {
					positions := map[string]string{
						"TESTO006": "selection_test.go:12:10: ",
						"TESTO012": "selection_test.go:27:20: ",
						"TESTO018": "selection_test.go:29:23: ",
						"TESTO019": "selection_test.go:28:22: ",
						"TESTO021": "selection_test.go:16:2: ",
						"TESTO022": "selection_test.go:18:19: ",
					}
					if !strings.Contains(string(output), positions[id]+id+":") {
						t.Errorf("missing %s\n%s", id, output)
					}
				}
			}
			clean := exec.Command(binary, "run", "--config", path, "./clean")
			clean.Dir = cmd.Dir
			if output, err := clean.CombinedOutput(); err != nil {
				t.Fatalf("clean: %v\n%s", err, output)
			}
		})
	}
}
