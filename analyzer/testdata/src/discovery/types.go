package discovery

import (
	"faketesto"
	"github.com/ozontech/testo"
)

type T = *testo.T
type Base = testo.Suite[T]

type Empty struct{ testo.Suite[T] }    // want "TESTO_DISCOVERY: suite Empty$"
type Pointer struct{ *testo.Suite[T] } // want "TESTO_DISCOVERY: suite Pointer$"
type AliasedBase struct{ Base }        // want "TESTO_DISCOVERY: suite AliasedBase$"
type EmptyAlias = Empty
type Alias = Active
type Active struct{ Base } // want "TESTO_DISCOVERY: suite Active$"
type Generic[P any] struct{ testo.Suite[T] }
type Indirect struct{ Empty }
type Unrelated struct{ faketesto.Suite[T] }
type NonStruct int

type Promoted struct { // want "TESTO_DISCOVERY: suite Promoted$"
	testo.Suite[T]
	Providers
}
type Providers struct{}

func (Providers) CasesAge() []int { return nil }
