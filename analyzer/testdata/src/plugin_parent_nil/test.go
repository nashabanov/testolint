package pluginparentnil

import (
	"github.com/ozontech/testo"
	p "github.com/ozontech/testo/testoplugin"
)

type Bad struct{}

func (*Bad) Plugin(parent p.Plugin, _ ...p.Option) p.Spec {
	if parent == nil { // want "TESTO024:.*asserted concrete pointer"
	}
	if nil != (parent) { // want "TESTO024:.*asserted concrete pointer"
	}
	return p.Spec{}
}

type Good struct{}

func (*Good) Plugin(parent p.Plugin, _ ...p.Option) p.Spec {
	previous := parent.(*Good)
	if previous == nil {
	}
	return p.Spec{}
}

type Changed struct{}

func (*Changed) Plugin(parent p.Plugin, _ ...p.Option) p.Spec {
	parent = nil
	if parent == nil {
	}
	return p.Spec{}
}

type Addressed struct{}

func (*Addressed) Plugin(parent p.Plugin, _ ...p.Option) p.Spec {
	_ = &parent
	if parent == nil {
	}
	return p.Spec{}
}

type Manual struct{}

func (*Manual) Plugin(parent p.Plugin, _ ...p.Option) p.Spec {
	if parent == nil {
	}
	return p.Spec{}
}

type Unused struct{}

func (*Unused) Plugin(parent p.Plugin, _ ...p.Option) p.Spec {
	if parent == nil {
	}
	return p.Spec{}
}

type Wrong struct{}

func (*Wrong) Plugin(parent any, _ ...p.Option) p.Spec {
	if parent == nil {
	}
	return p.Spec{}
}

type Value struct{}

func (Value) Plugin(parent p.Plugin, _ ...p.Option) p.Spec {
	if parent != nil { // want "TESTO024:"
	}
	return p.Spec{}
}

type Alias = Bad
type Context struct {
	*testo.T
	Bad       *Alias
	Good      *Good
	Changed   *Changed
	Addressed *Addressed
	Manual    *Manual
	Wrong     *Wrong
	Value     *Value
	hidden    *Unused
}

func execute() {
	testo.RunTest(nil, func(*Context) {})
	testo.RunTest(nil, func(*Context) {}) // no duplicate diagnostics
	new(Manual).Plugin(nil)
	_ = testo.Test(func(*UnusedContext) {}) // unexecuted adapter
}

type UnusedContext struct {
	*testo.T
	Unused *Unused
}

// The context type itself is excluded from Testo's plugin list, including
// contexts whose Plugin method has a value receiver.
type RootContext struct{ *testo.T }

func (RootContext) Plugin(parent p.Plugin, _ ...p.Option) p.Spec {
	if parent == nil {
	}
	return p.Spec{}
}

func rootIsNotPlugin() { testo.RunTest(nil, func(*RootContext) {}) }

type Nested struct{ Bad *Alias }
type NestedContext struct {
	*testo.T
	Nested *Nested
}

func nested() {
	run := testo.RunTest[*NestedContext]
	run(nil, func(*NestedContext) {})
}

type Generic[P any] struct{ hidden P }

func (*Generic[P]) Plugin(parent p.Plugin, _ ...p.Option) p.Spec {
	if parent == nil { // want "TESTO024:"
	}
	return p.Spec{}
}

type GenericContext struct {
	*testo.T
	Generic *Generic[int]
}

func generic() { testo.RunTest(nil, func(*GenericContext) {}) }

// A blocking wiring error suppresses dependent parent checks.
type Blocked struct{}

func (*Blocked) Plugin(parent p.Plugin, _ ...p.Option) p.Spec {
	if parent == nil {
	}
	return p.Spec{}
}

type BlockedContext struct {
	*testo.T
	Blocked *Blocked
	Invalid int // want "TESTO021:"
}

func blocked() { testo.RunTest(nil, func(*BlockedContext) {}) }

// The embedded plugin is traversed; its promoted method adds no duplicates.
type Promoted struct{ *Bad }
type PromotedContext struct {
	*testo.T
	Promoted *Promoted
}

func promoted() { testo.RunTest(nil, func(*PromotedContext) {}) }
