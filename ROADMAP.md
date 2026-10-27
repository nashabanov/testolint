# Roadmap

## Overview and project principles

testolint is a specialized static analyzer for
[Testo](https://github.com/ozontech/testo). This roadmap covers upcoming work from v0.8.0 to v1.0.0
for a project maintained by one developer. Versions
express an intended order, without dates or a promise to ship every candidate.
Small, useful releases take precedence over completing a rule quota.

Priorities are **P0** (correctness foundation or release blocker), **P1** (next
useful capability), and **P2** (optional improvement). **Planned** means intended
work subject to acceptance; **Experimental** requires a viable detection
strategy; **Deferred** has no delivery commitment. Required work, optional
enhancements, and research are distinguished below. Experimental candidates do
not block releases and may move to the backlog.

The project prioritizes Testo runtime failures, lifecycle misuse, test isolation,
reflection initialization, parameterization, annotations, and execution APIs.
Rules must earn their maintenance cost through Testo-specific knowledge.
Do not duplicate go vet, staticcheck, tparallel, paralleltest, thelper, standard
race detection, or generic style and naming checks. Existing Testo naming checks
remain relevant because the framework uses those names for discovery.

Prefer conservative AST and Go type analysis, explicit supported patterns, and
silence for unknown cases. Correctness matters more than rule count; complex
analysis cannot promise zero false positives. Extend an existing rule before
adding an overlapping one. Add model fields or discovery support only when a
concrete rule needs them; SSA, interprocedural analysis, and a general data-flow
framework require evidence that simpler approaches cannot solve a valuable case.

Rule classification and enablement are separate decisions. **Correctness**
means framework misuse, runtime errors or correctness-related hazards;
**policy** means optional conventions, restrictions or quality recommendations.
Classification must reflect the diagnostic's actual supported semantics, not
its name. A correctness rule can still have conservative scope; policy does not
mean an existing rule is silently disabled. Preserve current defaults unless a
documented migration is justified. New optional policies start disabled and
require explicit opt-in. Test quality may be a domain/group, not a third rule
class or an implicit enablement setting.

## Planning baseline

Implemented rules and limitations are documented in the [README](README.md);
completed work and release history belong in the [changelog](CHANGELOG.md).
The unreleased v0.7.0 implementation scope is recorded in the changelog.
The baseline for future work is Go 1.27.1, Testo v1.8.0 and golangci-lint v2.14.0.
Broader compatibility requires explicit verification.

Use the pinned Testo source when evaluating candidates, especially
[collector.go](https://github.com/ozontech/testo/blob/v1.8.0/collector.go),
[runner.go](https://github.com/ozontech/testo/blob/v1.8.0/runner.go),
[construct.go](https://github.com/ozontech/testo/blob/v1.8.0/construct.go),
[the plugin contract](https://github.com/ozontech/testo/blob/v1.8.0/testoplugin/plugin.go),
[options.go](https://github.com/ozontech/testo/blob/v1.8.0/options.go) and
[annotations.go](https://github.com/ozontech/testo/blob/v1.8.0/annotations.go).
Runtime reproductions remain an acceptance prerequisite. Distinguish runtime
failures, framework warnings and optional linter policy.

## Release overview

The implemented baseline remains TESTO001–TESTO024. No existing diagnostic ID
is reassigned. Future candidates receive no permanent ID before acceptance for
implementation; published IDs must never be recycled. The table describes
required outcomes, not rule quotas or promises to implement research candidates.

| Version | Focus | Required work | Optional / experimental work |
| --- | --- | --- | --- |
| v0.8.0 | Rule Registry, Configuration & Diagnostic UX | Static registry, shared reporting, validated rule selection, rule documentation | Safe SuggestedFixes; source suppression only with demonstrated need |
| v0.9.0 | Hardening, Architecture & Compatibility | Corpus audits, measured pipeline review, benchmarks, compatibility matrix | Evidence-driven architecture changes and justified additional support |
| v0.10.0 | Test Quality Policies & Extensibility | Research, candidate selection and documented decisions | Accepted opt-in quality checks; small custom-rules feasibility study |
| v1.0.0 | Stable Release | Stability contract and release-readiness evidence | No mandatory new rules; experimental work does not block release |

## Detailed version milestones

### v0.8.0 — Rule Registry, Configuration & Diagnostic UX

**P1 · Planned.** Standardize registration, reporting and configuration within
the existing analyzer. This is an incremental migration, not a complete rewrite
or a new rule quota. Preserve rule functions, execution groups and their order.

#### Static rule registry — required

- Keep the design inside the existing `analyzer` package. Add `registry.go` for
  a centralized, static registry of rule metadata and registration. No `Rule`
  interface, runtime registration, `init()` rule registration or abstraction
  framework. The existing golangci-lint adapter registration is separate from
  this internal registry and need not change merely to satisfy that constraint.
- Use the diagnostic ID as the lookup key. A small metadata model is sufficient:

  ```go
  type ruleMeta struct {
      Name   string
      Class  ruleClass
      DocURL string
  }
  ```

  Keep each rule's metadata in one authoritative source. Do not duplicate rule
  lists in configuration, documentation generators or the two integrations.
  Document defaults alongside the static definitions, adding a small explicit
  default-enablement field only if needed; do not infer defaults from Class.
- Preserve suite, parallel, final-suite, run and package execution groups and
  their existing rule functions. A checker can emit multiple IDs, as execution
  safety and construction currently do. Associate each emitted ID with metadata
  and filter individual diagnostics correctly without running the same checker
  twice. Shared prerequisites may run for selected checks without enabling their
  own disabled diagnostics. Do not force one function per ID as a prerequisite.
- Registry metadata describes the rule; diagnostic message templates stay with
  the rule implementation. Registration order must be deterministic and explicit.

#### Diagnostic reporting — required

Add `diagnostic.go` with shared helpers, conceptually:

```go
reportf(p, id, pos, format, args...)
report(p, id, diagnostic)
```

These are planned internal APIs; choose concrete Go signatures during
implementation. The ID comes from the registry, rather than being repeated in
message format strings. Standardize:

- Exactly one `TESTO###` prefix, a documented `Category` convention and rule
  documentation URLs resolved from registry metadata. Publish and test a URL
  scheme; do not invent links to pages that do not exist.
- Clear wording and actionable remediation that respect uncertainty. Establish
  a diagnostic baseline before migration. Preserve detection, severity of
  claims, suppression and defaults; intentional wording or position improvements
  need reviewed before/after fixtures and documentation.
- Consistent primary positions and useful ranges. Preserve caller-supplied
  `End`, `Related` and `SuggestedFixes` when reporting a full diagnostic; a helper
  must not replace it with a reduced message/position pair.
- Deterministic deduplication and cascading-diagnostic precedence. Account for
  related locations/fixes when designing deduplication so useful context is not
  silently discarded. Preserve established diagnostic ownership and reporting
  order unless a specific reviewed correction is justified.

Prioritize useful related source locations for TESTO004 (provider versus
parameter field), TESTO014–016 (initialization versus hazardous use), TESTO020
(teardown versus child use) and TESTO022 (cycle edges). Choose only locations
that actually clarify the proof; imported or unavailable bodies must not produce
fabricated ranges. `Related` and SuggestedFixes are not mandatory on every
report. SuggestedFixes are **P2 · Optional** and require a uniquely safe
transformation plus fix-output tests; lifecycle ownership changes usually need
human judgment.

#### Rule configuration — required

- Enable/disable individual IDs using the registry. Reject unknown IDs, malformed
  settings and invalid combinations with actionable errors before analysis.
- Document a small common precedence model: start from registry defaults, then
  apply the explicitly configured selections according to a defined conflict
  policy. Specify whether overlapping enable/disable lists are rejected; do not
  let CLI flag order or plugin map iteration accidentally decide the result.
  Document how driver-specific settings map to this common model.
- Expose equivalent settings through standalone go/analysis CLI flags and
  golangci-lint configuration. Verify default, valid and invalid examples through
  both real integration boundaries, including explicit opt-in policy selection.
- Keep configuration and mutable state isolated across analyzer instances and
  repeated runs. Avoid attaching process-global mutable selections to the
  exported Analyzer singleton. Introduce only the smallest construction boundary
  needed by the drivers, with tests for differently configured instances.
- Preserve current defaults, including existing policy behavior. TESTO006 is
  policy and TESTO012 includes intentional nil-safe suites; classify the complete
  implemented rule set based on actual semantics. Class alone does not change
  enablement. Any justified default change requires explicit migration guidance.

**P2 · Optional:** evaluate source suppression only if a demonstrated need
remains after rule selection. No custom DSL or general configuration framework.

#### Documentation and acceptance

**Required documentation:** every supported ID, name, correctness/policy class,
default enablement, violation and valid examples, limits and known uncertainty,
remediation, diagnostic links, and configuration examples for both integrations.
Use the registry as the metadata authority and automatically check documentation
coverage/consistency; prose examples and explanations may remain hand-written.

**Definition of Done:**

- Every implemented rule has one unique valid metadata entry and a documented
  default. Automated tests verify metadata/registration consistency, execution
  group coverage, duplicate/unknown IDs and required documentation links.
- All rule diagnostics use shared reporting. Unit tests preserve full diagnostic
  fields; end-to-end regressions preserve IDs, intentional messages/positions,
  suppression and deterministic counts without loss or duplication.
- Existing IDs remain stable. Any deliberate UX change is separately reviewed;
  registry/reporting migration introduces no unintended detection changes.
- Configuration tests cover defaults, empty selections, individual IDs emitted
  by shared checkers, invalid input, conflicts, explicit opt-in and instance
  isolation. CLI and a freshly rebuilt golangci-lint plugin agree on selected IDs,
  default diagnostics, positions and clean fixtures; invalid settings fail clearly.
- Ship working examples and migration notes. Apply the common rule acceptance
  and fixture-first workflow to diagnostic changes; helper/registry unit tests
  supplement rather than replace end-to-end coverage.

### v0.9.0 — Hardening, Architecture & Compatibility

**P0 · Planned:** establish reproducible correctness, performance and compatibility
evidence. **P1 · Planned:** review the pipeline architecture; implement targeted
changes only where that evidence demonstrates a benefit. A mandatory review is
not a requirement to refactor every area.

#### Architecture — targeted review and conditional improvements

Evaluate separation of discovery, semantic models, rule execution and reporting;
reuse of collected facts; redundant AST traversals; execution order/dependencies;
model construction only when selected checks require it; diagnostic ownership
and precedence; and pipeline testability and performance.

Record a concrete problem before a structural change: a reproduction, a
maintenance defect or a measured time/allocation cost. Prefer a small internal
change with a before/after result. Verify preserved reporting order, filtering,
shared-checker behavior, isolation and diagnostic field preservation. For example,
measure repeated construction-graph work before deciding whether shared facts
would help. Reuse v0.8.0 metadata/reporting rather than building parallel systems.

No generic plugin architecture, SSA framework, interprocedural analysis engine or
broad package decomposition without a concrete use case and comparative evidence
that a smaller solution is insufficient. Changes to discovery or diagnostics
still require fixture-first red/green coverage; structural refactoring proceeds
only with green tests and established behavioral equivalence.

#### Hardening and compatibility — required

- Audit at least two real Testo repositories with different suite/provider/plugin
  patterns. Record revisions, configuration, diagnostics and manual decisions;
  investigate both false positives and false negatives. A clean run alone does
  not establish adequate coverage. Add regression fixtures for confirmed defects
  and resolve release-blocking correctness issues.
- Audit generic declarations versus instantiated calls, embedding, aliases,
  pointer/value sets, promoted/imported methods, shadowing, ambiguity, invalid
  signatures and cascading diagnostics. Existing support is a regression target;
  broader discovery remains conditional on real usage and precision evidence.
- Check deterministic diagnostics across repeated uncached runs, including IDs,
  primary/related positions, ranges, deduplication and configuration. Verify CLI
  and freshly rebuilt plugin parity on clean and diagnostic fixtures.
- Establish representative small/large benchmarks. Record hardware, tool and
  Testo versions, fixture size, time and allocations. Measure model construction
  with relevant rule groups disabled as well as enabled. Compare substantial
  changes with that baseline; investigate regressions and set a justified budget
  after measurement, without invented throughput promises.
- Publish and exercise an explicit Go/Testo/golangci-lint compatibility matrix.
  Cover at least the declared minimum Go and the current selected supported
  toolchain. Claim only framework/tool combinations actually verified by runtime
  and integration checks. Update compatibility documentation and CI together.

**P2 · Optional:** add supported versions/discovery patterns or optimize the
pipeline when real usage or measurements justify the maintenance cost.
**Experimental:** document uncovered semantics without converting every unknown
scenario into a release blocker.

**Definition of Done:** save reproducible audit reports with repository revisions,
triaged findings, checked missed cases, exact commands/outcomes/skips, benchmark
baselines and comparisons, and an exercised compatibility matrix. Document
architecture decisions, including decisions not to refactor. Resolve confirmed
release-blocking correctness defects; keep unsupported/experimental scenarios
explicit. Meet measured performance budgets or explain accepted tradeoffs.
Cached success and skipped external integration are insufficient evidence.
No promise of zero false positives or universal framework coverage is made.

### v0.10.0 — Test Quality Policies & Extensibility

**P1 · Planned research and selection; implementation conditional.** This
pre-stable milestone evaluates optional checks before v1.0.0; unfinished
research must not delay the stable release. No rule quota and no permanent IDs
are assigned to the candidates below. An evidence-backed decision to defer or
drop a candidate is a valid research outcome.

#### Test quality policies

| Candidate | Priority | Goal |
| --- | --- | --- |
| NoFixedSleep | P1 · Experimental | Detect potentially harmful fixed waits in Testo execution contexts. |
| UnboundedPolling | P1 · Experimental | Identify polling without established termination bounds where Testo semantics add useful evidence. |
| SleepInLifecycle | P2 · Experimental | Investigate timing-related hazards inside Testo lifecycle hooks. |
| GlobalStateMutation | P2 · Experimental | Identify risky shared global state changes during Testo execution. |
| TimingDependentAssertion | P2 · Experimental | Investigate timing-dependent assertions in Testo contexts. |

All accepted checks start **opt-in and disabled by default**, regardless of
correctness/policy classification. Most quality recommendations are policy;
classify a particular check as correctness only when supported evidence warrants
that claim. An optional quality domain/group must not replace the primary
classification or implicitly select rules.

Require substantial Testo-specific value: reuse suite, test, hook and lifecycle
discovery and establish why execution context changes the conclusion. Check
against nosleep, forbidigo, staticcheck and other generic Go linters. A blanket
ban on time.Sleep or generic loop/global-variable style checks does not qualify.
Extend existing testolint rules when the same hazard is already modeled.

Distinguish fixed delays, intentional scheduling, bounded polling, cancellation,
and legitimate integration-test behavior. A loop is not unsafe merely because
it polls; inability to prove a bound is not proof of an unbounded harmful wait.
A Sleep call is not automatically a correctness error. Global-state checks need
ownership, overlap and restoration evidence rather than a ban on all assignments;
timing-dependent assertions need a concrete proof pattern rather than names or
heuristics alone. Skip unsupported/unknown flow and validate useful detection
against actual Testo repositories before accepting a rule.

#### Custom rules API — P2 · Experimental research only

Explore whether additional statically registered Go rule functions can use the
existing registry/execution groups without modifying core execution logic for
each extension. Begin with a concrete downstream use case that built-in rules
and configuration cannot solve. Compare a small internal/static composition
experiment with simply contributing another built-in rule.

No DSL, dynamic plugin loading, general-purpose framework, mandatory public
extension API or stable external API commitment in v0.10.0. Research may conclude
that no implementation is justified. Any implementation requires a separately
reviewed scope, demonstrated benefit and alignment with the planned v1.0 contract;
it must not weaken the supported built-in pipeline or require package
reorganization solely for a hypothetical extension.

#### Research gates and acceptance

1. Validate each candidate against actual pinned Testo semantics and establish
   its correctness/policy class and relevant execution context.
2. Demonstrate useful positive cases and legitimate-use negatives in real Testo
   repositories, recording revisions and manual decisions.
3. Establish reasonable false-positive behavior, unknown cases and explicit
   precision limits; investigate missed cases without promising full enforcement.
4. Compare existing Go linters and testolint rules. Defer/drop candidates with
   no substantial Testo-specific value or unjustified maintenance complexity.
5. Select only candidates that pass those gates; assign an ID only at acceptance
   for implementation, preserving all existing IDs.
6. Implement selected checks using v0.8.0 configuration, registry and reporting,
   with no change to stable defaults. Verify explicit opt-in and CLI/plugin parity.
7. Apply the mandatory end-to-end red/green workflow, real-Testo reproductions,
   registry consistency tests and documentation for every accepted rule.

**Definition of Done:** publish research/selection decisions, overlap analysis,
repository evidence and documented limits. Any accepted implementation has
end-to-end fixtures, remediation guidance and tested disabled-by-default behavior.
Report the custom-rules feasibility decision separately. Completing every quality
candidate or shipping a public extension API is not a release requirement.

### v1.0.0 — Stable Release

**P0 · Planned.** Stabilize the supported product and publish its contract; no
mandatory new rules or feature work. Experimental/deferred candidates and the
v0.10.0 research program do not block v1.0.0.

Define and verify the stable contract for:

- Diagnostic IDs, rule names/metadata, correctness/policy classifications and
  supported semantics. Preserve IDs and do not recycle removed IDs; document
  rule retirement. Document whether a metadata/class change affects selection.
- Configuration keys, validation, precedence and defaults across both drivers.
  Breaking supported semantics or configuration in v1.x requires a major version;
  document bug fixes and additive coverage in release notes. New optional policy
  remains explicitly selected, not silently enabled.
- Diagnostic formatting: the documented ID prefix, Category convention, rule URL
  scheme and preservation of supplied structured diagnostic fields. Guarantee
  meaningful source ownership for supported patterns, not immutable byte offsets.
- Documentation URLs and a migration/redirect strategy for moved rule pages.
  URLs must resolve to documentation of the corresponding ID.
- Supported Go/Testo/golangci-lint versions, compatibility policy and the standalone
  and module-plugin invocation/configuration contracts.

Publish a clear distinction between stable guarantees and patch improvements.
Exact message prose, the number/order of explanatory related locations, and
narrower/more useful source ranges may improve in patch releases with regression
coverage and release notes; they must retain the ID, meaning and documented
structured format. Metadata descriptions/links may be clarified without changing
rule identity, documented classification/defaults or supported behavior.
Do not promise immutability of every message, internal function, model, registry
layout or traversal algorithm. Detection fixes can change affected diagnostics
within the documented supported scope; explain these as fixes, not silent
configuration migrations.

**Definition of Done:** publish the stability/migration contract, meet the full
readiness checklist below and present evidence for maintainer review. Complete
Testo API coverage is unnecessary. Do not tag or publish without explicit
permission, or describe an unpublished version as already installable.

## Deferred / experimental backlog

Items here have no release commitment. Ideas receive no diagnostic ID until
their semantics are established and they are accepted for implementation.
The v0.10.0 quality candidates and custom-rules research are described in that
milestone only; they are not prerequisites for the stable release.

| Idea | Priority / status | Reason deferred and next evidence needed |
| --- | --- | --- |
| GlobalOptionsInTest | P2 · Deferred policy | Testo intentionally appends global configuration. Runtime evidence covers deliberate registration in a test for a later RunSubSuite/RunTest, while existing context snapshots remain unchanged. A later run alone does not establish misuse; require a justified opt-in policy and rule-selection support. |
| Broader TESTO023 coverage | P2 · Deferred | Current proof requires a fresh empty suite literal, a direct pointer dereference and a sole BeforeEach allocation. Gather examples justifying constructor/history or helper analysis without flagging valid zero values or BeforeAll setup. |
| Broader TESTO010 empty expressions | P2 · Deferred | Unknown capacities and unchanged empty locals need independent precision/escape evidence and fixture-first acceptance. |
| UninitializedPluginField | P2 · Deferred | Skipped nil fields can be intentional, initialized manually or used by nil-safe methods. Gather real-project cases and establish a useful proof of missing initialization at use before adding a rule. |
| SharedCollectionMutation | P2 · Deferred | Map/slice mutation needs alias, ownership and overlap evidence; first establish whether TESTO014/015 can cover useful direct patterns. |
| ParallelSuiteMutation | P2 · Deferred | Broadly reporting receiver writes would resemble generic race detection and flag safe code; require a distinct Testo lifecycle failure beyond the existing shared-field checks. |
| ExcessiveCaseCombinations | P2 · Deferred | Cartesian growth can be intentional; gather real costs and define opt-in thresholds before creating policy. |
| IgnoredRunResult | P2 · Deferred | Testo returns success booleans, but failures already propagate through testing; ignoring a result is commonly valid. Require a concrete lost behavior beyond generic unchecked-result linting. |
| Broader RunSubSuite recursion detection | P2 · Deferred | TESTO013 already covers narrow direct recurrence. Wrappers, multiple tests, options, filters and plugins complicate reachability; extend TESTO013 only for proven additional paths. |
| Generic suite discovery improvements | P2 · Deferred | Generic declarations are currently skipped while instantiated run checks exist. Define concrete instantiation and diagnostic ownership requirements before broadening suite discovery; prioritize real repository examples. |
| Interprocedural analysis research | P2 · Deferred | Helpers and cross-package ownership may matter, but facts/SSA increase cost. Collect missed high-value bugs and compare simpler local strategies first. |

## Rule acceptance criteria

Follow [AGENTS.md](AGENTS.md) for every rule or discovery change:

1. Reproduce the behavior against actual pinned Testo, not only the simplified
   analyzer stub. Record versioned source evidence and whether the outcome is a
   runtime failure, warning, or policy concern.
2. Check overlap with existing Go linters and existing testolint rules. Prefer an
   extension over a duplicate; identify the Testo-specific value.
3. Define supported and unsupported scenarios, precision assumptions, priority,
   correctness/policy class, diagnostic ownership/positions/ranges, useful related
   locations and explicit default enablement before implementation. From v0.8.0,
   register metadata and documentation in the authoritative registry and use shared
   reporting/configuration. Class and enablement remain independent. Research
   candidates receive no permanent ID before acceptance for implementation.
4. Add or extend a type-checking end-to-end fixture under
   `analyzer/testdata/src/`, with `// want "TESTOxxx..."` expectations and valid
   examples without diagnostics. Register new fixture packages in TestAnalyzer.
   CLI/plugin behavior changes also need their corresponding fixtures first.
5. Run `go test ./analyzer -run '^TestAnalyzer/<fixture>$' -count=1` before
   changing implementation and observe the intended diagnostic failure.
   Compilation, loading and environment errors do not count as the red phase.
6. Implement the smallest sufficient analyzer change only after that red phase.
   Refactor only with green tests; unit tests supplement rather than replace
   analyzer fixtures.
7. Verify positives and negatives, absence of unexpected/duplicate diagnostics,
   and relevant boundaries: aliases, receivers, unrelated types, signatures,
   promoted methods and context identity as appropriate. Rerun the targeted test,
   gofmt modified Go files, then `make test`, `go vet ./...` and `make lint` when
   available.
8. Add real-Testo cases for the affected semantics and exercise CLI/plugin
   integration. Rebuild the custom binary from current sources; run
   `TESTOLINT_GOLANGCI_BINARY="$PWD/bin/golangci-lint" go test ./plugin/golangci -run '^TestIntegration$' -count=1 -v`
   and `go test ./clean -count=1` from the real-Testo integration module.
   Record unavailable checks; a skipped external integration is not success.
9. Document behavior, limitations, potential false positives and remedies in
   README; record the change under the current unreleased changelog version.
   Report fixture names, observed red/green results and checks for review.
10. Accept only when the benefit justifies maintenance. Narrow or defer a rule
    that needs speculative runtime assumptions or excessive machinery. Review
    the concrete changes before requesting commit approval; commits and releases
    require explicit authorization.

## v1.0.0 readiness criteria

- Stable IDs, rule metadata/classifications and supported semantics are documented;
  configuration keys, precedence, defaults and migration rules are understandable.
  The registry, documentation and both drivers agree on the complete rule set.
- Diagnostic prefix, Category, structured field preservation, meaningful source
  ownership and documentation URL guarantees are published and verified. State
  which wording, related-location and range refinements may ship in patches;
  do not imply that every message or internal implementation detail is frozen.
- Limitations cover unknown flow, aliases, generic discovery, imported bodies,
  runtime plugins, filtering and concurrency assumptions for applicable rules.
  No rule claims complete race detection or universal API coverage.
- Runtime reproductions and positive/negative regression fixtures support every
  accepted rule. Real repository audit findings are triaged; no known
  release-blocking correctness defect remains. Experimental rules are excluded
  from stable defaults unless they independently pass acceptance. New optional
  policies remain opt-in; existing defaults retain their documented behavior. Deferred research does not block stability.
- Supported Go/Testo/golangci-lint combinations have recorded successful checks;
  CLI and freshly rebuilt plugin agree on fixtures and configuration. Benchmarks
  have a recorded baseline and explained material regressions. Targeted pipeline
  architecture changes have evidence and preserved behavior; a rewrite is not
  a readiness condition.
- Follow AGENTS.md pre-release review: record target version, reviewed commit,
  working-tree state and exact tool/framework versions; compare against the
  previous release and inspect actual pinned Testo semantics. Run `make test`,
  uncached `go test ./... -count=1`, `go vet ./...`, `make lint` when available,
  `go mod tidy -diff`, `go mod verify`, rebuilt plugin integration and the
  real-Testo clean module. Record failures, retries and skips accurately;
  resolve missing essential verification before recommending publication.
- Save a review report with severity-ordered findings, evidence, source links,
  checks, remaining conditions and readiness recommendation. README installation
  versions, rule count, compatibility, limitations and CI match the target.
  Changelog records changes since the previous release under the unreleased
  target without an invented release date. Present the report and documentation
  for review; do not describe an unpublished tag as installable.

## Sequencing toward v1.0.0

Deliver v0.8.0 incrementally: capture the existing diagnostic/default baseline,
introduce the static registry and reporting helpers, migrate existing checks,
then connect configuration and rule documentation to that authority. The first
registry migration must explicitly handle checkers emitting multiple IDs;
otherwise selection can lose diagnostics or run shared work twice. Establish
working documentation URLs before promising them as stable metadata.

Use the resulting registry/configuration/reporting baseline for v0.9.0 corpus
and architecture review. Measure costs before factoring shared traversals; do
not make v0.8.0 depend on a speculative v0.9.0 redesign. Corpus access and the
supported toolchain matrix are evidence dependencies, not grounds to invent
compatibility claims. Record unavailable checks and remaining release conditions.

Evaluate v0.10.0 candidates one at a time against real usage and the planned
stability contract; retain only benefits worth a single maintainer's ongoing
cost. Then stabilize the supported subset in v1.0.0 independently of unfinished
optional policy research. Configuration and metadata are prerequisites for
accepted opt-in checks, not reasons to delay stability until every candidate
is implemented.
Custom extension research must stay smaller than the concrete use case it serves.

This planning update introduces no implementation or dependency changes.
Versions express sequence rather than dates. Future additive work must preserve
stable behavior and defaults; experimental work may remain deferred indefinitely.
