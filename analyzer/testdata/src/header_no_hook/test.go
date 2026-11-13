package header_no_hook

import (
	"github.com/ozontech/testo"
	a "github.com/ozontech/testo-allure"
)

type Context struct{ *a.PluginAllure }
type T = *Context
type Suite struct{ testo.Suite[T] }

func (*Suite) BeforeEach(t T)  { t.Tags("shared") }
func (*Suite) TestMissing(t T) {} // want "TESTO028"
func (*Suite) TestValid(t T)   { t.Tags("local") }
