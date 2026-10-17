package parallel_coverage

import "github.com/ozontech/testo"

func helper() {}

type SerialLater struct {
	testo.Suite[*testo.T]
	resource *int
	ctx      *testo.T
}

func (s *SerialLater) BeforeEach(t *testo.T) { s.resource = new(int); s.ctx = t }
func (s *SerialLater) TestA(t *testo.T)      { t.Parallel(); t.Log(s.resource); s.ctx.Log("ctx") } // want "TESTO014:" "TESTO016:"
func (s *SerialLater) TestB(t *testo.T)      {}

type SerialEarlier struct {
	testo.Suite[*testo.T]
	resource *int
}

func (s *SerialEarlier) BeforeEach(t *testo.T) { s.resource = new(int) }
func (s *SerialEarlier) TestA(t *testo.T)      {}
func (s *SerialEarlier) TestB(t *testo.T)      { t.Parallel(); t.Log(s.resource) }

type Cleanup struct {
	testo.Suite[*testo.T]
	resource *int
}

func (s *Cleanup) BeforeEach(t *testo.T) {
	s.resource = new(int)
	t.Cleanup(func() { t.Log(s.resource) }) // want "TESTO015:"
}
func (s *Cleanup) TestA(t *testo.T) { t.Parallel() }
func (s *Cleanup) TestB(t *testo.T) {}

type Hooks struct{ testo.Suite[*testo.T] }

func (*Hooks) BeforeAll(t *testo.T) {
	helper()
	if true {
		helper()
	}
}
func (*Hooks) AfterAll(t *testo.T) {
	helper()
	for false {
		helper()
	}
}
func (*Hooks) BeforeEach(t *testo.T) { t.Parallel() }
func (*Hooks) TestA(t *testo.T)      { t.Parallel() } // want "TESTO017:"
type UnknownSibling struct{ testo.Suite[*testo.T] }

func (*UnknownSibling) BeforeEach(t *testo.T) { t.Parallel() }
func (*UnknownSibling) TestA(t *testo.T)      { t.Parallel() } // want "TESTO017:"
func (*UnknownSibling) TestB(t *testo.T)      { helper() }
func (*UnknownSibling) TestC(t *testo.T) {
	if true {
		helper()
	}
}
func (*UnknownSibling) TestD(t *testo.T, p struct{ N int }) { helper() }
func (*UnknownSibling) CasesN() []int                       { helper(); return nil }

type UnknownTarget struct{ testo.Suite[*testo.T] }

func (*UnknownTarget) BeforeEach(t *testo.T) { t.Parallel() }
func (*UnknownTarget) TestA(t *testo.T)      { helper(); t.Parallel() }

type Dynamic struct {
	testo.Suite[*testo.T]
	resource *int
}

func (s *Dynamic) BeforeEach(t *testo.T)               { s.resource = new(int) }
func (s *Dynamic) TestA(t *testo.T)                    { t.Parallel(); t.Log(s.resource) } // want "TESTO014:"
func (s *Dynamic) TestB(t *testo.T)                    {}
func (s *Dynamic) TestC(t *testo.T, p struct{ N int }) {}
func (*Dynamic) CasesN() []int                         { helper(); return nil }

// Regular methods precede parameterized methods even when their names sort later.
type CollectionOrder struct {
	testo.Suite[*testo.T]
	resource *int
}

func (s *CollectionOrder) BeforeEach(t *testo.T)               { s.resource = new(int) }
func (s *CollectionOrder) TestZ(t *testo.T)                    { t.Parallel(); t.Log(s.resource) } // want "TESTO014:"
func (s *CollectionOrder) TestA(t *testo.T, p struct{ N int }) {}
func (*CollectionOrder) CasesN() []int                         { return []int{1} }

type SafeOperations struct{ testo.Suite[*testo.T] }

func (*SafeOperations) BeforeEach(t *testo.T) { t.Parallel() }
func (*SafeOperations) TestA(t *testo.T)      { t.Helper(); t.Logf("%d", 1); t.Parallel() } // want "TESTO017:"
type UnknownBetween struct{ testo.Suite[*testo.T] }

func (*UnknownBetween) BeforeEach(t *testo.T) { t.Parallel() }
func (*UnknownBetween) TestA(t *testo.T)      { helper(); t.Parallel() }

type InvalidEach struct{ testo.Suite[*testo.T] }

func (*InvalidEach) BeforeEach(t int) {} // want "TESTO002:"
func (*InvalidEach) TestA(t *testo.T) { t.Parallel(); t.Parallel() }

type Cartesian struct {
	testo.Suite[*testo.T]
	resource *int
}

func (s *Cartesian) BeforeEach(t *testo.T)                  { s.resource = new(int) }
func (s *Cartesian) TestA(t *testo.T, p struct{ X, Y int }) { t.Parallel(); t.Log(s.resource) } // want "TESTO014:"
func (*Cartesian) CasesX() []int                            { return []int{1} }
func (*Cartesian) CasesY() []int                            { return []int{1, 2} }

// A definitely empty path cannot execute its otherwise unknown body.
type EmptySibling struct {
	testo.Suite[*testo.T]
	resource *int
}

func (s *EmptySibling) BeforeEach(t *testo.T)             { s.resource = new(int) }
func (s *EmptySibling) TestA(t *testo.T)                  { t.Parallel(); t.Log(s.resource) } // want "TESTO014:"
func (*EmptySibling) TestB(t *testo.T)                    {}
func (*EmptySibling) TestC(t *testo.T, p struct{ N int }) { helper() }
func (*EmptySibling) CasesN() []int                       { return nil } // want "TESTO010:"
