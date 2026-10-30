package analyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

func checkMissingCaseProvider(p *analysis.Pass, suite *suite) {
	for _, test := range suite.tests {
		if !test.validParams(suite) {
			continue
		}
		for _, param := range test.params {
			if !ast.IsExported(param.name) {
				continue // Not settable by Testo; reported by TESTO009.
			}

			if suite.casesByName[param.name] != nil {
				continue
			}

			pos := param.pos
			if !pos.IsValid() {
				pos = test.method.pos(suite)
			}
			reportf(p, "TESTO003", pos,
				"parameter %q requires Cases%s",
				param.name,
				param.name,
			)
		}
	}
}
