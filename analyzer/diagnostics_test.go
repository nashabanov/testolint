package analyzer

import (
	"go/token"
	"reflect"
	"testing"

	"golang.org/x/tools/go/analysis"
)

func TestReport(t *testing.T) {
	original := analysis.Diagnostic{
		Pos: 10, End: 20, Message: "message",
		Related:        []analysis.RelatedInformation{{Pos: 30, End: 40, Message: "related"}},
		SuggestedFixes: []analysis.SuggestedFix{{Message: "fix", TextEdits: []analysis.TextEdit{{Pos: 10, End: 20, NewText: []byte("replacement")}}}},
	}
	var got []analysis.Diagnostic
	pass := &analysis.Pass{Report: func(d analysis.Diagnostic) { got = append(got, d) }}
	report(pass, "TESTO004", original)
	want := original
	want.Message = "TESTO004: message"
	want.Category = "TESTO004"
	want.URL = registry["TESTO004"].DocURL
	if !reflect.DeepEqual(got, []analysis.Diagnostic{want}) {
		t.Fatalf("report = %+v, want %+v", got, want)
	}
	if original.Message != "message" {
		t.Fatal("input diagnostic mutated")
	}
	reportf(pass, "TESTO006", token.Pos(50), "unused %q", "provider")
	if got[1].Message != `TESTO006: unused "provider"` || got[1].Category != "TESTO006" || got[1].URL != registry["TESTO006"].DocURL || got[1].Pos != 50 {
		t.Fatalf("reportf = %+v", got[1])
	}
}

func TestReportUnknownID(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("unknown ID did not panic")
		}
	}()
	report(&analysis.Pass{Report: func(analysis.Diagnostic) { t.Error("unknown ID reported") }}, "TESTO999", analysis.Diagnostic{})
}

func TestDeduplicatingReporter(t *testing.T) {
	base := analysis.Diagnostic{Pos: 10, Message: "TESTO004: mismatch"}
	variants := []analysis.Diagnostic{base}
	end := base
	end.End = 20
	variants = append(variants, end)
	related := base
	related.Related = []analysis.RelatedInformation{{Pos: 30, Message: "provider"}}
	variants = append(variants, related)
	fix := base
	fix.SuggestedFixes = []analysis.SuggestedFix{{Message: "fix", TextEdits: []analysis.TextEdit{{Pos: 10, NewText: []byte("value")}}}}
	variants = append(variants, fix)
	category := base
	category.Category = "TESTO004"
	variants = append(variants, category)
	url := base
	url.URL = ruleDocURL
	variants = append(variants, url)
	var got []analysis.Diagnostic
	emit := deduplicatingReporter(func(d analysis.Diagnostic) { got = append(got, d) })
	for _, d := range variants {
		emit(d)
		emit(d)
	}
	if !reflect.DeepEqual(got, variants) {
		t.Fatalf("deduplication lost context or order: %+v", got)
	}
	var second []analysis.Diagnostic
	deduplicatingReporter(func(d analysis.Diagnostic) { second = append(second, d) })(base)
	if len(second) != 1 {
		t.Fatal("deduplication state shared between passes")
	}
}
