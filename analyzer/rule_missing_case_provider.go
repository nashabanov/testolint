package analyzer

import "golang.org/x/tools/go/analysis"

func checkMissingCaseProvider(p *analysis.Pass, suite *Suite) {
	for _, test := range suite.Tests {
		for _, param := range test.Params {
			if _, ok := suite.CasesByName[param.Name]; ok {
				continue
			}

			p.Reportf(
				param.Pos,
				"TESTO003: parameter %q requires Cases%s",
				param.Name,
				param.Name,
			)
		}
	}
}
