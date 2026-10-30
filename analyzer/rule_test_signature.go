package analyzer

import (
	"go/types"

	"golang.org/x/tools/go/analysis"
)

func checkTestSignature(p *analysis.Pass, suite *suite) {
	for _, test := range suite.tests {
		sig, ok := test.fn().Type().(*types.Signature)
		if !ok {
			continue
		}

		params := sig.Params()

		if params.Len() != 1 && params.Len() != 2 {
			reportf(p, "TESTO001", test.method.pos(suite),
				"invalid test signature: expected func(T) or func(T, struct{...})",
			)
			continue
		}

		if !types.Identical(params.At(0).Type(), suite.tType) {
			reportf(p, "TESTO001", test.method.pos(suite),
				"invalid test signature: first parameter must match suite T",
			)
			continue
		}

		if params.Len() == 2 && !isStructType(params.At(1).Type()) {
			reportf(p, "TESTO001", test.method.pos(suite),
				"invalid test signature: second parameter must be a struct",
			)
		}

		if sig.Results().Len() != 0 {
			reportf(p, "TESTO001", test.method.pos(suite),
				"test must not return values",
			)
		}
	}
}
