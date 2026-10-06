package analyzer

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

func checkEmptyCases(p *analysis.Pass, suite *Suite) {
	for _, provider := range suite.Cases {
		if !isValidPrefixedName(provider.Func.Name(), "Cases") {
			continue // TESTO008: not a runtime cases provider.
		}
		sig, ok := provider.Func.Type().(*types.Signature)
		if !ok || sig.Params().Len() != 0 || sig.Results().Len() != 1 ||
			!isSliceType(sig.Results().At(0).Type()) {
			continue // TESTO005: avoid secondary diagnostics.
		}
		decl := provider.Decl
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
			p.Reportf(provider.Pos(suite),
				"TESTO010: %s always returns an empty case set", provider.Func.Name())
		}
	}
}
