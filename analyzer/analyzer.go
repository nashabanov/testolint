package analyzer

import (
	"go/token"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "testolint",
	Doc:  "check testo usage",
	Run:  run,
}

func run(p *analysis.Pass) (any, error) {
	suites := discoverSuites(p)
	pass := *p
	type key struct {
		pos     token.Pos
		message string
	}
	seen := make(map[key]bool)
	pass.Report = func(d analysis.Diagnostic) {
		k := key{d.Pos, d.Message}
		if !seen[k] {
			seen[k] = true
			p.Report(d)
		}
	}

	for _, suite := range suites {
		for _, rule := range rules {
			rule(&pass, suite)
		}
	}

	return nil, nil
}
