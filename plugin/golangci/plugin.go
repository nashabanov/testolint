// Package golangci adapts testolint to the golangci-lint module plugin system.
package golangci

import (
	"github.com/golangci/plugin-module-register/register"
	"github.com/nashabanov/testolint/analyzer"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("testolint", newPlugin)
}

type plugin struct{ analyzer *analysis.Analyzer }

func newPlugin(settings any) (register.LinterPlugin, error) {
	config, err := register.DecodeSettings[analyzer.Config](settings)
	if err != nil {
		return nil, err
	}
	a, err := analyzer.New(config)
	if err != nil {
		return nil, err
	}
	return &plugin{analyzer: a}, nil
}

func (p *plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{p.analyzer}, nil
}

func (*plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
