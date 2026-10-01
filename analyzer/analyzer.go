package analyzer

import "golang.org/x/tools/go/analysis"

var Analyzer = &analysis.Analyzer{
	Name: "testolint",
	Doc:  "check testo usage",
	Run:  run,
}

func run(p *analysis.Pass) (any, error) {
	for _, file := range p.Files {
		p.Reportf(file.Pos(), "testolint: file analyzed")
	}
	return nil, nil
}
