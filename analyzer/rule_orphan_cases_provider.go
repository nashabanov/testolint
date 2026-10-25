package analyzer

import (
	"go/types"

	"golang.org/x/tools/go/analysis"
)

func checkOrphanCasesProvider(p *analysis.Pass, suite *suite) {
	// A malformed test may hide provider references; avoid secondary orphan reports.
	for _, test := range suite.tests {
		sig := test.fn().Type().(*types.Signature)
		if sig.Params().Len() > 1 && !test.validParams(suite) {
			return
		}
	}
	used := make(map[string]struct{})

	for _, test := range suite.tests {
		for _, param := range test.params {
			used[param.name] = struct{}{}
		}
	}

	for _, provider := range suite.cases {
		if !isValidPrefixedName(provider.fn.Name(), "Cases") {
			continue
		}

		name := casesName(provider)

		if _, ok := used[name]; ok {
			continue
		}

		p.Reportf(
			provider.pos(suite),
			"TESTO006: Cases%s is not used by any test parameter",
			name,
		)
	}
}
