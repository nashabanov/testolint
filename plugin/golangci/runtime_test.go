package golangci_test

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// Exercise the real framework's suspension and context identity, rather than
// treating analyzer stubs as runtime evidence.
func TestParallelRuntime(t *testing.T) {
	cmd := exec.Command("go", "test", "./testdata/parallelruntime", "-run", "^TestReplacement$", "-count=1", "-v")
	cmd.Dir = "testdata/integration"
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("replacement reproduction: %v\n%s", err, out)
	}
	for _, marker := range []string{"cleanup replacement confirmed", "field and context replacement confirmed"} {
		if !strings.Contains(string(out), marker) {
			t.Fatalf("missing %q\n%s", marker, out)
		}
	}
	cmd = exec.Command("go", "test", "./testdata/parallelruntime", "-run", "^TestDuplicate$", "-count=1", "-v")
	cmd.Dir = "testdata/integration"
	cmd.Env = append(os.Environ(), "TESTOLINT_REPRO_DUPLICATE=1")
	out, err = cmd.CombinedOutput()
	if err == nil || !strings.Contains(string(out), "Parallel called multiple times") {
		t.Fatalf("duplicate reproduction: %v\n%s", err, out)
	}
}
