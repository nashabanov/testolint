package analyzer

import (
	"golang.org/x/tools/go/analysis"
)

func checkCasesSignature(p *analysis.Pass, suite *suite) {
	for _, cases := range suite.cases {
		if _, problem := casesSignature(cases.fn.Type()); problem != "" {
			reportf(p, "TESTO005", cases.pos(suite),
				"%s", problem)
		}
	}
}
