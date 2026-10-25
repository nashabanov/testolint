package analyzer

import "golang.org/x/tools/go/analysis"

func checkSharedBeforeEachWrite(p *analysis.Pass, model *parallelModel) {
	for _, read := range model.sharedReads {
		// TESTO016 owns context stores; TESTO015 owns cleanup hazards for this
		// field, including suppression of ordinary reads with the same cause.
		if read.storesContext || read.event.inCleanup || read.alsoReadInCleanup {
			continue
		}
		p.Reportf(read.event.pos,
			"TESTO014: shared BeforeEach field %q can be overwritten before parallel use; capture a local value",
			read.event.field.Name())
	}
}
