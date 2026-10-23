package analyzer

import (
	"go/ast"
	"go/types"
	"slices"

	"golang.org/x/tools/go/analysis"
)

// discoverRuns identifies calls independently of suite declarations and rules.
// Local function aliases are accepted only with a direct initializer and no
// reassignment or address-taking anywhere in the package. This deliberately
// avoids making control-flow or execution-order assumptions.
func discoverRuns(p *analysis.Pass) []*RunCall {
	if p.TypesInfo == nil {
		return nil
	}
	return discoverRunsWithFacts(p, collectCallFacts(p))
}

func discoverRunsWithFacts(p *analysis.Pass, facts *callFacts) []*RunCall {
	var result []*RunCall
	for _, call := range facts.calls {
		fn := facts.function(p.TypesInfo, call.Fun)
		if fn == nil || fn.Pkg() == nil || fn.Pkg().Path() != testoImportPath {
			continue
		}
		sig, ok := fn.Type().(*types.Signature)
		if !ok || sig.Recv() != nil {
			continue
		}
		var kind RunKind
		switch fn.Name() {
		case "RunSuite":
			kind = RunSuite
		case "RunSubSuite":
			kind = RunSubSuite
		default:
			continue
		}
		if len(call.Args) < 2 {
			continue
		}
		suiteExpr := call.Args[1]
		suiteValue := facts.value(p.TypesInfo, suiteExpr)
		result = append(result, &RunCall{
			Kind: kind, Call: call, SuiteExpr: suiteExpr, EnclosingFunc: facts.functions[call],
			SuiteType: p.TypesInfo.TypeOf(suiteExpr), SuiteValue: suiteValue,
		})
	}
	return result
}

func runFunctionExpr(expr ast.Expr) ast.Expr {
	for {
		switch e := ast.Unparen(expr).(type) {
		case *ast.IndexExpr:
			expr = e.X
		case *ast.IndexListExpr:
			expr = e.X
		default:
			return e
		}
	}
}

// A call inside a function literal belongs to that literal, even when it is
// nested in a method. It cannot be treated as a direct call by the method.
func enclosingFunc(stack []ast.Node) *ast.FuncDecl {
	for _, node := range slices.Backward(stack) {
		switch node := node.(type) {
		case *ast.FuncDecl:
			return node
		case *ast.FuncLit:
			return nil
		}
	}
	return nil
}
