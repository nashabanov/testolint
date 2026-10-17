package integration

import "github.com/ozontech/testo"

type Coverage struct {
	testo.Suite[*testo.T]
	resource *int
	ctx      *testo.T
}

func (s *Coverage) BeforeEach(t *testo.T) { s.resource = new(int); s.ctx = t }
func (s *Coverage) TestA(t *testo.T) {
	t.Parallel()
	t.Log(s.resource)
	s.ctx.Log("ctx")
}
func (s *Coverage) TestB(t *testo.T)                    {}
func (s *Coverage) TestC(t *testo.T, p struct{ N int }) {}
func (*Coverage) CasesN() []int                         { n := 1; return []int{n} }

type CoverageCleanup struct {
	testo.Suite[*testo.T]
	resource *int
}

func (s *CoverageCleanup) BeforeEach(t *testo.T) {
	s.resource = new(int)
	t.Cleanup(func() { t.Log(s.resource) })
}
func (*CoverageCleanup) TestA(t *testo.T) { t.Parallel() }
func (*CoverageCleanup) TestB(t *testo.T) {}
func coverageHelper()                     {}

type CoverageDuplicate struct{ testo.Suite[*testo.T] }

func (*CoverageDuplicate) BeforeAll(t *testo.T)  { coverageHelper() }
func (*CoverageDuplicate) AfterAll(t *testo.T)   { coverageHelper() }
func (*CoverageDuplicate) BeforeEach(t *testo.T) { t.Parallel() }
func (*CoverageDuplicate) TestA(t *testo.T) {
	t.Helper()
	t.Logf("%d", 1)
	t.Parallel()
}
func (*CoverageDuplicate) TestB(t *testo.T) { coverageHelper() }
