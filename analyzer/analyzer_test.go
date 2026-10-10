package analyzer

import (
	"go/types"
	"reflect"
	"testing"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	tests := []string{
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
							if suite.TType == nil || suite.MethodSet == nil || suite.CasesByName == nil {
								t.Errorf("incomplete model for %s", suite.Type.Obj().Name())
							}
							p.Reportf(suite.Type.Obj().Pos(), "TESTO_DISCOVERY: suite %s", suite.Type.Obj().Name())
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
							if call.SuiteType == nil || call.SuiteExpr != call.Call.Args[1] {
								t.Fatal("incomplete run model")
							}
							name := "RunSuite"
							if call.Kind == RunSubSuite {
								name = "RunSubSuite"
							}
							p.Reportf(call.SuiteExpr.Pos(), "TESTO_RUN_DISCOVERY: %s %s", name, call.SuiteType)
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
						reportRun := func(p *analysis.Pass, call *RunCall) {
							p.Reportf(call.SuiteExpr.Pos(), "TESTO_RUN_PIPELINE: run")
							// The same diagnostic is also emitted by a suite rule.
							if named, ok := call.SuiteType.(*types.Named); ok && named.Obj().Name() == "Empty" {
								p.Reportf(named.Obj().Pos(), "TESTO011: suite %q contains no tests", "Empty")
							}
						}
						reportOther := func(p *analysis.Pass, call *RunCall) {
							p.Reportf(call.SuiteExpr.Pos(), "TESTO_RUN_PIPELINE: other")
						}
						var first []analysis.Diagnostic
						// Repeat the pass to check deterministic order and that
						// deduplication does not leak into the next invocation.
						for i := 0; i < 2; i++ {
							var diagnostics []analysis.Diagnostic
							pass := *p
							pass.Report = func(d analysis.Diagnostic) { diagnostics = append(diagnostics, d) }
							if _, err := runWithRules(&pass, suiteRules(), []RunRule{reportRun, reportRun, reportOther}); err != nil {
								return nil, err
							}
							if i == 0 {
								first = diagnostics
							} else if !reflect.DeepEqual(first, diagnostics) {
								t.Error("pipeline diagnostics changed between invocations")
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
