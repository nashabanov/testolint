package analyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

func checkUnexportedParamField(p *analysis.Pass, suite *Suite) {
	for _, test := range suite.Tests {
		for _, param := range test.Params {
			if ast.IsExported(param.Name) {
				continue
			}

			p.Reportf(
				param.Pos,
				"TESTO009: parameter field %q must be exported",
				param.Name,
			)
		}
	}
}
