# Roadmap

## Overview and project principles

testolint is a specialized static analyzer for
[Testo](https://github.com/ozontech/testo). This roadmap starts at v0.3.0 and
describes a path to v1.0.0 for a project maintained by one developer. Versions
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

## Current status

The repository's latest release tag and [changelog](CHANGELOG.md) are v0.4.2.
Release **v0.4.2** includes **TESTO001–TESTO017**,
all enabled by default. TESTO014–017 use the conservative subset documented in
the README; helper internals, custom contexts and unknown schedules remain excluded.
v0.4.2 extends partial lifecycle evidence across unsupported teardown or
receiver-independent sibling bodies while retaining the documented limits.
The [README](README.md) is the detailed reference for existing rules and limits.

Already implemented:

- Test and lifecycle signatures, Testo naming, parameter field visibility,
  field-name providers, assignability, unused providers, and narrowly provable
  empty providers and suites.
- Indirect suite embedding, aliases, pointer embedding, promoted and imported
  methods, and Go method-set shadowing and ambiguity handling. Generic suite
  declarations remain excluded from suite checks.
- Typed discovery of RunSuite/RunSubSuite, including generic calls and unchanged
  direct local function aliases; conservative nil-suite and direct receiver
  recursion checks, including instantiated run arguments.
- Shared suite/run diagnostic deduplication, regression fixtures, a reusable
  go/analysis analyzer, CLI, golangci-lint v2 module plugin, real-Testo fixtures,
  and CI integration. Determinism already has pipeline regression coverage.

These are foundations to preserve, not new deliverables. Current models retain
suite types, T, method sets and optional local method bodies, plus suite run
calls and limited local initializers. Parallel checks retain narrow local
lifecycle order, direct field ownership and Cleanup reads; they do not provide
general scheduling, alias analysis or standalone Testo execution-API modeling.
The plugin currently ignores settings; per-rule configuration is not implemented.
The declared baseline is Go 1.27.1, Testo v1.8.0 in integration fixtures, and
golangci-lint v2.14.0 in CI. Broader compatibility is not implied.

### Testo semantic baseline

Planning was checked against actual v1.8.0 source:

- [suite.go](https://github.com/ozontech/testo/blob/v1.8.0/suite.go) and
  [runner.go](https://github.com/ozontech/testo/blob/v1.8.0/runner.go): BeforeEach
  does not construct a fresh suite. A pointer suite shares its instance; value
  receivers copy values but can still share referenced data. Each test gets a
  constructed T; BeforeEach, test, and AfterEach use that test's context.
  BeforeAll/AfterAll use the suite context. AfterEach is deferred and can precede
  completion of parallel subtests.
- [collector.go](https://github.com/ozontech/testo/blob/v1.8.0/collector.go) and
  runner.go: providers are evaluated after BeforeAll and before BeforeEach;
  annotations are read during collection before the plugin plan is applied.
  Dynamic providers and dependencies on BeforeAll are supported.
- [construct.go](https://github.com/ozontech/testo/blob/v1.8.0/construct.go) and
  [reflectutil](https://github.com/ozontech/testo/blob/v1.8.0/internal/reflectutil/reflect.go):
  construction fills settable fields recursively and wires exported pointer
  fields, with special handling for testo.T. Exported non-pointer fields on
  traversed custom T/plugin structures can panic; unexported fields are skipped.
  Recursion follows traversed fields, not every reference in a Go type graph.
- runner.go accepts nil callbacks for **Run**, replacing them with a no-op.
  RunTest/Test use a different singleton path. Run during Cleanup is explicitly
  documented as panicking.
- [plugin contract](https://github.com/ozontech/testo/blob/v1.8.0/testoplugin/plugin.go):
  at the suite root, Plugin receives a non-nil interface containing a typed nil
  pointer of its own plugin type; nested contexts receive the parent instance.
- [options.go](https://github.com/ozontech/testo/blob/v1.8.0/options.go) and
  [annotations.go](https://github.com/ozontech/testo/blob/v1.8.0/annotations.go):
  registration uses synchronized process-global storage. Late calls can affect
  future runs; synchronization does not provide test isolation. BeforeAll is
  not automatically too late for annotations.

These source observations guide candidates; runtime reproductions remain an
acceptance prerequisite. Reports must distinguish runtime failures, framework
warnings (including strict-mode escalation), and optional linter policy.

## Release overview

Counts below are candidate counts, not promised totals. All IDs after TESTO017
are provisional and must be rechecked before implementation; implemented IDs
must never be reassigned.

| Version | Focus | Planned new rules | Experimental candidates | Other work |
| --- | --- | --- | --- | --- |
| v0.3.0 | Execution Discovery — released | 13 implemented overall | — | Preserve current behavior |
| v0.4.0 | Shared State & Parallel Safety — released | 4 implemented | 0 | Conservative concurrency evidence |
| v0.4.1 | Parallel Analysis Coverage Fixes — released | 0 | 0 | Serial interference, partial paths and case counts |
| v0.4.2 | Parallel Lifecycle Analysis Fixes — released | 0 | 0 | Independent hook and sibling lifecycle facts |
| v0.5.0 | Execution & Lifecycle Safety | 2 | 1 | Extend execution API discovery as needed |
| v0.6.0 | Custom T & Plugin Validation | 2 | 1 | Match reflection traversal |
| v0.7.0 | Providers & Annotations | 3 | 2 | Optional TESTO010 extension |
| v0.8.0 | Configuration & Developer Experience | 0 | 0 | Rule selection, defaults, diagnostic usability |
| v0.9.0 | Hardening & Compatibility | 0 | 0 | Corpus audits, compatibility, benchmarks |
| v1.0.0 | Stable Release | 0 | 0 | Explicit stability contract |

## Detailed version milestones

### v0.3.0 — Execution Discovery (released)

TESTO001–TESTO013 cover the existing suite validation and execution checks
summarized above. This release added typed execution discovery, statically nil
suite arguments, and narrowly defined direct RunSubSuite receiver recursion.

**Definition of Done:** historical baseline is the v0.3.0 tag and changelog;
future changes preserve its fixtures or explicitly document accepted behavior
changes. No additional work is assigned to this released version.

### v0.4.0 — Shared State & Parallel Safety

**Released.** TESTO014–017 cover local
straight-line lifecycle/test bodies with standard `*testo.T` context identity.
Shared-field checks require direct pointer receiver fields, BeforeEach writes
before suspension and two known parallel executions. Constant literal providers
can establish parameterized multiplicity. Cleanup diagnostics take precedence
over TESTO014; context fields use TESTO016. Unknown calls, branches, aliases,
imported/promoted bodies and additional test/AfterEach writes are conservatively
excluded. Unknown provider counts suppress shared-field checks; TESTO017 only
requires an established execution of its own test, independently of sibling
counts. Runtime reproductions and CLI/plugin fixtures
exercise the accepted subset. The 0.4.1 patch expanded this initial scope; see
README for current boundaries and CHANGELOG for the preserved release history.

**Required:** define context identity, receiver sharing, direct field accesses,
and execution ordering for each supported pattern before writing diagnostics.
Establish parallelism through resolved calls on the relevant Testo context;
do not infer it from names, pointer receivers, or a Parallel call on an unrelated
context. Begin with local bodies and simple unconditional paths. Skip unknown
aliases, synchronization, plugin overrides, helper behavior, and branching when
they prevent a defensible conclusion. Account for multiple parameterized
executions and the order of writes relative to Parallel suspension.

#### TESTO014 — SharedBeforeEachWrite

**P1 · Planned.** Detect a direct BeforeEach write to a shared receiver field
followed by use in a parallel test, where another execution can overwrite that
field before the use. Testo-specific value comes from combining shared suite
identity with hook/test ordering. Scope: proven shared storage, matching field
objects, and locally established overlapping executions. Risks: serial tests,
copies, captured per-test values, synchronization, and writes that occur outside
the overlap. A suite-field assignment alone is insufficient evidence.

#### TESTO015 — SharedFixtureCleanup

**P1 · Planned.** Detect a Cleanup callback that reads a mutable suite field
whose value another execution can replace before cleanup. Testo's reused
receiver makes an apparently per-test fixture reference point at another test's
resource. Scope: direct closures, resolved Cleanup, and a proven overwrite
window. Risks: immutable shared fixtures, protected fields, serial execution,
and callbacks using a captured local resource rather than re-reading the field.
Report the cleanup hazard without duplicating TESTO014 for the same cause.

#### TESTO016 — SharedTestContext

**P1 · Planned.** Detect storing a test-scoped T into shared suite state and
using that field from an execution that can observe a different test's context.
This depends on Testo's per-test T construction versus suite lifetime. Scope:
direct assignments of the hook/test parameter and direct field uses with proven
overlap or lifetime mismatch. Risks: a serial store used only by the same test,
suite-scoped BeforeAll T, per-test local copies, and intentional context storage.
Do not report every field of type T or every assignment to it.

#### TESTO017 — DuplicateParallelLifecycle

**P1 · Planned.** Detect duplicate Parallel calls across Testo lifecycle and
test bodies that execute on the same underlying test context. Scope: direct,
unconditional calls across BeforeEach/test/AfterEach, with context identity and
reachability established. Testo-specific value is cross-hook reasoning beyond
generic duplicate-call checks. Risks: separate suite/test/subtest contexts,
exclusive branches, early termination, and plugin overrides. Do not equate
BeforeAll and test T or count unrelated subtest calls together.

**Optional:** more local receiver/callback forms after precision is demonstrated.
**Research:** synchronization recognition or helper analysis; no framework is
required merely to enlarge coverage.

**Definition of Done:** accepted rules have runtime reproductions, serial and
parallel negative fixtures, documented scheduling assumptions, and integration
coverage. Ship only independently accepted rules; reclassify any candidate
without a precise strategy as Experimental rather than delaying the release.

### v0.5.0 — Execution & Lifecycle Safety

**Required:** extend typed call discovery only for APIs needed by accepted rules;
preserve import aliases, generics, unrelated-symbol exclusion and deduplication.

#### TESTO018 — RunInCleanup

**P1 · Planned.** Detect testo.Run executed directly inside a registered Cleanup
callback, matching Testo's documented cleanup restriction. Scope: resolved Testo
Run calls in direct cleanup closures with the relevant context established.
Risks: helper registration that does not execute the call, unrelated Cleanup
methods, and nested closures only invoked later. Do not broaden this to every
function whose name starts with Run or duplicate generic testing.Run checks.

#### TESTO019 — NilTestCallback

**P1 · Planned.** Detect statically nil function callbacks supplied to RunTest
or Test; their singleton execution path invokes the callback. Scope: explicit
nil, typed nil conversions, and unchanged local zero/nil function variables.
For Test, explain that failure occurs when the returned adapter is executed.
Risks: reassignment, unknown function values and adapters never invoked.
Explicitly exclude testo.Run, which supports nil as a no-op callback.

#### TESTO020 — PrematureTeardown

**P2 · Experimental.** Investigate AfterEach teardown of a resource used by
parallel subtests that finish later. The Testo-specific issue is deferred
AfterEach ordering. Candidate scope: a locally identified resource, explicit
parallel subtest use and direct teardown of that same resource. Risks: arbitrary
Close-like methods, independent resources, joins, cleanup ordering and plugin
behavior. Method names alone cannot establish teardown or dependence. A minimal
reproduction and precise resource identity are required before implementation.

**Definition of Done:** TESTO018/019 pass acceptance and execution-API boundary
tests; document the nil-Run exception. TESTO020 ends with either an accepted
narrow strategy or a recorded research outcome, not a required diagnostic.

### v0.6.0 — Custom T & Plugin Validation

**Required:** derive reachable construction roots from actual Testo use, then
mirror Fill/setPlugins traversal, visibility and testo.T exceptions. Avoid
validating arbitrary structs merely because they resemble a plugin.

#### TESTO021 — InvalidCustomTField

**P1 · Planned.** Detect exported non-pointer fields reached by Testo's custom
T/plugin wiring where setPlugins requires pointers. Scope: known concrete roots
and traversed nested structures, respecting aliases and special framework
types. This is a reflection initialization failure, not a general struct style
rule. Risks: unexported fields, structs outside the construction graph, and
incorrectly treating testo.T internals as user fields. Suppress dependent reports.

#### TESTO022 — RecursiveTGraph

**P1 · Planned.** Detect a type cycle that Testo's recursive initialization
actually traverses and cannot construct. Scope: concrete roots and an active
traversal path matching Fill; report a useful field/type cycle. Risks: repeated
types on independent branches, unexported fields and references through
containers that Fill does not descend into. General recursive Go types are valid
and are not sufficient evidence. Avoid a second report for the same cycle.

#### TESTO023 — UninitializedPluginField

**P2 · Experimental.** Investigate use of plugin fields skipped by automatic
initialization, such as inaccessible fields, where initialization is demonstrably
assumed. Testo's reflection visibility and plugin wiring provide the specific
context. Scope must include a proven use and missing initialization; a nil field
alone is insufficient. Risks: deliberate manual initialization, optional plugins,
nil-safe methods and ordinary application fields. Exported pointer fields are
normally filled; overlaps with TESTO021/022 should extend those rules instead.

**Definition of Done:** TESTO021/022 agree with real construction reproductions,
including visibility, nested pointers, harmless repeated types and framework
exceptions. TESTO023 ships only if a distinct, maintainable failure pattern is
established; otherwise retain research notes in the backlog.

### v0.7.0 — Providers & Annotations

**Required:** preserve support for dynamic Cases and BeforeAll initialization.
Use resolved provider fields and annotation targets rather than method names
alone to infer dependencies or timing.

#### TESTO024 — ProviderBeforeEachDependency

**P1 · Planned.** Detect a provider relying on state initialized only in
BeforeEach, although collection happens earlier. Scope: direct provider field
reads, direct hook initialization and a demonstrable required value, with no
known earlier initialization. Testo-specific value is collection order.
Risks: constructor initialization, valid zero values, BeforeAll writes and
intentional reads of initial state. Reading a field also written in BeforeEach
is not itself a violation; skip cases needing unknown call-site history.

#### TESTO025 — LateAnnotation

**P2 · Experimental.** Investigate For/ForEach calls whose target's annotations
have already been read for the current plan. Scope: locally resolved targets
and proven registration/collection ordering. Risks: calls configuring a future
suite run, other suites, and legitimate BeforeAll registration. Testo snapshots
annotations during collection; registration in a test cannot retroactively
change that plan but may intentionally affect a later run. Establish an
unambiguous misuse pattern before treating it as a correctness diagnostic.

#### TESTO026 — PluginParentNilCheck

**P1 · Planned.** Detect direct comparisons of a framework-supplied Plugin
parent interface with nil when used to distinguish the root. Testo supplies a
typed nil at the root, so that comparison cannot identify it. Scope: resolved
Plugin implementations reachable from Testo construction and the unchanged
parent parameter. Risks: reassigned parameters, manually invoked Plugin methods
and comparisons serving another purpose. Comparing an asserted concrete pointer
with nil is valid and must remain silent.

#### TESTO027 — GlobalOptionsInTest

**P1 · Planned.** Detect testo.Options registration during Testo execution when
its process-global append can leak configuration into a subsequent initialization.
Scope: resolved calls in tests/hooks and a supported, demonstrable later run;
skip uncertain ordering. This is Testo configuration isolation, not an
unsynchronized-memory warning: Options uses a mutex. Risks: deliberate setup for
a child suite, process-isolated tests and intentional global configuration.
Package initialization remains valid; decide whether ambiguous patterns belong
to opt-in policy before enabling diagnostics by default.

#### TESTO028 — StaticCasesOnly

**P2 · Experimental; optional policy, opt-in only.** Investigate reproducible
case discovery as a user-selected restriction; dynamic Cases are valid Testo.
Define a published expression whitelist: slice literals with supported elements,
compile-time constants used in cases, and simple local values proven unchanged
and composed of supported expressions. Go has no constant slices; constants may
supply elements or bounds. Explicitly support further provable expressions only
with fixtures.

Diagnose only established disallowed dependencies, such as recognized external
service/database calls or environment reads that determine cases. Arbitrary
helper calls are unknown, not automatically nondeterministic. Skip unknown
expressions and report the checked subset clearly. Risks include intentional
integration discovery, deterministic helpers and hidden mutation. Do not promise
complete enforcement without interprocedural analysis. Delivery depends on rule
selection support; move to v0.8.0 or later if that prerequisite is unavailable.

**Optional enhancement:** extend TESTO010, without a new ID, to additional
provably empty providers, starting with built-in make of a slice with constant
zero length and simple unchanged empty locals. Resolve built-ins by type/object
identity; skip mutations, escapes and unknown flow. Preserve current checks and
add false-positive fixtures before expanding accepted expressions.

**Definition of Done:** TESTO024/026/027 pass acceptance with initialization,
timing and intentional-use negatives. TESTO025/028 remain research unless precise
semantics are established; TESTO028 cannot ship enabled by default. Any TESTO010
extension has its own red/green evidence and documented limits.

### v0.8.0 — Configuration & Developer Experience

**Required (P1 · Planned):**

- Add per-rule enable/disable selection with a small validated configuration
  model exposed through go/analysis flags and golangci-lint settings. Reject
  unknown IDs and conflicting settings clearly; test configuration isolation
  across analyzer instances and integration parity.
- Classify rules as correctness or policy. TESTO006 is already policy;
  TESTO012 also includes intentional nil-safe suites. Classification must reflect
  existing behavior. Document migration if defaults change; do not silently
  disable existing rules. New optional policy such as TESTO028 stays opt-in.
- Document defaults, precedence, unsupported patterns and actionable diagnostics.
  Verify stable IDs, positions, suppression of cascading reports and CLI errors.

**Optional (P2):** SuggestedFixes only for edits with a uniquely safe meaning,
with fix-output tests; many lifecycle changes require human judgment. Evaluate
source suppression only if a concrete user need remains after rule selection.
No custom DSL or general configuration framework is planned.

**Definition of Done:** configuration examples work in both integrations;
defaults and invalid settings have end-to-end coverage; policy activation is
explicit and any changed defaults have migration guidance. No new rule quota.

### v0.9.0 — Hardening & Compatibility

**Required (P0 · Planned):**

- Audit at least two real Testo repositories with different suite/provider/plugin
  patterns, recording revisions, configuration, diagnostics and manual decisions.
  Record confirmed false positives and fixes or documented scope restrictions;
  leave no known release-blocking false positive unresolved.
- Add regression fixtures for confirmed defects. Audit generic declarations
  versus instantiated calls, embedding, aliases, pointer/value sets, promoted and
  imported methods, shadowing, ambiguity and diagnostic cascades. Existing
  support is a regression target; broader generic discovery remains conditional.
- Check deterministic diagnostics across repeated uncached runs and stable
  positions/deduplication across CLI and plugin outputs.
- Establish benchmarks for representative small and large packages. Record
  hardware, Go/Testo versions, fixture size, time and allocations; compare each
  substantial analysis addition with that baseline. Investigate regressions and
  set a justified budget after measurement, without invented throughput promises.
- Publish and exercise an explicit Go/Testo/golangci-lint compatibility matrix.
  Keep at least the declared minimum Go and current selected supported toolchain
  covered; claim only Testo versions verified by runtime and integration checks.
  Rebuild the custom binary from the reviewed checkout and compare clean and
  diagnostic fixtures, configuration and IDs with the CLI.

**Optional:** add supported versions and discovery patterns when real usage
justifies the cost. **Research:** assess uncovered semantics without turning
every gap into a release blocker.

**Definition of Done:** save a reproducible audit report with corpus results,
triaged findings, benchmark baseline, exercised compatibility matrix, exact check
outcomes and skipped checks. All confirmed correctness blockers are resolved;
unsupported cases are explicit. A skipped plugin test is not verification.

### v1.0.0 — Stable Release

**Required (P0 · Planned):** publish a stability contract covering diagnostic IDs
and supported semantics, configuration keys/defaults, limitations, compatibility,
and integrations. Preserve implemented IDs; do not recycle removed IDs. Within
v1.x, treat breaking supported semantics or configuration changes as requiring a
major version; document bug fixes and additive detection coverage in release
notes. New optional policy remains explicit.

**Definition of Done:** meet the readiness checklist below and publish only
after maintainer review. Complete coverage of every Testo API is not required.
Unresolved experimental proposals can remain outside the stable rule set.

## Deferred / experimental backlog

Items here have no release commitment. Unnumbered ideas receive no diagnostic
ID until their semantics are established. Milestone experiments TESTO020,
TESTO023, TESTO025 and TESTO028 can move here without blocking delivery.

| Idea | Priority / status | Reason deferred and next evidence needed |
| --- | --- | --- |
| SharedCollectionMutation | P2 · Deferred | Map/slice mutation needs alias, ownership and overlap evidence; first establish whether TESTO014/015 can cover useful direct patterns. |
| ParallelSuiteMutation | P2 · Deferred | Broadly reporting receiver writes would resemble generic race detection and flag safe code; require a distinct Testo lifecycle failure beyond v0.4.0. |
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
   diagnostic position and default enablement before implementation.
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

- Stable IDs and supported semantics are documented; configuration, defaults,
  correctness/policy classification and migration rules are understandable.
- Limitations cover unknown flow, aliases, generic discovery, imported bodies,
  runtime plugins, filtering and concurrency assumptions for applicable rules.
  No rule claims complete race detection or universal API coverage.
- Runtime reproductions and positive/negative regression fixtures support every
  accepted rule. Real repository audit findings are triaged; no known
  release-blocking correctness defect remains. Experimental rules are excluded
  from stable defaults unless they independently pass acceptance.
- Supported Go/Testo/golangci-lint combinations have recorded successful checks;
  CLI and freshly rebuilt plugin agree on fixtures and configuration. Benchmarks
  have a recorded baseline and explained material regressions.
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

The practical release strategy is to complete one rule at a time, defer research
that cannot meet precision requirements, and advance configuration work earlier
if an optional rule needs it. v1.0.0 depends on confidence in the supported
subset, not on shipping every provisional ID above.
