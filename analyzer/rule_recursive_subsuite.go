package analyzer

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

func checkRecursiveSubSuite(p *analysis.Pass, call *RunCall) {
	if call.Kind != RunSubSuite || call.SuiteType == nil || p.TypesInfo == nil || len(call.Call.Args) != 2 {
		return
	}
	decl := call.EnclosingFunc
	if decl == nil || decl.Recv == nil || decl.Body == nil || len(decl.Body.List) != 1 {
		return
	}
	stmt, ok := decl.Body.List[0].(*ast.ExprStmt)
	if !ok || ast.Unparen(stmt.X) != call.Call {
		return
	}
	fn, ok := p.TypesInfo.ObjectOf(decl.Name).(*types.Func)
	if !ok {
		return
	}
	sig, ok := fn.Type().(*types.Signature)
	if !ok || sig.Recv() == nil || sig.Params().Len() != 1 || sig.Results().Len() != 0 || sig.Variadic() {
		return
	}
	if fn.Name() != "BeforeAll" && !isValidPrefixedName(fn.Name(), "Test") {
		return
	}

	// The parent must be this method's Testo parameter. A receiver or parent
	// obtained from a function, field, or local alias needs additional reasoning.
	parent, ok := ast.Unparen(call.Call.Args[0]).(*ast.Ident)
	if !ok || p.TypesInfo.ObjectOf(parent) != sig.Params().At(0) {
		return
	}
	receiver := ast.Unparen(call.SuiteExpr)
	if unary, ok := receiver.(*ast.UnaryExpr); ok && (unary.Op == token.AND || unary.Op == token.MUL) {
		receiver = ast.Unparen(unary.X)
	}
	id, ok := receiver.(*ast.Ident)
	if !ok || p.TypesInfo.ObjectOf(id) != sig.Recv() {
		return
	}

	// Use the argument's actual method set, not the pointer method set used by
	// suite-level rules. Passing a value can remove the current pointer method.
	methodSet := types.NewMethodSet(call.SuiteType)
	selected := methodSet.Lookup(p.Pkg, fn.Name())
	if selected == nil {
		return
	}
	target, ok := selected.Obj().(*types.Func)
	if !ok || target.Origin() != fn.Origin() {
		return
	}
	tType := testoTType(methodSet)
	if tType == nil || !types.Identical(sig.Params().At(0).Type(), tType) {
		return
	}
	if !directRecurrenceSuite(methodSet, fn, tType) {
		return
	}
	p.Reportf(call.Call.Pos(), "TESTO013: recursive RunSubSuite invocation")
}

// directRecurrenceSuite rejects collection failures and execution paths that
// need provider, hook-body, or test-order analysis. BeforeAll repeats before
// other suite methods. A Test repeats only when it is the sole regular test
// and the two hooks that run before it are Testo's inherited defaults.
func directRecurrenceSuite(methodSet *types.MethodSet, current *types.Func, tType types.Type) bool {
	tests := 0
	for selection := range methodSet.Methods() {
		fn, ok := selection.Obj().(*types.Func)
		if !ok {
			return false
		}
		name := fn.Name()
		if strings.HasPrefix(name, "Cases") && name != "Cases" {
			return false
		}
		if !strings.HasPrefix(name, "Test") {
			continue
		}
		sig, ok := fn.Type().(*types.Signature)
		if !isValidPrefixedName(name, "Test") || !ok || sig.Variadic() || sig.Params().Len() != 1 || sig.Results().Len() != 0 || !types.Identical(sig.Params().At(0).Type(), tType) {
			return false
		}
		tests++
	}
	if tests == 0 {
		return false
	}
	if current.Name() == "BeforeAll" {
		return true
	}
	if tests != 1 {
		return false
	}
	for _, name := range []string{"BeforeAll", "BeforeEach"} {
		selection := methodSet.Lookup(nil, name)
		if selection == nil {
			return false
		}
		fn, ok := selection.Obj().(*types.Func)
		if !ok || fn.Pkg() == nil || fn.Pkg().Path() != testoImportPath {
			return false
		}
		sig, ok := fn.Type().(*types.Signature)
		if !ok || sig.Recv() == nil || !isTestoSuiteType(sig.Recv().Type()) {
			return false
		}
	}
	return true
}
