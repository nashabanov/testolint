package analyzer

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

// checkPrematureTeardown deliberately starts with a concrete resource contract:
// os.File.Close makes subsequent Read, Write and Stat fail. A single regular
// test, default surrounding hooks and complete tiny bodies establish identity
// and ordering without assuming what arbitrary Close methods or helpers do.
func checkPrematureTeardown(p *analysis.Pass, suite *suite) {
	if !standardTestoT(suite.tType) || len(suite.tests) != 1 || len(suite.cases) != 0 {
		return
	}
	var after *method
	for _, hook := range suite.hooks {
		if !validParallelSignature(hook.fn.Type().(*types.Signature), suite.tType, hook.fn.Name()) {
			return
		}
		if hook.fn.Name() == "AfterEach" {
			after = hook
		} else if hook.fn.Pkg().Path() != testoImportPath {
			return
		}
	}
	test := suite.tests[0].method
	if after == nil || !isValidPrefixedName(test.fn.Name(), "Test") ||
		!teardownMethod(suite, test) || !teardownMethod(suite, after) {
		return
	}
	closeCall := discardedCall(after.decl.Body.List[0])
	field := fileOperationField(p.TypesInfo, closeCall, after, "Close")
	if field == nil || len(closeCall.Args) != 0 {
		return
	}
	run := discardedCall(test.decl.Body.List[0])
	if run == nil || len(run.Args) != 3 || run.Ellipsis.IsValid() {
		return
	}
	// Resolve only the direct generic Testo function, not local function values.
	fn := (&callFacts{}).function(p.TypesInfo, run.Fun)
	if fn == nil || fn.Pkg() == nil || fn.Pkg().Path() != testoImportPath || fn.Name() != "Run" ||
		fn.Type().(*types.Signature).Recv() != nil {
		return
	}
	if !sameObject(p.TypesInfo, run.Args[0], test.fn.Type().(*types.Signature).Params().At(0)) ||
		p.TypesInfo.Types[run.Args[1]].Value == nil {
		return
	}
	child, ok := ast.Unparen(run.Args[2]).(*ast.FuncLit)
	if !ok || len(child.Body.List) != 2 {
		return
	}
	childSig, ok := p.TypesInfo.TypeOf(child).(*types.Signature)
	if !ok || childSig.Params().Len() != 1 || !standardTestoT(childSig.Params().At(0).Type()) {
		return
	}
	parallel := discardedCall(child.Body.List[0])
	if parallel == nil || len(parallel.Args) != 0 {
		return
	}
	sel, ok := ast.Unparen(parallel.Fun).(*ast.SelectorExpr)
	if !ok || !sameObject(p.TypesInfo, sel.X, childSig.Params().At(0)) {
		return
	}
	selection := p.TypesInfo.Selections[sel]
	if selection == nil || selection.Kind() != types.MethodVal || selection.Obj().Pkg() == nil ||
		selection.Obj().Pkg().Path() != testoImportPath || selection.Obj().Name() != "Parallel" {
		return
	}
	use := discardedCall(child.Body.List[1])
	if fileOperationField(p.TypesInfo, use, test, "Read", "Write", "Stat") != field {
		return
	}
	// Arguments must not mutate or replace the field before the operation.
	for _, arg := range use.Args {
		if !literalFileArgument(p.TypesInfo, arg) {
			return
		}
	}
	reportf(p, "TESTO020", closeCall.Pos(),
		"AfterEach closes file %q before its parallel subtest uses it; close the file in Cleanup", field.Name())
}

func teardownMethod(suite *suite, method *method) bool {
	if method.decl == nil || method.decl.Body == nil || len(method.decl.Body.List) != 1 {
		return false
	}
	sig := method.fn.Type().(*types.Signature)
	if sig.Recv() == nil || sig.Params().Len() != 1 || !validParallelSignature(sig, suite.tType, method.fn.Name()) {
		return false
	}
	receiver := sig.Recv().Type()
	if ptr, ok := types.Unalias(receiver).(*types.Pointer); ok {
		receiver = ptr.Elem()
	}
	return types.Identical(receiver, suite.typ)
}

// Only synchronous calls whose results are discarded have no intervening flow
// or assignments to state. In particular, defer/go and named results are unknown.
func discardedCall(stmt ast.Stmt) *ast.CallExpr {
	var expr ast.Expr
	switch stmt := stmt.(type) {
	case *ast.ExprStmt:
		expr = stmt.X
	case *ast.AssignStmt:
		if stmt.Tok != token.ASSIGN || len(stmt.Rhs) != 1 {
			return nil
		}
		for _, lhs := range stmt.Lhs {
			id, ok := lhs.(*ast.Ident)
			if !ok || id.Name != "_" {
				return nil
			}
		}
		expr = stmt.Rhs[0]
	default:
		return nil
	}
	call, _ := ast.Unparen(expr).(*ast.CallExpr)
	return call
}

func fileOperationField(info *types.Info, call *ast.CallExpr, method *method, names ...string) *types.Var {
	if call == nil {
		return nil
	}
	sel, ok := ast.Unparen(call.Fun).(*ast.SelectorExpr)
	if !ok {
		return nil
	}
	selection := info.Selections[sel]
	if selection == nil || selection.Kind() != types.MethodVal {
		return nil
	}
	fn := selection.Obj()
	if fn.Pkg() == nil || fn.Pkg().Path() != "os" {
		return nil
	}
	matched := false
	for _, name := range names {
		matched = matched || fn.Name() == name
	}
	if !matched {
		return nil
	}
	field, ok := ast.Unparen(sel.X).(*ast.SelectorExpr)
	if !ok || !sameObject(info, field.X, method.fn.Type().(*types.Signature).Recv()) {
		return nil
	}
	fs := info.Selections[field]
	if fs == nil || fs.Kind() != types.FieldVal || len(fs.Index()) != 1 {
		return nil
	}
	ptr, ok := types.Unalias(info.TypeOf(field)).(*types.Pointer)
	if !ok {
		return nil
	}
	named, ok := types.Unalias(ptr.Elem()).(*types.Named)
	if !ok || named.Obj().Pkg() == nil || named.Obj().Pkg().Path() != "os" || named.Obj().Name() != "File" {
		return nil
	}
	object, _ := fs.Obj().(*types.Var)
	return object
}

func literalFileArgument(info *types.Info, expr ast.Expr) bool {
	if info.Types[expr].Value != nil {
		return true
	}
	if id, ok := ast.Unparen(expr).(*ast.Ident); ok {
		return info.ObjectOf(id) == types.Universe.Lookup("nil")
	}
	lit, ok := ast.Unparen(expr).(*ast.CompositeLit)
	if !ok {
		return false
	}
	for _, element := range lit.Elts {
		if info.Types[element].Value == nil {
			return false
		}
	}
	return true
}
