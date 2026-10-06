package analyzer

import (
	"go/types"

	"golang.org/x/tools/go/analysis"
)

func checkCasesSignature(p *analysis.Pass, suite *Suite) {
	for _, cases := range suite.Cases {
		sig, ok := cases.Func.Type().(*types.Signature)
		if !ok {
			continue
		}

		if sig.Params().Len() != 0 {
			p.Reportf(
				cases.Pos(suite),
				"TESTO005: cases provider must not accept parameters",
			)
			continue
		}

		results := sig.Results()

		if results.Len() != 1 {
			p.Reportf(
				cases.Pos(suite),
				"TESTO005: cases provider must return exactly one slice",
			)
			continue
		}

		if !isSliceType(results.At(0).Type()) {
			p.Reportf(
				cases.Pos(suite),
				"TESTO005: cases provider must return a slice",
			)
		}
	}
}

func isSliceType(t types.Type) bool {
	_, ok := t.Underlying().(*types.Slice)
	return ok
}
