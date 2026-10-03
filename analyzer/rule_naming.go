package analyzer

import (
	"unicode"

	"golang.org/x/tools/go/analysis"
)

func checkMalformedTestName(p *analysis.Pass, suite *Suite) {
	for _, test := range suite.Tests {
		name := test.Method.Decl.Name.Name

		if isValidPrefixedName(name, "Test") {
			continue
		}

		p.Reportf(
			test.Method.Decl.Name.Pos(),
			"TESTO007: malformed test name %q: expected TestXxx",
			name,
		)
	}
}

func checkMalformedCasesName(p *analysis.Pass, suite *Suite) {
	for _, provider := range suite.Cases {
		name := provider.Decl.Name.Name

		if isValidPrefixedName(name, "Cases") {
			continue
		}

		p.Reportf(
			provider.Decl.Name.Pos(),
			"TESTO008: malformed cases provider name %q: expected CasesXxx",
			name,
		)
	}
}

func isValidPrefixedName(name, prefix string) bool {
	if len(name) <= len(prefix) {
		return false
	}

	if name[:len(prefix)] != prefix {
		return false
	}

	r := []rune(name[len(prefix):])
	if len(r) == 0 {
		return false
	}

	return unicode.IsUpper(r[0])
}
