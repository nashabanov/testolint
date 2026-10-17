package analyzer

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

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
func testExecutionCount(p *analysis.Pass, suite *Suite, test *Test) int {
	count := 1
	for _, param := range test.Params {
		size := constantProviderSize(p, suite.CasesByName[param.Name])
		if size == 0 {
			return 0
		}
		if size >= 2 {
			count = 2
		}
	}
	return count
}

// Recognize only a sole return of an unkeyed slice literal of constants.
// Zero means no execution count was established (including an empty provider).
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
		return 0
	}
	for _, element := range lit.Elts {
		// Constants have no helper calls or side effects during collection.
		if p.TypesInfo.Types[element].Value == nil {
			return 0
		}
	}
	return len(lit.Elts)
}
