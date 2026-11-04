package analyzer

import "testing"

func TestEmptyPipeline(t *testing.T) {
	// An empty selection needs neither package data nor a reporter.
	result, err := runWithRules(nil, ruleSet{})
	if result != nil || err != nil {
		t.Fatalf("empty pipeline = (%v, %v), want (nil, nil)", result, err)
	}
}
