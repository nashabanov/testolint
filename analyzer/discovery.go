package analyzer

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

func findTestMethods(p *analysis.Pass) []*ast.FuncDecl {
	var tests []*ast.FuncDecl

	for _, file := range p.Files {
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			if fn.Recv == nil {
				continue
			}

			if !strings.HasPrefix(fn.Name.Name, "Test") {
				continue
			}

			tests = append(tests, fn)
		}
	}

	return tests
}
