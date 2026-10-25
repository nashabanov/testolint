package construction

import (
	"testing"

	framework "github.com/ozontech/testo"
)

type Invalid struct {
	*framework.T
	Count int // want "TESTO021: exported field .*Count.* must be a pointer"
}

type InvalidAlias = Invalid
type BadSuite struct{ framework.Suite[*InvalidAlias] }

func (BadSuite) TestOK(*InvalidAlias) {}

type Nested struct {
	Label string // want "TESTO021: exported field .*Label.* must be a pointer"
}
type NestedT struct {
	*framework.T
	*Nested
}

type Loop struct {
	Next *Loop // want "TESTO022: recursive initialization.*Next"
}
type Recursive struct {
	*framework.T
	*Loop
	Bad int // no TESTO021: Fill fails before setPlugins runs.
}

type Left struct{ Right *Right }
type Right struct {
	Left *Left // want "TESTO022: recursive initialization.*Left"
}
type Mutual struct {
	*framework.T
	*Left
}

// Fill traverses value structs even though setPlugins would reject them.
type ValueCycle struct {
	Next *ValueCycle // want "TESTO022: recursive initialization.*Next"
}
type ValueRecursive struct {
	*framework.T
	Value ValueCycle
}

type Leaf struct{ Number *int }
type Safe struct {
	*framework.T
	First  *Leaf
	Second *Leaf // repeated types on independent paths are not cycles.
	hidden int
	loop   *Loop
	List   *[]*Loop
	Map    *map[string]*Loop
	Array  *[1]*Loop
	Fn     *func() *Loop
	Any    *any
	Deep   **Leaf
}

// Neither unused custom T declarations nor unrelated application structs are roots.
type Unused struct {
	*framework.T
	Bad int
}
type Suite struct{ Bad int }

func (Suite) TestOther(*Invalid)  {}
func RunTest(any, func(*Invalid)) {}

type Generic[T any] struct {
	*framework.T
	Field T // want "TESTO021: exported field .*Field.* must be a pointer"
}

func roots(t *testing.T, parent *framework.T) {
	framework.RunSuite(t, BadSuite{})
	framework.RunSuite(t, &BadSuite{}) // no duplicate field diagnostic.
	framework.RunSubSuite(parent, &BadSuite{})
	framework.RunTest(t, func(*NestedT) {})
	framework.RunTest(t, func(*Recursive) {})
	framework.RunTest(t, func(*Mutual) {})
	framework.RunTest(t, func(*ValueRecursive) {})
	framework.RunTest(t, func(*Safe) {})
	framework.RunTest(t, func(Safe) {}) // value context has the same traversed fields.
	framework.RunTest(t, func(*Generic[int]) {})
	framework.RunTest(t, func(*Generic[*int]) {})
	framework.RunTest(t, func(*framework.T) {})
	RunTest(t, func(*Invalid) {})
	_ = framework.Test(func(*Unused) {}) // adapter never executed.
}

type AdapterT struct {
	*framework.T
	Bad bool // want "TESTO021: exported field .*Bad.* must be a pointer"
}
type ChildT struct {
	*framework.T
	Bad float64 // want "TESTO021: exported field .*Bad.* must be a pointer"
}
type AliasCallT struct {
	*framework.T
	Bad complex64 // want "TESTO021: exported field .*Bad.* must be a pointer"
}

func adapters(t *testing.T, child *ChildT) {
	t.Run("adapter", framework.Test(func(*AdapterT) {}))
	framework.Test(func(*AdapterT) {})(t)
	framework.Run(child, "child", func(*ChildT) {})
	run := framework.RunTest[*AliasCallT]
	run(t, func(*AliasCallT) {})
}
