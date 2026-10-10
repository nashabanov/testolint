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
	return runWithRules(p, suiteRules(), nil)
}

// runWithRules shares discovery and diagnostic deduplication across both rule
// collections. All state belongs to this invocation, including the rule lists.
func runWithRules(p *analysis.Pass, suiteRules []Rule, runRules []RunRule) (any, error) {
	suites := discoverSuites(p)
	runs := discoverRuns(p)
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
		for _, rule := range suiteRules {
			rule(&pass, suite)
		}
	}

	for _, call := range runs {
		for _, rule := range runRules {
			rule(&pass, call)
		}
	}

	return nil, nil
}
