package clean

import (
	tt "github.com/ozontech/testo"
	"testing"
)

type Child struct{ tt.Suite[*tt.T] }

func (Child) TestOK(t *tt.T) {}

type Parent struct{ tt.Suite[*tt.T] }

func (*Parent) TestChild(t *tt.T) {
	tt.RunSubSuite[Child, *tt.T, *tt.T](t, Child{})
}

// BeforeAll limits recursion on this same instance; the test's sole self call
// alone cannot establish an unbounded execution path.
type Bounded struct {
	tt.Suite[*tt.T]
	Depth int
}

func (s *Bounded) BeforeAll(t *tt.T) {
	s.Depth++
	if s.Depth > 2 {
		t.Skip("bounded recursion")
	}
}
func (s *Bounded) TestRecursive(t *tt.T) {
	tt.RunSubSuite(t, s)
}

func TestRunAnalysis(t *testing.T) {
	var suite *Parent = nil
	if suite != nil {
		t.Fatal("expected an initially nil suite")
	}
	suite = &Parent{}
	tt.RunSuite(t, suite)
	bounded := &Bounded{}
	run := tt.RunSuite[*Bounded]
	if !run(t, bounded) {
		t.Fatal("bounded suite failed")
	}
	if bounded.Depth != 3 {
		t.Fatalf("depth=%d, want 3", bounded.Depth)
	}
	nil := &Child{}
	tt.RunSuite(t, nil)
}
