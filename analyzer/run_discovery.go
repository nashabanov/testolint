package analyzer

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

// discoverRuns identifies calls independently of suite declarations and rules.
// Local function aliases are accepted only with a direct initializer and no
// reassignment or address-taking anywhere in the package. This deliberately
// avoids making control-flow or execution-order assumptions.
func discoverRuns(p *analysis.Pass) []*RunCall {
	if p.TypesInfo == nil {
		return nil
	}
	localValues := make(map[types.Object]ast.Expr)
	invalid := make(map[types.Object]bool)
	var calls []*ast.CallExpr
	record := func(id *ast.Ident, expr ast.Expr) {
		obj, ok := p.TypesInfo.Defs[id].(*types.Var)
		if !ok || obj.Parent() == nil || p.Pkg == nil || obj.Parent() == p.Pkg.Scope() {
			return
		}
		localValues[obj] = expr
	}
	invalidate := func(expr ast.Expr) {
		if id, ok := ast.Unparen(expr).(*ast.Ident); ok {
			if obj := p.TypesInfo.ObjectOf(id); obj != nil {
				invalid[obj] = true
			}
		}
	}
	for _, file := range p.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			switch n := n.(type) {
			case *ast.CallExpr:
				calls = append(calls, n)
			case *ast.ValueSpec:
				if len(n.Values) == 0 {
					for _, id := range n.Names {
						record(id, nil)
					}
				} else if len(n.Names) == len(n.Values) {
					for i, id := range n.Names {
						record(id, n.Values[i])
					}
				}
			case *ast.AssignStmt:
				for i, lhs := range n.Lhs {
					id, ok := lhs.(*ast.Ident)
					if ok && n.Tok == token.DEFINE && p.TypesInfo.Defs[id] != nil {
						if len(n.Lhs) == len(n.Rhs) {
							record(id, n.Rhs[i])
						}
					} else {
						invalidate(lhs)
					}
				}
			case *ast.RangeStmt:
				if n.Tok == token.ASSIGN {
					if n.Key != nil {
						invalidate(n.Key)
					}
					if n.Value != nil {
						invalidate(n.Value)
					}
				}
			case *ast.UnaryExpr:
				if n.Op == token.AND {
					invalidate(n.X)
				}
			}
			return true
		})
	}
	var result []*RunCall
	for _, call := range calls {
		fun := runFunctionExpr(call.Fun)
		if id, ok := fun.(*ast.Ident); ok {
			obj := p.TypesInfo.ObjectOf(id)
			if expr := localValues[obj]; expr != nil && !invalid[obj] {
				fun = runFunctionExpr(expr)
			}
		}
		var obj types.Object
		switch fun := fun.(type) {
		case *ast.Ident:
			obj = p.TypesInfo.ObjectOf(fun)
		case *ast.SelectorExpr:
			obj = p.TypesInfo.ObjectOf(fun.Sel)
		}
		fn, ok := obj.(*types.Func)
		if !ok || fn.Pkg() == nil || fn.Pkg().Path() != testoImportPath {
			continue
		}
		sig, ok := fn.Type().(*types.Signature)
		if !ok || sig.Recv() != nil {
			continue
		}
		var kind RunKind
		switch fn.Name() {
		case "RunSuite":
			kind = RunSuite
		case "RunSubSuite":
			kind = RunSubSuite
		default:
			continue
		}
		if len(call.Args) < 2 {
			continue
		}
		suiteExpr := call.Args[1]
		suiteValue := suiteExpr
		if id, ok := ast.Unparen(suiteExpr).(*ast.Ident); ok {
			obj := p.TypesInfo.ObjectOf(id)
			if init, known := localValues[obj]; known && !invalid[obj] {
				suiteValue = init
			}
		}
		result = append(result, &RunCall{
			Kind: kind, Call: call, SuiteExpr: suiteExpr,
			SuiteType: p.TypesInfo.TypeOf(suiteExpr), SuiteValue: suiteValue,
		})
	}
	return result
}

func runFunctionExpr(expr ast.Expr) ast.Expr {
	for {
		switch e := ast.Unparen(expr).(type) {
		case *ast.IndexExpr:
			expr = e.X
		case *ast.IndexListExpr:
			expr = e.X
		default:
			return e
		}
	}
}
