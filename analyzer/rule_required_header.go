package analyzer

import (
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

func allureContext(t types.Type) bool {
	methods := types.NewMethodSet(t)
	for i := 0; i < methods.Len(); i++ {
		fn := methods.At(i).Obj()
		if fn.Name() == "unwrapAllure" && fn.Pkg() != nil && fn.Pkg().Path() == allureImportPath {
			return true
		}
	}
	return false
}

func validStyleTest(s *suite, test *test) bool {
	sig := test.fn().Type().(*types.Signature)
	return test.method.decl != nil && test.method.decl.Body != nil &&
		isValidPrefixedName(test.fn().Name(), "Test") && sig.Results().Len() == 0 &&
		sig.Params().Len() >= 1 && sig.Params().Len() <= 2 &&
		types.Identical(sig.Params().At(0).Type(), s.tType) &&
		(sig.Params().Len() == 1 || isStructType(sig.Params().At(1).Type()))
}

func checkRequiredHeader(p *analysis.Pass, s *suite) {
	checkRequiredHeaderWithConfig(p, s, HeaderConfig{})
}

func checkRequiredHeaderWithConfig(p *analysis.Pass, s *suite, config HeaderConfig) {
	if len(config.Required) == 0 || !allureContext(s.tType) {
		return
	}
	inherited := map[string]bool{}
	unknown := false
	if config.IncludeBeforeEach == nil || *config.IncludeBeforeEach {
		for _, hook := range s.hooks {
			if hook.fn.Name() != "BeforeEach" {
				continue
			}
			sig := hook.fn.Type().(*types.Signature)
			if hook.fn.Pkg() != nil && hook.fn.Pkg().Path() == testoImportPath && isTestoSuiteType(sig.Recv().Type()) {
				continue
			}
			if hook.decl == nil || hook.decl.Body == nil || sig.Params().Len() != 1 || sig.Results().Len() != 0 || !types.Identical(sig.Params().At(0).Type(), s.tType) {
				unknown = true
				break
			}
			unknown = collectHeader(p, hook.decl.Body, sig.Params().At(0), config, inherited)
		}
	}
	for _, test := range s.tests {
		if !validStyleTest(s, test) {
			continue
		}
		found := map[string]bool{}
		for name := range inherited {
			found[name] = true
		}
		sig := test.fn().Type().(*types.Signature)
		localUnknown := collectHeader(p, test.method.decl.Body, sig.Params().At(0), config, found)
		if unknown || localUnknown {
			continue
		}
		var missing []string
		for _, group := range config.Required {
			present := false
			for _, name := range group {
				present = present || found[name]
			}
			if !present {
				missing = append(missing, strings.Join(group, " or "))
			}
		}
		if len(missing) > 0 {
			reportf(p, "TESTO028", test.fn().Pos(), "missing required test metadata: %s", strings.Join(missing, ", "))
		}
	}
}

// Only unconditional direct metadata calls prove presence. Control flow and
// calls receiving the context may contribute metadata, but are unknown here.
// Merely creating a function literal does not execute its body.
func collectHeader(p *analysis.Pass, body *ast.BlockStmt, context types.Object, config HeaderConfig, found map[string]bool) bool {
	unknown := false
	for _, stmt := range body.List {
		if metadataStatementWithHeader(p, stmt, context, config) {
			call := statementCall(stmt)
			sel := ast.Unparen(call.Fun).(*ast.SelectorExpr)
			fn := p.TypesInfo.Selections[sel].Obj()
			if fn.Pkg().Path() == allureImportPath {
				found[fn.Name()] = true
			}
			continue
		}
		switch stmt.(type) {
		case *ast.IfStmt, *ast.SwitchStmt, *ast.TypeSwitchStmt, *ast.ForStmt, *ast.RangeStmt, *ast.SelectStmt, *ast.BranchStmt, *ast.ReturnStmt, *ast.LabeledStmt:
			return true
		}
		ast.Inspect(stmt, func(n ast.Node) bool {
			if _, ok := n.(*ast.FuncLit); ok {
				return false
			}
			if call, ok := n.(*ast.CallExpr); ok {
				for _, arg := range call.Args {
					if sameObject(p.TypesInfo, arg, context) {
						unknown = true
					}
				}
				if sel, ok := ast.Unparen(call.Fun).(*ast.SelectorExpr); ok && sameObject(p.TypesInfo, sel.X, context) {
					fn := p.TypesInfo.ObjectOf(sel.Sel)
					if fn == nil || fn.Pkg() == nil || (fn.Pkg().Path() != allureImportPath && fn.Pkg().Path() != testoImportPath) {
						unknown = true
					}
				}
			}
			return true
		})
	}
	return unknown
}
