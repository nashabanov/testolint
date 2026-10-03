package analyzer

import "golang.org/x/tools/go/analysis"

func checkOrphanCasesProvider(p *analysis.Pass, suite *Suite) {
	used := make(map[string]struct{})

	for _, test := range suite.Tests {
		for _, param := range test.Params {
			used[param.Name] = struct{}{}
		}
	}

	for _, provider := range suite.Cases {
		if !isValidPrefixedName(provider.Decl.Name.Name, "Cases") {
			continue
		}

		name := casesName(provider)

		if _, ok := used[name]; ok {
			continue
		}

		p.Reportf(
			provider.Decl.Name.Pos(),
			"TESTO006: Cases%s is not used by any test parameter",
			name,
		)
	}
}
