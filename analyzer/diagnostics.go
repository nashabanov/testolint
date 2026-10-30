package analyzer

import (
	"fmt"
	"go/token"
	"reflect"

	"golang.org/x/tools/go/analysis"
)

// reportf accepts a message without the diagnostic ID prefix.
func reportf(p *analysis.Pass, id ruleID, pos token.Pos, format string, args ...any) {
	report(p, id, analysis.Diagnostic{Pos: pos, Message: fmt.Sprintf(format, args...)})
}

// report assigns the registry identity while preserving ranges, related locations
// and suggested fixes. Category is the diagnostic ID; URL is the registry DocURL.
// Unknown IDs are programming errors, rather than silently unclassified reports.
func report(p *analysis.Pass, id ruleID, diagnostic analysis.Diagnostic) {
	meta, ok := registry[id]
	if !ok {
		panic(fmt.Sprintf("testolint: unknown diagnostic ID %q", id))
	}
	diagnostic.Message = string(id) + ": " + diagnostic.Message
	diagnostic.Category = string(id)
	diagnostic.URL = meta.DocURL
	p.Report(diagnostic)
}

// deduplicatingReporter retains the first occurrence of each complete diagnostic.
// Different ranges, related locations or fixes carry information and are retained.
func deduplicatingReporter(next func(analysis.Diagnostic)) func(analysis.Diagnostic) {
	type key struct {
		pos     token.Pos
		message string
	}
	seen := make(map[key][]analysis.Diagnostic)
	return func(d analysis.Diagnostic) {
		k := key{d.Pos, d.Message}
		for _, previous := range seen[k] {
			if reflect.DeepEqual(previous, d) {
				return
			}
		}
		seen[k] = append(seen[k], d)
		next(d)
	}
}
