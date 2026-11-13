package analyzer

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

// Symbol identifies either a package function or a method on a declared type.
// Pointer and value receivers have the same identity.
type Symbol struct {
	Package  string `json:"package"`
	Function string `json:"function"`
	Receiver string `json:"receiver"`
	Method   string `json:"method"`
}

func (s Symbol) validate() error {
	if s.Package == "" || (s.Function == "") == (s.Method == "") ||
		(s.Function != "" && (s.Receiver != "" || !token.IsIdentifier(s.Function))) ||
		(s.Method != "" && (!token.IsIdentifier(s.Receiver) || !token.IsIdentifier(s.Method))) {
		return fmt.Errorf("testolint: operation must specify package/function or package/receiver/method: %+v", s)
	}
	return nil
}

func symbolOf(fn *types.Func) Symbol {
	if fn == nil || fn.Pkg() == nil {
		return Symbol{}
	}
	s := Symbol{Package: fn.Pkg().Path()}
	sig := fn.Type().(*types.Signature)
	if sig.Recv() == nil {
		s.Function = fn.Name()
		return s
	}
	t := types.Unalias(sig.Recv().Type())
	if ptr, ok := t.(*types.Pointer); ok {
		t = types.Unalias(ptr.Elem())
	}
	if named, ok := t.(*types.Named); ok {
		s.Receiver = named.Obj().Name()
		s.Method = fn.Name()
	}
	return s
}

type operation struct {
	call    *ast.CallExpr
	site    *ast.CallExpr
	id      ruleID
	covered bool
}

type operationWalker struct {
	pass           *analysis.Pass
	facts          *callFacts
	declarations   map[*types.Func]*ast.FuncDecl
	kinds          map[Symbol]ruleID
	operations     []operation
	events         []phaseEvent
	suppressPhases bool
	helpers        map[*types.Func]helperFacts
}

func checkStepCoverage(p *analysis.Pass, s *suite) { checkStepCoverageWithConfig(p, s, StyleConfig{}) }

func checkStepCoverageWithConfig(p *analysis.Pass, s *suite, config StyleConfig) {
	if !allureContext(s.tType) {
		return
	}
	newOperationWalker(p, config).checkSuite(s, config)
}

func newOperationWalker(p *analysis.Pass, config StyleConfig) *operationWalker {
	w := &operationWalker{pass: p, facts: collectCallFacts(p), declarations: map[*types.Func]*ast.FuncDecl{}, kinds: builtinAssertions(), helpers: map[*types.Func]helperFacts{}}
	for _, symbol := range config.Assertions {
		w.kinds[symbol] = "TESTO031"
	}
	for _, symbol := range config.Actions {
		w.kinds[symbol] = "TESTO029"
	}
	for _, file := range p.Files {
		for _, decl := range file.Decls {
			if decl, ok := decl.(*ast.FuncDecl); ok {
				if fn, ok := p.TypesInfo.ObjectOf(decl.Name).(*types.Func); ok {
					w.declarations[fn.Origin()] = decl
				}
			}
		}
	}
	return w
}

func (w *operationWalker) checkSuite(s *suite, config StyleConfig) {
	p := w.pass
	for _, test := range s.tests {
		if !validStyleTest(s, test) {
			continue
		}
		w.operations = nil
		w.events = nil
		w.walk(test.method.decl.Body, false, nil, 0)
		reportPhases(p, w.events, config.Phases)
		reported := map[struct {
			pos token.Pos
			id  ruleID
		}]bool{}
		for _, op := range w.operations {
			if op.covered {
				continue
			}
			site := op.call
			message := "wrap this assertion in an Allure Step or Setup"
			if op.id == "TESTO029" {
				message = "wrap this action in an Allure Step or Setup"
			}
			var related []analysis.RelatedInformation
			if op.site != nil {
				site = op.site
				message = "wrap this helper call in a step in the test, or add the necessary steps inside the helper"
				related = []analysis.RelatedInformation{{Pos: op.call.Pos(), End: op.call.End(), Message: "uncovered operation inside helper"}}
			}
			key := struct {
				pos token.Pos
				id  ruleID
			}{site.Pos(), op.id}
			if reported[key] {
				continue
			}
			reported[key] = true
			report(p, op.id, analysis.Diagnostic{Pos: site.Pos(), End: site.End(), Message: message, Related: related})
		}
	}
}

func (w *operationWalker) walk(node ast.Node, covered bool, site *ast.CallExpr, depth int) {
	if node == nil {
		return
	}
	ast.Inspect(node, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.IfStmt, *ast.SwitchStmt, *ast.TypeSwitchStmt, *ast.ForStmt, *ast.RangeStmt, *ast.SelectStmt, *ast.LabeledStmt:
			// Coverage remains lexical; unsupported control flow interrupts
			// phase inference rather than concatenating mutually exclusive paths.
			w.event(phaseEvent{})
			saved := w.suppressPhases
			w.suppressPhases = true
			ast.Inspect(n, func(child ast.Node) bool {
				if child == n {
					return true
				}
				if child != nil {
					w.walk(child, covered, site, depth)
				}
				return false
			})
			w.suppressPhases = saved
			return false
		case *ast.ReturnStmt:
			for _, result := range n.Results {
				w.walk(result, covered, site, depth)
			}
			return false
		case *ast.BranchStmt:
			w.event(phaseEvent{})
			return false
		case *ast.FuncLit:
			return false
		case *ast.GoStmt, *ast.DeferStmt:
			// Asynchronous and deferred bodies are not synchronous report containers.
			return false
		case *ast.CallExpr:
			// Arguments execute before the call; function literal bodies do not.
			for _, arg := range n.Args {
				w.walk(arg, covered, site, depth)
			}
			fn := w.facts.function(w.pass.TypesInfo, n.Fun)
			symbol := symbolOf(fn)
			if symbol.Package == allureImportPath && symbol.Function != "" {
				switch symbol.Function {
				case "TearDown":
					return false
				case "Step", "Setup":
					saved := w.suppressPhases
					if symbol.Function == "Setup" {
						w.event(phaseEvent{call: n, site: site, phase: arrangePhase})
						w.suppressPhases = true
					}
					if len(n.Args) >= 3 {
						if callback, ok := ast.Unparen(n.Args[2]).(*ast.FuncLit); ok {
							w.walk(callback.Body, true, site, depth)
						} else {
							w.event(phaseEvent{})
						}
					}
					w.suppressPhases = saved
					return false
				}
			}
			if id := w.kinds[symbol]; id != "" {
				w.operations = append(w.operations, operation{n, site, id, covered})
				phase := actionPhase
				if id == "TESTO031" {
					phase = assertPhase
				}
				w.event(phaseEvent{call: n, site: site, phase: phase})
			} else if fn != nil && depth == 0 {
				if decl := w.declarations[fn.Origin()]; decl != nil && decl.Body != nil {
					w.helper(fn.Origin(), decl.Body, covered, n)
				}
			}
			return false
		}
		return true
	})
}

// Helper summaries are context-independent: coverage is relative to the helper,
// and phase events retain their operation positions. Callsite rebasing happens
// when consumed. No summary follows another helper.
type helperFacts struct {
	operations []operation
	events     []phaseEvent
}

func (w *operationWalker) helper(fn *types.Func, body *ast.BlockStmt, covered bool, site *ast.CallExpr) {
	facts, ok := w.helpers[fn]
	if !ok {
		operations, events, suppressed := w.operations, w.events, w.suppressPhases
		w.operations = nil
		w.events = nil
		w.suppressPhases = false
		w.walk(body, false, nil, 1)
		facts = helperFacts{w.operations, w.events}
		w.helpers[fn] = facts
		w.operations, w.events, w.suppressPhases = operations, events, suppressed
	}
	for _, op := range facts.operations {
		op.site = site
		op.covered = op.covered || covered
		w.operations = append(w.operations, op)
	}
	for _, event := range facts.events {
		event.site = site
		w.event(event)
	}
}
