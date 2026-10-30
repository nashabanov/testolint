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
	var checks []suiteRule
	for _, registration := range registrations {
		if registration.suite != nil {
			checks = append(checks, registration.suite)
		}
	}
	return checks
}

func afterParallelRules() []suiteRule {
	var checks []suiteRule
	for _, registration := range registrations {
		if registration.afterParallel != nil {
			checks = append(checks, registration.afterParallel)
		}
	}
	return checks
}

func parallelRules() []parallelRule {
	var checks []parallelRule
	for _, registration := range registrations {
		if registration.parallel != nil {
			checks = append(checks, registration.parallel)
		}
	}
	return checks
}

func runRules() []runRule {
	var checks []runRule
	for _, registration := range registrations {
		if registration.run != nil {
			checks = append(checks, registration.run)
		}
	}
	return checks
}

func packageRules() []packageRule {
	var checks []packageRule
	for _, registration := range registrations {
		if registration.packageCheck != nil {
			checks = append(checks, registration.packageCheck)
		}
	}
	return checks
}
