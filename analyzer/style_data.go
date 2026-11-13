package analyzer

import (
	"go/ast"
	"go/token"
)

// A literal request passed directly to a configured action is proven input
// preparation. Other assignments remain neutral: their purpose is unknown.
func (w *operationWalker) requestData(block *ast.BlockStmt) map[ast.Stmt]bool {
	result := map[ast.Stmt]bool{}
	for i, stmt := range block.List {
		assign, ok := stmt.(*ast.AssignStmt)
		if !ok || assign.Tok != token.DEFINE || len(assign.Lhs) != 1 || len(assign.Rhs) != 1 {
			continue
		}
		id, ok := assign.Lhs[0].(*ast.Ident)
		if !ok {
			continue
		}
		obj := w.pass.TypesInfo.Defs[id]
		if obj == nil || w.facts.invalid[obj] || !literalData(assign.Rhs[0]) {
			continue
		}
		for _, later := range block.List[i+1:] {
			call := statementCall(later)
			if call == nil {
				continue
			}
			fn := w.facts.function(w.pass.TypesInfo, call.Fun)
			if w.kinds[symbolOf(fn)] != "TESTO029" {
				continue
			}
			for _, arg := range call.Args {
				if sameObject(w.pass.TypesInfo, arg, obj) {
					result[stmt] = true
				}
			}
		}
	}
	return result
}

func literalData(expr ast.Expr) bool {
	switch expr := ast.Unparen(expr).(type) {
	case *ast.BasicLit:
		return true
	case *ast.UnaryExpr:
		return (expr.Op == token.AND || expr.Op == token.SUB || expr.Op == token.ADD) && literalData(expr.X)
	case *ast.CompositeLit:
		for _, element := range expr.Elts {
			if pair, ok := element.(*ast.KeyValueExpr); ok {
				element = pair.Value
			}
			if !literalData(element) {
				return false
			}
		}
		return true
	}
	return false
}
