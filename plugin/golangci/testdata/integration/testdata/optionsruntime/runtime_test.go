package optionsruntime

import (
	"testing"

	"github.com/ozontech/testo"
	"github.com/ozontech/testo/testoplugin"
)

type marker struct{}
type optionPlugin struct{ configured bool }

func (p *optionPlugin) Plugin(_ testoplugin.Plugin, options ...testoplugin.Option) testoplugin.Spec {
	for _, option := range options {
		if _, ok := option.Value.(marker); ok {
			p.configured = true
		}
	}
	return testoplugin.Spec{}
}

type context struct {
	*testo.T
	Plugin *optionPlugin
}
type setupSuite struct{ testo.Suite[*context] }

func (setupSuite) TestConfigure(t *context) {
	// This is deliberate process-wide configuration for subsequent initializations.
	testo.Options(testoplugin.Option{Value: marker{}}.Propagated())
	testo.Run(t, "existing snapshot", func(child *context) {
		if child.Plugin.configured {
			child.Fatal("existing context unexpectedly received later global options")
		}
	})
	testo.RunSubSuite(t, configuredSuite{})
}

type configuredSuite struct{ testo.Suite[*context] }

func (configuredSuite) TestConfigured(t *context) {
	if !t.Plugin.configured {
		t.Fatal("new subsuite did not receive global options")
	}
}
func TestIntentionalOptions(t *testing.T) {
	testo.RunSuite(t, setupSuite{})
	testo.RunTest(t, func(ctx *context) {
		if !ctx.Plugin.configured {
			ctx.Fatal("later run did not receive global options")
		}
	})
}
