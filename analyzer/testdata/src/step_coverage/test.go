package step_coverage

import (
	"github.com/ozontech/testo"
	a "github.com/ozontech/testo-allure"
)

type Context struct{ *a.PluginAllure }
type T = *Context
type Suite struct{ testo.Suite[T] }

func verify() {}
func action() {}
func (*Suite) TestDirect(t T) {
	verify() // want "TESTO031"
	action() // want "TESTO029"
}
func (*Suite) TestCovered(t T) {
	a.Step(t, "step", func(t T) { verify(); action(); helper(t) })
	a.Setup(t, "setup", func(t T) { verify(); action() })
	a.Step(t, "outer", func(t T) { a.Step(t, "inner", func(t T) { verify() }) })
}
func helper(t T) { a.Step(t, "covered", func(t T) { verify(); action() }); verify(); action() }
func covered(t T) {
	a.Step(t, "one", func(t T) { verify(); action() })
	a.Step(t, "two", func(t T) { verify() })
}
func data() int { return 1 }
func chain(t T) { helper(t) }
func (*Suite) TestHelpers(t T) {
	helper(t) // want "TESTO031" "TESTO029"
	covered(t)
	_ = data()
	chain(t)
	a.Step(t, "wrapped", func(t T) { helper(t) })
	helper(t) // want "TESTO031" "TESTO029"
}
func (*Suite) TestClosures(t T) {
	_ = func() { verify(); action() }
	a.TearDown(t, "cleanup", func(t T) { verify(); action() })
}

type Plain struct{ testo.Suite[*testo.T] }

func (*Plain) TestIgnored(t *testo.T) { verify(); action() }
