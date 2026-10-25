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
	return runWithRules(p, suiteRules(), runRules())
}

// runWithRules shares call facts and diagnostic deduplication across suite/run
// rules and standalone execution checks. All state belongs to this invocation.
func runWithRules(p *analysis.Pass, suiteRules []Rule, runRules []RunRule) (any, error) {
	suites := discoverSuites(p)
	facts := collectCallFacts(p)
	runs := discoverRunsWithFacts(p, facts)
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
		suite.parallel = buildParallelModel(&pass, suite)
		for _, rule := range suiteRules {
			rule(&pass, suite)
		}
	}

	for _, call := range runs {
		for _, rule := range runRules {
			rule(&pass, call)
		}
	}

	checkExecutionSafety(&pass, facts)
	checkConstruction(&pass, facts)
	return nil, nil
}
