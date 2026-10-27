package analyzer

import (
	"go/ast"
	"go/constant"
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
		if emptySliceExpr(p.TypesInfo, ret.Results[0]) {
			p.Reportf(provider.pos(suite),
				"TESTO010: %s always returns an empty case set", provider.fn.Name())
		}
	}
}

// Keep empty-case diagnostics and lifecycle multiplicity in agreement. Unknown
// capacities may execute arbitrary code or panic, so only constants qualify.
func emptySliceExpr(info *types.Info, expression ast.Expr) bool {
	switch expr := ast.Unparen(expression).(type) {
	case *ast.Ident:
		return info.ObjectOf(expr) == types.Universe.Lookup("nil")
	case *ast.CompositeLit:
		t := info.TypeOf(expr)
		return len(expr.Elts) == 0 && t != nil && isSliceType(t)
	case *ast.CallExpr:
		if len(expr.Args) < 2 || len(expr.Args) > 3 || expr.Ellipsis.IsValid() ||
			!sameObject(info, expr.Fun, types.Universe.Lookup("make")) {
			return false
		}
		t := info.TypeOf(expr)
		length := info.Types[expr.Args[1]].Value
		if t == nil || !isSliceType(t) || length == nil || constant.Sign(length) != 0 {
			return false
		}
		if len(expr.Args) == 3 {
			capacity := info.Types[expr.Args[2]].Value
			return capacity != nil && constant.Sign(capacity) >= 0
		}
		return true
	}
	return false
}
