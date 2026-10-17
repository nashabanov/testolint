package integration

import "github.com/ozontech/testo"

type Shared struct {
	testo.Suite[*testo.T]
	value int
	ctx   *testo.T
}

func (s *Shared) BeforeEach(t *testo.T) { s.value = 1; s.ctx = t }
func (s *Shared) TestA(t *testo.T) {
	t.Parallel()
	t.Log(s.value)
	s.ctx.Log("context")
}
func (s *Shared) TestB(t *testo.T) { t.Parallel() }

type Cleanup struct {
	testo.Suite[*testo.T]
	value int
}

func (s *Cleanup) BeforeEach(t *testo.T) {
	s.value = 1
	t.Cleanup(func() { t.Log(s.value) })
}
func (s *Cleanup) TestA(t *testo.T) { t.Parallel(); t.Log(s.value) }
func (s *Cleanup) TestB(t *testo.T) { t.Parallel() }

type Duplicate struct{ testo.Suite[*testo.T] }

func (*Duplicate) BeforeEach(t *testo.T) { t.Parallel() }
func (*Duplicate) TestA(t *testo.T)      { t.Parallel() }
func (*Duplicate) AfterEach(t *testo.T)  { t.Parallel() }

type Serial struct {
	testo.Suite[*testo.T]
	value int
	ctx   *testo.T
}

func (s *Serial) BeforeEach(t *testo.T) { s.value = 1; s.ctx = t; t.Cleanup(func() { t.Log(s.value) }) }
func (s *Serial) TestA(t *testo.T)      { t.Log(s.value); s.ctx.Log("ok") }
func (s *Serial) TestB(t *testo.T)      { t.Log(s.value) }

type Captured struct {
	testo.Suite[*testo.T]
	value int
}

func (s *Captured) BeforeEach(t *testo.T) { s.value = 1; v := s.value; t.Cleanup(func() { t.Log(v) }) }
func (s *Captured) TestA(t *testo.T)      { v := s.value; t.Parallel(); t.Log(v) }
func (s *Captured) TestB(t *testo.T)      { t.Parallel() }

type Single struct {
	testo.Suite[*testo.T]
	value int
}

func (s *Single) BeforeEach(t *testo.T) { s.value = 1 }
func (s *Single) TestA(t *testo.T)      { t.Parallel(); t.Log(s.value) }

type Branch struct{ testo.Suite[*testo.T] }

func (*Branch) BeforeEach(t *testo.T) {
	if false {
		t.Parallel()
	}
}
func (*Branch) TestA(t *testo.T) { t.Parallel() }

type Separate struct{ testo.Suite[*testo.T] }

func (*Separate) BeforeAll(t *testo.T) { t.Parallel() }
func (*Separate) TestA(t *testo.T)     { t.Parallel() }

// A test's own write is not a BeforeEach initialization hazard.
type OwnWrite struct {
	testo.Suite[*testo.T]
	value int
}

func (s *OwnWrite) TestA(t *testo.T) { s.value = 1; t.Parallel(); t.Log(s.value) }
func (s *OwnWrite) TestB(t *testo.T) { t.Parallel() }

// Replacing the field after suspension needs a different concurrency model.
type Replaced struct {
	testo.Suite[*testo.T]
	value int
}

func (s *Replaced) BeforeEach(t *testo.T) { s.value = 1 }
func (s *Replaced) TestA(t *testo.T)      { t.Parallel(); s.value = 2; t.Log(s.value) }
func (s *Replaced) TestB(t *testo.T)      { t.Parallel() }

// One parameterized method still creates multiple executions.
type Parameterized struct {
	testo.Suite[*testo.T]
	value int
}

func (s *Parameterized) BeforeEach(t *testo.T) { s.value = 1 }
func (s *Parameterized) CasesN() []int         { return []int{1, 2} }
func (s *Parameterized) TestA(t *testo.T, p struct{ N int }) {
	t.Parallel()
	t.Log(s.value)
}

type Alias = testo.T
type Aliased struct {
	testo.Suite[*Alias]
	value int
}

func (s *Aliased) BeforeEach(t *Alias) { s.value = 1 }
func (s *Aliased) TestA(t *Alias)      { t.Parallel(); t.Log(s.value) }
func (s *Aliased) TestB(t *Alias)      { t.Parallel() }

// Unknown helper behavior and early termination invalidate the timeline.
func helper() {}

type Unknown struct {
	testo.Suite[*testo.T]
	value int
}

func (s *Unknown) BeforeEach(t *testo.T) { s.value = 1; helper() }
func (s *Unknown) TestA(t *testo.T)      { t.Parallel(); t.Log(s.value) }
func (s *Unknown) TestB(t *testo.T)      { t.Parallel() }

// Resource methods inside Cleanup read their receiver before invocation.
type Resource struct{}

func (*Resource) Close() {}

type FixtureCleanup struct {
	testo.Suite[*testo.T]
	resource *Resource
}

func (s *FixtureCleanup) BeforeEach(t *testo.T) {
	s.resource = new(Resource)
	t.Cleanup(func() { s.resource.Close() })
}
func (s *FixtureCleanup) TestA(t *testo.T) { t.Parallel() }
func (s *FixtureCleanup) TestB(t *testo.T) { t.Parallel() }
