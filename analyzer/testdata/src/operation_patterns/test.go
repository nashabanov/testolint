package operation_patterns

import (
	"github.com/ozontech/testo"
	a "github.com/ozontech/testo-allure"
	alias "operation_patterns/remote"
)

type Context struct{ *a.PluginAllure }
type T = *Context
type Suite struct{ testo.Suite[T] }
type Client struct{}
type Alias = Client

func (Client) Send()    {}
func (*Client) Update() {}
func (Client) Close()   {}
func Check()            {}
func (*Suite) TestCalls(t T) {
	var c Alias
	c.Close()
	c.Send()     // want "TESTO029"
	c.Update()   // want "TESTO029"
	alias.Send() // want "TESTO029"
	Check()      // want "TESTO031"
	a.Step(t, "covered", func(t T) { c.Send(); alias.Send(); Check() })
}
