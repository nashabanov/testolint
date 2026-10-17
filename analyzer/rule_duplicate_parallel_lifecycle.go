package analyzer

import "golang.org/x/tools/go/analysis"

func checkDuplicateParallelLifecycle(p *analysis.Pass, suite *Suite) {
	if suite.parallel == nil {
		return
	}
	for _, path := range suite.parallel.paths {
		// Only this test must have an established execution. Sibling case
		// counts do not affect the identity of its lifecycle context.
		if path.executions <= 0 {
			continue
		}
		seen := false
		for _, event := range path.events {
			if event.kind != parallelCall {
				continue
			}
			if seen {
				p.Reportf(event.pos,
					"TESTO017: Parallel called again on the same test context across lifecycle/test bodies")
			}
			seen = true
		}
	}
}
