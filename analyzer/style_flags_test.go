package analyzer

import "testing"

func TestStyleFlags(t *testing.T) {
	for _, value := range []string{
		`{"actions":[{"package":"x/...","function":"*"}],"assertions":[{"package":"x/sub/...","function":"*"}]}`,
		`{"actions":[{"package":"x","receiver":"Client","method":"Send*"}]}`, `null`, `[]`, `{"unknown":true}`, `{"phases":{"strict":"false"}}`, `{"header":{"required":[[]]}}`, `{"actions":[{"package":"x","function":"x.y"}]}`, `{} {}`} {
		a, err := New(Config{})
		if err != nil {
			t.Fatal(err)
		}
		if err := a.Flags.Set("style-config", value); err == nil {
			t.Errorf("accepted %s", value)
		}
	}
	for _, strictFirst := range []bool{true, false} {
		a, err := New(Config{})
		if err != nil {
			t.Fatal(err)
		}
		settings := [][2]string{{"style-config", `{"phases":{"strict":true}}`}, {"phase-order-strict", "false"}}
		if strictFirst {
			settings[0], settings[1] = settings[1], settings[0]
		}
		for _, setting := range settings {
			if err := a.Flags.Set(setting[0], setting[1]); err != nil {
				t.Fatal(err)
			}
		}
		if got := a.Flags.Lookup("phase-order-strict").Value.String(); got != "false" {
			t.Fatalf("strict=%s", got)
		}
		if err := a.Flags.Set("style-config", `{}`); err == nil {
			t.Fatal("accepted duplicate style flag")
		}
		if err := a.Flags.Set("phase-order-strict", "true"); err == nil {
			t.Fatal("accepted duplicate phase flag")
		}
	}
}

func TestOperationPatternConflicts(t *testing.T) {
	for _, pair := range [][2]string{{"x/...", "x/sub/..."}, {"x/sub/...", "x/..."}} {
		for _, method := range []bool{false, true} {
			a, b := Symbol{Package: pair[0], Function: "*"}, Symbol{Package: pair[1], Function: "*"}
			if method {
				a.Function, b.Function = "", ""
				a.Receiver, b.Receiver = "*", "*"
				a.Method, b.Method = "*", "*"
			}
			for _, kind := range []string{"assertions", "ignore"} {
				cfg := Config{Style: StyleConfig{Actions: []Symbol{a}}}
				if kind == "assertions" {
					cfg.Style.Assertions = []Symbol{b}
				} else {
					cfg.Style.Ignore = []Symbol{b}
				}
				if _, err := New(cfg); err == nil {
					t.Errorf("accepted %s conflict: %+v / %+v", kind, a, b)
				}
				cfg.Style.Actions, cfg.Style.Assertions = nil, []Symbol{a}
				cfg.Style.Ignore = []Symbol{b}
				if _, err := New(cfg); err == nil {
					t.Errorf("accepted assertion/ignore conflict: %+v / %+v", a, b)
				}
			}
		}
	}
	for _, pair := range [][2]string{{"x/sub/...", "x/submarine/..."}, {"x", "x/sub/..."}} {
		if _, err := New(Config{Style: StyleConfig{
			Actions:    []Symbol{{Package: pair[0], Function: "*"}},
			Assertions: []Symbol{{Package: pair[1], Function: "*"}},
		}}); err != nil {
			t.Fatal(err)
		}
	}
}
