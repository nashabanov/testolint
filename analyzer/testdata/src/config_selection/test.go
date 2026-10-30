package selection

import (
	"github.com/ozontech/testo"
	"testing"
)

type T = *testo.T
type S struct{ testo.Suite[T] }

func (S) TestOK(T)           {}
func (S) CasesUnused() []int { return []int{1} } // want "TESTO006"

type Invalid struct {
	*testo.T
	Count int // want "TESTO021"
}
type Loop struct{ Next *Loop }
type Recursive struct {
	*testo.T
	*Loop
	Bad int // A cycle suppresses this field even when TESTO022 is disabled.
}

func calls(t *testing.T, ctx T) {
	var s *S
	testo.RunSuite(t, s)
	testo.RunTest[T](t, nil) // want "TESTO019"
	ctx.Cleanup(func() { testo.Run(ctx, "cleanup", nil) })
	testo.RunTest(t, func(*Invalid) {})
	testo.RunTest(t, func(*Recursive) {})
	testo.RunTest(t, func(T) {})
	testo.RunSuite(t, S{})
}
