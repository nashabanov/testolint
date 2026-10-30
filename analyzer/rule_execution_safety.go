package analyzer

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

// Execution checks are independent of suite discovery: standalone APIs also
// appear in ordinary testing functions. Names alone never identify an API.
func checkExecutionSafety(p *analysis.Pass, facts *callFacts) {
	if p.TypesInfo == nil {
		return
	}
	for _, call := range facts.calls {
		fn := facts.function(p.TypesInfo, call.Fun)
		if fn == nil || fn.Pkg() == nil || fn.Pkg().Path() != testoImportPath {
			continue
		}
		sig, ok := fn.Type().(*types.Signature)
		if !ok || sig.Recv() != nil {
			continue
		}
		switch fn.Name() {
		case "Run":
			checkRunInCleanup(p, facts, call)
		case "RunTest", "Test":
			index := 1
			if fn.Name() == "Test" {
				if !executedAdapter(p.TypesInfo, facts, call) {
					continue
				}
				index = 0
			}
			if len(call.Args) <= index {
				continue
			}
			arg := call.Args[index]
			value := facts.value(p.TypesInfo, arg)
			if value != nil && !nilCallbackExpr(p.TypesInfo, value) {
				continue
			}
			message := "RunTest callback is statically nil; supply a non-nil function"
			if index == 0 {
				message = "Test callback is statically nil; executing the adapter invokes it; supply a non-nil function"
			}
			reportf(p, "TESTO019", arg.Pos(), "%s", message)
		}
	}
}

func nilCallbackExpr(info *types.Info, expr ast.Expr) bool {
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
		_, function := typ.Underlying().(*types.Signature)
		return function && nilCallbackExpr(info, expr.Args[0])
	}
	return false
}

func checkRunInCleanup(p *analysis.Pass, facts *callFacts, call *ast.CallExpr) {
	if len(call.Args) < 3 {
		return
	}
	// Only direct, synchronous expression statements in a direct cleanup
	// literal qualify. Deferred/go calls and nested function bodies stay unknown.
	node, parent := parentWithoutParens(facts, call)
	stmt, ok := parent.(*ast.ExprStmt)
	if !ok || stmt.X != node {
		return
	}
	block, ok := facts.parents[stmt].(*ast.BlockStmt)
	if !ok {
		return
	}
	lit, ok := facts.parents[block].(*ast.FuncLit)
	if !ok {
		return
	}
	node, parent = parentWithoutParens(facts, lit)
	registration, ok := parent.(*ast.CallExpr)
	if !ok || len(registration.Args) != 1 || registration.Args[0] != node {
		return
	}
	sel, ok := ast.Unparen(registration.Fun).(*ast.SelectorExpr)
	if !ok || !standardTestoT(p.TypesInfo.TypeOf(sel.X)) {
		return
	}
	selection := p.TypesInfo.Selections[sel]
	if selection == nil || selection.Kind() != types.MethodVal || selection.Obj().Pkg() == nil ||
		selection.Obj().Pkg().Path() != testoImportPath || selection.Obj().Name() != "Cleanup" {
		return
	}
	context, ok := ast.Unparen(sel.X).(*ast.Ident)
	runContext, runOK := ast.Unparen(call.Args[0]).(*ast.Ident)
	if !ok || !runOK {
		return
	}
	obj, ok := p.TypesInfo.ObjectOf(context).(*types.Var)
	if !ok || obj.Parent() == p.Pkg.Scope() || facts.invalid[obj] || p.TypesInfo.ObjectOf(runContext) != obj {
		return
	}
	// Branches, returns and other control-flow constructs before the call are
	// outside the supported straight-line prefix.
	for _, preceding := range block.List {
		if preceding == stmt {
			reportf(p, "TESTO018", call.Pos(),
				"testo.Run on this context during Cleanup panics; run subtests before cleanup")
			return
		}
		switch preceding.(type) {
		case *ast.ExprStmt, *ast.AssignStmt, *ast.DeclStmt, *ast.EmptyStmt:
		default:
			return
		}
	}
}
