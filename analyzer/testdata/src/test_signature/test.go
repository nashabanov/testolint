package testsignature

import "github.com/ozontech/testo"

type T struct{ *testo.T }

type Suite struct {
	testo.Suite[T]
}

func (Suite) TestValid(t T) {}

func (Suite) TestValidParams(
	t T,
	p struct {
		ID int
	},
) {
}

func (Suite) TestMissingT() {} // want "TESTO001"

func (Suite) TestTooMany(t T, p struct{}, extra int) {} // want "TESTO001"

func (Suite) TestWrongParams(t T, p string) {} // want "TESTO001"

func (Suite) TestReturns(t T) any { return 0 } // want "TESTO001"

type WrongT struct{ *testo.T }

func (Suite) TestWrongT( // want "TESTO001"
	t WrongT,
	p struct{},
) {
}
