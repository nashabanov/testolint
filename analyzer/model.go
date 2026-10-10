package analyzer

import (
	"go/ast"
	"go/token"
	"go/types"
)

type Suite struct {
	Type      *types.Named
	TType     types.Type
	MethodSet *types.MethodSet

	Tests []*Test
	Cases []*Method
	Hooks []*Method

	CasesByName map[string]*Method
}

type Method struct {
	Func *types.Func
	// Decl is optional and available only for declarations in this pass.
	// Func remains the semantic identity, including instantiated methods.
	Decl *ast.FuncDecl
}

type Param struct {
	Name string
	Type types.Type
	Pos  token.Pos
}

type Test struct {
	Method *Method
	Params []Param
}

// Pos reports local methods at their declaration, and imported methods at
// the local suite that exposes them. Dependency files are outside this pass.
func (m *Method) Pos(suite *Suite) token.Pos {
	if m.Func.Pkg() != suite.Type.Obj().Pkg() || !m.Func.Pos().IsValid() {
		return suite.Type.Obj().Pos()
	}
	return m.Func.Pos()
}

func (t *Test) Func() *types.Func {
	return t.Method.Func
}

// validParams reports whether parameter checks have a meaningful runtime target.
func (t *Test) validParams(suite *Suite) bool {
	sig := t.Func().Type().(*types.Signature)
	return isValidPrefixedName(t.Func().Name(), "Test") &&
		sig.Params().Len() == 2 && !sig.Variadic() && sig.Results().Len() == 0 &&
		types.Identical(sig.Params().At(0).Type(), suite.TType) &&
		isStructType(sig.Params().At(1).Type())
}

// RunKind identifies a Testo suite execution entry point.
type RunKind uint8

const (
	RunSuite RunKind = iota
	RunSubSuite
)

// RunCall retains source expressions and the suite's static argument type.
// SuiteType may be nil when type information is incomplete.
type RunCall struct {
	Kind      RunKind
	Call      *ast.CallExpr
	SuiteExpr ast.Expr
	SuiteType types.Type
}
