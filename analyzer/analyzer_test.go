package analyzer

import (
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
			analysistest.Run(
				t,
				analysistest.TestData(),
				analyzer,
				tt,
			)
		})
	}
}
