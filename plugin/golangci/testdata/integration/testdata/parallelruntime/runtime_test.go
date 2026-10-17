package parallelruntime

import (
	"github.com/ozontech/testo"
	"os"
	"testing"
)

type shared struct {
	testo.Suite[*testo.T]
	resource *int
	context  *testo.T
	observed *testing.T
	count    int
}

func (s *shared) BeforeEach(t *testo.T) {
	s.count++
	s.resource = new(int)
	*s.resource = s.count
	s.context = t
	localResource, localContext := s.resource, t
	t.Cleanup(func() {
		// Both hooks finish before either parallel body resumes. At least the first
		// callback therefore sees the second execution's resource and context.
		if localResource != s.resource && localContext != s.context {
			s.observed.Log("cleanup replacement confirmed")
		}
	})
}
func (s *shared) TestA(t *testo.T) {
	t.Parallel()
	if *s.resource != 2 {
		t.Fatal("expected last BeforeEach resource")
	}
	if s.context != t {
		s.observed.Log("field and context replacement confirmed")
	}
}
func (s *shared) TestB(t *testo.T) { s.TestA(t) }
func TestReplacement(t *testing.T) { testo.RunSuite(t, &shared{observed: t}) }

type duplicate struct{ testo.Suite[*testo.T] }

func (*duplicate) BeforeEach(t *testo.T) { t.Parallel() }
func (*duplicate) TestA(t *testo.T)      { t.Parallel() }
func TestDuplicate(t *testing.T) {
	if os.Getenv("TESTOLINT_REPRO_DUPLICATE") != "1" {
		t.Skip("intentional runtime failure; subprocess only")
	}
	testo.RunSuite(t, &duplicate{})
}
