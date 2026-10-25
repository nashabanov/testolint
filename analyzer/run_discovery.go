package analyzer

import (
	"go/types"

	"golang.org/x/tools/go/analysis"
)

// discoverRuns identifies calls independently of suite declarations and rules.
// Local function aliases are accepted only with a direct initializer and no
// reassignment or address-taking anywhere in the package. This deliberately
// avoids making control-flow or execution-order assumptions.
func discoverRuns(p *analysis.Pass) []*runCall {
	if p.TypesInfo == nil {
		return nil
	}
	return discoverRunsWithFacts(p, collectCallFacts(p))
}

func discoverRunsWithFacts(p *analysis.Pass, facts *callFacts) []*runCall {
	var result []*runCall
	for _, call := range facts.calls {
		fn := facts.function(p.TypesInfo, call.Fun)
		if fn == nil || fn.Pkg() == nil || fn.Pkg().Path() != testoImportPath {
			continue
		}
		sig, ok := fn.Type().(*types.Signature)
		if !ok || sig.Recv() != nil {
			continue
		}
		var kind runKind
		switch fn.Name() {
		case "RunSuite":
			kind = runSuite
		case "RunSubSuite":
			kind = runSubSuite
		default:
			continue
		}
		if len(call.Args) < 2 {
			continue
		}
		suiteExpr := call.Args[1]
		suiteValue := facts.value(p.TypesInfo, suiteExpr)
		result = append(result, &runCall{
			kind: kind, call: call, suiteExpr: suiteExpr, enclosingFunc: facts.functions[call],
			suiteType: p.TypesInfo.TypeOf(suiteExpr), suiteValue: suiteValue,
		})
	}
	return result
}
