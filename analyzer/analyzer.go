package analyzer

import "golang.org/x/tools/go/analysis"

var Analyzer = &analysis.Analyzer{
	Name: "testolint",
	Doc:  "check testo usage",
	Run:  run,
}

func run(p *analysis.Pass) (any, error) {
	tests := findTestMethods(p)

	for _, test := range tests {
		p.Reportf(
			test.Name.Pos(),
			"testolint: found new test method %s",
			test.Name.Name,
		)
	}

	return nil, nil
}
