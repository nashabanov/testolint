package analyzer

import "golang.org/x/tools/go/analysis"

type suiteRule func(*analysis.Pass, *suite)

type runRule func(*analysis.Pass, *runCall)

// packageRule checks package-wide calls independently of suite declarations.
type packageRule func(*analysis.Pass, *callFacts)

type parallelRule func(*analysis.Pass, *parallelModel)

type ruleSet struct {
	suites   []suiteRule
	parallel []parallelRule
	// Final suite checks retain their reporting order after parallel checks.
	afterParallel []suiteRule
	runs          []runRule
	packages      []packageRule
}

func suiteRules() []suiteRule {
	return []suiteRule{
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

func afterParallelRules() []suiteRule {
	return []suiteRule{
		checkPrematureTeardown,
	}
}

func parallelRules() []parallelRule {
	return []parallelRule{
		checkSharedBeforeEachWrite,
		checkSharedFixtureCleanup,
		checkSharedTestContext,
		checkDuplicateParallelLifecycle,
	}
}

func runRules() []runRule {
	return []runRule{
		checkNilSuite,
		checkRecursiveSubSuite,
		checkProviderBeforeEachDependency,
	}
}

func packageRules() []packageRule {
	return []packageRule{
		checkExecutionSafety,
		checkConstruction,
		checkPluginParentNil,
	}
}
