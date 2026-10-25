package analyzer

import (
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

func discoverSuites(p *analysis.Pass) []*suite {
	var result []*suite
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
				suite := &suite{
					typ:         named,
					tType:       tType,
					methodSet:   methodSet,
					casesByName: make(map[string]*method),
				}
				result = append(result, suite)
			}
		}
	}

	// The method set contains exactly the selectable declared and promoted
	// methods, with instantiated signatures and Go's shadowing rules applied.
	for _, suite := range result {
		for selection := range suite.methodSet.Methods() {
			f := selection.Obj().(*types.Func)
			method := &method{fn: f, decl: decls[f.Origin()]}
			switch name := f.Name(); {
			case strings.HasPrefix(name, "Test"):
				suite.tests = append(suite.tests, testFromMethod(method, suite))
			case strings.HasPrefix(name, "Cases") && name != "Cases":
				suite.cases = append(suite.cases, method)
				suite.casesByName[casesName(method)] = method
			case isHookMethod(name):
				suite.hooks = append(suite.hooks, method)
			}
		}
	}

	return result
}

func testFromMethod(method *method, suite *suite) *test {
	test := &test{
		method: method,
	}

	sig, ok := method.fn.Type().(*types.Signature)
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
		if !pos.IsValid() || field.Pkg() != suite.typ.Obj().Pkg() {
			pos = method.pos(suite)
		}

		test.params = append(test.params, param{
			name: field.Name(),
			typ:  field.Type(),
			pos:  pos,
		})
	}

	return test
}
