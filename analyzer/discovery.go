package analyzer

import (
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

const testoImportPath = "github.com/ozontech/testo"

func discoverSuites(p *analysis.Pass) []*Suite {
	var result []*Suite
	decls := make(map[*types.Func]*ast.FuncDecl)
	for _, file := range p.Files {
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv == nil {
				continue
			}
			if f, ok := p.TypesInfo.ObjectOf(fn.Name).(*types.Func); ok {
				decls[f.Origin()] = fn
			}
		}
	}

	// Discover package-level named types before examining any methods.
	for _, file := range p.Files {
		for _, decl := range file.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}
			for _, spec := range gen.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok || typeSpec.Assign.IsValid() {
					continue
				}
				obj, ok := p.TypesInfo.ObjectOf(typeSpec.Name).(*types.TypeName)
				if !ok {
					continue
				}
				named, ok := obj.Type().(*types.Named)
				if !ok || named.TypeParams().Len() != 0 {
					continue
				}
				if _, ok := named.Underlying().(*types.Struct); !ok {
					continue
				}
				methodSet := types.NewMethodSet(types.NewPointer(named))
				tType := testoTType(methodSet)
				if tType == nil {
					continue
				}
				suite := &Suite{
					Type:        named,
					TType:       tType,
					MethodSet:   methodSet,
					CasesByName: make(map[string]*Method),
				}
				result = append(result, suite)
			}
		}
	}

	// The method set contains exactly the selectable declared and promoted
	// methods, with instantiated signatures and Go's shadowing rules applied.
	for _, suite := range result {
		for selection := range suite.MethodSet.Methods() {
			f := selection.Obj().(*types.Func)
			method := &Method{Func: f, Decl: decls[f.Origin()]}
			switch name := f.Name(); {
			case strings.HasPrefix(name, "Test"):
				suite.Tests = append(suite.Tests, testFromMethod(method, suite))
			case strings.HasPrefix(name, "Cases") && name != "Cases":
				suite.Cases = append(suite.Cases, method)
				suite.CasesByName[casesName(method)] = method
			case isHookMethod(name):
				suite.Hooks = append(suite.Hooks, method)
			}
		}
	}

	return result
}

func testFromMethod(method *Method, suite *Suite) *Test {
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

	for field := range paramStruct.Fields() {
		pos := field.Pos()
		if !pos.IsValid() || field.Pkg() != suite.Type.Obj().Pkg() {
			pos = method.Pos(suite)
		}

		test.Params = append(test.Params, Param{
			Name: field.Name(),
			Type: field.Type(),
			Pos:  pos,
		})
	}

	return test
}

func isTestoSuiteType(t types.Type) bool {
	t = types.Unalias(t)
	if ptr, ok := t.(*types.Pointer); ok {
		t = types.Unalias(ptr.Elem())
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

func isHookMethod(name string) bool {
	switch name {
	case "BeforeAll", "BeforeEach", "AfterEach", "AfterAll":
		return true
	default:
		return false
	}
}

// testoTType uses Testo's private suite-interface marker rather than
// reconstructing Go's embedding and promotion rules. The selected method's
// receiver is the instantiated testo.Suite[T], even for indirect embedding.
func testoTType(methodSet *types.MethodSet) types.Type {
	for selection := range methodSet.Methods() {
		f := selection.Obj().(*types.Func)
		if f.Name() != "private" || f.Pkg() == nil || f.Pkg().Path() != testoImportPath {
			continue
		}
		sig := f.Type().(*types.Signature)
		if !isTestoSuiteType(sig.Recv().Type()) {
			continue
		}
		t := types.Unalias(sig.Recv().Type())
		if ptr, ok := t.(*types.Pointer); ok {
			t = types.Unalias(ptr.Elem())
		}
		args := t.(*types.Named).TypeArgs()
		if args.Len() == 1 {
			return args.At(0)
		}
	}
	return nil
}

func casesName(method *Method) string {
	return strings.TrimPrefix(
		method.Func.Name(),
		"Cases",
	)
}
