package analyzer

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

func checkNilSuite(p *analysis.Pass, call *runCall) {
	if p.TypesInfo == nil || call.suiteType == nil {
		return
	}
	// Only pointer suites and explicit untyped nil are supported. In particular,
	// a typed nil pointer boxed in an interface is not a nil interface.
	_, pointer := call.suiteType.Underlying().(*types.Pointer)
	if !pointer && call.suiteType != types.Typ[types.UntypedNil] {
		return
	}
	if call.suiteValue == nil || nilPointerExpr(p.TypesInfo, call.suiteValue) {
		p.Reportf(call.suiteExpr.Pos(), "TESTO012: suite argument is statically nil")
	}
}

// nilPointerExpr accepts the predeclared nil and pointer type conversions of
// it, never function calls or variable values. Discovery resolves a single
// unchanged local initializer before the rule runs.
func nilPointerExpr(info *types.Info, expr ast.Expr) bool {
	switch expr := ast.Unparen(expr).(type) {
	case *ast.Ident:
		return info.ObjectOf(expr) == types.Universe.Lookup("nil")
	case *ast.CallExpr:
		if len(expr.Args) != 1 || expr.Ellipsis.IsValid() || !info.Types[ast.Unparen(expr.Fun)].IsType() {
			return false
		}
		typ := info.TypeOf(expr)
		if typ == nil {
			return false
		}
		if _, ok := typ.Underlying().(*types.Pointer); !ok {
			return false
		}
		return nilPointerExpr(info, expr.Args[0])
	}
	return false
}
