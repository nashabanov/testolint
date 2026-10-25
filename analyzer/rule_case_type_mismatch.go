package analyzer

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

func checkCaseTypeMismatch(p *analysis.Pass, suite *suite) {
	for _, test := range suite.tests {
		if !test.validParams(suite) {
			continue
		}
		for _, param := range test.params {
			if !ast.IsExported(param.name) {
				continue // Not settable by Testo; reported by TESTO009.
			}

			provider := suite.casesByName[param.name]
			if provider == nil {
				continue // TESTO003
			}

			slice, _ := casesSignature(provider.fn.Type())
			if slice == nil {
				continue // TESTO005
			}

			if types.AssignableTo(slice.Elem(), param.typ) {
				continue
			}

			p.Reportf(
				provider.pos(suite),
				"TESTO004: Cases%s provides %s, but parameter %q expects %s",
				param.name,
				slice.Elem(),
				param.name,
				param.typ,
			)
		}
	}
}
