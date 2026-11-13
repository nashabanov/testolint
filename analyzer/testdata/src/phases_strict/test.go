package phases_strict

import (
	"github.com/ozontech/testo"
	a "github.com/ozontech/testo-allure"
)

type Context struct{ *a.PluginAllure }
type T = *Context
type Suite struct{ testo.Suite[T] }

func action(...any) {}
func verify()       {}
func (*Suite) TestGood(t T) {
	a.Setup(t, "setup", func(t T) { action(); verify() })
	action()
	action()
	verify()
	verify()
}
func (*Suite) TestRepeat(t T) {
	action()
	verify()
	action() // want "TESTO030"
	verify()
}
func (*Suite) TestLateSetup(t T) {
	action()
	verify()
	a.Setup(t, "late", func(t T) { action(); verify() }) // want "TESTO030"
}
func (*Suite) TestMixedStep(t T) {
	a.Step(t, "mixed", func(t T) {
		action()
		verify()
		action() // want "TESTO030"
	})
}
func helper(t T) { action(); verify(); action() }
func (*Suite) TestHelper(t T) {
	helper(t) // want "TESTO030"
}
func (*Suite) TestUnknownBranch(t T) {
	if condition {
		verify()
	} else {
		action()
	}
	action()
	verify()
	action() // want "TESTO030"
}

var condition bool

func (*Suite) TestOnlyAssert(t T) { verify() }
func (*Suite) TestNeutral(t T)    { action(); verify(); x := 1; _ = x + 1 }

func maybeVerify() bool { verify(); return true }
func (*Suite) TestShortCircuit(t T) {
	_ = condition && maybeVerify()
	action()
	verify()
	action() // want "TESTO030"
}
func (*Suite) TestShortCircuitOr(t T) {
	_ = condition || maybeVerify()
	action()
}

func (*Suite) TestLateRequestData(t T) {
	action()
	request := struct{ ID int }{ID: 42}
	action(request)
}
func (*Suite) TestRequestData(t T) {
	request := struct{ ID int }{ID: 42}
	action(request)
	verify()
}

func (*Suite) TestNativeRequire(t T) {
	t.Require().Equal(1, 1)
	action() // want "TESTO030"
}
func (*Suite) TestNativeAssert(t T) {
	checks := t.Assert().Masked()
	checks.Equal(1, 1)
	action() // want "TESTO030"
}
func (*Suite) TestNativeConstructor(t T) {
	_ = t.Require().Hidden()
	_ = t.Assert().Masked()
	action()
}
func (*Suite) TestNativeSetup(t T) {
	a.Setup(t, "precondition", func(t T) { t.Require().Equal(1, 1) })
	action()
}
