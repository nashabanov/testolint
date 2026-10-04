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
	Decl *ast.FuncDecl
	Func *types.Func
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

func (t *Test) Decl() *ast.FuncDecl {
	return t.Method.Decl
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

// promotedTests makes orphan detection inconclusive: discovery only collects
// declared tests, while Testo also executes promoted methods.
func (s *Suite) promotedTests() bool {
	for i := 0; i < s.MethodSet.Len(); i++ {
		method := s.MethodSet.At(i)
		if len(method.Index()) > 1 && isValidPrefixedName(method.Obj().Name(), "Test") {
			return true
		}
	}
	return false
}
