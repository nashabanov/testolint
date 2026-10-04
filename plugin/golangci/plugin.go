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

type plugin struct{}

func newPlugin(_ any) (register.LinterPlugin, error) {
	return &plugin{}, nil
}

func (*plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{analyzer.Analyzer}, nil
}

func (*plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
