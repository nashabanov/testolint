package analyzer

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
	"testing"
)

func TestStructureSelection(t *testing.T) {
	style := StyleConfig{Actions: []Symbol{{Package: "step_coverage", Function: "action"}, {Package: "step_coverage/remote", Receiver: "Client", Method: "Send"}, {Package: "step_coverage/remote", Function: "Send"}}, Assertions: []Symbol{{Package: "step_coverage", Function: "verify"}}}
	wrapper := &analysis.Analyzer{Name: "structureselection", Doc: "verify independent coverage selections"}
	wrapper.Run = func(p *analysis.Pass) (any, error) {
		walker := newOperationWalker(p, style)
		walker.collectPhases = false
		for _, s := range discoverSuites(p) {
			for _, test := range s.tests {
				if validStyleTest(s, test) {
					walker.walk(test.method.decl.Body, false, nil, 0)
				}
			}
		}
		if len(walker.events) != 0 {
			t.Fatal("phase events collected with phase analysis disabled")
		}
		for _, facts := range walker.helpers {
			if len(facts.events) != 0 {
				t.Fatal("cached disabled phase events")
			}
		}

		for _, tc := range []struct {
			ids                 []string
			actions, assertions int
		}{
			{nil, 0, 0},
			{[]string{"TESTO029"}, 7, 0},
			{[]string{"TESTO031"}, 0, 5},
			{[]string{"TESTO029", "TESTO031"}, 7, 5},
		} {
			a, err := New(Config{DisableAll: true, Enable: tc.ids, Style: style})
			if err != nil {
				return nil, err
			}
			counts := map[string]int{}
			pass := *p
			pass.Report = func(d analysis.Diagnostic) { counts[d.Category]++ }
			if _, err := a.Run(&pass); err != nil {
				return nil, err
			}
			if counts["TESTO029"] != tc.actions || counts["TESTO031"] != tc.assertions || counts["TESTO030"] != 0 {
				t.Errorf("%v: %v", tc.ids, counts)
			}
		}
		a, err := New(Config{DisableAll: true, Enable: []string{"TESTO029", "TESTO031"}, Style: style})
		if err != nil {
			return nil, err
		}
		return a.Run(p)
	}
	analysistest.Run(t, analysistest.TestData(), wrapper, "step_coverage")
}
