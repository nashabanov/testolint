package run_discovery

import (
	"faketesto"
	"github.com/ozontech/testo"
)

type S struct{ testo.Suite[*testo.T] }

func (S) TestOK(t *testo.T) {}

type Alias = S

func ordinary(t *testo.T) {
	testo.RunSuite(t, S{})                                  // want `TESTO_RUN_DISCOVERY: RunSuite run_discovery.S`
	testo.RunSuite[*S](t, &S{})                             // want `TESTO_RUN_DISCOVERY: RunSuite \*run_discovery.S`
	testo.RunSuite[Alias, *testo.T](t, Alias{})             // want `TESTO_RUN_DISCOVERY: RunSuite run_discovery.(Alias|S)`
	testo.RunSubSuite(t, &S{})                              // want `TESTO_RUN_DISCOVERY: RunSubSuite \*run_discovery.S`
	testo.RunSubSuite[*S, *testo.T, *testo.T](t, (*S)(nil)) // want `TESTO_RUN_DISCOVERY: RunSubSuite \*run_discovery.S` "TESTO012: suite argument is statically nil"
	(testo.RunSuite[S])(t, S{})                             // want `TESTO_RUN_DISCOVERY: RunSuite run_discovery.S`
	run := testo.RunSuite[S]
	run(t, S{}) // want `TESTO_RUN_DISCOVERY: RunSuite run_discovery.S`
	var sub = testo.RunSubSuite[S, *testo.T, *testo.T]
	sub(t, S{}) // want `TESTO_RUN_DISCOVERY: RunSubSuite run_discovery.S`
	var inferred func(any, S, ...any) bool = testo.RunSuite
	inferred(t, S{}) // want `TESTO_RUN_DISCOVERY: RunSuite run_discovery.S`
	RunSuite(t, S{})
	faketesto.RunSuite(t, S{})
	(Other{}).RunSubSuite(t, S{})
	changed := testo.RunSuite[S]
	changed = RunSuite
	changed(t, S{})
	escaped := testo.RunSuite[S]
	_ = &escaped
	escaped(t, S{})
	var parameter func(any, S, ...any) bool = RunSuite
	parameter(t, S{})
	global(t, S{})
	// Writes in closures and range clauses invalidate the binding too.
	captured := testo.RunSuite[S]
	change := func() { captured = RunSuite }
	change()
	captured(t, S{})
	ranged := testo.RunSuite[S]
	for _, ranged = range []func(any, S, ...any) bool{RunSuite} {
		ranged(t, S{})
	}
	reused := testo.RunSuite[S]
	reused, x := RunSuite, 1
	_ = x
	reused(t, S{})
	// A shadowed variable does not invalidate the outer binding.
	stable := testo.RunSuite[S]
	{
		stable := RunSuite
		stable(t, S{})
	}
	stable(t, S{}) // want "TESTO_RUN_DISCOVERY: RunSuite run_discovery.S"
	// Aliases of aliases are intentionally outside the direct-alias subset.
	indirect := stable
	indirect(t, S{})
	defer testo.RunSuite(t, S{}) // want "TESTO_RUN_DISCOVERY: RunSuite run_discovery.S"
	go testo.RunSubSuite(t, S{}) // want "TESTO_RUN_DISCOVERY: RunSubSuite run_discovery.S"
}

var global = testo.RunSuite[S]

func RunSuite(t any, s S, options ...any) bool { return true }

type Other struct{}

func (Other) RunSubSuite(t any, s S) {}
