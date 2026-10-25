package analyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

func checkUnexportedParamField(p *analysis.Pass, suite *suite) {
	for _, test := range suite.tests {
		if !test.validParams(suite) {
			continue
		}
		for _, param := range test.params {
			if ast.IsExported(param.name) {
				continue
			}

			pos := param.pos
			if !pos.IsValid() {
				pos = test.method.pos(suite)
			}
			p.Reportf(
				pos,
				"TESTO009: parameter field %q must be exported",
				param.name,
			)
		}
	}
}
