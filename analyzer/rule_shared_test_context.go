package analyzer

import "golang.org/x/tools/go/analysis"

func checkSharedTestContext(p *analysis.Pass, suite *Suite) {
	if suite.parallel == nil {
		return
	}
	for _, read := range suite.parallel.sharedReads {
		if read.storesContext {
			p.Reportf(read.event.pos,
				"TESTO016: shared field %q can refer to another test's context",
				read.event.field.Name())
		}
	}
}
