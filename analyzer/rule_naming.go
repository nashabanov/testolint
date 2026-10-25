package analyzer

import (
	"golang.org/x/tools/go/analysis"
)

func checkMalformedTestName(p *analysis.Pass, suite *suite) {
	for _, test := range suite.tests {
		name := test.fn().Name()

		if isValidPrefixedName(name, "Test") {
			continue
		}

		p.Reportf(
			test.method.pos(suite),
			"TESTO007: malformed test name %q: suffix must be empty or start with a non-lowercase rune",
			name,
		)
	}
}

func checkMalformedCasesName(p *analysis.Pass, suite *suite) {
	for _, provider := range suite.cases {
		name := provider.fn.Name()

		if isValidPrefixedName(name, "Cases") {
			continue
		}

		p.Reportf(
			provider.pos(suite),
			"TESTO008: malformed cases provider name %q: suffix must be empty or start with a non-lowercase rune",
			name,
		)
	}
}
