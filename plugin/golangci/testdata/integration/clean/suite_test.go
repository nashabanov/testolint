package clean

import (
	"github.com/ozontech/testo"
	"testing"
)

type T = *testo.T
type Base = testo.Suite[T]
type Number int
type Numbers []Number
type Params struct{ Age Number }
type Providers struct{}

func (Providers) CasesAge() Numbers { return Numbers{1} }

type Suite struct {
	Base
	Providers
}

func (*Suite) TestAge(t T, p Params) {
	if p.Age != 1 {
		t.Fatal(p.Age)
	}
}
func (*Suite) Cases(x int) string { return "ignored by Testo" }

type Tests struct{}

func (Tests) TestAge(t T, p Params) {}

type Inherited struct {
	testo.Suite[T]
	Tests
}

func (Inherited) CasesAge() Numbers { return Numbers{1} }
func TestSuites(t *testing.T) {
	testo.RunSuite(t, new(Suite))
	testo.RunSuite(t, new(Inherited))
}
