package constructionruntime

import (
	"testing"

	"github.com/ozontech/testo"
	"github.com/ozontech/testo/testoplugin"
)

type invalidContext struct {
	*testo.T
	Count int
}

func TestInvalidField(t *testing.T) { testo.RunTest(t, func(*invalidContext) {}) }

type invalidPlugin struct{ Value int }

func (*invalidPlugin) Plugin(testoplugin.Plugin, ...testoplugin.Option) testoplugin.Spec {
	return testoplugin.Spec{}
}

type pluginContext struct {
	*testo.T
	Plugin *invalidPlugin
}

func TestInvalidPluginField(t *testing.T) { testo.RunTest(t, func(*pluginContext) {}) }

type recursive struct{ Next *recursive }
type recursiveContext struct {
	*testo.T
	Part *recursive
	Bad  int // Fill fails first, before setPlugins validates this field.
}

func TestRecursive(t *testing.T) { testo.RunTest(t, func(*recursiveContext) {}) }

type left struct{ Right *right }
type right struct{ Left *left }
type mutualContext struct {
	*testo.T
	Part *left
}

func TestMutual(t *testing.T) { testo.RunTest(t, func(*mutualContext) {}) }

type valueRecursive struct{ Next *valueRecursive }
type valueContext struct {
	*testo.T
	Part valueRecursive
}

func TestValueRecursive(t *testing.T) { testo.RunTest(t, func(*valueContext) {}) }
