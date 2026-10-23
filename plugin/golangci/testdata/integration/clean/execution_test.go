package clean

import (
	"github.com/ozontech/testo"
	"os/exec"
	"strings"
	"testing"
)

func TestExecutionSafe(t *testing.T) {
	testo.RunTest(t, func(ctx *testo.T) {
		testo.Run(ctx, "nil is a no-op", nil)
		ctx.Cleanup(func() { ctx.Log("cleanup") })
	})
	t.Run("adapter", testo.Test(func(*testo.T) {}))
	_ = testo.Test[*testo.T](nil)
}

func TestExecutionRuntime(t *testing.T) {
	for _, tc := range []struct {
		name, evidence string
		fails          bool
	}{
		{"TestNilRunTest", "invalid memory address or nil pointer dereference", true},
		{"TestNilAdapter", "invalid memory address or nil pointer dereference", true},
		{"TestCleanupRun", "testing: t.Run called during t.Cleanup", true},
		{"TestClosedFile", "TESTO020 confirmed: Read, Write and Stat used the file after AfterEach closed it", true},
		{"TestPremature", "AfterEach ran before parallel child resumed", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cmd := exec.Command("go", "test", "../testdata/executionruntime", "-run", "^"+tc.name+"$", "-count=1", "-v")
			out, err := cmd.CombinedOutput()
			if (err != nil) != tc.fails || !strings.Contains(string(out), tc.evidence) {
				t.Fatalf("unexpected runtime result: %v\n%s", err, out)
			}
		})
	}
}
