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
	return nil
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
	return nil
}

// TESTO006: orphan provider
func (Suite) CasesUnused() []bool { // want "TESTO006"
	return nil
}
