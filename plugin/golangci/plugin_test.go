package golangci_test

import (
	"testing"

	"github.com/golangci/plugin-module-register/register"
	"github.com/nashabanov/testolint/analyzer"
	_ "github.com/nashabanov/testolint/plugin/golangci"
)

func TestRegisteredAnalyzer(t *testing.T) {
	newPlugin, err := register.GetPlugin("testolint")
	if err != nil {
		t.Fatal(err)
	}
	plugin, err := newPlugin(nil)
	if err != nil {
		t.Fatal(err)
	}
	analyzers, err := plugin.BuildAnalyzers()
	if err != nil {
		t.Fatal(err)
	}
	if len(analyzers) != 1 || analyzers[0] != analyzer.Analyzer {
		t.Fatal("plugin must return the existing analyzer.Analyzer")
	}
	if plugin.GetLoadMode() != register.LoadModeTypesInfo {
		t.Fatal("plugin must request type information")
	}
}
