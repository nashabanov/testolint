package integration

import (
	tt "github.com/ozontech/testo"
	"testing"
)

type Recursive struct{ tt.Suite[*tt.T] }
type RunAlias = Recursive

func (s *Recursive) TestRecursive(t *tt.T) {
	tt.RunSubSuite(t, s)
}

type RecursiveHook struct{ tt.Suite[*tt.T] }

func (RecursiveHook) TestOK(t *tt.T) {}
func (s *RecursiveHook) BeforeAll(t *tt.T) {
	tt.RunSubSuite(t, s)
}

// This helper is analyzed without executing its intentionally nil arguments.
func RunNilSuiteArguments(t *testing.T, parent *tt.T) {
	var suite *Recursive
	tt.RunSuite(t, suite)
	run := tt.RunSuite[*RunAlias, *tt.T]
	run(t, (*RunAlias)(nil))
	tt.RunSubSuite[*Recursive, *tt.T, *tt.T](parent, (*Recursive)(nil))
}
