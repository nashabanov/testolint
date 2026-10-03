package casessignature

import "github.com/ozontech/testo"

type T struct{ *testo.T }

type Suite struct {
	testo.Suite[T]
}

func (Suite) CasesValid() []int { // want "TESTO006"
	return nil
}

func (Suite) CasesWithArg(x int) []int { // want "TESTO005" "TESTO006"
	return nil
}

func (Suite) CasesNoResult() {} // want "TESTO005" "TESTO006"

func (Suite) CasesTwoResults() ([]int, error) { // want "TESTO005" "TESTO006"
	return nil, nil
}

func (Suite) CasesNotSlice() int { // want "TESTO005" "TESTO006"
	return 0
}
