package emptycases

import (
	"emptycaseslib"
	"github.com/ozontech/testo"
)

type T = *testo.T
type Names []string
type NamesAlias = Names
type Suite struct{ testo.Suite[T] }

func (Suite) TestAll(t T, p struct {
	Nil, Empty, NonEmpty, Call, Variable, Named, Alias, Made, Flow, Shadow string
}) {
}

func (Suite) CasesNil() []string      { return nil }            // want "TESTO010: CasesNil always returns an empty case set"
func (*Suite) CasesEmpty() []string   { return []string{} }     // want "TESTO010: CasesEmpty always returns an empty case set"
func (Suite) CasesNamed() Names       { return Names{} }        // want "TESTO010: CasesNamed always returns an empty case set"
func (Suite) CasesAlias() NamesAlias  { return (NamesAlias{}) } // want "TESTO010: CasesAlias always returns an empty case set"
func (Suite) CasesNonEmpty() []string { return []string{"admin"} }
func (Suite) CasesCall() []string     { return loadRoles() }
func loadRoles() []string             { return nil }
func (Suite) CasesVariable() []string {
	var roles []string
	return roles
}
func (Suite) CasesMade() []string { return make([]string, 0) }
func (Suite) CasesFlow() []string {
	if true {
		return nil
	}
	return []string{}
}
func (Suite) CasesShadow() []string {
	nil := []string{"admin"}
	return nil
}

// Invalid signatures and malformed names do not receive cascading TESTO010.
func (Suite) CasesArg(x int) []string         { return nil }      // want "TESTO005" "TESTO006"
func (Suite) CasesResults() ([]string, error) { return nil, nil } // want "TESTO005" "TESTO006"
func (Suite) CasesMap() map[string]int        { return nil }      // want "TESTO005" "TESTO006"
func (Suite) CasesNone()                      {}                  // want "TESTO005" "TESTO006"
func (Suite) Casesbad() []string              { return nil }      // want "TESTO008"
func (Suite) Cases() []string                 { return nil }

// Local promoted providers retain their AST; duplicates are suppressed.
type Providers struct{}

func (Providers) CasesRole() []string { return nil } // want "TESTO010: CasesRole always returns an empty case set"
type Local struct {
	testo.Suite[T]
	Providers
}

func (Local) TestRole(t T, p struct{ Role string }) {}

type LocalChild struct{ *Local }

// Instantiated local methods use their origin declaration.
type GenericProviders[P any] struct{}

func (GenericProviders[P]) CasesValue() []P { return []P{} } // want "TESTO010: CasesValue always returns an empty case set"
type GenericSuite struct {
	testo.Suite[T]
	GenericProviders[string]
}

func (GenericSuite) TestValue(t T, p struct{ Value string }) {}

// Imported providers have no AST in this pass and are skipped by TESTO010.
type Imported struct {
	testo.Suite[T]
	emptycaseslib.Providers
}

func (Imported) TestRole(t T, p struct{ Role string }) {}
