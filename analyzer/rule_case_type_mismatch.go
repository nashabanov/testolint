package analyzer

import (
	"go/types"

	"golang.org/x/tools/go/analysis"
)

func checkCaseTypeMismatch(p *analysis.Pass, suite *Suite) {
	for _, test := range suite.Tests {
		for _, param := range test.Params {
			provider := suite.CasesByName[param.Name]
			if provider == nil {
				continue // TESTO003
			}

			sig, ok := provider.Func.Type().(*types.Signature)
			if !ok {
				continue
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
				"TESTO004: Cases%s returns []%s, expected []%s",
				param.Name,
				slice.Elem(),
				param.Type,
			)
		}
	}
}
