package step_coverage

import (
	"github.com/ozontech/testo"
	a "github.com/ozontech/testo-allure"
	api "step_coverage/remote"
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

type Worker struct{}

func (Worker) work(t T) { action(); verify() }

type Promoted struct {
	testo.Suite[T]
	Worker
}

func (s *Promoted) TestHelper(t T) {
	s.work(t) // want "TESTO029" "TESTO031"
}
func genericHelper[Context any](t Context) { action(); verify() }
func (*Suite) TestGenericHelper(t T) {
	genericHelper[T](t) // want "TESTO029" "TESTO031"
}
func (*Suite) TestAsyncDeferredUnknown(t T) {
	go helper(t)
	defer helper(t)
	go a.Step(t, "async", func(t T) { action(); verify() })
	defer a.Step(t, "deferred", func(t T) { action(); verify() })
}

func (*Suite) TestImported(t T) {
	new(api.Client).Send() // want "TESTO029"
	api.Send[T](t)         // want "TESTO029"
	api.Helper()           // opaque: no configured operation at this callsite
	a.Step(t, "covered imported", func(t T) { new(api.Client).Send() })
}
