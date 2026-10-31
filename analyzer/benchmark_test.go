package analyzer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/packages"
)

// BenchmarkAnalyzer measures the complete default analyzer on preloaded packages.
// Source generation, loading, type checking and validation are outside the timer.
func BenchmarkAnalyzer(b *testing.B) {
	for _, scenario := range []string{"NoTesto", "Valid", "Diagnostics", "Complex"} {
		for _, size := range []int{1, 10, 100} {
			b.Run(fmt.Sprintf("%s/%d", scenario, size), func(b *testing.B) {
				b.StopTimer()
				pass := benchmarkPass(b, benchmarkSource(scenario, size))
				count := 0
				pass.Report = func(analysis.Diagnostic) { count++ }
				if _, err := Analyzer.Run(pass); err != nil {
					b.Fatal(err)
				}
				want := 0
				if scenario == "Diagnostics" {
					want = size
				}
				if count != want {
					b.Fatalf("diagnostics = %d, want %d", count, want)
				}
				count = 0
				b.ReportAllocs()
				b.ResetTimer()
				b.StartTimer()
				for i := 0; i < b.N; i++ {
					if _, err := Analyzer.Run(pass); err != nil {
						b.Fatal(err)
					}
				}
				b.StopTimer()
				if count != b.N*want {
					b.Fatalf("diagnostics = %d, want %d", count, b.N*want)
				}
			})
		}
	}
}

func benchmarkPass(b *testing.B, source string) *analysis.Pass {
	b.Helper()
	dir := b.TempDir()
	stub, err := filepath.Abs("testdata/src/github.com/ozontech/testo")
	if err != nil {
		b.Fatal(err)
	}
	module := fmt.Sprintf("module benchmarkfixture\n\ngo 1.27.1\n\nrequire github.com/ozontech/testo v0.0.0\nreplace github.com/ozontech/testo => %q\n", stub)
	for name, content := range map[string]string{"go.mod": module, "suite.go": source} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600); err != nil {
			b.Fatal(err)
		}
	}
	pkgs, err := packages.Load(&packages.Config{Dir: dir, Mode: packages.LoadAllSyntax}, ".")
	if err != nil {
		b.Fatal(err)
	}
	if packages.PrintErrors(pkgs) != 0 {
		b.Fatal("benchmark fixture failed to load")
	}
	if len(pkgs) != 1 {
		b.Fatalf("loaded %d packages, want 1", len(pkgs))
	}
	pkg := pkgs[0]
	return &analysis.Pass{Analyzer: Analyzer, Fset: pkg.Fset, Files: pkg.Syntax,
		Pkg: pkg.Types, TypesInfo: pkg.TypesInfo, TypesSizes: pkg.TypesSizes}
}

func benchmarkSource(scenario string, size int) string {
	var src strings.Builder
	src.WriteString("package benchmarkfixture\n")
	if scenario != "NoTesto" {
		src.WriteString("import \"github.com/ozontech/testo\"\ntype T = *testo.T\n")
	}
	for i := 0; i < size; i++ {
		switch scenario {
		case "NoTesto":
			fmt.Fprintf(&src, "type Plain%d struct { Value int }\nfunc (s *Plain%d) TestValue(v int) int { return s.Value + v }\n", i, i)
		case "Valid":
			fmt.Fprintf(&src, "type Suite%d struct { testo.Suite[T] }\nfunc (s *Suite%d) TestValue(t T) {}\n", i, i)
		case "Diagnostics":
			fmt.Fprintf(&src, "type Suite%d struct { testo.Suite[T] }\nfunc (s *Suite%d) TestValue() {}\n", i, i)
		case "Complex":
			fmt.Fprintf(&src, `
type Methods%d struct{}
func (*Methods%d) TestValue(t T, p struct { Value int }) {}
func (*Methods%d) CasesValue() []int { return []int{1, 2, 3} }
type Base%d struct { testo.Suite[T]; Methods%d }
type Suite%d struct { *Base%d }
func (*Suite%d) BeforeEach(t T) {}
`, i, i, i, i, i, i, i, i)
		}
	}
	return src.String()
}
