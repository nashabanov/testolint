# Changelog

## 0.7.0

### Added

- TESTO023: report a used provider dereferencing a fresh suite’s nil pointer
  field before its direct BeforeEach allocation; require a default, zero-safe
  BeforeAll path, one provider and valid collection signatures. Preserve dynamic providers,
  constructor initialization, earlier provider setup, unknown instance history
  and BeforeAll setup.
- TESTO024: report direct nil comparisons on unchanged Plugin parent interfaces
  in local implementations reached by valid Testo construction roots. Support
  concrete generics and pointer/value implementations; skip explicit local
  method references, reassigned/address-taken parameters and unknown graphs.

### Changed

- Extend TESTO010 to built-in make of slices with constant zero length and an
  omitted or constant nonnegative capacity. Share proven emptiness with parallel
  lifecycle multiplicity; preserve silence for unknown lengths/capacities.
- Link run calls to exact local suite models without broadening suite discovery.

### Tests and documentation

- Add provider_before_each, plugin_parent_nil and empty_make_shadow analyzer
  fixtures; extend empty_cases and parallel_coverage with regression boundaries.
- Add real-Testo provider panic, constructor/BeforeAll initialization, typed-nil
  root/child parent and intentional future global Options reproductions. Extend
  CLI/plugin assertions to 29 diagnostics with exact positions and clean cases.
- Document supported patterns, remedies and limitations; defer the ambiguous
  global-options policy without assigning it a diagnostic ID.
- Remove LateAnnotation and StaticCasesOnly from the roadmap. Number the new,
  unpublished rules consecutively as TESTO023/024; leave deferred ideas unnumbered
  and preserve all published diagnostic IDs.

## 0.6.1

### Refactoring

- Make analyzer models, rule types, fields, methods and run-kind constants
  private; keep `analyzer.Analyzer` as the package's only exported declaration.
- Consolidate shared Testo type, naming and signature semantics in
  `testo_semantics.go` and call facts and interpretation in `calls.go`. Keep
  individual rule files separate.
- Keep suite discovery models independent of parallel-analysis state; build one
  separate lifecycle model per suite and pass it directly to parallel rules.
- Register package-wide execution and construction checks alongside suite and
  run checks, sharing invocation-local call facts and diagnostic deduplication.
  Preserve diagnostic messages, positions and default reporting order.
- Pass all rule groups explicitly through `ruleSet`; build the parallel model
  once per suite only when parallel checks are selected. Keep final suite checks
  after parallel diagnostics to preserve the existing reporting order.

### Tests

- Extend the run_pipeline analyzer fixture with parallel and standalone checks,
  safe examples, teardown diagnostics, phase-order assertions, independent
  package/parallel selection, omitted groups, duplicate registration and
  repeated-pass checks.

## 0.6.0

### Added

- TESTO021: report exported non-pointer fields reached by custom T/plugin
  wiring from typed Testo execution calls, including nested and imported types.
- TESTO022: report active-path type cycles traversed by context initialization;
  ignore inaccessible fields and container contents, and allow repeated types
  on independent branches. Cycles suppress dependent field diagnostics.
- Share construction-root discovery and type traversal for both checks, with
  aliases, pointer/value contexts, concrete generics and local function aliases.
  Skip unknown traversed type parameters and unexecuted Test adapters.

### Tests and documentation

- Add construction analyzer fixtures, generic/imported regression boundaries,
  CLI/plugin diagnostic count and position assertions, and real-Testo runtime
  failures plus safe construction cases.
- Document construction semantics, diagnostic precedence and supported roots.
- Defer experimental TESTO023 after feasibility review; it is outside v0.6.0.
- Limit the roadmap to upcoming v0.7.0–v1.0.0 milestones and deferred work;
  keep completed implementation and release history in the changelog.

## 0.5.0

### Added

- TESTO018: detect direct testo.Run calls on the same unchanged standard context
  inside a registered Cleanup closure.
- TESTO019: detect statically nil callbacks to RunTest and immediately executed
  Test adapters, including direct testing.T.Run arguments. Preserve Run's valid
  nil-as-no-op behavior; do not flag unused/stored Test adapters.
- TESTO020: detect an os.File closed in AfterEach before a direct parallel child
  calls Read, Write or Stat on the same field. Restrict analysis to one regular
  test, default surrounding hooks and complete direct bodies; recommend Cleanup.

### Refactoring

- Share typed call facts, unchanged-local resolution and assignment/address
  invalidation across suite execution discovery and standalone execution checks.
  Keep diagnostic deduplication in the existing analyzer pipeline.
- Centralize suite-run rule registration in runRules.

### Tests and documentation

- Add execution_safety analyzer fixtures, real-Testo runtime reproductions,
  safe examples and CLI/custom-plugin assertions for IDs and source positions.
- Document supported execution patterns, unknown cases and runtime assumptions.
- Add premature_teardown analyzer fixtures, real-Testo closed-file reproduction
  and safe serial/Cleanup runtime checks. Verify os.ErrClosed for Read, Write
  and Stat, plus CLI/plugin diagnostic position and count.
- Update the roadmap for the v0.5.0 implementation.

## 0.4.2

### Fixed

- TESTO017 retains duplicates established before an unsupported AfterEach,
  while treating unknown lifecycle stages as barriers to later diagnostics.
  Lifecycle signature validation remains independent of body support.
- TESTO014–016 use a later execution's understood BeforeEach even when its test
  body is unsupported and does not expose the suite receiver or test context.
  Unknown receiver/context access and known test/AfterEach field replacements
  still suppress unsafe proofs.

### Tests and documentation

- Extend existing analyzer and real-Testo CLI/plugin fixtures, assert diagnostic
  positions and repeated CLI output, and reproduce opaque sibling lifecycle
  ordering and duplicate Parallel with opaque teardown at runtime.
- Document lifecycle-stage boundaries and receiver-ownership assumptions;
  correct the roadmap release baseline and the obsolete review-file reference.

## 0.4.1

### Fixed

- TESTO014–016 recognize later serial executions as possible BeforeEach
  overwriters, respecting Testo's regular-before-parameterized collection order.
- Unknown sibling provider counts no longer hide independently proven shared
  field hazards; distinguish proven empty paths from unknown counts and inspect
  each provider size once per suite.
- Validate BeforeAll/AfterAll signatures without parsing unrelated hook bodies.
- TESTO017 analyzes known test paths independently of unsupported siblings;
  unsupported operations on the target path still suppress the diagnostic.
- Recognize resolved standard context Helper and Logf calls in local timelines.

### Tests and documentation

- Add analyzer coverage regressions, real-Testo serial-order and Cartesian-case
  runtime reproductions, and CLI/custom-plugin expectations.
- Document partial analysis, ordering evidence and retained conservative limits
  in README.md.

## 0.4.0

### Added

- `TESTO014`: shared BeforeEach field reads after proven Parallel suspension.
- `TESTO015`: Cleanup callbacks re-reading replaceable shared fixture fields,
  with precedence over TESTO014 for the same field.
- `TESTO016`: shared BeforeEach storage exposing another test context.
- `TESTO017`: repeated Parallel calls on a single lifecycle/test context.
- Conservative local timelines, resolved context identity, direct receiver-field
  identity and constant literal provider multiplicity for parallel checks.
- Fixture-first regression coverage, real-Testo runtime reproductions, CLI and
  custom golangci-lint coverage, and executable safe serial/local-capture cases.

### Changed

- Split TESTO014–017 into four separately registered rules, backed by a shared
  model built once per suite. Separate AST parsing, provider counts and
  shared-field analysis from diagnostic reporting; remove nested parser closures.

### Fixed

- Unknown sibling case counts no longer suppress TESTO017 on a test with its
  own established execution. Preserve suppression for collection errors and
  unknown or empty cases of the checked test.
- Disable golangci-lint's identical-message cap in the integration harness so
  it verifies every expected diagnostic, including four TESTO017 reports.

### Documentation

- Document supported scheduling patterns, remedies and precision boundaries.
- Include the version roadmap and a pre-commit 0.4.0 review report.

## 0.3.0

### Added

- Discover `RunSuite` and `RunSubSuite` using Go types, including generic calls,
  import aliases and unchanged direct local function aliases.
- `TESTO012`: report explicit nil and unchanged local nil suite pointers.
- `TESTO013`: report direct receiver recursion in `BeforeAll` or a sole regular
  test with default preceding hooks; skip execution paths requiring inference.
- Regression fixtures for discovery, both rules, false positives, diagnostic
  positions, deduplication and deterministic execution alongside existing rules.
- CLI and custom golangci-lint integration coverage for both new rules and
  clean execution with real Testo, including bounded receiver recursion.
- Document detection boundaries and Testo v1.8.0 runtime behavior.

### Changed

- Add separate run-level rules with shared discovery and diagnostic deduplication,
  keep analysis state local, and modernize loops and test error handling.
- Use the latest published release in installation examples and record pre-release checks.

## 0.2.0

### Added

- `TESTO010`: report providers whose single return statement returns `nil` or an
  empty slice literal, including named slices and aliases. Local promoted
  providers are checked; imported bodies and more complex provider logic are skipped.
- `TESTO011`: report suites with no tests recognized by Testo's naming rules,
  including suites with only hooks or providers.
- Discover suites that indirectly embed `testo.Suite[T]`.
- Check promoted and imported tests, providers, and hooks using Go's pointer
  method set, respecting shadowing and ambiguous methods. Imported method
  diagnostics point to the local suite declaration.
- Repository guidelines for fixture-first TDD, pre-release review, and changelog
  updates before requesting commit approval.

### Fixed

- Discover package-level suite types even when they declare no methods.
  Type aliases do not create duplicate suites.

### Changed

- Share provider signature validation between `TESTO004`, `TESTO005`, and
  `TESTO010`, preserving diagnostic messages and behavior.
- Update README installation examples to target v0.2.0 once its tag is published.

## 0.1.0

Initial release.

### Added

- Standalone `testolint` CLI and reusable `go/analysis` analyzer.
- golangci-lint v2 module plugin using the same analyzer and rules as the CLI.
- `TESTO001` and `TESTO002`: validate test and lifecycle hook signatures.
- `TESTO003`, `TESTO004`, and `TESTO005`: validate required field-name providers,
  assignable slice element types, and provider signatures.
- `TESTO006`: report providers unused by test parameters as a linter policy.
- `TESTO007` and `TESTO008`: validate Testo test and provider naming, including
  Unicode suffixes; ignore the method named exactly `Cases`.
- `TESTO009`: report unexported test parameter fields.
- Direct `testo.Suite[T]` discovery, including aliases and pointer embedding,
  with checks for declared pointer and value receiver methods.
- Support named parameter structs, named provider slices, and type aliases.
- Suppress dependent diagnostics for invalid signatures and unexported fields;
  deduplicate identical diagnostics at the same source position.
- Analyzer fixtures, CLI exit-code tests, real-Testo integration fixtures,
  custom golangci-lint integration tests, and CI checks.
- Installation, usage, compatibility, and limitations documentation, with runtime
  semantics checked against Testo v1.8.0.
