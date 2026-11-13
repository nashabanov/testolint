package required_header

import (
	"github.com/ozontech/testo"
	a "github.com/ozontech/testo-allure"
)

type Context struct{ *a.PluginAllure }
type T = *Context
type Suite struct{ testo.Suite[T] }

func (*Suite) TestMissing(t T)     {} // want "TESTO028"
func (*Suite) TestValid(t T)       { t.Title("title"); t.Tags("smoke") }
func (*Suite) TestAlternative(t T) { t.Titlef("%s", "title"); t.Tags("smoke") }

type Base struct{ testo.Suite[T] }

func (*Base) BeforeEach(t T) { t.Tags("shared") }

type Inherited struct{ Base }

func (*Inherited) TestGood(t T)    { t.Title("title") }
func (*Inherited) TestMissing(t T) {} // want "TESTO028"
type Shadow struct{ Base }

func (*Shadow) BeforeEach(t T)  {}
func (*Shadow) TestMissing(t T) { t.Title("title") } // want "TESTO028"
type Conditional struct{ testo.Suite[T] }

func (*Conditional) BeforeEach(t T) {
	if condition {
		t.Tags("maybe")
	}
}
func (*Conditional) TestUnknown(t T) { t.Title("title") }

var condition bool

type Closure struct{ testo.Suite[T] }

func (*Closure) BeforeEach(t T)  { _ = func() { t.Tags("never") } }
func (*Closure) TestMissing(t T) { t.Title("title") } // want "TESTO028"
type Plain struct{ testo.Suite[*testo.T] }

func (*Plain) TestIgnored(t *testo.T) {}
