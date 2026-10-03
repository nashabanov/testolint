package naming

import "github.com/ozontech/testo"

type T struct{ *testo.T }

type Suite struct {
	testo.Suite[T]
}

func (Suite) TestValid(t T) {}

func (Suite) Testinvalid(t T) {} // want "TESTO007"

func (Suite) CasesValid() []int { // want "TESTO006"
	return nil
}

func (Suite) Casesinvalid() []int { // want "TESTO008"
	return nil
}

// A provider without a suffix is reported only by the naming rule.
func (Suite) Cases() []int { // want "TESTO008"
	return nil
}
