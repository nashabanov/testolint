package required_header

import (
	"github.com/ozontech/testo"
	a "github.com/ozontech/testo-allure"
	"required_header/lib"
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

func (*Conditional) TestProvenMissingTitle(t T) {} // want "TESTO028.*Title"

type ConditionCall struct{ testo.Suite[T] }

func unknownMetadata(t T) bool { t.Title("title"); t.Tags("tag"); return true }
func (*ConditionCall) BeforeEach(t T) {
	if unknownMetadata(t) {
	}
}
func (*ConditionCall) TestUnknown(t T) {}

type Imported struct{ lib.Base[T] }

func (*Imported) TestUnknown(t T) {}

func (*Suite) TestAlias(t T) {
	alias := t
	alias.Title("alias title")
	alias.Tags("smoke")
}
func (*Suite) TestMethodValue(t T) {
	title := t.Title
	title("method value title")
	t.Tags("smoke")
}
func (*Suite) TestConditionalAlias(t T) {
	if alias := t; condition {
		alias.Title("title")
		alias.Tags("smoke")
	}
}
func (*Suite) TestStoredClosure(t T) { // want "TESTO028.*Title or Titlef, Tags"
	_ = func() { alias := t; alias.Title("unused"); alias.Tags("unused") }
}

// Known limitation: passing the context to Step suppresses missing metadata.
func (*Suite) TestMissingWithStep(t T) {
	a.Step(t, "work", func(t T) { t.Log("nothing") })
}

type HeaderContext = a.CommonT
type GenericBase[C HeaderContext] struct{ testo.Suite[C] }

func (*GenericBase[C]) BeforeEach(t C) { t.Tags("smoke") }

type GenericSuite struct{ GenericBase[T] }

func (*GenericSuite) TestGenericHook(t T)  { t.Title("generic hook title") }
func (*GenericSuite) TestMissingTitle(t T) {} // want "TESTO028.*Title or Titlef"

type GenericTests[C HeaderContext] struct{ testo.Suite[C] }

func (GenericTests[C]) TestGeneric(t C)     { t.Title("generic test"); t.Tags("smoke") }
func (GenericTests[C]) TestMissingTags(t C) { t.Title("title") } // want "TESTO028.*missing required test metadata: Tags"
type PromotedTests struct{ GenericTests[T] }
