package analyzer

import (
	"go/token"
	"go/types"
	"sort"

	"golang.org/x/tools/go/analysis"
)

type parallelEventKind uint8

const (
	parallelFieldRead parallelEventKind = iota
	parallelFieldWrite
	parallelCall
)

type parallelEvent struct {
	kind          parallelEventKind
	pos           token.Pos
	field         *types.Var
	storesContext bool // Field write stores the method's T parameter.
	inCleanup     bool // Field read occurs in a deferred Cleanup callback.
	beforeEach    bool // Event belongs to BeforeEach, not the test or AfterEach.
}

// Body support is independent of signature validity and execution count.
type parallelBody struct {
	events                []parallelEvent
	understood            bool // All statements in this body are supported.
	unknownReceiverAccess bool // Unknown code can access the suite directly or through T.
}

// Each path combines BeforeEach, one test method and AfterEach on the same T.
// Cleanup reads retain their deferred nature rather than implying inline use.
type parallelPath struct {
	events                      []parallelEvent
	executions                  int // Capped at two; zero is unknown, noExecutions is proven empty.
	parallel                    bool
	understood                  bool // Complete lifecycle timeline, required for shared-field consumers.
	beforeEach, test, afterEach parallelBody
}

// Shared by the four rules. Building this model never emits diagnostics.
type parallelModel struct {
	paths       []parallelPath
	sharedReads []sharedFieldRead
}

func buildParallelModel(p *analysis.Pass, suite *suite) *parallelModel {
	if !standardTestoT(suite.tType) || !validParallelProviders(suite) {
		return nil
	}
	hooks, ok := parseParallelHooks(p, suite)
	if !ok {
		return nil
	}
	model := &parallelModel{}
	tests := append([]*test(nil), suite.tests...)
	// Testo Collect appends regular methods first, then parameterized methods.
	sort.SliceStable(tests, func(i, j int) bool {
		a, b := tests[i].fn().Type().(*types.Signature).Params().Len(), tests[j].fn().Type().(*types.Signature).Params().Len()
		if a != b {
			return a < b
		}
		return tests[i].fn().Name() < tests[j].fn().Name()
	})
	sizes := make(map[*method]int)
	for _, provider := range suite.cases {
		sizes[provider] = constantProviderSize(p, provider)
	}
	for _, test := range tests {
		if !validParallelSignature(test.fn().Type().(*types.Signature), suite.tType, test.fn().Name()) || !isValidPrefixedName(test.fn().Name(), "Test") || !validParallelParameters(suite, test) {
			return nil
		}
		body := parallelBodyFacts(p, suite, test.method)
		path := newParallelPath(hooks["BeforeEach"], body, hooks["AfterEach"])
		path.executions = testExecutionCount(suite, test, sizes)
		model.paths = append(model.paths, path)
	}
	model.sharedReads = findSharedFieldReads(model.paths)
	return model
}

func parseParallelHooks(p *analysis.Pass, suite *suite) (map[string]parallelBody, bool) {
	hooks := make(map[string]parallelBody)
	for _, method := range suite.hooks {
		if !validParallelSignature(method.fn.Type().(*types.Signature), suite.tType, method.fn.Name()) {
			return nil, false
		}
		if method.fn.Name() == "BeforeAll" || method.fn.Name() == "AfterAll" {
			continue
		}
		hooks[method.fn.Name()] = parallelBodyFacts(p, suite, method)
	}
	return hooks, true
}

// Keep only complete, reachable lifecycle stages for duplicate-call checking.
// An unknown stage is a barrier, not an empty body. Later stages retain their
// independent facts but do not contribute events to this known prefix.
func newParallelPath(before, test, after parallelBody) parallelPath {
	path := parallelPath{
		beforeEach: before, test: test, afterEach: after,
		understood: before.understood && test.understood && after.understood,
	}
	if !before.understood {
		return path
	}
	for _, event := range before.events {
		event.beforeEach = true
		path.events = append(path.events, event)
	}
	if test.understood {
		path.events = append(path.events, test.events...)
		if after.understood {
			path.events = append(path.events, after.events...)
		}
	}
	path.parallel = hasParallel(before.events) || (test.understood && hasParallel(test.events))
	return path
}

func hasParallel(events []parallelEvent) bool {
	for _, event := range events {
		if event.kind == parallelCall {
			return true
		}
	}
	return false
}
