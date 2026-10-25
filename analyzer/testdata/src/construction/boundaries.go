package construction

import (
	"testing"

	"constructionlib"
	"github.com/ozontech/testo"
)

// A generic call has no concrete construction graph until instantiated.
// It may legally receive pointer fields; do not call an unknown T non-pointer.
type Deferred[X any] struct {
	*testo.T
	Field X
}

func generic[X any](t *testing.T) {
	testo.RunTest(t, func(*Deferred[X]) {})
}

type EmbeddedValue struct {
	testo.T // want "TESTO021: exported field .*T.* must be a pointer"
}
type ValueStruct struct {
	*testo.T
	Value Leaf // want "TESTO021: exported field .*Value.* must be a pointer"
}
type NamedPointer *Leaf
type NamedPointerT struct {
	*testo.T
	Value NamedPointer
}

type Imported struct {
	*testo.T
	Value *constructionlib.Invalid
}
type ImportedAlias = Imported
type ImportedCycle struct {
	*testo.T
	Value *constructionlib.Cycle
}

type GenericCycle[X any] struct {
	Next *GenericCycle[X] // want "TESTO022: recursive initialization.*Next"
}
type GenericCycleT[X any] struct {
	*testo.T
	Cycle *GenericCycle[X]
}

// The complete concrete graph is known even when root declarations are imported.
func boundaries(t *testing.T) {
	testo.RunTest(t, func(*GenericCycleT[int]) {})
	testo.RunTest(t, func(*GenericCycleT[string]) {}) // one diagnostic at the shared field declaration.
	testo.RunTest(t, func(*EmbeddedValue) {})
	testo.RunTest(t, func(*ValueStruct) {})
	testo.RunTest(t, func(*NamedPointerT) {})
	testo.RunTest(t, func(*Imported) {})                // want "TESTO021: exported field .*Bad.* must be a pointer"
	testo.RunTest(t, func(*Imported) {})                // do not repeat the same root's imported diagnostic.
	testo.RunTest(t, func(*ImportedAlias) {})           // aliases must not repeat imported diagnostics either.
	testo.RunTest(t, func(*ImportedCycle) {})           // want "TESTO022: recursive initialization.*Next"
	testo.RunTest(t, func(*constructionlib.Context) {}) // want "TESTO021: exported field .*Bad.* must be a pointer"
}
