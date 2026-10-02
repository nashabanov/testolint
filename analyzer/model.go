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
	Type *types.Named

	Tests []*Method
	Cases []*Method
	Hooks []*Method
}
