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
	if len(analyzers) != 1 || analyzers[0] == analyzer.Analyzer || analyzers[0].Name != analyzer.Analyzer.Name {
		t.Fatal("plugin must return an independent testolint analyzer")
	}
	if plugin.GetLoadMode() != register.LoadModeTypesInfo {
		t.Fatal("plugin must request type information")
	}
}

func TestIndependentPlugins(t *testing.T) {
	factory, err := register.GetPlugin("testolint")
	if err != nil {
		t.Fatal(err)
	}
	first, err := factory(map[string]any{"disable-all": true, "enable": []string{"TESTO019"}})
	if err != nil {
		t.Fatal(err)
	}
	second, err := factory(nil)
	if err != nil {
		t.Fatal(err)
	}
	a, err := first.BuildAnalyzers()
	if err != nil {
		t.Fatal(err)
	}
	b, err := second.BuildAnalyzers()
	if err != nil {
		t.Fatal(err)
	}
	if a[0] == b[0] || a[0] == analyzer.Analyzer || b[0] == analyzer.Analyzer {
		t.Fatal("analyzer instances shared")
	}
	if a[0].Flags.Lookup("enable").Value.String() != "TESTO019" || b[0].Flags.Lookup("enable").Value.String() != "" {
		t.Fatal("plugin configuration leaked")
	}
	if err := a[0].Flags.Set("disable", "TESTO006"); err != nil {
		t.Fatal(err)
	}
	if b[0].Flags.Lookup("disable").Value.String() != "" || analyzer.Analyzer.Flags.Lookup("disable") != nil {
		t.Fatal("flag configuration leaked")
	}
}
