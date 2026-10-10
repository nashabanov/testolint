package nil_suite

import (
	"faketesto"
	tt "github.com/ozontech/testo"
	"promotionlib"
)

type S struct{ tt.Suite[*tt.T] }

func (*S) TestOK(t *tt.T) {}

type Alias = S

func explicit(t *tt.T) {
	tt.RunSuite(t,
		(*S)(nil), // want "TESTO012: suite argument is statically nil"
	)
	tt.RunSuite(t, (*S)((*Alias)(nil)))       // want "TESTO012: suite argument is statically nil"
	tt.RunSuite(t, (*promotionlib.Base)(nil)) // want "TESTO012: suite argument is statically nil"
	tt.RunSuite(t, (*Generic[int])(nil))      // want "TESTO012: suite argument is statically nil"
	tt.RunSuite(t, convert(nil))
	tt.RunSuite(t, (*S)(nil))                    // want "TESTO012: suite argument is statically nil"
	tt.RunSuite[*S](t, nil)                      // want "TESTO012: suite argument is statically nil"
	tt.RunSuite[*Alias, *tt.T](t, (*Alias)(nil)) // want "TESTO012: suite argument is statically nil"
	tt.RunSubSuite(t, ((*S)(nil)))               // want "TESTO012: suite argument is statically nil"
	run := tt.RunSuite[*S]
	run(t, (*S)(nil)) // want "TESTO012: suite argument is statically nil"
	tt.RunSuite(t, &S{})
	tt.RunSuite(t, new(S))
	tt.RunSuite(t, makeSuite())
	faketesto.RunSuite(t, (*S)(nil))
}

func locals(t *tt.T) {
	var declared *S = nil
	tt.RunSuite(t, declared) // want "TESTO012: suite argument is statically nil"
	short := (*Alias)(nil)
	tt.RunSubSuite(t, short) // want "TESTO012: suite argument is statically nil"
	var zero *S
	tt.RunSuite(t, zero) // want "TESTO012: suite argument is statically nil"
	var a, b *S
	tt.RunSuite(t, a) // want "TESTO012: suite argument is statically nil"
	tt.RunSuite(t, b) // want "TESTO012: suite argument is statically nil"
	yes, no := (*S)(nil), &S{}
	tt.RunSuite(t, yes) // want "TESTO012: suite argument is statically nil"
	tt.RunSuite(t, no)
	var first, second = &S{}, (*S)(nil)
	tt.RunSuite(t, first)
	tt.RunSuite(t, second)     // want "TESTO012: suite argument is statically nil"
	defer tt.RunSuite(t, zero) // want "TESTO012: suite argument is statically nil"
	go tt.RunSuite(t, zero)    // want "TESTO012: suite argument is statically nil"
	func() { tt.RunSuite(t, zero) /* want "TESTO012: suite argument is statically nil" */ }()
	stable := (*S)(nil)
	{
		stable := &S{}
		tt.RunSuite(t, stable)
	}
	tt.RunSuite(t, stable) // want "TESTO012: suite argument is statically nil"
}

var global *S

func negatives(t *tt.T, parameter *S, cond bool) {
	tt.RunSuite(t, parameter)
	tt.RunSuite(t, global)
	s := (*S)(nil)
	s = &S{}
	tt.RunSuite(t, s)
	branch := (*S)(nil)
	if cond {
		branch = &S{}
	}
	tt.RunSuite(t, branch)
	captured := (*S)(nil)
	fill := func() { captured = &S{} }
	fill()
	tt.RunSuite(t, captured)
	escaped := (*S)(nil)
	mutate(&escaped)
	tt.RunSuite(t, escaped)
	var assigned *S
	assigned = nil // Assignments after declarations are deliberately not tracked.
	tt.RunSuite(t, assigned)
	var call = makeSuite()
	tt.RunSuite(t, call)
	// Any reassignment suppresses inference, even after the call.
	later := (*S)(nil)
	tt.RunSuite(t, later)
	later = &S{}
	_ = later
	// Shadowing the predeclared nil must not create a diagnostic.
	nil := &S{}
	tt.RunSuite(t, (*S)(nil))
	shadow := nil
	tt.RunSuite(t, shadow)
}

func sideEffects(t *tt.T) {
	s := (*S)(nil)
	parent := func() *tt.T { s = &S{}; return t }
	tt.RunSuite(parent(), s)
	ranged := (*S)(nil)
	for _, ranged = range []*S{&S{}} {
		tt.RunSuite(t, ranged)
	}
	reused := (*S)(nil)
	reused, n := &S{}, 1
	_ = n
	tt.RunSuite(t, reused)
}

func makeSuite() *S { return &S{} }
func mutate(s **S)  { *s = &S{} }

// Zero values of suites passed by value are not nil.
type Value struct{ tt.Suite[*tt.T] }

func (Value) TestOK(t *tt.T) {}
func values(t *tt.T) {
	var s Value
	tt.RunSuite(t, s)
	tt.RunSubSuite(t, Value{})
}

func convert(s *S) *S { return &S{} }

type Generic[X any] struct{ tt.Suite[*tt.T] }

func (*Generic[X]) TestOK(t *tt.T) {}
