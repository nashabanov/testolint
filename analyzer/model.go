package analyzer

import (
	"go/ast"
	"go/types"
)

type Method struct {
	Decl *ast.FuncDecl
	Func *types.Func
}

type Suite struct {
	Type  *types.Named
	TType types.Type

	Tests []*Method
	Cases []*Method
	Hooks []*Method
}
