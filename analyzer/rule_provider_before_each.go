package analyzer

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

// A fresh zero-initialized suite and default BeforeAll establish the field's
// initial state without guessing about constructors or previous executions.
func checkProviderBeforeEachDependency(p *analysis.Pass, call *runCall) {
	suite := call.model
	if suite == nil || !standardTestoT(suite.tType) || len(call.call.Args) != 2 || len(suite.cases) != 1 || !validParallelProviders(suite) {
		return
	}
	address, ok := ast.Unparen(call.suiteExpr).(*ast.UnaryExpr)
	if !ok || address.Op != token.AND {
		return
	}
	literal, ok := ast.Unparen(address.X).(*ast.CompositeLit)
	if !ok || len(literal.Elts) != 0 || !types.Identical(p.TypesInfo.TypeOf(literal), suite.typ) {
		return
	}
	var before *method
	for _, hook := range suite.hooks {
		sig := hook.fn.Type().(*types.Signature)
		if !validParallelSignature(sig, suite.tType, hook.fn.Name()) {
			return
		}
		switch hook.fn.Name() {
		case "BeforeAll":
			if !isTestoSuiteType(sig.Recv().Type()) || !zeroSafeHookPath(suite, hook.fn.Name()) {
				return
			}
		case "BeforeEach":
			before = hook
		}
	}
	if before == nil || !directPointerMethod(before, suite) || len(before.decl.Body.List) != 1 {
		return
	}
	assignment, ok := before.decl.Body.List[0].(*ast.AssignStmt)
	if !ok || assignment.Tok != token.ASSIGN || len(assignment.Lhs) != 1 || len(assignment.Rhs) != 1 {
		return
	}
	field := receiverPointerField(p.TypesInfo, assignment.Lhs[0], before)
	allocation, ok := ast.Unparen(assignment.Rhs[0]).(*ast.CallExpr)
	if field == nil || !ok || len(allocation.Args) != 1 || !sameObject(p.TypesInfo, allocation.Fun, types.Universe.Lookup("new")) {
		return
	}
	used := make(map[*method]bool)
	for _, test := range suite.tests {
		if !validParallelSignature(test.fn().Type().(*types.Signature), suite.tType, test.fn().Name()) || !validParallelParameters(suite, test) {
			return
		}
		for _, param := range test.params {
			used[suite.casesByName[param.name]] = true
		}
	}
	// Follow declaration/method-set order; map iteration must not order reports.
	for _, provider := range suite.cases {
		if !used[provider] || !directPointerMethod(provider, suite) || len(provider.decl.Body.List) != 1 {
			continue
		}
		ret, ok := provider.decl.Body.List[0].(*ast.ReturnStmt)
		if !ok || len(ret.Results) != 1 {
			continue
		}
		cases, ok := ast.Unparen(ret.Results[0]).(*ast.CompositeLit)
		if !ok || len(cases.Elts) != 1 {
			continue
		}
		read, ok := ast.Unparen(cases.Elts[0]).(*ast.StarExpr)
		if ok && receiverPointerField(p.TypesInfo, read.X, provider) == field {
			p.Reportf(read.Pos(), "TESTO023: provider dereferences nil field %q before BeforeEach initializes it; initialize it before collection, for example in BeforeAll", field.Name())
		}
	}
}

// Invoking an inherited value hook through a nil embedded pointer can panic
// before collection. A fresh literal establishes only value embedding paths.
func zeroSafeHookPath(suite *suite, name string) bool {
	selection := suite.methodSet.Lookup(nil, name)
	if selection == nil {
		return false
	}
	t := types.Type(suite.typ)
	path := selection.Index()
	for _, index := range path[:len(path)-1] {
		st, ok := t.Underlying().(*types.Struct)
		if !ok {
			return false
		}
		t = types.Unalias(st.Field(index).Type())
		if _, pointer := t.Underlying().(*types.Pointer); pointer {
			return false
		}
	}
	return true
}

func directPointerMethod(method *method, suite *suite) bool {
	if method.decl == nil || method.decl.Body == nil {
		return false
	}
	sig := method.fn.Type().(*types.Signature)
	ptr, ok := types.Unalias(sig.Recv().Type()).(*types.Pointer)
	return ok && types.Identical(types.Unalias(ptr.Elem()), suite.typ)
}

func receiverPointerField(info *types.Info, expr ast.Expr, method *method) *types.Var {
	sel, ok := ast.Unparen(expr).(*ast.SelectorExpr)
	if !ok || !sameObject(info, sel.X, method.fn.Type().(*types.Signature).Recv()) {
		return nil
	}
	selection := info.Selections[sel]
	if selection == nil || selection.Kind() != types.FieldVal || len(selection.Index()) != 1 {
		return nil
	}
	field, ok := selection.Obj().(*types.Var)
	if !ok {
		return nil
	}
	if _, ok := types.Unalias(field.Type()).Underlying().(*types.Pointer); !ok {
		return nil
	}
	return field
}
