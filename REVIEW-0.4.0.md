# v0.4.0 pre-commit review

Status: implemented, unreleased, awaiting user review. No commit, tag or
publication has been created. Review date: 2026-10-17.

## Scope and baseline

- Target: v0.4.0, compared with v0.3.0.
- HEAD: `8261c66ba25e00bdc16df8334e844ec5a7769951`; changes below are uncommitted.
  HEAD differs from v0.3.0 by removal of generated review documents.
- The incoming tree already contained staged `ROADMAP.md`. Its staged contents
  were preserved; its working copy now records implementation status.
- Added TESTO014–017 through one shared local lifecycle timeline; existing
  discovery, diagnostic IDs, CLI implementation and plugin adapter are unchanged.
- Changed: rule registration, per-suite model preparation, analyzer fixtures and Testo stub,
  CLI/plugin expectations, README, CHANGELOG and ROADMAP.
- Added: `analyzer/parallel_{body,model,fields,cases}.go`, four separately
  registered `rule_shared_*.go` / `rule_duplicate_parallel_lifecycle.go` files,
  analyzer fixtures `parallel_safety` and `duplicate_parallel`,
  `plugin/golangci/runtime_test.go`, integration `parallel_test.go` and
  `duplicate_parallel_test.go`, clean `parallel_test.go`, runtime reproduction
  under integration `testdata/parallelruntime/runtime_test.go`, and this report.
  The original combined `rule_parallel_safety.go` was removed during refactoring.
- Go: `go1.27.1 darwin/arm64`.
- System golangci-lint: `2.14.0`, built with `go1.27.0`, commit `114493f9`.
- Rebuilt plugin: `v2.14.0-custom-gcl-2lQx9obpTZpdGFmqqUBOCX73dI38zHDOjDA2tk`,
  built with `go1.27.1`, timestamp `2026-10-17 01:07:43.813661 +0000 UTC`.
- Actual integration dependency: `github.com/ozontech/testo v1.8.0`.

## Refactoring requested during review

The first implementation combined AST parsing, scheduling facts and four
policies in large functions. It has been replaced with these responsibilities:

- `parallel_body.go`: a local parser with explicit state and small methods for
  signatures, expressions, assignments, context calls and Cleanup. No nested
  parser closures; events have a kind instead of independent write/parallel flags.
- `parallel_model.go`: build each suite's hook/test paths once, before running
  rules. No global cache, repeated body parsing or diagnostic reporting.
- `parallel_cases.go`: collection validation and case multiplicity are separate.
  Invalid collection and unknown counts no longer share one rejection condition.
- `parallel_fields.go`: compute overwrite/read facts once for the three field
  rules; retain the scheduling constraints and exclusion of extra writes.
- Four rule files, each registered in `rules.go`: each emits exactly its own ID.
  TESTO014/015 explicitly yield to the more specific cleanup/context diagnostics.
  TESTO017 scans context calls without depending on shared-field analysis.
- The largest function in these eight implementation files is 28 lines;
  responsibilities are explicit without adding a generic visitor/framework.

The structural refactor first passed the existing targeted fixture and full
uncached suite without changing expected diagnostics. A separate fixture-first
change then removed TESTO017's dependence on sibling counts; this is the only
intended detection change in the refactoring pass. Parser support and existing
messages/positions remain unchanged. Reporting now follows rule registration
order, with existing pass-wide deduplication and deterministic execution.

## Findings, in severity order

No unresolved release-blocking defect was found in the supported subset.
This is a local source/fixture review, not a repository corpus audit.

### Resolved correctness findings

1. Test-body writes were initially attributed to BeforeEach; later writes could
   also invalidate the modeled cleanup resource. Regression examples OwnWrite,
   Replaced and CleanupReplaced in `analyzer/testdata/src/parallel_safety/test.go`
   produced unexpected TESTO014/015 diagnostics before fixes. The implementation
   now tracks BeforeEach provenance and excludes fields written by a test or
   AfterEach (`analyzer/parallel_fields.go`, replaced-field map and event
   processing). This intentionally sacrifices coverage for precision.
2. Invalid BeforeAll signatures initially allowed cascading concurrency
   diagnostics. InvalidHook reproduced the unwanted TESTO014 alongside TESTO002.
   All lifecycle hook timelines/signatures now gate the new checks (`parallel_body.go:21`).
3. `new(expr)` can invoke a helper in modern Go. NewExpression reproduced an
   unjustified TESTO014 for an unknown initializer. Only type arguments or
   constant initializers are now accepted by the timeline's builtin handling.
4. Cleanup and normal reads initially emitted TESTO015 plus TESTO014 for the
   same field/cause. Cleanup's negative normal-read expectation reproduced this;
   cleanup takes precedence. Existing pass-wide deduplication handles repeated
   hook positions without new global state.
5. Unknown sibling case counts originally suppressed TESTO017 for regular tests
   and tests with their own known cases. `duplicate_parallel` reproduced two
   missing TESTO017 diagnostics. The rule now checks each path's own count;
   collection validation still prevents secondary messages for missing/type-
   incompatible providers and unexported parameter fields.

### Coverage boundaries / non-blocking follow-up

- `analyzer/parallel_body.go:21` deliberately uses a strict local
  parser. Unknown calls, branching, early exits, aliases and imported/promoted
  bodies skip the suite. This misses common helper-based tests; it is documented
  rather than treated as proof of safety. Future expansion needs separate
  precision evidence and fixture-first changes.
- `analyzer/parallel_cases.go:59` only counts non-empty, unkeyed constant
  literal provider returns. Dynamic counts, empty sets, runtime filtering,
  annotations and plugin plans are not inferred. TESTO017 can inspect known
  test paths even with unknown sibling counts; TESTO014–016 remain conservative
  for the whole suite. Constant multiplicity has
  analyzer and real-Testo diagnostic coverage, but the deterministic runtime
  overwrite reproduction uses two regular methods rather than a parameterized
  method. A dedicated parameterized runtime reproduction is future coverage.
- Plugin overrides/synchronization, custom T, shared referenced data in value
  receivers and context stores inside tests are excluded. No helper/SSA/race
  analysis or platform compatibility expansion was attempted.

## Testo source comparison

Read actual files from the pinned module cache, not only the analyzer stub:

- [suite.go, v1.8.0](https://github.com/ozontech/testo/blob/v1.8.0/suite.go):
  the private marker and generic T support existing discovery. BeforeEach,
  test and AfterEach share a test context; BeforeAll/AfterAll use suite context.
  Default hooks do not add Parallel. Value receivers copy their storage; new
  field diagnostics require the suite's own pointer receiver.
- [runner.go, v1.8.0](https://github.com/ozontech/testo/blob/v1.8.0/runner.go#L343):
  the collected plan runs through sibling testing.T.Run calls; each constructs
  its own T and passes the same suite value. BeforeEach precedes test.Run and
  AfterEach is deferred. Parallel suspension permits another BeforeEach to
  replace a shared field before the first test resumes. Plugins/filtering can
  alter execution and remain an explicit assumption.
- [t.go, v1.8.0](https://github.com/ozontech/testo/blob/v1.8.0/t.go#L141):
  Parallel resolves on actual *testo.T and normally delegates to testing;
  [Cleanup](https://github.com/ozontech/testo/blob/v1.8.0/t.go#L451)
  executes after the test and its subtests. Both support runtime overrides,
  which cannot be inferred statically.
- [collector.go, v1.8.0](https://github.com/ozontech/testo/blob/v1.8.0/collector.go):
  Go method sets, Test/Cases suffix rules, exact T signatures, field-name provider
  matching, slice element assignability and reflected field setting remain the
  baseline for existing rules. Providers execute after BeforeAll; their case
  permutations create multiple independently constructed test contexts. Empty
  provider/suite conditions are framework warnings (potential strict-mode
  failures); unused providers remain linter policy. New rules do not change
  discovery or these existing diagnostics.

TESTO014–016 identify a Testo lifecycle/ownership hazard, not a guaranteed race
or runtime failure. The deterministic reproduction confirms replacement of both
resource and T and a cleanup re-read, without requiring a data race. TESTO017's
second Parallel call produces an actual runtime test failure with default Testo.
Runtime harness checks that failure deliberately in a subprocess; it is not a
failed verification check.

Reviewed discovery/method sets, rule gating, positions and deduplication, CLI
singlechecker, plugin registration/load mode and CI. Existing alias, promotion,
shadowing, ambiguity, signature, unrelated-type and deterministic pipeline
fixtures remain green. New checks use field objects and resolved methods,
never names alone to establish Testo context identity. No similar-looking
existing discovery/rule code was refactored merely for release cleanup.

## Fixture-first evidence

Fixture: `analyzer/testdata/src/parallel_safety/test.go`, registered in
`TestAnalyzer`. It includes positives plus serial, single-execution, local
capture, pointer/value, alias, unrelated-context, branching, unknown helper,
replacement and invalid-hook boundaries.

Command for every red/green cycle:

```sh
go test ./analyzer -run '^TestAnalyzer/parallel_safety$' -count=1
```

The refactoring regression used `^TestAnalyzer/duplicate_parallel$`; the final
targeted run selected both fixtures with `^TestAnalyzer/(parallel_safety|duplicate_parallel)$`.

- Initial red: five missing diagnostics covering all four new IDs, after the
  fixture loaded and type-checked. Initial sandbox/toolchain/cache access errors
  were environmental and were not counted as red. Retried with approved access.
- OwnWrite/Replaced/Parameterized red: two unexpected TESTO014 diagnostics plus
  a missing parameterized TESTO014; green after provenance/count changes.
- Cleanup precedence red: unexpected TESTO014; green after precedence handling.
- CleanupReplaced red: unexpected TESTO015; green after excluding extra writes.
- InvalidHook red: unexpected TESTO014; green after hook gating.
- FixtureCleanup red: missing TESTO015; green after terminal resource-method
  receiver support. Unknown method bodies do not contribute further scheduling
  evidence; only the receiver read is used.
- NewExpression red: unexpected TESTO014; green after restricting new arguments.
  ValueParallel also verifies the same-context rule with value receivers.
- Refactoring regression `duplicate_parallel`: two missing TESTO017 diagnostics
  on a regular method and a known-case method with dynamic siblings. Fixtures
  loaded/type-checked. Negative cases cover one Parallel, unknown/empty own
  cases, missing/incompatible providers and private parameter fields. Green
  after per-path applicability and separate collection validation; both targeted
  fixtures and the full suite passed.
- Final targeted run: PASS. No diagnostic behavior change preceded its red fixture; structural refactoring
  was performed with green tests.
  An initial implementation compilation typo was corrected; compilation failures
  were not counted as diagnostic red evidence.

## Verification

Final code checks:

| Command | Outcome |
| --- | --- |
| `gofmt` on modified Go files | Applied |
| Targeted analyzer fixture, uncached | PASS |
| `make test` | PASS |
| `go test ./... -count=1` | PASS |
| `TESTOLINT_GOLANGCI_BINARY="$PWD/bin/golangci-lint" go test ./... -count=1` | PASS, including external integration |
| `go vet ./...` | PASS |
| `make lint` | PASS, 0 issues |
| `go mod tidy -diff` | PASS, no changes |
| `go mod verify` | PASS, all modules verified |
| `bin/bootstrap/golangci-lint custom -v` | PASS, rebuilt after final analyzer edit |
| `TESTOLINT_GOLANGCI_BINARY="$PWD/bin/golangci-lint" go test ./plugin/golangci -run '^TestIntegration$' -count=1 -v` | PASS with final rebuilt binary |
| `go test ./clean -count=1` in integration module | PASS, actual Testo serial/captured parallel executions |
| `go test ./plugin/golangci -run '^TestParallelRuntime$' -count=1 -v` | PASS; also executed uncached through the full suite |
| `git diff --check` | PASS |

The ordinary full suite skips TestIntegration without its binary variable; the
separate explicit run above passed and is the external integration evidence.
CLI now checks 16 total diagnostics; the plugin verifies per-ID counts,
including TESTO014=3, TESTO015=2, TESTO016=1, TESTO017=4 and 16 testolint outputs.

Intermediate failures/retries: lint reported QF1001 (boolean expression), fixed
and retried successfully. First plugin integration counted `want` comments in
rendered source snippets twice; a copied value-receiver example also triggered
staticcheck SA4005. Removed analyzer-specific expectation comments and that
redundant real-Testo negative example; retained its analyzer fixture. Plugin
integration then passed, including a run after the refactoring code rebuild.
That refactoring integration run first exposed golangci-lint's default limit of
three identical messages: TESTO017=3 instead of 4, although the CLI reported all
16 diagnostics. Its local `run --help` confirmed `--max-same-issues` defaults to
3. The integration harness now passes `--max-same-issues=0`; the repeated run
passed with all 16 diagnostics. Production linter configuration is unchanged.
During the refactoring baseline, a sandboxed full run failed CLI package loading
with "matched no packages". The identical uncached command outside the sandbox
passed before any implementation edits; the environment failure was not counted
as a TDD red phase. No essential check remains unavailable or only skipped.

## Release conditions

README documents 17 enabled rules, actual tested versions and detection limits.
Installation still uses the published `latest`, not an unpublished v0.4.0 tag.
CHANGELOG records 0.4.0 as Unreleased without a release date. CI already builds
and exercises the custom plugin and clean real-Testo module; no CI change was
needed. ROADMAP distinguishes implemented scope from broader candidates.

Ready for user review of this working-tree implementation. Before committing,
obtain explicit approval of code, fixtures and documentation. Before publishing,
review the resulting commit and authorize the release/tag separately. Broader
compatibility and arbitrary concurrent Go code are not claimed.
