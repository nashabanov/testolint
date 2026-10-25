package analyzer

import (
	"go/types"
	"reflect"
	"strings"
	"testing"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	tests := []string{
		"construction",
		"execution_safety",
		"premature_teardown",
		"parallel_safety",
		"parallel_coverage",
		"duplicate_parallel",
		"wrongpackage",
		"robustness",
		"no_suites",
		"test_signature",
		"hook_signature",
		"cases_signature",
		"parametrization",
		"naming",
		"param_visibility",
		"discovery",
		"run_discovery",
		"run_pipeline",
		"nil_suite",
		"recursive_subsuite",
		"promotion",
		"empty_cases",
		"empty_suite",
	}

	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			analyzer := Analyzer
			if tt == "discovery" {
				// Observe suites with no declared methods through analysistest
				// diagnostics, without adding a public rule for empty suites.
				analyzer = &analysis.Analyzer{
					Name: "discoverytest",
					Doc:  "check discovered suite models and ordinary diagnostics",
					Run: func(p *analysis.Pass) (any, error) {
						for _, suite := range discoverSuites(p) {
							if suite.tType == nil || suite.methodSet == nil || suite.casesByName == nil {
								t.Errorf("incomplete model for %s", suite.typ.Obj().Name())
							}
							p.Reportf(suite.typ.Obj().Pos(), "TESTO_DISCOVERY: suite %s", suite.typ.Obj().Name())
						}
						return run(p)
					},
				}
			}
			if tt == "run_discovery" {
				analyzer = &analysis.Analyzer{
					Name: "rundiscoverytest",
					Doc:  "check discovered run models and ordinary diagnostics",
					Run: func(p *analysis.Pass) (any, error) {
						for _, call := range discoverRuns(p) {
							if call.suiteType == nil || call.suiteExpr != call.call.Args[1] {
								t.Fatal("incomplete run model")
							}
							name := "RunSuite"
							if call.kind == runSubSuite {
								name = "RunSubSuite"
							}
							p.Reportf(call.suiteExpr.Pos(), "TESTO_RUN_DISCOVERY: %s %s", name, call.suiteType)
						}
						return run(p)
					},
				}
			}
			if tt == "run_pipeline" {
				analyzer = &analysis.Analyzer{
					Name: "runpipelinetest",
					Doc:  "check suite and run rules together",
					Run: func(p *analysis.Pass) (any, error) {
						reportRun := func(p *analysis.Pass, call *runCall) {
							p.Reportf(call.suiteExpr.Pos(), "TESTO_RUN_PIPELINE: run")
							// The same diagnostic is also emitted by a suite rule.
							if named, ok := call.suiteType.(*types.Named); ok && named.Obj().Name() == "Empty" {
								p.Reportf(named.Obj().Pos(), "TESTO011: suite %q contains no tests", "Empty")
							}
						}
						reportOther := func(p *analysis.Pass, call *runCall) {
							p.Reportf(call.suiteExpr.Pos(), "TESTO_RUN_PIPELINE: other")
						}
						reportBeforeParallel := func(p *analysis.Pass, suite *suite) {
							if suite.typ.Obj().Name() == "Parallel" {
								p.Reportf(suite.typ.Obj().Pos(), "TESTO_RUN_PIPELINE: before parallel")
							}
						}
						reportAfterParallel := func(p *analysis.Pass, suite *suite) {
							if suite.typ.Obj().Name() == "Parallel" {
								p.Reportf(suite.typ.Obj().Pos(), "TESTO_RUN_PIPELINE: after parallel")
							}
						}
						suiteChecks := append([]suiteRule{reportBeforeParallel}, suiteRules()...)
						afterParallelChecks := append(afterParallelRules(), reportAfterParallel)
						// Register package and parallel checks twice to exercise the
						// same global deduplication used by suite and run rules.
						packageChecks := packageRules()
						packageChecks = append(packageChecks, packageRules()...)
						parallelChecks := parallelRules()
						parallelChecks = append(parallelChecks, parallelRules()...)
						checks := ruleSet{
							suites:        suiteChecks,
							parallel:      parallelChecks,
							afterParallel: afterParallelChecks,
							runs:          []runRule{reportRun, reportRun, reportOther},
							packages:      packageChecks,
						}
						var first []analysis.Diagnostic
						// Repeat the pass to check deterministic order and that
						// deduplication does not leak into the next invocation.
						for i := range 2 {
							var diagnostics []analysis.Diagnostic
							pass := *p
							pass.Report = func(d analysis.Diagnostic) { diagnostics = append(diagnostics, d) }
							if _, err := runWithRules(&pass, checks); err != nil {
								return nil, err
							}
							if i == 0 {
								first = diagnostics
							} else if !reflect.DeepEqual(first, diagnostics) {
								t.Error("pipeline diagnostics changed between invocations")
							}
						}
						var parallelOrder []string
						for _, d := range first {
							switch {
							case strings.HasPrefix(d.Message, "TESTO014:"):
								parallelOrder = append(parallelOrder, "TESTO014")
							case strings.HasPrefix(d.Message, "TESTO017:"):
								parallelOrder = append(parallelOrder, "TESTO017")
							case d.Message == "TESTO_RUN_PIPELINE: before parallel", d.Message == "TESTO_RUN_PIPELINE: after parallel":
								parallelOrder = append(parallelOrder, d.Message)
							}
						}
						wantParallelOrder := []string{"TESTO_RUN_PIPELINE: before parallel", "TESTO014", "TESTO017", "TESTO_RUN_PIPELINE: after parallel"}
						if !reflect.DeepEqual(parallelOrder, wantParallelOrder) {
							t.Errorf("parallel diagnostic order = %v, want %v", parallelOrder, wantParallelOrder)
						}
						// Selected groups keep the same diagnostics and order as the
						// complete pipeline; omitting parallel rules leaves all others.
						var wantPackages, wantParallel, wantWithoutParallel []analysis.Diagnostic
						for _, d := range first {
							if strings.HasPrefix(d.Message, "TESTO019:") || strings.HasPrefix(d.Message, "TESTO021:") {
								wantPackages = append(wantPackages, d)
							}
							if strings.HasPrefix(d.Message, "TESTO014:") || strings.HasPrefix(d.Message, "TESTO017:") {
								wantParallel = append(wantParallel, d)
							} else {
								wantWithoutParallel = append(wantWithoutParallel, d)
							}
						}
						withoutParallel := checks
						withoutParallel.parallel = nil
						for _, selected := range []struct {
							name  string
							rules ruleSet
							want  []analysis.Diagnostic
						}{
							{"packages only", ruleSet{packages: packageChecks}, wantPackages},
							{"parallel only", ruleSet{parallel: parallelChecks}, wantParallel},
							{"without parallel", withoutParallel, wantWithoutParallel},
							{"no rules", ruleSet{}, nil},
						} {
							var diagnostics []analysis.Diagnostic
							pass := *p
							pass.Report = func(d analysis.Diagnostic) { diagnostics = append(diagnostics, d) }
							if _, err := runWithRules(&pass, selected.rules); err != nil {
								return nil, err
							}
							if !reflect.DeepEqual(diagnostics, selected.want) {
								t.Errorf("%s diagnostics = %v, want %v", selected.name, diagnostics, selected.want)
							}
						}
						for _, d := range first {
							p.Report(d)
						}
						return nil, nil
					},
				}
			}
			analysistest.Run(
				t,
				analysistest.TestData(),
				analyzer,
				tt,
			)
		})
	}
}
