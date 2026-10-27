package analyzer

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

const noExecutions = -1 // A proven empty provider or Cartesian product.

func validParallelProviders(suite *suite) bool {
	for _, provider := range suite.cases {
		if _, problem := casesSignature(provider.fn.Type()); problem != "" || !isValidPrefixedName(provider.fn.Name(), "Cases") {
			return false
		}
	}
	return true
}

// Collection errors invalidate every test in the suite, unlike an unknown
// provider count. Keep these two reasons for excluding analysis separate.
func validParallelParameters(suite *suite, test *test) bool {
	if test.fn().Type().(*types.Signature).Params().Len() == 1 {
		return true
	}
	if !test.validParams(suite) {
		return false
	}
	for _, param := range test.params {
		provider := suite.casesByName[param.name]
		if !ast.IsExported(param.name) || provider == nil {
			return false
		}
		slice, _ := casesSignature(provider.fn.Type())
		if slice == nil || !types.AssignableTo(slice.Elem(), param.typ) {
			return false
		}
	}
	return true
}

// Called after collection validation. Cap multiplicity at two, which is
// sufficient to establish a second execution of the same test method.
func testExecutionCount(suite *suite, test *test, sizes map[*method]int) int {
	count := 1
	for _, param := range test.params {
		size := sizes[suite.casesByName[param.name]]
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
func constantProviderSize(p *analysis.Pass, provider *method) int {
	if provider.decl == nil || provider.decl.Body == nil || len(provider.decl.Body.List) != 1 {
		return 0
	}
	ret, ok := provider.decl.Body.List[0].(*ast.ReturnStmt)
	if !ok || len(ret.Results) != 1 {
		return 0
	}
	if emptySliceExpr(p.TypesInfo, ret.Results[0]) {
		return noExecutions
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
