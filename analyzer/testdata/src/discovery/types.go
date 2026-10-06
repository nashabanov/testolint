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
type EmptyIndirect struct{ *EmptyAlias } // want "TESTO_DISCOVERY: suite EmptyIndirect$"
type EmptyRecursive struct {             // want "TESTO_DISCOVERY: suite EmptyRecursive$"
	*EmptyRecursive
	EmptyIndirect
}
type NotEmbedded struct{ Base Empty }
type OtherBase struct{ Base } // want "TESTO_DISCOVERY: suite OtherBase$"
type AmbiguousBase struct {
	Empty
	OtherBase
}
type Alias = Active
type Active struct{ Base } // want "TESTO_DISCOVERY: suite Active$"
type Generic[P any] struct{ testo.Suite[T] }
type Indirect struct{ Empty } // want "TESTO_DISCOVERY: suite Indirect$"
type Unrelated struct{ faketesto.Suite[T] }
type NonStruct int

type Promoted struct { // want "TESTO_DISCOVERY: suite Promoted$"
	testo.Suite[T]
	Providers
}
type Providers struct{}

func (Providers) CasesAge() []int { return nil } // want "TESTO006"
