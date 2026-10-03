package analyzer

import "golang.org/x/tools/go/analysis"

type Rule func(*analysis.Pass, *Suite)

var rules = []Rule{
	checkTestSignature,
	checkHookSignature,
	checkMissingCaseProvider,
	checkCaseTypeMismatch,
	checkCasesSignature,
	checkOrphanCasesProvider,
}
