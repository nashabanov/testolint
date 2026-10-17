package analyzer

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

// This parser accepts only straight-line local code. An unsupported statement
// invalidates the whole body; partially collected events are never returned.
type parallelBodyParser struct {
	info           *types.Info
	receiver       *types.Var
	context        *types.Var
	sharedReceiver bool
	events         []parallelEvent
}

func parseParallelBody(p *analysis.Pass, suite *Suite, method *Method) ([]parallelEvent, bool) {
	sig := method.Func.Type().(*types.Signature)
	if !validParallelSignature(sig, suite.TType, method.Func.Name()) {
		return nil, false
	}
	if method.Func.Pkg().Path() == testoImportPath {
		return nil, true // Testo's inherited hooks have no relevant events.
	}
	if method.Decl == nil || method.Decl.Body == nil || sig.Recv() == nil {
		return nil, false
	}
	ptr, ok := types.Unalias(sig.Recv().Type()).(*types.Pointer)
	shared := ok && types.Identical(ptr.Elem(), suite.Type)
	if !shared && !types.Identical(sig.Recv().Type(), suite.Type) {
		return nil, false // A promoted body operates on a different receiver.
	}
	parser := parallelBodyParser{
		info: p.TypesInfo, receiver: sig.Recv(), context: sig.Params().At(0),
		sharedReceiver: shared,
	}
	if !parser.statements(method.Decl.Body.List, false) {
		return nil, false
	}
	return parser.events, true
}

func validParallelSignature(sig *types.Signature, context types.Type, name string) bool {
	if sig.Variadic() || sig.Results().Len() != 0 {
		return false
	}
	switch sig.Params().Len() {
	case 1:
	case 2:
		if !isValidPrefixedName(name, "Test") || !isStructType(sig.Params().At(1).Type()) {
			return false
		}
	default:
		return false
	}
	return types.Identical(sig.Params().At(0).Type(), context)
}

func (b *parallelBodyParser) isObject(expr ast.Expr, object types.Object) bool {
	id, ok := ast.Unparen(expr).(*ast.Ident)
	return ok && b.info.ObjectOf(id) == object
}

func (b *parallelBodyParser) field(expr ast.Expr) *types.Var {
	sel, ok := ast.Unparen(expr).(*ast.SelectorExpr)
	if !ok || !b.sharedReceiver || !b.isObject(sel.X, b.receiver) {
		return nil
	}
	selection := b.info.Selections[sel]
	if selection == nil || selection.Kind() != types.FieldVal || len(selection.Index()) != 1 {
		return nil
	}
	field, _ := selection.Obj().(*types.Var)
	return field
}

func (b *parallelBodyParser) statements(statements []ast.Stmt, inCleanup bool) bool {
	for i, statement := range statements {
		switch statement := statement.(type) {
		case *ast.AssignStmt:
			if !b.assignment(statement, inCleanup) {
				return false
			}
		case *ast.ExprStmt:
			if !b.callStatement(statement, inCleanup, i == len(statements)-1) {
				return false
			}
		default:
			return false
		}
	}
	return true
}

func (b *parallelBodyParser) assignment(statement *ast.AssignStmt, inCleanup bool) bool {
	if len(statement.Lhs) != 1 || len(statement.Rhs) != 1 {
		return false
	}
	lhs, rhs := statement.Lhs[0], statement.Rhs[0]
	if field := b.field(lhs); field != nil {
		if inCleanup || statement.Tok != token.ASSIGN {
			return false
		}
		storesContext := b.isObject(rhs, b.context)
		if !storesContext && !b.expression(rhs, inCleanup) {
			return false
		}
		b.events = append(b.events, parallelEvent{
			kind: parallelFieldWrite, pos: statement.Pos(), field: field,
			storesContext: storesContext,
		})
		return true
	}
	id, ok := lhs.(*ast.Ident)
	if !ok || b.isObject(id, b.receiver) || b.isObject(id, b.context) {
		return false
	}
	return (statement.Tok == token.DEFINE || id.Name == "_") && b.expression(rhs, inCleanup)
}

func (b *parallelBodyParser) callStatement(statement *ast.ExprStmt, inCleanup, last bool) bool {
	call, ok := statement.X.(*ast.CallExpr)
	if !ok {
		return false
	}
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if ok && inCleanup && last && b.field(sel.X) != nil {
		return b.cleanupResourceCall(call, sel)
	}
	if ok && b.isObject(sel.X, b.context) && (sel.Sel.Name == "Parallel" || sel.Sel.Name == "Cleanup") {
		return !inCleanup && b.contextCall(call, sel)
	}
	return b.expression(call, inCleanup)
}

func (b *parallelBodyParser) contextCall(call *ast.CallExpr, sel *ast.SelectorExpr) bool {
	if !b.isTestoMethod(sel) {
		return false
	}
	if sel.Sel.Name == "Parallel" {
		if len(call.Args) != 0 {
			return false
		}
		b.events = append(b.events, parallelEvent{kind: parallelCall, pos: call.Pos()})
		return true
	}
	if len(call.Args) != 1 {
		return false
	}
	callback, ok := call.Args[0].(*ast.FuncLit)
	return ok && b.statements(callback.Body.List, true)
}

// A terminal Cleanup resource call evaluates its receiver before invoking the
// unknown method. No later event relies on the method returning or synchronizing.
func (b *parallelBodyParser) cleanupResourceCall(call *ast.CallExpr, sel *ast.SelectorExpr) bool {
	selection := b.info.Selections[sel]
	if selection == nil || selection.Kind() != types.MethodVal {
		return false
	}
	if _, ok := types.Unalias(b.field(sel.X).Type()).(*types.Pointer); !ok {
		return false
	}
	return b.expression(sel.X, true) && b.arguments(call.Args, true)
}

func (b *parallelBodyParser) expression(expr ast.Expr, inCleanup bool) bool {
	if field := b.field(expr); field != nil {
		b.events = append(b.events, parallelEvent{
			kind: parallelFieldRead, pos: expr.Pos(), field: field, inCleanup: inCleanup,
		})
		return true
	}
	switch expr := ast.Unparen(expr).(type) {
	case *ast.BasicLit:
		return true
	case *ast.Ident:
		return !b.isObject(expr, b.receiver) && !b.isObject(expr, b.context)
	case *ast.CallExpr:
		if id, ok := expr.Fun.(*ast.Ident); ok {
			return b.allocation(expr, id)
		}
		return b.logCall(expr, inCleanup)
	default:
		return false
	}
}

func (b *parallelBodyParser) allocation(call *ast.CallExpr, name *ast.Ident) bool {
	builtin, ok := b.info.ObjectOf(name).(*types.Builtin)
	if !ok || builtin.Name() != "new" || len(call.Args) != 1 {
		return false
	}
	// new(expr) can execute arbitrary code; accept only types or constants.
	arg := b.info.Types[call.Args[0]]
	return arg.IsType() || arg.Value != nil
}

func (b *parallelBodyParser) logCall(call *ast.CallExpr, inCleanup bool) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != "Log" || !b.isTestoMethod(sel) {
		return false
	}
	if !b.isObject(sel.X, b.context) {
		if b.field(sel.X) == nil || !standardTestoT(b.info.TypeOf(sel.X)) || !b.expression(sel.X, inCleanup) {
			return false
		}
	}
	return b.arguments(call.Args, inCleanup)
}

func (b *parallelBodyParser) isTestoMethod(sel *ast.SelectorExpr) bool {
	selection := b.info.Selections[sel]
	if selection == nil {
		return false
	}
	method, ok := selection.Obj().(*types.Func)
	return ok && method.Pkg() != nil && method.Pkg().Path() == testoImportPath
}

func (b *parallelBodyParser) arguments(args []ast.Expr, inCleanup bool) bool {
	for _, arg := range args {
		if !b.expression(arg, inCleanup) {
			return false
		}
	}
	return true
}
