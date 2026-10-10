package analyzer

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"testing"

	"golang.org/x/tools/go/analysis"
)

func TestDiscoverRunsMissingTypeInfo(t *testing.T) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "runs.go", `package p
 func f() { run := testo.RunSuite[S]; run(t, s); testo.RunSubSuite(t, s) }
 `, 0)
	if err != nil {
		t.Fatal(err)
	}
	for _, info := range []*types.Info{nil, {}} {
		p := &analysis.Pass{Fset: fset, Files: []*ast.File{file}, TypesInfo: info}
		if got := discoverRuns(p); len(got) != 0 {
			t.Fatalf("unexpected runs without type information: %v", got)
		}
	}
}
