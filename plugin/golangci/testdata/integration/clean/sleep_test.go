package clean

import (
	"os/exec"
	"strings"
	"testing"
)

func TestMalformedSleepName(t *testing.T) {
	cmd := exec.Command("go", "test", "../testdata/sleep", "-run", "^TestMalformedName$", "-count=1", "-v")
	out, err := cmd.CombinedOutput()
	if err == nil || !strings.Contains(string(out), "Testfoo has malformed name") {
		t.Fatalf("expected Testo to reject the malformed name: %v\n%s", err, out)
	}
}
