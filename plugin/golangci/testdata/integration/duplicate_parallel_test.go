package integration

import "github.com/ozontech/testo"

// A dynamic sibling count cannot change the T identity of a regular test.
// It still prevents the shared-field rules from establishing an overwrite window.
type DynamicSibling struct {
	testo.Suite[*testo.T]
	value int
	ctx   *testo.T
}

func (s *DynamicSibling) BeforeEach(t *testo.T) {
	s.value = 1
	s.ctx = t
	t.Cleanup(func() { t.Log(s.value) })
	t.Parallel()
}
func (s *DynamicSibling) TestKnown(t *testo.T) {
	t.Parallel()
	t.Log(s.value)
	s.ctx.Log("context")
}
func (*DynamicSibling) TestDynamic(t *testo.T, p struct{ N int }) {}
func (*DynamicSibling) CasesN() []int {
	n := 1
	return []int{n}
}

// The same independence applies to a test with its own known non-empty cases.
type KnownCases struct{ testo.Suite[*testo.T] }

func (*KnownCases) BeforeEach(t *testo.T) { t.Parallel() }
func (*KnownCases) TestKnown(t *testo.T, p struct{ X int }) {
	t.Parallel()
}
func (*KnownCases) TestDynamic(t *testo.T, p struct{ N int }) {}
func (*KnownCases) CasesX() []int                             { return []int{1} }
func (*KnownCases) CasesN() []int {
	n := 1
	return []int{n}
}
