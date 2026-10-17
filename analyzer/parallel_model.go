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

func buildParallelModel(p *analysis.Pass, suite *Suite) *parallelModel {
	if !standardTestoT(suite.TType) || !validParallelProviders(suite) {
		return nil
	}
	hooks, ok := parseParallelHooks(p, suite)
	if !ok {
		return nil
	}
	model := &parallelModel{}
	tests := append([]*Test(nil), suite.Tests...)
	// Testo Collect appends regular methods first, then parameterized methods.
	sort.SliceStable(tests, func(i, j int) bool {
		a, b := tests[i].Func().Type().(*types.Signature).Params().Len(), tests[j].Func().Type().(*types.Signature).Params().Len()
		if a != b {
			return a < b
		}
		return tests[i].Func().Name() < tests[j].Func().Name()
	})
	sizes := make(map[*Method]int)
	for _, provider := range suite.Cases {
		sizes[provider] = constantProviderSize(p, provider)
	}
	for _, test := range tests {
		if !validParallelSignature(test.Func().Type().(*types.Signature), suite.TType, test.Func().Name()) || !isValidPrefixedName(test.Func().Name(), "Test") || !validParallelParameters(suite, test) {
			return nil
		}
		body := parallelBodyFacts(p, suite, test.Method)
		path := newParallelPath(hooks["BeforeEach"], body, hooks["AfterEach"])
		path.executions = testExecutionCount(suite, test, sizes)
		model.paths = append(model.paths, path)
	}
	model.sharedReads = findSharedFieldReads(model.paths)
	return model
}

func parseParallelHooks(p *analysis.Pass, suite *Suite) (map[string]parallelBody, bool) {
	hooks := make(map[string]parallelBody)
	for _, method := range suite.Hooks {
		if !validParallelSignature(method.Func.Type().(*types.Signature), suite.TType, method.Func.Name()) {
			return nil, false
		}
		if method.Func.Name() == "BeforeAll" || method.Func.Name() == "AfterAll" {
			continue
		}
		hooks[method.Func.Name()] = parallelBodyFacts(p, suite, method)
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

func standardTestoT(t types.Type) bool {
	ptr, ok := types.Unalias(t).(*types.Pointer)
	if !ok {
		return false
	}
	named, ok := types.Unalias(ptr.Elem()).(*types.Named)
	return ok && named.Obj().Pkg() != nil &&
		named.Obj().Pkg().Path() == testoImportPath && named.Obj().Name() == "T"
}
