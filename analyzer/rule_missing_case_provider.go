package analyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

func checkMissingCaseProvider(p *analysis.Pass, suite *Suite) {
	for _, test := range suite.Tests {
		if !test.validParams(suite) {
			continue
		}
		for _, param := range test.Params {
			if !ast.IsExported(param.Name) {
				continue // Not settable by Testo; reported by TESTO009.
			}

			if suite.MethodSet.Lookup(nil, "Cases"+param.Name) != nil {
				continue
			}

			pos := param.Pos
			if !pos.IsValid() {
				pos = test.Decl().Name.Pos()
			}
			p.Reportf(
				pos,
				"TESTO003: parameter %q requires Cases%s",
				param.Name,
				param.Name,
			)
		}
	}
}
