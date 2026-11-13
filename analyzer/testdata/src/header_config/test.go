package header_config

import (
	"github.com/ozontech/testo"
	a "github.com/ozontech/testo-allure"
)

type Context struct{ *a.PluginAllure }
type T = *Context
type Suite struct{ testo.Suite[T] }

func (*Suite) TestCustom(t T) {
	t.Title("title")
	t.Description("body") // want "TESTO025"
}
func (*Suite) TestValid(t T) {
	t.Title("title")

	t.Description("body")
}
func (*Suite) TestParallel(t T) {
	t.Parallel()
	t.Description("body") // want "TESTO025"
}
