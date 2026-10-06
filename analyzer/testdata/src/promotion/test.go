package promotion

import (
	"github.com/ozontech/testo"
	"promotionlib"
)

type T = *testo.T
type Base struct{ testo.Suite[T] }
type BaseAlias = Base
type Indirect struct{ BaseAlias }
type Pointer struct{ *BaseAlias }

func (Indirect) TestInvalid()  {} // want "TESTO001"
func (*Pointer) TestValid(t T) {}
func (*Pointer) BeforeEach()   {} // want "TESTO002"

type Methods struct{}

func (Methods) TestValid(t T, p struct{ Role string }) {}
func (*Methods) CasesRole() []string                   { return nil }
func (Methods) BeforeEach(t T)                         {}
func (Methods) TestBad()                               {}             // want "TESTO001"
func (Methods) AfterEach()                             {}             // want "TESTO002"
func (Methods) Testbad(t T)                            {}             // want "TESTO007"
func (Methods) Casesbad() []int                        { return nil } // want "TESTO008"
func (Methods) CasesUnused() []int                     { return nil } // want "TESTO006"

type Value struct {
	Base
	Methods
}
type Ptr struct {
	*Base
	*Methods
}
type Deep struct{ *Ptr }

type BadProviders struct{}

func (BadProviders) CasesRole(x int) []string { return nil } // want "TESTO005"
type Bad struct {
	Base
	BadProviders
}

func (Bad) TestUser(t T, p struct{ Role string }) {}

type MismatchProviders struct{}

func (MismatchProviders) CasesRole() []int { return nil } // want "TESTO004"
type Mismatch struct {
	Base
	MismatchProviders
}

func (Mismatch) TestUser(t T, p struct{ Role string }) {}

type ParamsTests struct{}

func (ParamsTests) TestParams(t T, p struct {
	Missing int // want "TESTO003"
	private int // want "TESTO009"
}) {
}

type ParamsSuite struct {
	Base
	ParamsTests
}

// Selection must follow Go shadowing and ambiguity, not collect all embedded methods.
type Shadow struct {
	Base
	Methods
}

func (Shadow) TestBad(t T)        {}
func (Shadow) AfterEach(t T)      {}
func (Shadow) Testbad(t T)        {}             // want "TESTO007"
func (Shadow) Casesbad() []int    { return nil } // want "TESTO008"
func (Shadow) CasesUnused() []int { return nil } // want "TESTO006"
type Left struct{}
type Right struct{}

func (Left) TestAmbiguous()  {}
func (Right) TestAmbiguous() {}

type Ambiguous struct {
	Base
	Left
	Right
}

// Recursive embedding is legal through pointers and must terminate.
type Recursive struct {
	*Recursive
	Base
}

func (Recursive) TestRecursive(t T) {}

type CycleA struct{ *CycleB }
type CycleB struct{ *CycleA }

func (CycleA) TestIgnored() {}

// Imported and instantiated methods have no local FuncDecl.
type Imported struct{ *promotionlib.Base }             // want "TESTO005"
func (Imported) TestRole(t T, p struct{ Role string }) {}

type GenericMethods struct {
	Base
	promotionlib.Methods[T]
}
