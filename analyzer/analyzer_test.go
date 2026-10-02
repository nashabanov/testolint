package analyzer

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	tests := []string{
		"wrongpackage",
		"test_signature",
		"hook_signature",
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
