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
