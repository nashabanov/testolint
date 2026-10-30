package selection

import (
	"github.com/ozontech/testo"
	"testing"
)

type T = *testo.T
type S struct{ testo.Suite[T] }

func (S) TestOK(T)           {}
func (S) CasesUnused() []int { return []int{1} }

type Invalid struct {
	*testo.T
	Count int
}
type Loop struct{ Next *Loop } // want "TESTO022"
type Recursive struct {
	*testo.T
	*Loop
	Bad int // A cycle suppresses this field even when TESTO022 is disabled.
}

func calls(t *testing.T, ctx T) {
	var s *S
	testo.RunSuite(t, s) // want "TESTO012"
	testo.RunTest[T](t, nil)
	ctx.Cleanup(func() { testo.Run(ctx, "cleanup", nil) }) // want "TESTO018"
	testo.RunTest(t, func(*Invalid) {})
	testo.RunTest(t, func(*Recursive) {})
	testo.RunTest(t, func(T) {})
	testo.RunSuite(t, S{})
}
