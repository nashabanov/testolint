package analyzer

import (
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/analysis"
)

type testPhase uint8

const (
	unknownPhase testPhase = iota
	arrangePhase
	actionPhase
	assertPhase
)

type phaseEvent struct {
	call  ast.Node
	site  *ast.CallExpr
	phase testPhase
}

func (w *operationWalker) event(event phaseEvent) {
	if !w.suppressPhases {
		w.events = append(w.events, event)
	}
}

func reportPhases(p *analysis.Pass, events []phaseEvent, config PhaseConfig) {
	strict := config.Strict == nil || *config.Strict
	previous := unknownPhase
	reported := map[token.Pos]bool{}
	for _, event := range events {
		if event.phase == unknownPhase {
			previous = unknownPhase
			continue
		}
		invalid := event.phase == arrangePhase && previous > arrangePhase || strict && event.phase == actionPhase && previous == assertPhase
		if invalid {
			call := event.call
			var related []analysis.RelatedInformation
			if event.site != nil {
				call = event.site
				related = []analysis.RelatedInformation{{Pos: event.call.Pos(), End: event.call.End(), Message: "phase transition inside helper"}}
			}
			if !reported[call.Pos()] {
				report(p, "TESTO030", analysis.Diagnostic{Pos: call.Pos(), End: call.End(), Message: "test phases must follow Data/Setup, Action, Assert order", Related: related})
				reported[call.Pos()] = true
			}
		}
		previous = event.phase
	}
}
