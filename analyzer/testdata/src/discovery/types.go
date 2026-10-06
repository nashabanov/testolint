package discovery

import (
	"faketesto"
	"github.com/ozontech/testo"
)

type T = *testo.T
type Base = testo.Suite[T]

type Empty struct{ testo.Suite[T] }    // want "TESTO_DISCOVERY: suite Empty$" "TESTO011"
type Pointer struct{ *testo.Suite[T] } // want "TESTO_DISCOVERY: suite Pointer$" "TESTO011"
type AliasedBase struct{ Base }        // want "TESTO_DISCOVERY: suite AliasedBase$" "TESTO011"
type EmptyAlias = Empty
type EmptyIndirect struct{ *EmptyAlias } // want "TESTO_DISCOVERY: suite EmptyIndirect$" "TESTO011"
type EmptyRecursive struct {             // want "TESTO_DISCOVERY: suite EmptyRecursive$" "TESTO011"
	*EmptyRecursive
	EmptyIndirect
}
type NotEmbedded struct{ Base Empty }
type OtherBase struct{ Base } // want "TESTO_DISCOVERY: suite OtherBase$" "TESTO011"
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

type Promoted struct { // want "TESTO_DISCOVERY: suite Promoted$" "TESTO011"
	testo.Suite[T]
	Providers
}
type Providers struct{}

func (Providers) CasesAge() []int { return []int{0} } // want "TESTO006"
