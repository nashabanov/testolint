package providerruntime

import (
	"github.com/ozontech/testo"
	"testing"
)

type Suite struct {
	testo.Suite[*testo.T]
	value *int
}

func (s *Suite) BeforeEach(*testo.T)                       { s.value = new(int) }
func (*Suite) TestValue(t *testo.T, p struct{ Value int }) {}
func (s *Suite) CasesValue() []int                         { return []int{*s.value} }
func TestLateProvider(t *testing.T)                        { testo.RunSuite(t, &Suite{}) }

func TestConstructorProvider(t *testing.T) {
	value := 42
	testo.RunSuite(t, &Suite{value: &value})
}

type emptySuite struct{ testo.Suite[*testo.T] }

func (emptySuite) CasesValue() []int { return make([]int, 0, 4) }
func (emptySuite) TestValue(t *testo.T, p struct{ Value int }) {
	t.Fatal("empty provider must not schedule this test")
}
func TestEmptyMakeProvider(t *testing.T) { testo.RunSuite(t, emptySuite{}) }
