package execution_safety

import (
	. "github.com/ozontech/testo"
	"testing"
)

type contextAlias = *T
type functionAlias = func(contextAlias)

type owner struct{ Suite[*T] }

func (*owner) TestA(t *T) {
	t.Cleanup(func() { Run(t, "pointer receiver", nil) }) // want "TESTO018"
}

type valueOwner struct{ Suite[*T] }

func (valueOwner) TestA(t *T) {
	t.Cleanup(func() { Run(t, "value receiver", nil) }) // want "TESTO018"
}

func moreCallbacks(t *testing.T) {
	var f functionAlias = nil
	RunTest(t, f) // want "TESTO019"
	g := functionAlias(nil)
	RunTest(t, g) // want "TESTO019"
	adapter := Test[*T]
	adapter(nil)(t)          // want "TESTO019.*adapter"
	(Test[*T](nil))(t)       // want "TESTO019.*adapter"
	RunTest[*custom](t, nil) // want "TESTO019"
	RunTest(t, functionAlias(func(contextAlias) {}))
	fakeRunTest := func(any, func(*T)) {}
	fakeRunTest(t, nil)
	var globalCopy = globalCallback
	RunTest(t, globalCopy)
	var f2 func(*T)
	for _, f2 = range []func(*T){func(*T) {}} {
	}
	RunTest(t, f2)
}

var globalCallback func(*T)

func cleanupBoundaries(t *T, yes bool) {
	callback := func() { Run(t, "stored", nil) }
	t.Cleanup(callback)
	t.Cleanup(func() {
		if yes {
			Run(t, "branch", nil)
		}
		Run(t, "after branch", nil)
	})
	t.Cleanup(func() {
		defer Run(t, "deferred", nil)
		go Run(t, "goroutine", nil)
	})
	t.Cleanup(func() {
		(func() { Run(t, "nested invocation", nil) })()
	})
	alias := t
	t.Cleanup(func() { Run(alias, "context alias", nil) })
	(*T).Cleanup(t, func() { Run(t, "method expression", nil) })
}
