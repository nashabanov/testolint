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
		for _, test := range suite.Tests {
			p.Reportf(
				test.Decl.Name.Pos(),
				"suite %s: found test %s",
				suite.Type.Obj().Name(),
				test.Decl.Name.Name,
			)
		}
	}

	return nil, nil
}
