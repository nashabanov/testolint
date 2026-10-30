package analyzer

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/types/typeutil"
)

// Construction starts at concrete contexts used by Testo, never at arbitrary
// type declarations. Fill runs before setPlugins, so cycles suppress field
// diagnostics for that root. All caches belong to this analyzer invocation.
func checkConstruction(p *analysis.Pass, facts *callFacts) {
	var checked typeutil.Map
	cycles := make(map[types.Type]bool)
	cycleFields := make(map[token.Pos]bool)
	for _, call := range facts.calls {
		root := constructionRoot(p.TypesInfo, facts, call)
		if root == nil {
			continue
		}
		root = types.Unalias(root)
		if checked.At(root) != nil {
			continue
		}
		checked.Set(root, true)
		graph := constructionGraph{active: make(map[types.Type]int), done: make(map[types.Type]bool)}
		if graph.findCycle(root, nil) {
			duplicate := false
			for _, typ := range graph.cycle {
				duplicate = duplicate || cycles[typ]
				cycles[typ] = true
			}
			if graph.field.Pkg() == p.Pkg {
				duplicate = duplicate || cycleFields[graph.field.Pos()]
				cycleFields[graph.field.Pos()] = true
			}
			if !duplicate {
				reportf(p, "TESTO022", constructionPos(p, graph.field, call.Pos()),
					"recursive initialization through field %q: %s", graph.field.Name(), graph.path)
			}
			continue
		}
		if graph.unknown {
			continue
		}
		checkConstructionFields(p, root, call.Pos(), make(map[types.Type]bool))
	}
}

func constructionRoot(info *types.Info, facts *callFacts, call *ast.CallExpr) types.Type {
	if info == nil {
		return nil
	}
	fn := facts.function(info, call.Fun)
	if fn == nil || fn.Pkg() == nil || fn.Pkg().Path() != testoImportPath {
		return nil
	}
	decl, ok := fn.Type().(*types.Signature)
	if !ok || decl.Recv() != nil {
		return nil
	}
	// The expression type contains the inferred/explicit generic instantiation,
	// unlike the function's declaration object. Local aliases retain this type.
	sig, ok := info.TypeOf(call.Fun).Underlying().(*types.Signature)
	if !ok {
		return nil
	}
	switch fn.Name() {
	case "RunSuite", "RunSubSuite":
		if len(call.Args) >= 2 {
			return testoTType(types.NewMethodSet(info.TypeOf(call.Args[1])))
		}
	case "Run":
		return sig.Params().At(0).Type()
	case "RunTest", "Test":
		index := 1
		if fn.Name() == "Test" {
			if !executedAdapter(info, facts, call) {
				return nil
			}
			index = 0
		}
		callback, ok := sig.Params().At(index).Type().Underlying().(*types.Signature)
		if ok && callback.Params().Len() == 1 {
			return callback.Params().At(0).Type()
		}
	}
	return nil
}

type constructionGraph struct {
	active  map[types.Type]int
	done    map[types.Type]bool
	stack   []types.Type
	fields  []*types.Var
	cycle   []types.Type
	field   *types.Var
	path    string
	unknown bool
}

// Match Fill's active path, not a global visited set: the same type on two
// independent branches is safe. Containers and inaccessible fields are leaves.
func (g *constructionGraph) findCycle(t types.Type, field *types.Var) bool {
	t = types.Unalias(t)
	if _, ok := t.(*types.TypeParam); ok {
		g.unknown = true
		return false
	}
	if g.done[t] {
		return false
	}
	if index, ok := g.active[t]; ok {
		g.cycle = append([]types.Type(nil), g.stack[index:]...)
		g.field = field
		names := []string{types.TypeString(t, func(pkg *types.Package) string { return pkg.Name() })}
		var previous *types.Var
		for _, f := range append(g.fields[index+1:], field) {
			if f != nil && f != previous {
				names = append(names, f.Name())
				previous = f
			}
		}
		names = append(names, names[0])
		g.path = strings.Join(names, " -> ")
		return true
	}
	g.active[t] = len(g.stack)
	g.stack = append(g.stack, t)
	g.fields = append(g.fields, field)
	defer func() {
		delete(g.active, t)
		g.stack = g.stack[:len(g.stack)-1]
		g.fields = g.fields[:len(g.fields)-1]
		g.done[t] = true
	}()
	switch typ := t.Underlying().(type) {
	case *types.Pointer:
		return g.findCycle(typ.Elem(), field)
	case *types.Struct:
		for f := range typ.Fields() {
			if f.Exported() && g.findCycle(f.Type(), f) {
				return true
			}
		}
	}
	return false
}

func checkConstructionFields(p *analysis.Pass, t types.Type, fallback token.Pos, seen map[types.Type]bool) {
	t = types.Unalias(t)
	if seen[t] {
		return
	}
	seen[t] = true
	if ptr, ok := t.Underlying().(*types.Pointer); ok {
		checkConstructionFields(p, ptr.Elem(), fallback, seen)
		return
	}
	if named, ok := t.(*types.Named); ok && named.Obj().Pkg() != nil &&
		named.Obj().Pkg().Path() == testoImportPath && named.Obj().Name() == "T" {
		return
	}
	st, ok := t.Underlying().(*types.Struct)
	if !ok {
		return
	}
	for field := range st.Fields() {
		if !field.Exported() {
			continue
		}
		if _, ok := types.Unalias(field.Type()).Underlying().(*types.Pointer); !ok {
			reportf(p, "TESTO021", constructionPos(p, field, fallback),
				"exported field %q must be a pointer for Testo initialization", field.Name())
			continue
		}
		checkConstructionFields(p, field.Type(), fallback, seen)
	}
}

func constructionPos(p *analysis.Pass, field *types.Var, fallback token.Pos) token.Pos {
	if field.Pkg() == p.Pkg && field.Pos().IsValid() {
		return field.Pos()
	}
	return fallback
}
