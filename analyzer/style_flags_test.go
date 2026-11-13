package analyzer

import "testing"

func TestStyleFlags(t *testing.T) {
	for _, value := range []string{`null`, `[]`, `{"unknown":true}`, `{"phases":{"strict":"false"}}`, `{"header":{"required":[[]]}}`, `{"actions":[{"package":"x","function":"x.y"}]}`, `{} {}`} {
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
