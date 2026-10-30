package analyzer

import (
	"go/types"

	"golang.org/x/tools/go/analysis"
)

func checkHookSignature(p *analysis.Pass, suite *suite) {
	for _, hook := range suite.hooks {
		sig, ok := hook.fn.Type().(*types.Signature)
		if !ok {
			continue
		}

		params := sig.Params()

		if params.Len() != 1 {
			reportf(p, "TESTO002", hook.pos(suite),
				"invalid hook signature: expected func(T)",
			)
			continue
		}

		if suite.tType != nil &&
			!types.Identical(params.At(0).Type(), suite.tType) {
			reportf(p, "TESTO002", hook.pos(suite),
				"invalid hook signature: parameter must match suite T",
			)
		}

		if sig.Results().Len() != 0 {
			reportf(p, "TESTO002", hook.pos(suite),
				"hook must not return values",
			)
		}
	}
}
