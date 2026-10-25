package analyzer

import "golang.org/x/tools/go/analysis"

func checkSharedTestContext(p *analysis.Pass, model *parallelModel) {
	for _, read := range model.sharedReads {
		if read.storesContext {
			p.Reportf(read.event.pos,
				"TESTO016: shared field %q can refer to another test's context",
				read.event.field.Name())
		}
	}
}
