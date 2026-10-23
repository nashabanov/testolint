package execution_safety

import (
	tt "github.com/ozontech/testo"
	"testing"
)

type Ctx = tt.T
type Callback func(*Ctx)
type unrelated struct{}

func (unrelated) Cleanup(f func())            {}
func (unrelated) RunTest(t any, f func(*Ctx)) {}

type custom struct{ *Ctx }

func (*custom) Cleanup(f func()) {}

func callbacks(t *testing.T, unknown func(*Ctx)) {
	tt.RunTest[*Ctx](t, nil)     // want "TESTO019"
	tt.RunTest(t, Callback(nil)) // want "TESTO019"
	var zero func(*Ctx)
	tt.RunTest(t, zero)
	local := (func(*Ctx))(nil)
	run := tt.RunTest[*Ctx]
	run(t, local)
	tt.Test[*Ctx](nil)(t)                // want "TESTO019.*adapter"
	t.Run("nil", tt.Test(Callback(nil))) // want "TESTO019.*adapter"
	tt.RunTest(t, unknown)
	tt.RunTest(t, func(*Ctx) {})
	zero = func(*Ctx) {}
	tt.RunTest(t, zero) // reassignment anywhere suppresses all uses of zero
	changed := (func(*Ctx))(nil)
	_ = &changed
	tt.RunTest(t, changed)
	unrelated{}.RunTest(t, nil)
	_ = tt.Test[*Ctx](nil) // unused adapter does not execute its callback
	adapter := tt.Test[*Ctx](nil)
	_ = adapter // adapter tracking is outside scope
	run = func(any, func(*Ctx), ...any) bool { return true }
	run(t, nil)
}

func stable(t *testing.T) {
	var zero func(*Ctx)
	tt.RunTest(t, zero) // want "TESTO019"
	run := tt.RunTest[*Ctx]
	run(t, nil) // want "TESTO019"
}

func cleanup(t, other *Ctx, fake unrelated, c *custom) {
	t.Cleanup(func() {
		tt.Run(t, "bad", nil)                 // want "TESTO018"
		tt.Run[*Ctx](t, "bad", func(*Ctx) {}) // want "TESTO018"
		run := tt.Run[*Ctx]
		run(t, "bad", nil) // want "TESTO018"
		tt.Run(other, "other", nil)
		_ = func() { tt.Run(t, "later", nil) }
	})
	tt.Run(t, "nil is supported", nil)
	fake.Cleanup(func() { tt.Run(t, "unrelated", nil) })
	c.Cleanup(func() { tt.Run(c, "override", nil) })
	register := t.Cleanup
	register(func() { tt.Run(t, "unknown registration", nil) })
}

func changedContext(t, other *Ctx) {
	t.Cleanup(func() { tt.Run(t, "unknown identity", nil) })
	t = other
}
func escapedContext(t *Ctx) {
	_ = &t
	t.Cleanup(func() { tt.Run(t, "unknown identity", nil) })
}
