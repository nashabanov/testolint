package analyzer

import "golang.org/x/tools/go/analysis"

func checkEmptySuite(p *analysis.Pass, suite *Suite) {
	for _, test := range suite.Tests {
		if isValidPrefixedName(test.Func().Name(), "Test") {
			return
		}
	}
	obj := suite.Type.Obj()
	p.Reportf(obj.Pos(), "TESTO011: suite %q contains no tests", obj.Name())
}
