package analyzer

type ruleID string

type ruleClass string

const (
	correctness ruleClass = "correctness"
	policy      ruleClass = "policy"
	style       ruleClass = "style"
)

type ruleMeta struct {
	Name            string
	Class           ruleClass
	EnableByDefault bool
	DocURL          string
}

// Existing correctness and policy rules are enabled by default. Class does not
// determine enablement. Each rule links to its own documentation page.
var registry = map[ruleID]ruleMeta{
	"TESTO029": {Name: "unreported-action", Class: style, DocURL: ruleDocURL + "/TESTO029.md"},
	"TESTO031": {Name: "assert-inside-step", Class: style, DocURL: ruleDocURL + "/TESTO031.md"},
	"TESTO028": {Name: "required-test-header", Class: style, DocURL: ruleDocURL + "/TESTO028.md"},
	"TESTO001": {Name: "test-signature", Class: correctness, EnableByDefault: true, DocURL: ruleDocURL + "/TESTO001.md"},
	"TESTO002": {Name: "hook-signature", Class: correctness, EnableByDefault: true, DocURL: ruleDocURL + "/TESTO002.md"},
	"TESTO003": {Name: "missing-case-provider", Class: correctness, EnableByDefault: true, DocURL: ruleDocURL + "/TESTO003.md"},
	"TESTO004": {Name: "case-type-mismatch", Class: correctness, EnableByDefault: true, DocURL: ruleDocURL + "/TESTO004.md"},
	"TESTO005": {Name: "cases-signature", Class: correctness, EnableByDefault: true, DocURL: ruleDocURL + "/TESTO005.md"},
	"TESTO006": {Name: "orphan-cases-provider", Class: policy, EnableByDefault: true, DocURL: ruleDocURL + "/TESTO006.md"},
	"TESTO007": {Name: "test-name", Class: correctness, EnableByDefault: true, DocURL: ruleDocURL + "/TESTO007.md"},
	"TESTO008": {Name: "cases-name", Class: correctness, EnableByDefault: true, DocURL: ruleDocURL + "/TESTO008.md"},
	"TESTO009": {Name: "unexported-param-field", Class: correctness, EnableByDefault: true, DocURL: ruleDocURL + "/TESTO009.md"},
	"TESTO010": {Name: "empty-cases", Class: correctness, EnableByDefault: true, DocURL: ruleDocURL + "/TESTO010.md"},
	"TESTO011": {Name: "empty-suite", Class: correctness, EnableByDefault: true, DocURL: ruleDocURL + "/TESTO011.md"},
	"TESTO012": {Name: "nil-suite", Class: policy, EnableByDefault: true, DocURL: ruleDocURL + "/TESTO012.md"},
	"TESTO013": {Name: "recursive-subsuite", Class: correctness, EnableByDefault: true, DocURL: ruleDocURL + "/TESTO013.md"},
	"TESTO014": {Name: "shared-before-each-write", Class: policy, EnableByDefault: true, DocURL: ruleDocURL + "/TESTO014.md"},
	"TESTO015": {Name: "shared-fixture-cleanup", Class: policy, EnableByDefault: true, DocURL: ruleDocURL + "/TESTO015.md"},
	"TESTO016": {Name: "shared-test-context", Class: policy, EnableByDefault: true, DocURL: ruleDocURL + "/TESTO016.md"},
	"TESTO017": {Name: "duplicate-parallel-lifecycle", Class: correctness, EnableByDefault: true, DocURL: ruleDocURL + "/TESTO017.md"},
	"TESTO018": {Name: "run-during-cleanup", Class: correctness, EnableByDefault: true, DocURL: ruleDocURL + "/TESTO018.md"},
	"TESTO019": {Name: "nil-callback", Class: correctness, EnableByDefault: true, DocURL: ruleDocURL + "/TESTO019.md"},
	"TESTO020": {Name: "premature-teardown", Class: correctness, EnableByDefault: true, DocURL: ruleDocURL + "/TESTO020.md"},
	"TESTO021": {Name: "construction-field", Class: correctness, EnableByDefault: true, DocURL: ruleDocURL + "/TESTO021.md"},
	"TESTO022": {Name: "construction-cycle", Class: correctness, EnableByDefault: true, DocURL: ruleDocURL + "/TESTO022.md"},
	"TESTO023": {Name: "provider-before-each", Class: correctness, EnableByDefault: true, DocURL: ruleDocURL + "/TESTO023.md"},
	"TESTO024": {Name: "plugin-parent-nil", Class: correctness, EnableByDefault: true, DocURL: ruleDocURL + "/TESTO024.md"},
	"TESTO025": {Name: "blank-after-header", Class: style, DocURL: ruleDocURL + "/TESTO025.md"},
	"TESTO026": {Name: "blank-between-steps", Class: style, DocURL: ruleDocURL + "/TESTO026.md"},
	"TESTO027": {Name: "no-sleep", Class: policy, EnableByDefault: true, DocURL: ruleDocURL + "/TESTO027.md"},
}

const ruleDocURL = "https://github.com/nashabanov/testolint/blob/main/docs/rules"

// Each checker appears once, even when it emits multiple diagnostic IDs. The
// order within each execution group preserves the existing reporting order.
type ruleRegistration struct {
	ids           []ruleID
	suite         suiteRule
	parallel      parallelRule
	afterParallel suiteRule
	run           runRule
	packageCheck  packageRule
}

var registrations = []ruleRegistration{
	{ids: []ruleID{"TESTO029", "TESTO031"}, suite: checkStepCoverage},
	{ids: []ruleID{"TESTO028"}, suite: checkRequiredHeader},
	{ids: []ruleID{"TESTO001"}, suite: checkTestSignature},
	{ids: []ruleID{"TESTO002"}, suite: checkHookSignature},
	{ids: []ruleID{"TESTO003"}, suite: checkMissingCaseProvider},
	{ids: []ruleID{"TESTO004"}, suite: checkCaseTypeMismatch},
	{ids: []ruleID{"TESTO005"}, suite: checkCasesSignature},
	{ids: []ruleID{"TESTO006"}, suite: checkOrphanCasesProvider},
	{ids: []ruleID{"TESTO007"}, suite: checkMalformedTestName},
	{ids: []ruleID{"TESTO008"}, suite: checkMalformedCasesName},
	{ids: []ruleID{"TESTO009"}, suite: checkUnexportedParamField},
	{ids: []ruleID{"TESTO010"}, suite: checkEmptyCases},
	{ids: []ruleID{"TESTO011"}, suite: checkEmptySuite},
	{ids: []ruleID{"TESTO027"}, packageCheck: checkSleep},
	{ids: []ruleID{"TESTO014"}, parallel: checkSharedBeforeEachWrite},
	{ids: []ruleID{"TESTO015"}, parallel: checkSharedFixtureCleanup},
	{ids: []ruleID{"TESTO016"}, parallel: checkSharedTestContext},
	{ids: []ruleID{"TESTO017"}, parallel: checkDuplicateParallelLifecycle},
	{ids: []ruleID{"TESTO020"}, afterParallel: checkPrematureTeardown},
	{ids: []ruleID{"TESTO012"}, run: checkNilSuite},
	{ids: []ruleID{"TESTO013"}, run: checkRecursiveSubSuite},
	{ids: []ruleID{"TESTO023"}, run: checkProviderBeforeEachDependency},
	{ids: []ruleID{"TESTO018", "TESTO019"}, packageCheck: checkExecutionSafety},
	{ids: []ruleID{"TESTO021", "TESTO022"}, packageCheck: checkConstruction},
	{ids: []ruleID{"TESTO024"}, packageCheck: checkPluginParentNil},
	{ids: []ruleID{"TESTO025", "TESTO026"}, suite: checkStyle},
}
