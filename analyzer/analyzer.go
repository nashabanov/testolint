package analyzer

import (
	"go/types"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "testolint",
	Doc:  "check testo usage",
	Run:  run,
}

func run(p *analysis.Pass) (any, error) {
	selected, err := resolveConfig(Config{})
	if err != nil {
		return nil, err
	}
	return runSelected(p, selected)
}

// runWithRules shares call facts and diagnostic deduplication across all rule
// groups. All state belongs to this invocation.
func runWithRules(p *analysis.Pass, rules ruleSet) (any, error) {
	if len(rules.suites)+len(rules.parallel)+len(rules.afterParallel)+len(rules.runs)+len(rules.packages) == 0 {
		return nil, nil
	}
	suites := discoverSuites(p)
	var facts *callFacts
	if len(rules.packages) > 0 || len(rules.runs) > 0 {
		facts = collectCallFacts(p)
	}
	var runs []*runCall
	if len(rules.runs) > 0 {
		runs = discoverRunsWithFacts(p, facts)
	}
	models := make(map[types.Type]*suite, len(suites))
	for _, suite := range suites {
		models[suite.typ] = suite
	}
	for _, call := range runs {
		if call.suiteType != nil {
			t := types.Unalias(call.suiteType)
			if ptr, ok := t.(*types.Pointer); ok {
				t = types.Unalias(ptr.Elem())
			}
			call.model = models[t]
		}
	}
	pass := *p
	pass.Report = deduplicatingReporter(p.Report)

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
		rule(&pass, facts, suites)
	}
	return nil, nil
}
