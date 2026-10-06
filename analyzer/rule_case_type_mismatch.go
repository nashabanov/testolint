package analyzer

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

func checkCaseTypeMismatch(p *analysis.Pass, suite *Suite) {
	for _, test := range suite.Tests {
		if !test.validParams(suite) {
			continue
		}
		for _, param := range test.Params {
			if !ast.IsExported(param.Name) {
				continue // Not settable by Testo; reported by TESTO009.
			}

			provider := suite.CasesByName[param.Name]
			if provider == nil {
				continue // TESTO003
			}

			slice, _ := casesSignature(provider.Func.Type())
			if slice == nil {
				continue // TESTO005
			}

			if types.AssignableTo(slice.Elem(), param.Type) {
				continue
			}

			p.Reportf(
				provider.Pos(suite),
				"TESTO004: Cases%s provides %s, but parameter %q expects %s",
				param.Name,
				slice.Elem(),
				param.Name,
				param.Type,
			)
		}
	}
}
