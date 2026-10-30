package analyzer

import "golang.org/x/tools/go/analysis"

func checkSharedTestContext(p *analysis.Pass, model *parallelModel) {
	for _, read := range model.sharedReads {
		if read.storesContext {
			reportf(p, "TESTO016", read.event.pos,
				"shared field %q can refer to another test's context",
				read.event.field.Name())
		}
	}
}
