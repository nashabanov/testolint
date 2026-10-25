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
	return runWithRules(p, ruleSet{
		suites:        suiteRules(),
		parallel:      parallelRules(),
		afterParallel: afterParallelRules(),
		runs:          runRules(),
		packages:      packageRules(),
	})
}

// runWithRules shares call facts and diagnostic deduplication across all rule
// groups. All state belongs to this invocation.
func runWithRules(p *analysis.Pass, rules ruleSet) (any, error) {
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
		for _, rule := range rules.suites {
			rule(&pass, suite)
		}
		if len(rules.parallel) > 0 {
			if model := buildParallelModel(&pass, suite); model != nil {
				for _, rule := range rules.parallel {
					rule(&pass, model)
				}
			}
		}
		for _, rule := range rules.afterParallel {
			rule(&pass, suite)
		}
	}

	for _, call := range runs {
		for _, rule := range rules.runs {
			rule(&pass, call)
		}
	}

	for _, rule := range rules.packages {
		rule(&pass, facts)
	}
	return nil, nil
}
