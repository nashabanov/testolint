package analyzer

import "go/types"

// A read of a BeforeEach field with an established overwrite window. These
// facts let each rule decide applicability and diagnostic precedence explicitly.
type sharedFieldRead struct {
	event             parallelEvent
	storesContext     bool
	alsoReadInCleanup bool
}

func findSharedFieldReads(paths []parallelPath) []sharedFieldRead {
	if !knownParallelCounts(paths) || parallelExecutionCount(paths) < 2 {
		return nil
	}
	replaced := fieldsWrittenOutsideBeforeEach(paths)
	var reads []sharedFieldRead
	for _, path := range paths {
		if path.parallel {
			reads = append(reads, sharedReadsInPath(path, replaced)...)
		}
	}
	return reads
}

func parallelExecutionCount(paths []parallelPath) int {
	count := 0
	for _, path := range paths {
		if path.parallel {
			count += path.executions
		}
	}
	return count
}

// Additional writes require ordering analysis beyond the BeforeEach model.
func fieldsWrittenOutsideBeforeEach(paths []parallelPath) map[*types.Var]bool {
	fields := make(map[*types.Var]bool)
	for _, path := range paths {
		for _, event := range path.events {
			if event.kind == parallelFieldWrite && !event.beforeEach {
				fields[event.field] = true
			}
		}
	}
	return fields
}

func cleanupFields(events []parallelEvent) map[*types.Var]bool {
	fields := make(map[*types.Var]bool)
	for _, event := range events {
		if event.inCleanup {
			fields[event.field] = true
		}
	}
	return fields
}

func sharedReadsInPath(path parallelPath, replaced map[*types.Var]bool) []sharedFieldRead {
	cleanup := cleanupFields(path.events)
	writes := make(map[*types.Var]parallelEvent)
	var reads []sharedFieldRead
	suspended := false
	for _, event := range path.events {
		switch event.kind {
		case parallelCall:
			suspended = true
		case parallelFieldWrite:
			delete(writes, event.field)
			if event.beforeEach && !suspended {
				writes[event.field] = event
			}
		case parallelFieldRead:
			write, ok := writes[event.field]
			if !ok || replaced[event.field] || (!suspended && !event.inCleanup) {
				continue
			}
			reads = append(reads, sharedFieldRead{
				event:             event,
				storesContext:     write.storesContext,
				alsoReadInCleanup: cleanup[event.field],
			})
		}
	}
	return reads
}
