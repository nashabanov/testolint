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
	cmd := exec.Command(binary, "run", "--config", ".golangci.yml", "./...")
	cmd.Dir = "testdata/integration"
	output, err := cmd.CombinedOutput()
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) || exitErr.ExitCode() != 1 {
		t.Fatalf("expected diagnostic exit code 1, got %v\n%s", err, output)
	}
	for _, diagnostic := range []struct {
		id    string
		count int
	}{{"TESTO001:", 1}, {"TESTO012:", 3}, {"TESTO013:", 2}, {"(testolint)", 6}} {
		if count := strings.Count(string(output), diagnostic.id); count != diagnostic.count {
			t.Errorf("%s count=%d, want %d\n%s", diagnostic.id, count, diagnostic.count, output)
		}
	}
	clean := exec.Command(binary, "run", "--config", ".golangci.yml", "./clean")
	clean.Dir = "testdata/integration"
	if output, err := clean.CombinedOutput(); err != nil || strings.Contains(string(output), "TESTO") {
		t.Fatalf("clean suite: %v\n%s", err, output)
	}
}
