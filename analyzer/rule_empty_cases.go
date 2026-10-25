package analyzer

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

func checkEmptyCases(p *analysis.Pass, suite *suite) {
	for _, provider := range suite.cases {
		if !isValidPrefixedName(provider.fn.Name(), "Cases") {
			continue // TESTO008: not a runtime cases provider.
		}
		if slice, _ := casesSignature(provider.fn.Type()); slice == nil {
			continue // TESTO005: avoid secondary diagnostics.
		}
		decl := provider.decl
		if decl == nil || decl.Body == nil || len(decl.Body.List) != 1 {
			continue
		}
		ret, ok := decl.Body.List[0].(*ast.ReturnStmt)
		if !ok || len(ret.Results) != 1 {
			continue
		}
		empty := false
		switch expr := ast.Unparen(ret.Results[0]).(type) {
		case *ast.Ident:
			empty = p.TypesInfo.ObjectOf(expr) == types.Universe.Lookup("nil")
		case *ast.CompositeLit:
			t := p.TypesInfo.TypeOf(expr)
			empty = len(expr.Elts) == 0 && t != nil && isSliceType(t)
		}
		if empty {
			p.Reportf(provider.pos(suite),
				"TESTO010: %s always returns an empty case set", provider.fn.Name())
		}
	}
}
