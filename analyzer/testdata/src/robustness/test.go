package robustness

import "github.com/ozontech/testo"

type T = *testo.T
type Base = testo.Suite[T]
type Number int
type Numbers []Number
type Params struct{ Age Number }
type ParamsAlias = Params

type Providers struct{}

func (Providers) CasesAge() Numbers { return Numbers{0} }

type Suite struct {
	Base
	Providers
}

func (*Suite) TestPromoted(t T, p ParamsAlias) {}

// Testo ignores Cases entirely, including its signature.
func (*Suite) Cases(x int) string { return "" }

type Tests struct{}

func (Tests) TestInherited(t T, p Params) {}

type Inherited struct {
	testo.Suite[T]
	Tests
}

func (Inherited) CasesAge() Numbers { return Numbers{0} }

type PointerSuite struct{ *testo.Suite[T] }

func (PointerSuite) TestNamed(t T, p Params) {}
func (*PointerSuite) CasesAge() Numbers      { return Numbers{0} }

type Embedded struct{ Number }
type EmbeddedSuite struct{ testo.Suite[T] }

func (EmbeddedSuite) TestEmbedded(t T, p Embedded) {}
func (EmbeddedSuite) CasesNumber() Numbers         { return Numbers{0} }

// Importing Testo does not turn unrelated types into suites.
type Unrelated struct{}

func (Unrelated) TestInvalid()       {}
func (Unrelated) CasesInvalid(x int) {}

type Shared struct{ Missing int } // want `TESTO003: parameter "Missing" requires CasesMissing`
type Broken struct{ testo.Suite[T] }

func (Broken) TestOne(t T, p Shared)                      {}
func (Broken) TestTwo(t T, p Shared)                      {}
func (Broken) TestMismatchOne(t T, p struct{ Age int })   {}
func (Broken) TestMismatchTwo(t T, p struct{ Age int })   {}
func (Broken) CasesAge() []string                         { return []string{"admin"} } // want `TESTO004: CasesAge provides string, but parameter "Age" expects int`
func (Broken) TestBadProvider(t T, p struct{ Value int }) {}
func (Broken) CasesValue(x int) []string                  { return []string{"admin"} } // want "TESTO005"
func (Broken) TestMalformed(t T, p struct { // want "TESTO001"
	Absent  int
	private int
}, extra int) {
}
func (Broken) TestWrongT(t string, p struct{ Absent int }) {} // want "TESTO001"

// Generic declarations are skipped rather than comparing different receiver
// instantiations as separate suites.
type Generic[P any] struct{ testo.Suite[T] }

func (Generic[P]) TestGeneric(t T, p struct{ Value P }) {}
func (Generic[P]) CasesValue() []P                      { return []P{*new(P)} }

// Variadic arguments are not valid test or hook signatures.
type Variadic struct{ testo.Suite[T] }

func (Variadic) TestVariadic(t ...T) {} // want "TESTO001"
func (Variadic) BeforeEach(t ...T)   {} // want "TESTO002"

type PrivateParams struct{ private int } // want "TESTO009"
type PrivateSuite struct{ testo.Suite[T] }

func (PrivateSuite) TestOne(t T, p PrivateParams) {}
func (PrivateSuite) TestTwo(t T, p PrivateParams) {}
