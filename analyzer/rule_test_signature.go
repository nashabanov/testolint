package analyzer

import (
	"go/types"

	"golang.org/x/tools/go/analysis"
)

func checkTestSignature(p *analysis.Pass, suite *Suite) {
	for _, test := range suite.Tests {
		sig, ok := test.Func().Type().(*types.Signature)
		if !ok {
			continue
		}

		params := sig.Params()

		if params.Len() != 1 && params.Len() != 2 {
			p.Reportf(
				test.Method.Pos(suite),
				"TESTO001: invalid test signature: expected func(T) or func(T, struct{...})",
			)
			continue
		}

		if !types.Identical(params.At(0).Type(), suite.TType) {
			p.Reportf(
				test.Method.Pos(suite),
				"TESTO001: invalid test signature: first parameter must match suite T",
			)
			continue
		}

		if params.Len() == 2 && !isStructType(params.At(1).Type()) {
			p.Reportf(
				test.Method.Pos(suite),
				"TESTO001: invalid test signature: second parameter must be a struct",
			)
		}

		if sig.Results().Len() != 0 {
			p.Reportf(
				test.Method.Pos(suite),
				"TESTO001: test must not return values",
			)
		}
	}
}

func isStructType(t types.Type) bool {
	_, ok := t.Underlying().(*types.Struct)
	return ok
}
