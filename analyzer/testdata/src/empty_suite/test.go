package emptysuite

import (
	"github.com/ozontech/testo"
	"promotionlib"
)

type T = *testo.T
type Empty struct{ testo.Suite[T] } // want `TESTO011: suite "Empty" contains no tests`
type EmptyAlias = Empty
type Hooks struct{ testo.Suite[T] } // want `TESTO011: suite "Hooks" contains no tests`
func (Hooks) BeforeAll(t T)         {}

type ProvidersOnly struct{ testo.Suite[T] } // want `TESTO011: suite "ProvidersOnly" contains no tests`
func (ProvidersOnly) CasesRole() []string   { return []string{"admin"} } // want "TESTO006"

type Valid struct{ testo.Suite[T] }

func (Valid) TestFoo(t T) {}

type Methods struct{}

func (*Methods) TestFoo(t T) {}

type Promoted struct {
	testo.Suite[T]
	Methods
}
type Indirect struct{ *Promoted }
type IndirectAlias = Indirect
type IndirectValue struct{ Promoted }
type IndirectEmpty struct{ *EmptyAlias } // want `TESTO011: suite "IndirectEmpty" contains no tests`
type Imported struct {
	testo.Suite[T]
	promotionlib.Methods[T]
}

type Malformed struct{ testo.Suite[T] } // want `TESTO011: suite "Malformed" contains no tests`
func (Malformed) Testfoo(t T)           {} // want "TESTO007"
func (Malformed) Testé(t T)             {} // want "TESTO007"

// Each naming boundary gets its own suite so another method cannot mask it.
type Exact struct{ testo.Suite[T] }

func (Exact) Test(t T) {}

type Digit struct{ testo.Suite[T] }

func (Digit) Test1(t T) {}

type Underscore struct{ testo.Suite[T] }

func (Underscore) Test_Foo(t T) {}

type Unicode struct{ testo.Suite[T] }

func (Unicode) TestÉ(t T) {}

type Chinese struct{ testo.Suite[T] }

func (Chinese) Test中(t T) {}

// Signature problems are handled by TESTO001; Testo still recognizes the name.
type BadSignature struct{ testo.Suite[T] }

func (BadSignature) TestFoo() {} // want "TESTO001"

type Left struct{}
type Right struct{}

func (Left) TestFoo(t T)  {}
func (Right) TestFoo(t T) {}

type Ambiguous struct { // want `TESTO011: suite "Ambiguous" contains no tests`
	testo.Suite[T]
	Left
	Right
}
type Shadowed struct { // want `TESTO011: suite "Shadowed" contains no tests`
	testo.Suite[T]
	Methods
	TestFoo int
}
