package analyzer

import (
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

const testoImportPath = "github.com/ozontech/testo"

func discoverSuites(p *analysis.Pass) []*Suite {
	byType := make(map[*types.Named]*Suite)

	for _, file := range p.Files {
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv == nil {
				continue
			}

			recv := receiverType(p, fn)
			if recv == nil {
				continue
			}

			suite, ok := byType[recv]
			if !ok {
				suite = &Suite{
					Type:  recv,
					TType: testoTType(recv),
				}
				byType[recv] = suite
			}

			method := methodFromDecl(p, fn)
			if method == nil {
				continue
			}

			switch {
			case isTestMethod(fn):
				suite.Tests = append(suite.Tests, testFromMethod(method))

			case isCasesMethod(fn):
				suite.Cases = append(suite.Cases, method)

			case isHookMethod(fn):
				suite.Hooks = append(suite.Hooks, method)
			}
		}
	}

	result := make([]*Suite, 0, len(byType))

	for _, suite := range byType {
		if !isTestoSuite(suite.Type) {
			continue
		}

		result = append(result, suite)
	}

	return result
}

func testFromMethod(method *Method) *Test {
	test := &Test{
		Method: method,
	}

	sig, ok := method.Func.Type().(*types.Signature)
	if !ok {
		return test
	}

	params := sig.Params()

	if params.Len() < 2 {
		return test
	}

	paramStruct, ok := params.At(1).Type().Underlying().(*types.Struct)
	if !ok {
		return test
	}

	for i := 0; i < paramStruct.NumFields(); i++ {
		field := paramStruct.Field(i)

		test.Param = append(test.Param, Param{
			Name: field.Name(),
			Type: field.Type(),
			Pos:  field.Pos(),
		})
	}

	return test
}

func isTestoSuite(named *types.Named) bool {
	st, ok := named.Underlying().(*types.Struct)
	if !ok {
		return false
	}

	for i := 0; i < st.NumFields(); i++ {
		field := st.Field(i)

		if !field.Embedded() {
			continue
		}

		if isTestoSuiteType(field.Type()) {
			return true
		}
	}

	return false
}

func isTestoSuiteType(t types.Type) bool {
	if ptr, ok := t.(*types.Pointer); ok {
		t = ptr.Elem()
	}

	named, ok := t.(*types.Named)
	if !ok {
		return false
	}

	obj := named.Obj()
	if obj == nil || obj.Pkg() == nil {
		return false
	}

	return obj.Pkg().Path() == testoImportPath &&
		obj.Name() == "Suite"
}

func methodFromDecl(
	p *analysis.Pass,
	fn *ast.FuncDecl,
) *Method {
	obj := p.TypesInfo.ObjectOf(fn.Name)

	f, ok := obj.(*types.Func)
	if !ok {
		return nil
	}

	return &Method{
		Decl: fn,
		Func: f,
	}
}

func receiverType(
	p *analysis.Pass,
	fn *ast.FuncDecl,
) *types.Named {
	if fn.Recv == nil || len(fn.Recv.List) == 0 {
		return nil
	}

	recvExpr := fn.Recv.List[0].Type

	t := p.TypesInfo.TypeOf(recvExpr)
	if t == nil {
		return nil
	}

	if ptr, ok := t.(*types.Pointer); ok {
		t = ptr.Elem()
	}

	named, ok := t.(*types.Named)
	if !ok {
		return nil
	}

	return named
}

func isTestMethod(fn *ast.FuncDecl) bool {
	return fn.Recv != nil &&
		strings.HasPrefix(fn.Name.Name, "Test")
}

func isCasesMethod(fn *ast.FuncDecl) bool {
	return fn.Recv != nil &&
		strings.HasPrefix(fn.Name.Name, "Cases")
}

func isHookMethod(fn *ast.FuncDecl) bool {
	switch fn.Name.Name {
	case "BeforeAll", "BeforeEach", "AfterEach", "AfterAll":
		return true
	default:
		return false
	}
}

func testoTType(named *types.Named) types.Type {
	st, ok := named.Underlying().(*types.Struct)
	if !ok {
		return nil
	}

	for i := 0; i < st.NumFields(); i++ {
		field := st.Field(i)
		if !field.Embedded() {
			continue
		}

		t := field.Type()

		if ptr, ok := t.(*types.Pointer); ok {
			t = ptr.Elem()
		}

		testoSuite, ok := t.(*types.Named)
		if !ok {
			continue
		}

		obj := testoSuite.Obj()
		if obj == nil || obj.Pkg() == nil {
			continue
		}

		if obj.Pkg().Path() != testoImportPath ||
			obj.Name() != "Suite" {
			continue
		}

		args := testoSuite.TypeArgs()
		if args.Len() != 1 {
			return nil
		}

		return args.At(0)
	}

	return nil
}
