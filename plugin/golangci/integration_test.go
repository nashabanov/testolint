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
	if !strings.Contains(string(output), "TESTO001:") || !strings.Contains(string(output), "(testolint)") {
		t.Fatalf("expected TESTO001 from testolint\n%s", output)
	}
}
