package duplicate_parallel

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
	t.Parallel() // want "TESTO017:"
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
	t.Parallel() // want "TESTO017:"
}
func (*KnownCases) TestDynamic(t *testo.T, p struct{ N int }) {}
func (*KnownCases) CasesX() []int                             { return []int{1} }
func (*KnownCases) CasesN() []int {
	n := 1
	return []int{n}
}

// A single Parallel remains valid even when a sibling's count is unknown.
type SingleCall struct{ testo.Suite[*testo.T] }

func (*SingleCall) BeforeEach(t *testo.T)                     { t.Parallel() }
func (*SingleCall) TestKnown(t *testo.T)                      {}
func (*SingleCall) TestDynamic(t *testo.T, p struct{ N int }) {}
func (*SingleCall) CasesN() []int {
	n := 1
	return []int{n}
}

// The checked test must still have an established execution of its own.
type UnknownOwn struct{ testo.Suite[*testo.T] }

func (*UnknownOwn) BeforeEach(t *testo.T)                     { t.Parallel() }
func (*UnknownOwn) TestDynamic(t *testo.T, p struct{ N int }) { t.Parallel() }
func (*UnknownOwn) CasesN() []int {
	n := 1
	return []int{n}
}

type EmptyOwn struct{ testo.Suite[*testo.T] }

func (*EmptyOwn) BeforeEach(t *testo.T)                   { t.Parallel() }
func (*EmptyOwn) TestEmpty(t *testo.T, p struct{ N int }) { t.Parallel() }
func (*EmptyOwn) CasesN() []int                           { return nil } // want "TESTO010:"

// Invalid collection must not produce a secondary TESTO017 on a valid test.
type MissingProvider struct{ testo.Suite[*testo.T] }

func (*MissingProvider) BeforeEach(t *testo.T)                     { t.Parallel() }
func (*MissingProvider) TestKnown(t *testo.T)                      { t.Parallel() }
func (*MissingProvider) TestInvalid(t *testo.T, p struct{ N int }) {} // want "TESTO003:"

type WrongProvider struct{ testo.Suite[*testo.T] }

func (*WrongProvider) BeforeEach(t *testo.T)                     { t.Parallel() }
func (*WrongProvider) TestKnown(t *testo.T)                      { t.Parallel() }
func (*WrongProvider) TestInvalid(t *testo.T, p struct{ N int }) {}
func (*WrongProvider) CasesN() []string                          { return []string{"1"} } // want "TESTO004:"

type PrivateParameter struct{ testo.Suite[*testo.T] }

func (*PrivateParameter) BeforeEach(t *testo.T)                     { t.Parallel() }
func (*PrivateParameter) TestKnown(t *testo.T)                      { t.Parallel() }
func (*PrivateParameter) TestInvalid(t *testo.T, p struct{ n int }) {} // want "TESTO009:"

// Unknown teardown cannot hide a duplicate reached before teardown starts.
func cleanupDatabase() {}

type UnknownAfter struct{ testo.Suite[*testo.T] }

func (*UnknownAfter) BeforeAll(t *testo.T)  { cleanupDatabase() }
func (*UnknownAfter) AfterAll(t *testo.T)   { cleanupDatabase() }
func (*UnknownAfter) BeforeEach(t *testo.T) { t.Parallel() }
func (*UnknownAfter) TestA(t *testo.T)      { t.Parallel() } // want "TESTO017:"
func (*UnknownAfter) AfterEach(t *testo.T)  { cleanupDatabase() }

type SingleUnknownAfter struct{ testo.Suite[*testo.T] }

func (*SingleUnknownAfter) BeforeEach(t *testo.T) { t.Parallel() }
func (*SingleUnknownAfter) TestA(t *testo.T)      {}
func (*SingleUnknownAfter) AfterEach(t *testo.T)  { cleanupDatabase() }

type UnknownBefore struct{ testo.Suite[*testo.T] }

func (*UnknownBefore) BeforeEach(t *testo.T) { cleanupDatabase(); t.Parallel() }
func (*UnknownBefore) TestA(t *testo.T)      { t.Parallel() }
func (*UnknownBefore) AfterEach(t *testo.T)  { t.Parallel() }

type UnknownTest struct{ testo.Suite[*testo.T] }

func (*UnknownTest) BeforeEach(t *testo.T) { t.Parallel() }
func (*UnknownTest) TestA(t *testo.T)      { cleanupDatabase() }
func (*UnknownTest) AfterEach(t *testo.T)  { t.Parallel() }

type UnknownAfterBetween struct{ testo.Suite[*testo.T] }

func (*UnknownAfterBetween) TestA(t *testo.T)     { t.Parallel() }
func (*UnknownAfterBetween) AfterEach(t *testo.T) { cleanupDatabase(); t.Parallel() }

type SupportedAfter struct{ testo.Suite[*testo.T] }

func (*SupportedAfter) TestA(t *testo.T)     { t.Parallel() }
func (*SupportedAfter) AfterEach(t *testo.T) { t.Parallel() } // want "TESTO017:"

type InvalidAfter struct{ testo.Suite[*testo.T] }

func (*InvalidAfter) BeforeEach(t *testo.T) { t.Parallel() }
func (*InvalidAfter) TestA(t *testo.T)      { t.Parallel() }
func (*InvalidAfter) AfterEach(t int)       {} // want "TESTO002:"

type SuiteContext struct{ testo.Suite[*testo.T] }

func (*SuiteContext) BeforeAll(t *testo.T) { t.Parallel() }
func (*SuiteContext) TestA(t *testo.T)     { t.Parallel() }
func (*SuiteContext) AfterEach(t *testo.T) { cleanupDatabase() }
