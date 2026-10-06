package analyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

func checkUnexportedParamField(p *analysis.Pass, suite *Suite) {
	for _, test := range suite.Tests {
		if !test.validParams(suite) {
			continue
		}
		for _, param := range test.Params {
			if ast.IsExported(param.Name) {
				continue
			}

			pos := param.Pos
			if !pos.IsValid() {
				pos = test.Method.Pos(suite)
			}
			p.Reportf(
				pos,
				"TESTO009: parameter field %q must be exported",
				param.Name,
			)
		}
	}
}
