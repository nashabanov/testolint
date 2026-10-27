package integration

import "github.com/ozontech/testo"

type lateProviderSuite struct {
	testo.Suite[*testo.T]
	value *int
}

func (s *lateProviderSuite) BeforeEach(*testo.T)                       { s.value = new(int) }
func (*lateProviderSuite) TestValue(t *testo.T, p struct{ Value int }) {}
func (s *lateProviderSuite) CasesValue() []int {
	return []int{*s.value}
}

func lateProviderRoot() { testo.RunSuite(nil, &lateProviderSuite{}) }
