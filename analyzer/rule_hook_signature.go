package analyzer

import (
	"go/types"

	"golang.org/x/tools/go/analysis"
)

func checkHookSignature(p *analysis.Pass, suite *Suite) {
	for _, hook := range suite.Hooks {
		sig, ok := hook.Func.Type().(*types.Signature)
		if !ok {
			continue
		}

		params := sig.Params()

		if params.Len() != 1 {
			p.Reportf(
				hook.Decl.Name.Pos(),
				"TESTO002: invalid hook signature: expected func(T)",
			)
			continue
		}

		if suite.TType != nil &&
			!types.Identical(params.At(0).Type(), suite.TType) {
			p.Reportf(
				hook.Decl.Name.Pos(),
				"TESTO002: invalid hook signature: parameter must match suite T",
			)
		}

		if sig.Results().Len() != 0 {
			p.Reportf(
				hook.Decl.Name.Pos(),
				"TESTO002: hook must not return values",
			)
		}
	}
}
