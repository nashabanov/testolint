package analyzer

import "golang.org/x/tools/go/analysis"

var Analyzer = &analysis.Analyzer{
	Name: "testolint",
	Doc:  "check testo usage",
	Run:  run,
}

func run(p *analysis.Pass) (any, error) {
	suites := discoverSuites(p)

	for _, suite := range suites {
		for _, rule := range rules {
			rule(p, suite)
		}
	}

	return nil, nil
}
