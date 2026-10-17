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
	// BeforeEach supplies the write independently of an unsupported sibling
	// body. Receiver access in unknown code could restore/replace that field;
	// never interpret those effects as either an observed write or a safe body.
	for _, path := range paths {
		if path.executions == noExecutions {
			continue
		}
		if !path.beforeEach.understood || !path.afterEach.understood ||
			(!path.test.understood && path.test.unknownReceiverAccess) {
			return nil
		}
	}
	replaced := fieldsWrittenOutsideBeforeEach(paths)
	var reads []sharedFieldRead
	for i, path := range paths {
		if path.understood && path.parallel && path.executions > 0 && hasLaterExecution(paths, i) {
			reads = append(reads, sharedReadsInPath(path, replaced)...)
		}
	}
	return reads
}

// A suspended consumer allows every later scheduled execution to enter
// BeforeEach, regardless of whether that sibling itself calls Parallel.
func hasLaterExecution(paths []parallelPath, index int) bool {
	if paths[index].executions >= 2 {
		return true
	}
	for _, later := range paths[index+1:] {
		if later.executions > 0 {
			return true
		}
	}
	return false
}

// Additional writes require ordering analysis beyond the BeforeEach model.
func fieldsWrittenOutsideBeforeEach(paths []parallelPath) map[*types.Var]bool {
	fields := make(map[*types.Var]bool)
	for _, path := range paths {
		if path.executions == noExecutions {
			continue
		}
		// These facts remain relevant even when an earlier lifecycle stage is
		// unknown and consequently absent from the reachable prefix.
		for _, body := range []parallelBody{path.test, path.afterEach} {
			for _, event := range body.events {
				if event.kind == parallelFieldWrite {
					fields[event.field] = true
				}
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
