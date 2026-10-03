package analyzer

import (
	"go/ast"
	"go/token"
	"go/types"
)

type Suite struct {
	Type  *types.Named
	TType types.Type

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
