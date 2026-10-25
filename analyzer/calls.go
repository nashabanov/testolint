package analyzer

import (
	"go/ast"
	"go/token"
	"go/types"
	"slices"

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

// A call inside a function literal belongs to that literal, even when it is
// nested in a method. It cannot be treated as a direct call by the method.
func enclosingFunc(stack []ast.Node) *ast.FuncDecl {
	for _, node := range slices.Backward(stack) {
		switch node := node.(type) {
		case *ast.FuncDecl:
			return node
		case *ast.FuncLit:
			return nil
		}
	}
	return nil
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

func parentWithoutParens(facts *callFacts, node ast.Node) (ast.Node, ast.Node) {
	for {
		parent := facts.parents[node]
		if _, ok := parent.(*ast.ParenExpr); !ok {
			return node, parent
		}
		node = parent
	}
}

func executedAdapter(info *types.Info, facts *callFacts, call *ast.CallExpr) bool {
	node, parent := parentWithoutParens(facts, call)
	outer, ok := parent.(*ast.CallExpr)
	if !ok {
		return false
	}
	if outer.Fun == node {
		return true
	}
	if len(outer.Args) != 2 || outer.Args[1] != node {
		return false
	}
	sel, ok := ast.Unparen(outer.Fun).(*ast.SelectorExpr)
	if !ok {
		return false
	}
	selection := info.Selections[sel]
	if selection == nil || selection.Kind() != types.MethodVal {
		return false
	}
	fn := selection.Obj()
	return fn.Pkg() != nil && fn.Pkg().Path() == "testing" && fn.Name() == "Run"
}

func sameObject(info *types.Info, expr ast.Expr, object types.Object) bool {
	id, ok := ast.Unparen(expr).(*ast.Ident)
	return ok && info.ObjectOf(id) == object
}
