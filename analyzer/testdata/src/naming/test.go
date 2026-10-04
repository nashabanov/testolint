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

// Testo ignores the empty provider suffix.
func (Suite) Cases() []int {
	return nil
}

func (Suite) Test(t T)     {}
func (Suite) TestFoo(t T)  {}
func (Suite) Test1(t T)    {}
func (Suite) Test_Foo(t T) {}
func (Suite) Testfoo(t T)  {} // want "TESTO007"
func (Suite) TestÉ(t T)    {}
func (Suite) Testé(t T)    {} // want "TESTO007"
func (Suite) Test中(t T)    {}

func (Suite) CasesFoo() []int  { return nil } // want "TESTO006"
func (Suite) Cases1() []int    { return nil } // want "TESTO006"
func (Suite) Cases_Foo() []int { return nil } // want "TESTO006"
func (Suite) Casesfoo() []int  { return nil } // want "TESTO008"
func (Suite) CasesÉ() []int    { return nil } // want "TESTO006"
func (Suite) Casesé() []int    { return nil } // want "TESTO008"
func (Suite) Cases中() []int    { return nil } // want "TESTO006"
