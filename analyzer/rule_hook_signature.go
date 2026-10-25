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
			p.Reportf(
				hook.pos(suite),
				"TESTO002: invalid hook signature: expected func(T)",
			)
			continue
		}

		if suite.tType != nil &&
			!types.Identical(params.At(0).Type(), suite.tType) {
			p.Reportf(
				hook.pos(suite),
				"TESTO002: invalid hook signature: parameter must match suite T",
			)
		}

		if sig.Results().Len() != 0 {
			p.Reportf(
				hook.pos(suite),
				"TESTO002: hook must not return values",
			)
		}
	}
}
