# Rule catalog

<a id="rule-examples-and-remedies"></a>

Both the CLI and golangci-lint run the same 27 rules. All 25 correctness/policy
rules are enabled by default; the two style rules are opt-in.

- **Correctness:** Testo contract violations, missing execution and concrete lifecycle hazards.
- **Policy:** checks that can flag intentional code, including unused providers and shared-field ownership.
- **Style:** source formatting, with automatic fixes where safe.

See [configuration](configuration.md) to select rules and [shared analysis scope](analysis-scope.md)
for common limitations. Rule pages explain each diagnostic, its remedy and boundaries.
Examples are fragments; “Valid” means valid for the illustrated rule. These checks
do not execute tests or prove race freedom; runtime plugins and filters can affect execution.

The table's names, classes and defaults are checked against the rule registry.

| ID | Name | Class | Default | Autofix | Purpose |
| --- | --- | --- | --- | --- | --- |
| <a id="testo001"></a>[`TESTO001`](rules/TESTO001.md) | `test-signature` | correctness | enabled | No | Tests accept exactly the suite's `T`, optionally a struct, and return no values. |
| <a id="testo002"></a>[`TESTO002`](rules/TESTO002.md) | `hook-signature` | correctness | enabled | No | `BeforeAll`, `BeforeEach`, `AfterEach`, `AfterAll` accept exactly `T` and return no values. |
| <a id="testo003"></a>[`TESTO003`](rules/TESTO003.md) | `missing-case-provider` | correctness | enabled | No | Exported parameter fields require a `Cases<Field>` method. |
| <a id="testo004"></a>[`TESTO004`](rules/TESTO004.md) | `case-type-mismatch` | correctness | enabled | No | Provider slice elements must be assignable to the parameter field type. |
| <a id="testo005"></a>[`TESTO005`](rules/TESTO005.md) | `cases-signature` | correctness | enabled | No | Providers accept no arguments and return exactly one slice. |
| <a id="testo006"></a>[`TESTO006`](rules/TESTO006.md) | `orphan-cases-provider` | policy | enabled | No | Policy: providers should be referenced by a test parameter; Testo permits unused providers. |
| <a id="testo007"></a>[`TESTO007`](rules/TESTO007.md) | `test-name` | correctness | enabled | No | After `Test`, the suffix is empty or begins with a non-lowercase Unicode rune. |
| <a id="testo008"></a>[`TESTO008`](rules/TESTO008.md) | `cases-name` | correctness | enabled | No | After `Cases`, the suffix is empty or begins with a non-lowercase Unicode rune. |
| <a id="testo009"></a>[`TESTO009`](rules/TESTO009.md) | `unexported-param-field` | correctness | enabled | No | Parameter fields must be exported so reflection can set them. |
| <a id="testo010"></a>[`TESTO010`](rules/TESTO010.md) | `empty-cases` | correctness | enabled | No | Providers returning `nil`, an empty slice literal or built-in `make` with constant zero length always return an empty case set. |
| <a id="testo011"></a>[`TESTO011`](rules/TESTO011.md) | `empty-suite` | correctness | enabled | No | Suite must contain at least one runnable Testo test. |
| <a id="testo012"></a>[`TESTO012`](rules/TESTO012.md) | `nil-suite` | policy | enabled | No | Report a statically nil suite argument to `RunSuite` or `RunSubSuite`. |
| <a id="testo013"></a>[`TESTO013`](rules/TESTO013.md) | `recursive-subsuite` | correctness | enabled | No | Report a directly recursive `RunSubSuite` invocation on the current receiver. |
| <a id="testo014"></a>[`TESTO014`](rules/TESTO014.md) | `shared-before-each-write` | policy | enabled | No | BeforeEach writes a shared receiver field read after Parallel suspension. |
| <a id="testo015"></a>[`TESTO015`](rules/TESTO015.md) | `shared-fixture-cleanup` | policy | enabled | No | Cleanup re-reads a shared BeforeEach field that another execution can replace. |
| <a id="testo016"></a>[`TESTO016`](rules/TESTO016.md) | `shared-test-context` | policy | enabled | No | A shared BeforeEach field can expose another execution's test context. |
| <a id="testo017"></a>[`TESTO017`](rules/TESTO017.md) | `duplicate-parallel-lifecycle` | correctness | enabled | No | Repeated Parallel calls on the same context across BeforeEach, test and AfterEach. |
| <a id="testo018"></a>[`TESTO018`](rules/TESTO018.md) | `run-during-cleanup` | correctness | enabled | No | Direct `testo.Run` on the cleanup context inside its Cleanup callback. |
| <a id="testo019"></a>[`TESTO019`](rules/TESTO019.md) | `nil-callback` | correctness | enabled | No | Statically nil callbacks to `RunTest` or an immediately executed `Test` adapter. |
| <a id="testo020"></a>[`TESTO020`](rules/TESTO020.md) | `premature-teardown` | correctness | enabled | No | AfterEach closes an `*os.File` before a parallel child uses the same file. |
| <a id="testo021"></a>[`TESTO021`](rules/TESTO021.md) | `construction-field` | correctness | enabled | No | Exported fields traversed during custom T/plugin wiring must be pointers. |
| <a id="testo022"></a>[`TESTO022`](rules/TESTO022.md) | `construction-cycle` | correctness | enabled | No | Report a type cycle traversed during custom T/plugin initialization. |
| <a id="testo023"></a>[`TESTO023`](rules/TESTO023.md) | `provider-before-each` | correctness | enabled | No | A provider dereferences a fresh suite’s nil field before BeforeEach initializes it. |
| <a id="testo024"></a>[`TESTO024`](rules/TESTO024.md) | `plugin-parent-nil` | correctness | enabled | No | A Plugin parent interface compared with nil cannot identify the suite root. |
| <a id="testo025"></a>[`TESTO025`](rules/TESTO025.md) | `blank-after-header` | style | disabled | Yes, where safe | Separate initial Allure metadata from the test body with a blank line. |
| <a id="testo026"></a>[`TESTO026`](rules/TESTO026.md) | `blank-between-steps` | style | disabled | Yes, where safe | Separate adjacent Allure Step statements with a blank line. |
| <a id="testo027"></a>[`TESTO027`](rules/TESTO027.md) | `no-sleep` | policy | enabled | No | Avoid time.Sleep in Testo tests and hooks; use polling or condition-based waiting. |
| <a id="testo028"></a>[`TESTO028`](rules/TESTO028.md) | `required-test-header` | style | disabled | No | Required configured metadata, including effective BeforeEach. |
| <a id="testo029"></a>[`TESTO029`](rules/TESTO029.md) | `unreported-action` | style | disabled | No | Coverage of configured actions by Step/Setup. |
| <a id="testo031"></a>[`TESTO031`](rules/TESTO031.md) | `assert-inside-step` | style | disabled | No | Coverage of configured assertions by Step/Setup. |
| <a id="testo030"></a>[`TESTO030`](rules/TESTO030.md) | `test-phase-order` | style | disabled | No | Ordered Setup, configured actions and assertions. |
