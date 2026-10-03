package analyzer

import (
	"strings"
	"unicode"
	"unicode/utf8"

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
			"TESTO007: malformed test name %q: suffix must be empty or start with a non-lowercase rune",
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
			"TESTO008: malformed cases provider name %q: suffix must be empty or start with a non-lowercase rune",
			name,
		)
	}
}

func isValidPrefixedName(name, prefix string) bool {
	if !strings.HasPrefix(name, prefix) {
		return false
	}

	if len(name) == len(prefix) {
		return true
	}

	r, _ := utf8.DecodeRuneInString(name[len(prefix):])
	return !unicode.IsLower(r)
}
