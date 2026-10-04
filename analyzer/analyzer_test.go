package analyzer

import (
	"testing"

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
