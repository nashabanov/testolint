package analyzer

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

const noExecutions = -1 // A proven empty provider or Cartesian product.

func validParallelProviders(suite *Suite) bool {
	for _, provider := range suite.Cases {
		if _, problem := casesSignature(provider.Func.Type()); problem != "" || !isValidPrefixedName(provider.Func.Name(), "Cases") {
			return false
		}
	}
	return true
}

// Collection errors invalidate every test in the suite, unlike an unknown
// provider count. Keep these two reasons for excluding analysis separate.
func validParallelParameters(suite *Suite, test *Test) bool {
	if test.Func().Type().(*types.Signature).Params().Len() == 1 {
		return true
	}
	if !test.validParams(suite) {
		return false
	}
	for _, param := range test.Params {
		provider := suite.CasesByName[param.Name]
		if !ast.IsExported(param.Name) || provider == nil {
			return false
		}
		slice, _ := casesSignature(provider.Func.Type())
		if slice == nil || !types.AssignableTo(slice.Elem(), param.Type) {
			return false
		}
	}
	return true
}

// Called after collection validation. Cap multiplicity at two, which is
// sufficient to establish a second execution of the same test method.
func testExecutionCount(suite *Suite, test *Test, sizes map[*Method]int) int {
	count := 1
	for _, param := range test.Params {
		size := sizes[suite.CasesByName[param.Name]]
		if size == noExecutions {
			return noExecutions
		}
		if size == 0 {
			count = 0
		}
		if size >= 2 && count != 0 {
			count = 2
		}
	}
	return count
}

// Recognize only a sole return of an unkeyed slice literal of constants.
// Zero means unknown; noExecutions distinguishes a proven empty provider.
func constantProviderSize(p *analysis.Pass, provider *Method) int {
	if provider.Decl == nil || provider.Decl.Body == nil || len(provider.Decl.Body.List) != 1 {
		return 0
	}
	ret, ok := provider.Decl.Body.List[0].(*ast.ReturnStmt)
	if !ok || len(ret.Results) != 1 {
		return 0
	}
	lit, ok := ast.Unparen(ret.Results[0]).(*ast.CompositeLit)
	if !ok {
		if id, ok := ast.Unparen(ret.Results[0]).(*ast.Ident); ok && p.TypesInfo.ObjectOf(id) == types.Universe.Lookup("nil") {
			return noExecutions
		}
		return 0
	}
	if len(lit.Elts) == 0 {
		return noExecutions
	}
	for _, element := range lit.Elts {
		// Constants have no helper calls or side effects during collection.
		if p.TypesInfo.Types[element].Value == nil {
			return 0
		}
	}
	return len(lit.Elts)
}
