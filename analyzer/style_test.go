package analyzer

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestStyleFixes(t *testing.T) {
	a, err := New(Config{Enable: []string{"TESTO025", "TESTO026"}})
	if err != nil {
		t.Fatal(err)
	}
	analysistest.RunWithSuggestedFixes(t, analysistest.TestData(), a, "style")
}

func TestStyleSelection(t *testing.T) {
	wrapper := &analysis.Analyzer{Name: "styleselection", Doc: "check independent style selections"}
	wrapper.Run = func(p *analysis.Pass) (any, error) {
		for _, tc := range []struct {
			config        Config
			header, steps int
		}{
			{Config{}, 0, 0},
			{Config{Enable: []string{"TESTO025"}}, 5, 0},
			{Config{Enable: []string{"TESTO026"}}, 0, 4},
			{Config{Enable: []string{"TESTO025", "TESTO026"}, DisableAll: true}, 5, 4},
		} {
			a, err := New(tc.config)
			if err != nil {
				return nil, err
			}
			header, steps := 0, 0
			pass := *p
			pass.Report = func(d analysis.Diagnostic) {
				if len(d.SuggestedFixes) > 0 {
					file := p.Fset.File(d.Pos)
					source, err := p.ReadFile(file.Name())
					if err != nil {
						t.Fatal(err)
					}
					want := []byte("\n")
					if bytes.Contains(source, []byte("\r\n")) {
						want = []byte("\r\n")
					}
					for _, edit := range d.SuggestedFixes[0].TextEdits {
						if !bytes.Equal(edit.NewText, want) {
							t.Errorf("line ending: got %q, want %q", edit.NewText, want)
						}
					}
				}
				switch d.Category {
				case "TESTO025":
					header++
				case "TESTO026":
					steps++
				default:
					t.Errorf("unexpected diagnostic: %v", d)
				}
			}
			if _, err := a.Run(&pass); err != nil {
				return nil, err
			}
			if header != tc.header || steps != tc.steps {
				t.Errorf("%+v: counts %d/%d, want %d/%d", tc.config, header, steps, tc.header, tc.steps)
			}
		}
		enabled, _ := New(Config{Enable: []string{"TESTO025", "TESTO026"}})
		return enabled.Run(p)
	}
	analysistest.Run(t, analysistest.TestData(), wrapper, "style")
}

// Reanalyze the golden output to ensure fixes leave no style diagnostics.
func TestStyleFixed(t *testing.T) {
	dir := t.TempDir()
	if err := os.CopyFS(dir, os.DirFS(analysistest.TestData())); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"test.go", "crlf.go"} {
		path := filepath.Join(dir, "src", "style", name)
		source, err := os.ReadFile(path + ".golden")
		if err != nil {
			t.Fatal(err)
		}
		source = regexp.MustCompile(` // want "TESTO0(25|26)"`).ReplaceAll(source, nil)
		if err := os.WriteFile(path, source, 0600); err != nil {
			t.Fatal(err)
		}
	}
	a, err := New(Config{Enable: []string{"TESTO025", "TESTO026"}})
	if err != nil {
		t.Fatal(err)
	}
	analysistest.Run(t, dir, a, "style")
}
