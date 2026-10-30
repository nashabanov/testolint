package analyzer

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

const pluginImportPath = testoImportPath + "/testoplugin"

// Only inspect locally declared plugin methods reached by a concrete, valid
// construction graph. Explicit method references may supply a genuine nil
// interface, so they disable this inference for that method in this package.
func checkPluginParentNil(p *analysis.Pass, facts *callFacts) {
	methods := make(map[*types.Func]bool)
	for _, call := range facts.calls {
		root := constructionRoot(p.TypesInfo, facts, call)
		if root == nil {
			continue
		}
		graph := constructionGraph{active: make(map[types.Type]int), done: make(map[types.Type]bool)}
		if graph.findCycle(root, nil) || graph.unknown {
			continue
		}
		var found []*types.Func
		if wiredPluginMethods(root, root, p.Pkg, make(map[types.Type]bool), &found) {
			for _, fn := range found {
				methods[fn] = true
			}
		}
	}
	for _, obj := range p.TypesInfo.Uses {
		if fn, ok := obj.(*types.Func); ok {
			delete(methods, fn.Origin())
		}
	}
	for _, file := range p.Files {
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Body == nil || !methods[functionObject(p.TypesInfo, fn)] {
				continue
			}
			sig := functionObject(p.TypesInfo, fn).Type().(*types.Signature)
			parent := sig.Params().At(0)
			if facts.invalid[parent] {
				continue
			}
			// A direct if condition gives the comparison a control-flow role.
			// Aliases, nested closures and compound expressions remain unknown.
			ast.Inspect(fn.Body, func(n ast.Node) bool {
				if _, ok := n.(*ast.FuncLit); ok {
					return false
				}
				statement, ok := n.(*ast.IfStmt)
				if !ok {
					return true
				}
				comparison, ok := ast.Unparen(statement.Cond).(*ast.BinaryExpr)
				if !ok || (comparison.Op != token.EQL && comparison.Op != token.NEQ) {
					return true
				}
				isNil := func(expr ast.Expr) bool {
					return sameObject(p.TypesInfo, expr, types.Universe.Lookup("nil"))
				}
				if (sameObject(p.TypesInfo, comparison.X, parent) && isNil(comparison.Y)) ||
					(sameObject(p.TypesInfo, comparison.Y, parent) && isNil(comparison.X)) {
					reportf(p, "TESTO024", comparison.Pos(),
						"Plugin parent interface cannot identify the root with nil; compare the asserted concrete pointer with nil")
				}
				return true
			})
		}
	}
}

func functionObject(info *types.Info, decl *ast.FuncDecl) *types.Func {
	fn, _ := info.Defs[decl.Name].(*types.Func)
	return fn
}

// Follow the pointer-only wiring supported by setPlugins, stopping at testo.T.
// Construction errors suppress this dependent check for the entire root.
func wiredPluginMethods(t, root types.Type, pkg *types.Package, seen map[types.Type]bool, found *[]*types.Func) bool {
	t = types.Unalias(t)
	if seen[t] {
		return true
	}
	seen[t] = true
	_, pointer := t.Underlying().(*types.Pointer)
	if pointer && !types.Identical(t, types.Unalias(root)) {
		obj, index, _ := types.LookupFieldOrMethod(t, false, pkg, "Plugin")
		if fn, ok := obj.(*types.Func); ok && len(index) == 1 && fn.Pkg() == pkg {
			sig := fn.Type().(*types.Signature)
			if sig.Params().Len() > 0 {
				parent, ok := types.Unalias(sig.Params().At(0).Type()).(*types.Named)
				if ok && parent.Obj().Pkg() != nil && parent.Obj().Pkg().Path() == pluginImportPath && parent.Obj().Name() == "Plugin" {
					if iface, ok := parent.Underlying().(*types.Interface); ok && types.Implements(t, iface) {
						*found = append(*found, fn.Origin())
					}
				}
			}
		}
	}
	if ptr, ok := t.Underlying().(*types.Pointer); ok {
		return wiredPluginMethods(ptr.Elem(), root, pkg, seen, found)
	}
	if named, ok := t.(*types.Named); ok && named.Obj().Pkg() != nil && named.Obj().Pkg().Path() == testoImportPath && named.Obj().Name() == "T" {
		return true
	}
	st, ok := t.Underlying().(*types.Struct)
	if !ok {
		return true
	}
	for field := range st.Fields() {
		if !field.Exported() {
			continue
		}
		if _, ok := types.Unalias(field.Type()).Underlying().(*types.Pointer); !ok || !wiredPluginMethods(field.Type(), root, pkg, seen, found) {
			return false
		}
	}
	return true
}
