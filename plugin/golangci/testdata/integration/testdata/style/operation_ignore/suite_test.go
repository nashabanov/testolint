package operation_ignore

import (
	"github.com/ozontech/testo"
	a "github.com/ozontech/testo-allure"
	"testolint-style-integration/operation_ignore/remote"
)

type Context struct {
	*testo.T
	*a.PluginAllure
}
type T = *Context
type Suite struct{ testo.Suite[T] }

func action() int { return 1 }
func verify()     {}
func ignored(int) { action(); verify() }
func helper()     { ignored(0) }
func uncovered()  { action() }
func (*Suite) TestIgnored(t T) {
	a.Step(t, "assert", func(t T) { verify() })
	ignored(0)
	remote.Ignored()
	helper()
}
func (*Suite) TestArgument(t T) {
	ignored(action()) // want "TESTO029"
}
func (*Suite) TestArgumentPhase(t T) {
	a.Step(t, "assert", func(t T) { verify() })
	a.Step(t, "argument", func(t T) {
		ignored(action()) // want "TESTO030"
	})
}
func (*Suite) TestUncovered(t T) {
	uncovered()     // want "TESTO029"
	remote.Action() // want "TESTO029"
}
