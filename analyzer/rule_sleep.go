package analyzer

import (
	"go/ast"
	"go/types"
	"golang.org/x/tools/go/analysis"
)

func checkSleepBody(p *analysis.Pass, body *ast.BlockStmt) {
	if body == nil {
		return
	}
	ast.Inspect(body, func(n ast.Node) bool {
		c, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		var obj types.Object
		switch f := ast.Unparen(c.Fun).(type) {
		case *ast.SelectorExpr:
			obj = p.TypesInfo.ObjectOf(f.Sel)
		case *ast.Ident:
			obj = p.TypesInfo.ObjectOf(f)
		}
		if fn, ok := obj.(*types.Func); ok && fn.Pkg() != nil && fn.Pkg().Path() == "time" && fn.Name() == "Sleep" {
			reportf(p, "TESTO027", c.Pos(), "avoid time.Sleep in tests; use polling or another condition-based wait")
		}
		return true
	})
}
func checkSleepSuite(p *analysis.Pass, s *suite) {
	methods := append([]*method{}, s.hooks...)
	for _, t := range s.tests {
		methods = append(methods, t.method)
	}
	for _, m := range methods {
		if m.decl == nil {
			continue
		}
		sig := m.fn.Type().(*types.Signature)
		if sig.Results().Len() != 0 || sig.Params().Len() < 1 || !types.Identical(sig.Params().At(0).Type(), s.tType) {
			continue
		}
		if isValidPrefixedName(m.fn.Name(), "Test") {
			if sig.Params().Len() > 2 || (sig.Params().Len() == 2 && !isStructType(sig.Params().At(1).Type())) {
				continue
			}
		} else if sig.Params().Len() != 1 {
			continue
		}
		checkSleepBody(p, m.decl.Body)
	}
}
func checkSleepCallbacks(p *analysis.Pass, facts *callFacts) {
	for _, c := range facts.calls {
		var obj types.Object
		switch f := runFunctionExpr(c.Fun).(type) {
		case *ast.SelectorExpr:
			obj = p.TypesInfo.ObjectOf(f.Sel)
		case *ast.Ident:
			obj = p.TypesInfo.ObjectOf(f)
		}
		fn, ok := obj.(*types.Func)
		if !ok || fn.Pkg() == nil || fn.Pkg().Path() != testoImportPath {
			continue
		}
		if fn.Name() != "Run" && fn.Name() != "RunTest" && fn.Name() != "Test" {
			continue
		}
		for _, arg := range c.Args {
			if lit, ok := ast.Unparen(arg).(*ast.FuncLit); ok {
				checkSleepBody(p, lit.Body)
			}
		}
	}
}

func checkSleep(p *analysis.Pass, facts *callFacts) {
	for _, s := range discoverSuites(p) {
		checkSleepSuite(p, s)
	}
	checkSleepCallbacks(p, facts)
}
