package analyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

func checkMissingCaseProvider(p *analysis.Pass, suite *Suite) {
	for _, test := range suite.Tests {
		for _, param := range test.Params {
			if !ast.IsExported(param.Name) {
				continue // Not settable by Testo; reported by TESTO009.
			}

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
