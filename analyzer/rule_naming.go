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

		reportf(p, "TESTO007", test.method.pos(suite),
			"malformed test name %q: suffix must be empty or start with a non-lowercase rune",
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

		reportf(p, "TESTO008", provider.pos(suite),
			"malformed cases provider name %q: suffix must be empty or start with a non-lowercase rune",
			name,
		)
	}
}
