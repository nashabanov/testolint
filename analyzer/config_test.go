package analyzer

import (
	"io"
	"reflect"
	"strings"
	"testing"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestConfigValidation(t *testing.T) {
	for _, tc := range []struct {
		name   string
		config Config
		want   string
	}{
		{"unknown enable", Config{Enable: []string{"TESTO999"}}, "unknown rule"},
		{"unknown disable", Config{Disable: []string{"TESTO999"}}, "unknown rule"},
		{"empty", Config{Enable: []string{""}}, "empty rule ID"},
		{"case", Config{Enable: []string{"testo019"}}, "unknown rule"},
		{"whitespace", Config{Enable: []string{" TESTO019"}}, "unknown rule"},
		{"duplicate enable", Config{Enable: []string{"TESTO019", "TESTO019"}}, "duplicate rule"},
		{"duplicate disable", Config{Disable: []string{"TESTO019", "TESTO019"}}, "duplicate rule"},
		{"conflict", Config{Enable: []string{"TESTO019"}, Disable: []string{"TESTO019"}}, "both enable and disable"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := New(tc.config); err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("expected %q, got %v", tc.want, err)
			}
		})
	}
}

func TestSelectionFlags(t *testing.T) {
	for _, flags := range [][]string{
		{"-enable=TESTO019", "-disable=TESTO019"},
		{"-disable=TESTO019", "-enable=TESTO019"},
		{"-enable=TESTO019,"}, {"-disable=TESTO999"},
		{"-enable=TESTO019", "-enable=TESTO018"},
		{"-disable-all=invalid"},
	} {
		t.Run(strings.Join(flags, " "), func(t *testing.T) {
			a, err := New(Config{})
			if err != nil {
				t.Fatal(err)
			}
			a.Flags.SetOutput(io.Discard)
			if err := a.Flags.Parse(flags); err == nil {
				t.Error("invalid flags accepted")
			}
		})
	}
	var zero selectionFlag
	if zero.String() != "" {
		t.Error("zero flag must support flag package help rendering")
	}
}

// Reuse the no-diagnostic fixture to assert each selection directly against a
// baseline, including suppressed construction errors and both multi-ID checkers.
func TestConfiguredInstances(t *testing.T) {
	wrapper := &analysis.Analyzer{Name: "configinstances", Doc: "check configuration isolation"}
	wrapper.Run = func(p *analysis.Pass) (any, error) {
		input := Config{DisableAll: true, Enable: []string{"TESTO019"}, Disable: []string{"TESTO006"}}
		callback, err := New(input)
		if err != nil {
			return nil, err
		}
		input.Enable[0] = "TESTO018"
		input.Disable[0] = "TESTO019" // Mutation must not affect the constructed analyzer.
		flags, err := New(Config{})
		if err != nil {
			return nil, err
		}
		if err := flags.Flags.Parse([]string{"-enable= TESTO019 ", "-disable-all", "-disable="}); err != nil {
			return nil, err
		}
		emptyFlags, err := New(Config{})
		if err != nil {
			return nil, err
		}
		if err := emptyFlags.Flags.Parse([]string{"-enable=", "-disable="}); err != nil {
			return nil, err
		}
		capture := func(a *analysis.Analyzer) []analysis.Diagnostic {
			var ds []analysis.Diagnostic
			pass := *p
			pass.Report = func(d analysis.Diagnostic) { ds = append(ds, d) }
			if _, err := a.Run(&pass); err != nil {
				t.Fatal(err)
			}
			return ds
		}
		baseline := capture(Analyzer)
		if len(baseline) != 6 {
			t.Fatalf("baseline diagnostic count = %d, want 6", len(baseline))
		}
		for _, id := range []string{"TESTO006", "TESTO012", "TESTO018", "TESTO019", "TESTO021", "TESTO022"} {
			a, err := New(Config{DisableAll: true, Enable: []string{id}})
			if err != nil {
				return nil, err
			}
			var want []analysis.Diagnostic
			for _, d := range baseline {
				if d.Category == id {
					want = append(want, d)
				}
			}
			for range 2 {
				if got := capture(a); !reflect.DeepEqual(got, want) {
					t.Errorf("%s: got %+v, want %+v", id, got, want)
				}
			}
		}
		for _, a := range []*analysis.Analyzer{callback, flags} {
			got := capture(a)
			if len(got) != 1 || got[0].Category != "TESTO019" {
				t.Errorf("callback instance selection leaked: %+v", got)
			}
		}
		if got := capture(emptyFlags); !reflect.DeepEqual(got, baseline) {
			t.Error("empty lists changed defaults")
		}
		// Default singleton remains unchanged after configuring and running instances.
		if got := capture(Analyzer); !reflect.DeepEqual(got, baseline) {
			t.Error("default analyzer mutated")
		}
		return nil, nil
	}
	analysistest.Run(t, analysistest.TestData(), wrapper, "config_none")
}

func TestDefaultRegistrySelection(t *testing.T) {
	selected, err := resolveConfig(Config{})
	if err != nil {
		t.Fatal(err)
	}
	for id, meta := range registry {
		if selected[id] != meta.EnableByDefault {
			t.Errorf("default mismatch for %s", id)
		}
	}
	none, err := resolveConfig(Config{DisableAll: true})
	if err != nil {
		t.Fatal(err)
	}
	rules := selectedRules(none)
	if !reflect.DeepEqual(rules, ruleSet{}) {
		t.Fatal("disabled checks still scheduled")
	}
	both, err := resolveConfig(Config{DisableAll: true, Enable: []string{"TESTO018", "TESTO019", "TESTO021", "TESTO022"}})
	if err != nil {
		t.Fatal(err)
	}
	if got := len(selectedRules(both).packages); got != 2 {
		t.Errorf("multi-ID checkers registered %d times, want 2", got)
	}
}
