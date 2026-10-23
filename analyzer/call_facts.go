package analyzer

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

// callFacts keeps package-local initializers and invalidates every reassigned or
// address-taken variable, including captured variables. No flow is inferred.
type callFacts struct {
	values    map[types.Object]ast.Expr
	invalid   map[types.Object]bool
	calls     []*ast.CallExpr
	functions map[*ast.CallExpr]*ast.FuncDecl
	parents   map[ast.Node]ast.Node
}

func collectCallFacts(p *analysis.Pass) *callFacts {
	if p.TypesInfo == nil {
		return &callFacts{}
	}
	localValues := make(map[types.Object]ast.Expr)
	invalid := make(map[types.Object]bool)
	parents := make(map[ast.Node]ast.Node)
	var calls []*ast.CallExpr
	functions := make(map[*ast.CallExpr]*ast.FuncDecl)
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
		ast.PreorderStack(file, nil, func(n ast.Node, stack []ast.Node) bool {
			if len(stack) > 0 {
				parents[n] = stack[len(stack)-1]
			}
			switch n := n.(type) {
			case *ast.CallExpr:
				calls = append(calls, n)
				functions[n] = enclosingFunc(stack)
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
	return &callFacts{localValues, invalid, calls, functions, parents}
}

func (f *callFacts) function(info *types.Info, expr ast.Expr) *types.Func {
	fun := runFunctionExpr(expr)
	if id, ok := fun.(*ast.Ident); ok {
		obj := info.ObjectOf(id)
		if value := f.values[obj]; value != nil && !f.invalid[obj] {
			fun = runFunctionExpr(value)
		}
	}
	var obj types.Object
	switch fun := fun.(type) {
	case *ast.Ident:
		obj = info.ObjectOf(fun)
	case *ast.SelectorExpr:
		obj = info.ObjectOf(fun.Sel)
	}
	fn, _ := obj.(*types.Func)
	return fn
}

func (f *callFacts) value(info *types.Info, expr ast.Expr) ast.Expr {
	if id, ok := ast.Unparen(expr).(*ast.Ident); ok {
		obj := info.ObjectOf(id)
		if value, known := f.values[obj]; known && !f.invalid[obj] {
			return value
		}
	}
	return expr
}
