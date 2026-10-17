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

// Each path combines BeforeEach, one test method and AfterEach on the same T.
// Cleanup reads retain their deferred nature rather than implying inline use.
type parallelPath struct {
	events     []parallelEvent
	executions int // Capped at two; zero is unknown, noExecutions is proven empty.
	parallel   bool
	understood bool // Complete local timeline; unknown bodies contribute no scheduling evidence.
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
		body, ok := parseParallelBody(p, suite, test.Method)
		if !ok {
			body = nil
		}
		path := newParallelPath(hooks["BeforeEach"], body, hooks["AfterEach"])
		path.understood = ok
		path.executions = testExecutionCount(suite, test, sizes)
		model.paths = append(model.paths, path)
	}
	model.sharedReads = findSharedFieldReads(model.paths)
	return model
}

func parseParallelHooks(p *analysis.Pass, suite *Suite) (map[string][]parallelEvent, bool) {
	hooks := make(map[string][]parallelEvent)
	for _, method := range suite.Hooks {
		if !validParallelSignature(method.Func.Type().(*types.Signature), suite.TType, method.Func.Name()) {
			return nil, false
		}
		if method.Func.Name() == "BeforeAll" || method.Func.Name() == "AfterAll" {
			continue
		}
		body, ok := parseParallelBody(p, suite, method)
		if !ok {
			return nil, false
		}
		hooks[method.Func.Name()] = body
	}
	return hooks, true
}

func newParallelPath(before, test, after []parallelEvent) parallelPath {
	events := make([]parallelEvent, 0, len(before)+len(test)+len(after))
	for _, event := range before {
		event.beforeEach = true
		events = append(events, event)
	}
	events = append(events, test...)
	events = append(events, after...)
	return parallelPath{events: events, parallel: hasParallel(before) || hasParallel(test)}
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
