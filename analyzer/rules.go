package analyzer

import "golang.org/x/tools/go/analysis"

type Rule func(*analysis.Pass, *Suite)

type RunRule func(*analysis.Pass, *RunCall)

func suiteRules() []Rule {
	return []Rule{
		checkTestSignature,
		checkHookSignature,
		checkMissingCaseProvider,
		checkCaseTypeMismatch,
		checkCasesSignature,
		checkOrphanCasesProvider,
		checkMalformedTestName,
		checkMalformedCasesName,
		checkUnexportedParamField,
		checkEmptyCases,
		checkEmptySuite,
	}
}
