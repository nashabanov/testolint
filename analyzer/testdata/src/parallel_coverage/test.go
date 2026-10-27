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

// A later body need not be understood to establish entry into BeforeEach.
type UnsupportedSerial struct {
	testo.Suite[*testo.T]
	resource *int
	ctx      *testo.T
}

func (s *UnsupportedSerial) BeforeEach(t *testo.T) { s.resource = new(int); s.ctx = t }
func (s *UnsupportedSerial) TestA(t *testo.T) {
	t.Parallel()
	t.Log(s.resource)    // want "TESTO014:"
	s.ctx.Log("context") // want "TESTO016:"
}
func (*UnsupportedSerial) TestB(t *testo.T)                    { helper() }
func (*UnsupportedSerial) TestC(t *testo.T, p struct{ N int }) { helper() }
func (*UnsupportedSerial) CasesN() []int                       { helper(); return nil }

type UnsupportedParallel struct {
	testo.Suite[*testo.T]
	resource *int
}

func (s *UnsupportedParallel) BeforeEach(t *testo.T) { s.resource = new(int) }
func (s *UnsupportedParallel) TestA(t *testo.T)      { t.Parallel(); t.Log(s.resource) } // want "TESTO014:"
func (*UnsupportedParallel) TestB(t *testo.T)        { t.Parallel(); helper() }

type UnsupportedCleanup struct {
	testo.Suite[*testo.T]
	resource *int
}

func (s *UnsupportedCleanup) BeforeEach(t *testo.T) {
	s.resource = new(int)
	t.Cleanup(func() { t.Log(s.resource) }) // want "TESTO015:"
}
func (s *UnsupportedCleanup) TestA(t *testo.T) { t.Parallel(); t.Log(s.resource) }
func (*UnsupportedCleanup) TestB(t *testo.T) {
	if true {
		helper()
	}
}

type KnownCasesUnknownSibling struct {
	testo.Suite[*testo.T]
	resource *int
}

func (s *KnownCasesUnknownSibling) BeforeEach(t *testo.T) { s.resource = new(int) }
func (s *KnownCasesUnknownSibling) TestA(t *testo.T, p struct{ X int }) {
	t.Parallel()
	t.Log(s.resource) // want "TESTO014:"
}
func (*KnownCasesUnknownSibling) CasesX() []int                       { return []int{1, 2} }
func (*KnownCasesUnknownSibling) TestB(t *testo.T, p struct{ N int }) { helper() }
func (*KnownCasesUnknownSibling) CasesN() []int                       { helper(); return nil }

type ReceiverEscape struct {
	testo.Suite[*testo.T]
	resource *int
}

func (s *ReceiverEscape) BeforeEach(t *testo.T) { s.resource = new(int) }
func (s *ReceiverEscape) TestA(t *testo.T)      { t.Parallel(); t.Log(s.resource) }
func (s *ReceiverEscape) TestB(t *testo.T)      { modify(s) }

type UnknownWrite struct {
	testo.Suite[*testo.T]
	resource *int
}

func (s *UnknownWrite) BeforeEach(t *testo.T) { s.resource = new(int) }
func (s *UnknownWrite) TestA(t *testo.T)      { t.Parallel(); t.Log(s.resource) }
func (s *UnknownWrite) TestB(t *testo.T)      { s.resource = new(int); helper() }

type ReceiverAlias struct {
	testo.Suite[*testo.T]
	resource *int
}

func (s *ReceiverAlias) BeforeEach(t *testo.T) { s.resource = new(int) }
func (s *ReceiverAlias) TestA(t *testo.T)      { t.Parallel(); t.Log(s.resource) }
func (s *ReceiverAlias) TestB(t *testo.T)      { alias := s; modify(alias) }

type ConditionalWrite struct {
	testo.Suite[*testo.T]
	resource *int
}

func (s *ConditionalWrite) BeforeEach(t *testo.T) { s.resource = new(int) }
func (s *ConditionalWrite) TestA(t *testo.T)      { t.Parallel(); t.Log(s.resource) }
func (s *ConditionalWrite) TestB(t *testo.T) {
	if true {
		s.resource = new(int)
	}
}

type AddressEscape struct {
	testo.Suite[*testo.T]
	resource *int
}

func (s *AddressEscape) BeforeEach(t *testo.T) { s.resource = new(int) }
func (s *AddressEscape) TestA(t *testo.T)      { t.Parallel(); t.Log(s.resource) }
func (s *AddressEscape) TestB(t *testo.T)      { modify(&s.resource) }

func modify(any) {}

type UnsupportedEarlier struct {
	testo.Suite[*testo.T]
	resource *int
}

func (s *UnsupportedEarlier) BeforeEach(t *testo.T) { s.resource = new(int) }
func (s *UnsupportedEarlier) TestA(t *testo.T)      { helper() }
func (s *UnsupportedEarlier) TestB(t *testo.T)      { t.Parallel(); t.Log(s.resource) }

type CaptureUnknownSibling struct {
	testo.Suite[*testo.T]
	resource *int
}

func (s *CaptureUnknownSibling) BeforeEach(t *testo.T) { s.resource = new(int) }
func (s *CaptureUnknownSibling) TestA(t *testo.T)      { v := s.resource; t.Parallel(); t.Log(v) }
func (s *CaptureUnknownSibling) TestB(t *testo.T)      { helper() }

type ValueUnknownSibling struct {
	testo.Suite[*testo.T]
	resource *int
}

func (s ValueUnknownSibling) BeforeEach(t *testo.T) { s.resource = new(int) }
func (s ValueUnknownSibling) TestA(t *testo.T)      { t.Parallel(); t.Log(s.resource) }
func (ValueUnknownSibling) TestB(t *testo.T)        { helper() }

type UnknownOnlyInterferer struct {
	testo.Suite[*testo.T]
	resource *int
}

func (s *UnknownOnlyInterferer) BeforeEach(t *testo.T)             { s.resource = new(int) }
func (s *UnknownOnlyInterferer) TestA(t *testo.T)                  { t.Parallel(); t.Log(s.resource) }
func (*UnknownOnlyInterferer) TestB(t *testo.T, p struct{ N int }) { helper() }
func (*UnknownOnlyInterferer) CasesN() []int                       { helper(); return nil }

// An empty provider cannot supply the only later execution.
type EmptyOnlyInterferer struct {
	testo.Suite[*testo.T]
	resource *int
}

func (s *EmptyOnlyInterferer) BeforeEach(t *testo.T)             { s.resource = new(int) }
func (s *EmptyOnlyInterferer) TestA(t *testo.T)                  { t.Parallel(); t.Log(s.resource) }
func (*EmptyOnlyInterferer) TestB(t *testo.T, p struct{ N int }) { helper() }
func (*EmptyOnlyInterferer) CasesN() []int                       { return []int{} } // want "TESTO010:"

// A known AfterEach replacement remains relevant after an unknown test body.
type TeardownWrite struct {
	testo.Suite[*testo.T]
	resource *int
}

func (s *TeardownWrite) BeforeEach(t *testo.T) { s.resource = new(int) }
func (s *TeardownWrite) TestA(t *testo.T)      { t.Parallel(); t.Log(s.resource) }
func (*TeardownWrite) TestB(t *testo.T)        { helper() }
func (s *TeardownWrite) AfterEach(t *testo.T)  { s.resource = nil }

// Each suite has only its own execution; models must not pool instances.
type IndependentOne struct {
	testo.Suite[*testo.T]
	resource *int
}

func (s *IndependentOne) BeforeEach(t *testo.T) { s.resource = new(int) }
func (s *IndependentOne) TestA(t *testo.T)      { t.Parallel(); t.Log(s.resource) }

type IndependentTwo struct {
	testo.Suite[*testo.T]
	resource *int
}

func (s *IndependentTwo) BeforeEach(t *testo.T) { s.resource = new(int) }
func (s *IndependentTwo) TestA(t *testo.T)      { t.Parallel(); t.Log(s.resource) }

// An opaque helper can recover the suite from Testo's context reflection.
// Passing T (or an alias) is not receiver-independent evidence.
type ContextEscape struct {
	testo.Suite[*testo.T]
	resource *int
}

func (s *ContextEscape) BeforeEach(t *testo.T) { s.resource = new(int) }
func (s *ContextEscape) TestA(t *testo.T)      { t.Parallel(); t.Log(s.resource) }
func (*ContextEscape) TestB(t *testo.T)        { modify(t) }

type ContextAlias struct {
	testo.Suite[*testo.T]
	resource *int
}

func (s *ContextAlias) BeforeEach(t *testo.T) { s.resource = new(int) }
func (s *ContextAlias) TestA(t *testo.T)      { t.Parallel(); t.Log(s.resource) }
func (*ContextAlias) TestB(t *testo.T)        { alias := t; modify(alias) }

// An empty make provider excludes an opaque receiver-using sibling entirely.
type EmptyMakeSibling struct {
	testo.Suite[*testo.T]
	resource *int
}

func (s *EmptyMakeSibling) BeforeEach(t *testo.T)               { s.resource = new(int) }
func (s *EmptyMakeSibling) TestA(t *testo.T)                    { t.Parallel(); t.Log(s.resource) } // want "TESTO014:"
func (*EmptyMakeSibling) TestB(t *testo.T)                      {}
func (s *EmptyMakeSibling) TestC(t *testo.T, p struct{ N int }) { opaqueResource(s.resource) }
func (*EmptyMakeSibling) CasesN() []int                         { return make([]int, 0, 7) } // want "TESTO010:"
func opaqueResource(*int)                                       {}
