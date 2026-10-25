package analyzer

import "golang.org/x/tools/go/analysis"

func checkSharedFixtureCleanup(p *analysis.Pass, model *parallelModel) {
	for _, read := range model.sharedReads {
		// Stored test contexts are reported by TESTO016, including in Cleanup.
		if read.storesContext || !read.event.inCleanup {
			continue
		}
		p.Reportf(read.event.pos,
			"TESTO015: Cleanup reads shared field %q after another execution can replace it; capture a local value",
			read.event.field.Name())
	}
}
