package analyzer

import (
	"go/ast"
	"go/token"
	"go/types"
)

type suite struct {
	typ       *types.Named
	tType     types.Type
	methodSet *types.MethodSet

	tests []*test
	cases []*method
	hooks []*method

	casesByName map[string]*method
}

type method struct {
	fn *types.Func
	// decl is optional and available only for declarations in this pass.
	// fn remains the semantic identity, including instantiated methods.
	decl *ast.FuncDecl
}

type param struct {
	name string
	typ  types.Type
	pos  token.Pos
}

type test struct {
	method *method
	params []param
}

// pos reports local methods at their declaration, and imported methods at
// the local suite that exposes them. Dependency files are outside this pass.
func (m *method) pos(suite *suite) token.Pos {
	if m.fn.Pkg() != suite.typ.Obj().Pkg() || !m.fn.Pos().IsValid() {
		return suite.typ.Obj().Pos()
	}
	return m.fn.Pos()
}

func (t *test) fn() *types.Func {
	return t.method.fn
}

// runKind identifies a Testo suite execution entry point.
type runKind uint8

const (
	runSuite runKind = iota
	runSubSuite
)

// runCall retains source expressions and the suite's static argument type.
// suiteType may be nil when type information is incomplete.
type runCall struct {
	// model is linked to an exact local non-generic suite by the rule pipeline.
	model     *suite
	kind      runKind
	call      *ast.CallExpr
	suiteExpr ast.Expr
	suiteType types.Type
	// enclosingFunc is nil for calls in function literals or package initializers.
	enclosingFunc *ast.FuncDecl
	// suiteValue is suiteExpr, or the initializer of an unchanged local
	// variable used directly as the argument. Nil denotes its zero value.
	// Initializers are resolved only one level; their type may differ from
	// suiteType because of implicit conversions in the declaration.
	suiteValue ast.Expr
}
