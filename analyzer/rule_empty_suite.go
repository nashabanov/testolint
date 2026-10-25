package analyzer

import "golang.org/x/tools/go/analysis"

func checkEmptySuite(p *analysis.Pass, suite *suite) {
	for _, test := range suite.tests {
		if isValidPrefixedName(test.fn().Name(), "Test") {
			return
		}
	}
	obj := suite.typ.Obj()
	p.Reportf(obj.Pos(), "TESTO011: suite %q contains no tests", obj.Name())
}
