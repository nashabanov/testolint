package parameterization

import "github.com/ozontech/testo"

type T struct{ *testo.T }

type Suite struct {
	testo.Suite[T]
}

// valid: provider exists and type matches
func (Suite) TestValid(
	t T,
	p struct {
		UserID int
	},
) {
}

func (Suite) CasesUserID() []int {
	return []int{0}
}

// TESTO003: missing provider
func (Suite) TestMissingCases(
	t T,
	p struct {
		Role string // want "TESTO003"
	},
) {
}

// TESTO004: provider type mismatch
func (Suite) TestWrongCasesType(
	t T,
	p struct {
		Age int
	},
) {
}

func (Suite) CasesAge() []string { // want "TESTO004"
	return []string{"admin"}
}

// TESTO006: orphan provider
func (Suite) CasesUnused() []bool { // want "TESTO006"
	return []bool{true}
}

// A concrete element type is assignable to an interface field, but not identical.
func (Suite) TestAssignable(t T, p struct{ Value any }) {}
func (Suite) CasesValue() []int                         { return []int{0} }

// The reverse direction is not assignable: any cannot be assigned to int.
func (Suite) TestNotAssignable(t T, p struct{ Count int }) {}
func (Suite) CasesCount() []any                            { return []any{0} } // want "TESTO004"
