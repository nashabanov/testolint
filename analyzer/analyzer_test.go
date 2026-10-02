package analyzer

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	tests := []string{
		"basic",
		"aliasimport",
		"pointersuite",
		"wrongpackage",
	}

	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			analysistest.Run(
				t,
				analysistest.TestData(),
				Analyzer,
				tt,
			)
		})
	}
}

func TestDiscoverSuites(t *testing.T) {
	results := analysistest.Run(
		t,
		analysistest.TestData(),
		Analyzer,
		"basic",
	)
	if len(results) != 1 {
		t.Fatalf("got %d analysis results, want 1", len(results))
	}
	if results[0].Pass == nil {
		t.Fatal("analysis pass is nil")
	}

	suites := discoverSuites(results[0].Pass)
	got := make(map[string]*Suite, len(suites))
	for _, suite := range suites {
		got[suite.Type.Obj().Name()] = suite
	}

	if len(got) != 3 {
		t.Fatalf("got suites %v, want FirstSuite, EmptySuite, and SecondSuite", suiteNames(got))
	}
	assertMethods(t, got["FirstSuite"], "FirstSuite", []string{"TestFoo", "TestBar"}, []string{"CasesFoo"}, []string{"BeforeEach", "AfterAll"})
	assertMethods(t, got["EmptySuite"], "EmptySuite", nil, nil, nil)
	assertMethods(t, got["SecondSuite"], "SecondSuite", []string{"TestBaz"}, nil, nil)
}

func assertMethods(t *testing.T, suite *Suite, name string, tests, cases, hooks []string) {
	t.Helper()
	if suite == nil {
		t.Fatalf("suite %s was not discovered", name)
	}
	assertMethodNames(t, name+" tests", suite.Tests, tests)
	assertMethodNames(t, name+" cases", suite.Cases, cases)
	assertMethodNames(t, name+" hooks", suite.Hooks, hooks)
}

func assertMethodNames(t *testing.T, label string, methods []*Method, want []string) {
	t.Helper()
	if len(methods) != len(want) {
		t.Fatalf("%s: got %d methods, want %d", label, len(methods), len(want))
	}
	got := make(map[string]bool, len(methods))
	for _, method := range methods {
		if method == nil || method.Decl == nil || method.Func == nil {
			t.Fatalf("%s contains an incomplete method", label)
		}
		got[method.Decl.Name.Name] = true
		if method.Func.Name() != method.Decl.Name.Name {
			t.Errorf("%s: types.Func name %q does not match declaration %q", label, method.Func.Name(), method.Decl.Name.Name)
		}
	}
	for _, name := range want {
		if !got[name] {
			t.Errorf("%s does not contain %s", label, name)
		}
	}
}

func suiteNames(suites map[string]*Suite) []string {
	names := make([]string, 0, len(suites))
	for name := range suites {
		names = append(names, name)
	}
	return names
}
