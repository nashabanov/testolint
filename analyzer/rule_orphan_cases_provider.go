package analyzer

import (
	"go/types"

	"golang.org/x/tools/go/analysis"
)

func checkOrphanCasesProvider(p *analysis.Pass, suite *Suite) {
	// A malformed test may hide provider references; avoid secondary orphan reports.
	for _, test := range suite.Tests {
		sig := test.Func().Type().(*types.Signature)
		if sig.Params().Len() > 1 && !test.validParams(suite) {
			return
		}
	}
	used := make(map[string]struct{})

	for _, test := range suite.Tests {
		for _, param := range test.Params {
			used[param.Name] = struct{}{}
		}
	}

	for _, provider := range suite.Cases {
		if !isValidPrefixedName(provider.Func.Name(), "Cases") {
			continue
		}

		name := casesName(provider)

		if _, ok := used[name]; ok {
			continue
		}

		p.Reportf(
			provider.Pos(suite),
			"TESTO006: Cases%s is not used by any test parameter",
			name,
		)
	}
}
