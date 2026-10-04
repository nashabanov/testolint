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

			sig, ok := provider.Func.Type().(*types.Signature)
			if !ok || sig.Params().Len() != 0 {
				continue // TESTO005
			}

			results := sig.Results()
			if results.Len() != 1 {
				continue // TESTO005
			}

			slice, ok := results.At(0).Type().Underlying().(*types.Slice)
			if !ok {
				continue // TESTO005
			}

			if types.AssignableTo(slice.Elem(), param.Type) {
				continue
			}

			p.Reportf(
				provider.Decl.Name.Pos(),
				"TESTO004: Cases%s provides %s, but parameter %q expects %s",
				param.Name,
				slice.Elem(),
				param.Name,
				param.Type,
			)
		}
	}
}
