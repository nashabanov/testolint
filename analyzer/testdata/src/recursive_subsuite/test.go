package recursive_subsuite

import tt "github.com/ozontech/testo"

type T = *tt.T

type Pointer struct{ tt.Suite[T] }

func (s *Pointer) TestRecursive(t T) {
	tt.RunSubSuite(t, s) // want "TESTO013: recursive RunSubSuite invocation"
}

// Promotion must not duplicate the diagnostic in the method body.
type Promoted struct{ *Pointer }

type Value struct{ tt.Suite[T] }

func (s Value) TestRecursive(t T) {
	tt.RunSubSuite[Value, T, T](t, s) // want "TESTO013: recursive RunSubSuite invocation"
}

type Address struct{ tt.Suite[T] }

func (s Address) TestRecursive(t T) {
	tt.RunSubSuite(t, &s) // want "TESTO013: recursive RunSubSuite invocation"
}

type Aliased struct{ tt.Suite[T] }
type Alias = Aliased

func (s *Alias) TestRecursive(t T) {
	(tt.RunSubSuite[*Alias, T, T])(t, s) // want "TESTO013: recursive RunSubSuite invocation"
}

type Generic[P any] struct{ tt.Suite[T] }

func (s *Generic[P]) TestRecursive(t T) {
	tt.RunSubSuite(t, s) // want "TESTO013: recursive RunSubSuite invocation"
}

type Hook struct{ tt.Suite[T] }

func (Hook) TestOK(t T) {}
func (s *Hook) BeforeAll(t T) {
	tt.RunSubSuite( // want "TESTO013: recursive RunSubSuite invocation"
		t, s,
	)
}

// The same type does not mean the same receiver value.
type OtherInstance struct{ tt.Suite[T] }

func (s *OtherInstance) TestRecursive(t T) {
	tt.RunSubSuite(t, &OtherInstance{})
}

// Pointer receiver tests disappear from a suite's value method set.
type Dereferenced struct{ tt.Suite[T] }

func (s *Dereferenced) TestRecursive(t T) {
	tt.RunSubSuite(t, *s)
}

// Guarded or state-changing bodies are outside the direct recurrence subset.
type Guarded struct {
	tt.Suite[T]
	Depth int
}

func (s *Guarded) TestRecursive(t T) {
	if s.Depth < 2 {
		s.Depth++
		tt.RunSubSuite(t, s)
	}
}

type EarlyReturn struct{ tt.Suite[T] }

func (s *EarlyReturn) TestRecursive(t T) {
	if t == nil {
		return
	}
	tt.RunSubSuite(t, s)
}

type Changed struct{ tt.Suite[T] }

func (s *Changed) TestRecursive(t T) {
	s = &Changed{}
	tt.RunSubSuite(t, s)
}

// A Testo hook can terminate otherwise unconditional receiver recursion.
type Bounded struct {
	tt.Suite[T]
	Depth int
}

func (s *Bounded) BeforeAll(t T) {
	s.Depth++
	if s.Depth > 2 {
		t.Skip("bounded recursion")
	}
}
func (s *Bounded) TestRecursive(t T) {
	tt.RunSubSuite(t, s)
}

type EachHook struct{ tt.Suite[T] }

func (*EachHook) BeforeEach(t T) {}
func (s *EachHook) TestRecursive(t T) {
	tt.RunSubSuite(t, s)
}

// Helpers and deferred/closure calls do not establish a direct execution path.
type Helper struct{ tt.Suite[T] }

func (Helper) TestOK(t T)          {}
func (s *Helper) Helper(t T)       { tt.RunSubSuite(t, s) }
func (s *Helper) TestDeferred(t T) { defer tt.RunSubSuite(t, s) }
func (s *Helper) TestClosure(t T)  { func() { tt.RunSubSuite(t, s) }() }

// Only the method's own Testo parameter is accepted as the parent.
type OtherParent struct{ tt.Suite[T] }

var parent T

func (s *OtherParent) TestRecursive(t T) { tt.RunSubSuite(parent, s) }

// Receiver aliases and explicit options require additional reasoning.
type ReceiverAlias struct{ tt.Suite[T] }

func (s *ReceiverAlias) TestRecursive(t T) {
	self := s
	tt.RunSubSuite(t, self)
}

type Options struct{ tt.Suite[T] }

func (s *Options) TestRecursive(t T) { tt.RunSubSuite(t, s, "option") }

// Multiple or parameterized tests and providers are initially excluded.
type Multiple struct{ tt.Suite[T] }

func (Multiple) TestOK(t T)           {}
func (s *Multiple) TestRecursive(t T) { tt.RunSubSuite(t, s) }

type Parameterized struct{ tt.Suite[T] }

func (s *Parameterized) TestRecursive(t T, p struct{}) { tt.RunSubSuite(t, s) }

type Provider struct{ tt.Suite[T] }

func (s *Provider) TestRecursive(t T) { tt.RunSubSuite(t, s) }
func (Provider) CasesUnused() []int   { return []int{1} } // want "TESTO006"

// Invalid or ignored runtime methods must not produce cascading TESTO013.
type Invalid struct{ tt.Suite[T] }

func (s *Invalid) TestRecursive(t T) { tt.RunSubSuite(t, s) }
func (Invalid) TestBad()             {} // want "TESTO001"

type Malformed struct{ tt.Suite[T] } // want "TESTO011"
func (s *Malformed) Testbad(t T) { // want "TESTO007"
	tt.RunSubSuite(t, s)
}

type EmptyHook struct{ tt.Suite[T] } // want "TESTO011"
func (s *EmptyHook) BeforeAll(t T)   { tt.RunSubSuite(t, s) }

// Other lifecycle hooks and top-level wrappers are intentionally deferred.
type After struct{ tt.Suite[T] }

func (After) TestOK(t T)      {}
func (s *After) AfterAll(t T) { tt.RunSubSuite(t, s) }
func wrapper(t T, s *Pointer) { tt.RunSubSuite(t, s) }

// Unrelated methods with the same name are not Testo execution functions.
type Unrelated struct{ tt.Suite[T] }

func (s *Unrelated) TestRecursive(t T)           { s.RunSubSuite(t, s) }
func (*Unrelated) RunSubSuite(t T, s *Unrelated) {}
