package analyzer

import (
	"bytes"
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

const allureImportPath = "github.com/ozontech/testo-allure"

// Style checks only local, runnable test declarations. Imported bodies and
// helper functions are deliberately outside their scope.
func checkStyle(p *analysis.Pass, s *suite) {
	for _, test := range s.tests {
		decl := test.method.decl
		sig := test.fn().Type().(*types.Signature)
		if decl == nil || decl.Body == nil || !isValidPrefixedName(test.fn().Name(), "Test") ||
			sig.Params().Len() < 1 || sig.Params().Len() > 2 || sig.Results().Len() != 0 ||
			!types.Identical(sig.Params().At(0).Type(), s.tType) ||
			(sig.Params().Len() == 2 && !isStructType(sig.Params().At(1).Type())) {
			continue
		}
		body := decl.Body.List
		header := 0
		for header < len(body) && metadataStatement(p, body[header], sig.Params().At(0)) {
			header++
		}
		if header > 0 && header < len(body) {
			requireBlank(p, "TESTO025", body[header-1], body[header])
		}
		ast.Inspect(decl.Body, func(n ast.Node) bool {
			// Inspect nested lexical blocks too, including step callbacks.
			block, ok := n.(*ast.BlockStmt)
			if !ok {
				return true
			}
			for i := 1; i < len(block.List); i++ {
				if stepStatement(p, block.List[i-1]) && stepStatement(p, block.List[i]) {
					requireBlank(p, "TESTO026", block.List[i-1], block.List[i])
				}
			}
			return true
		})
	}
}

func statementCall(stmt ast.Stmt) *ast.CallExpr {
	expr, ok := stmt.(*ast.ExprStmt)
	if !ok {
		return nil
	}
	call, _ := ast.Unparen(expr.X).(*ast.CallExpr)
	return call
}

func metadataStatement(p *analysis.Pass, stmt ast.Stmt, context types.Object) bool {
	call := statementCall(stmt)
	if call == nil {
		return false
	}
	sel, ok := ast.Unparen(call.Fun).(*ast.SelectorExpr)
	if !ok || !sameObject(p.TypesInfo, sel.X, context) {
		return false
	}
	selection := p.TypesInfo.Selections[sel]
	if selection == nil || selection.Kind() != types.MethodVal {
		return false
	}
	fn := selection.Obj()
	if fn.Pkg() == nil || fn.Pkg().Path() != allureImportPath {
		return false
	}
	switch fn.Name() {
	case "Title", "Titlef", "Tags", "Description", "Descriptionf", "ID", "Labels", "Links", "Owner", "Severity", "Epic", "Feature", "Story", "Parameters", "Flaky", "Muted", "Known":
		return true
	}
	return false
}

func stepStatement(p *analysis.Pass, stmt ast.Stmt) bool {
	call := statementCall(stmt)
	if call == nil {
		return false
	}
	sel, ok := runFunctionExpr(call.Fun).(*ast.SelectorExpr)
	if !ok { // Also support dot imports, without relying on the local name.
		id, ok := runFunctionExpr(call.Fun).(*ast.Ident)
		if !ok {
			return false
		}
		fn, ok := p.TypesInfo.ObjectOf(id).(*types.Func)
		return ok && fn.Pkg() != nil && fn.Pkg().Path() == allureImportPath && fn.Name() == "Step"
	}
	fn, ok := p.TypesInfo.ObjectOf(sel.Sel).(*types.Func)
	return ok && fn.Pkg() != nil && fn.Pkg().Path() == allureImportPath && fn.Name() == "Step" && fn.Type().(*types.Signature).Recv() == nil
}

// Read physical source lines (not //line-adjusted positions). Comments do not
// count as whitespace. Insert before all leading comments of the next statement,
// preserving a trailing comment on the previous statement and the line ending.
func requireBlank(p *analysis.Pass, id ruleID, previous, next ast.Stmt) {
	file := p.Fset.File(previous.Pos())
	if file == nil || p.Fset.File(next.Pos()) != file {
		return
	}
	source, err := p.ReadFile(file.Name())
	if err != nil {
		return
	}
	end := file.Offset(previous.End())
	start := file.Offset(next.Pos())
	if end > start || start > len(source) {
		return
	}
	var comments []*ast.Comment
	for _, syntax := range p.Files {
		if p.Fset.File(syntax.Pos()) != file {
			continue
		}
		for _, group := range syntax.Comments {
			for _, comment := range group.List {
				if comment.Pos() >= previous.End() && comment.End() <= next.Pos() {
					comments = append(comments, comment)
					// A comment starting on the previous statement's final line belongs
					// to that statement, even when the comment spans several lines.
					if file.Line(comment.Pos()) == file.Line(file.Pos(end)) {
						end = file.Offset(comment.End())
					}
				}
			}
		}
	}
	firstNewline := bytes.IndexByte(source[end:start], '\n')
	var fixes []analysis.SuggestedFix
	if firstNewline >= 0 {
		insertion := end + firstNewline + 1
		lines := bytes.Split(source[insertion:start], []byte("\n"))
		offset := insertion
		for _, line := range lines[:len(lines)-1] {
			insideComment := false
			for _, comment := range comments {
				if file.Offset(comment.Pos()) <= offset && offset < file.Offset(comment.End()) {
					insideComment = true
					break
				}
			}
			if !insideComment && len(bytes.TrimSpace(line)) == 0 {
				return
			}
			offset += len(line) + 1
		}
		newline := []byte("\n")
		if insertion >= 2 && source[insertion-2] == '\r' {
			newline = []byte("\r\n")
		}
		fixes = []analysis.SuggestedFix{{Message: "Insert blank line", TextEdits: []analysis.TextEdit{{Pos: file.Pos(insertion), End: file.Pos(insertion), NewText: newline}}}}
	}
	// Same-line statements receive a diagnostic without an unsafe source edit.
	message := "insert a blank line between adjacent steps"
	if id == "TESTO025" {
		message = "insert a blank line after test metadata"
	}
	report(p, id, analysis.Diagnostic{Pos: next.Pos(), Message: message, SuggestedFixes: fixes})
}
