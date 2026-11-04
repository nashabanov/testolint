# Roadmap

## Current baseline

v0.8.2 provides 27 rules (TESTO001–TESTO027), a shared rule registry,
validated configuration for the CLI and golangci-lint plugin, structured
reporting, Allure style fixes, and analyzer benchmarks. TESTO001–TESTO024 are
enabled by default; TESTO025–TESTO027 are opt-in. Completed implementation work
and release history belong in [CHANGELOG.md](CHANGELOG.md), with supported
behavior and limitations in [README.md](README.md) and [docs/rules.md](docs/rules.md).

The planning baseline is Go 1.27.1, Testo v1.8.0, Testo Allure v1.5.0 and
golangci-lint v2.14.0. These versions identify the starting configuration;
broader compatibility must be verified before it is claimed.

## Principles

- Prioritize precision, useful diagnostics and maintenance cost over rule count.
- Use conservative AST and Go type analysis, with explicit unsupported cases.
  Unknown behavior is not evidence of an error.
- Extend existing rules before adding overlapping checks. Require Testo-specific
  value beyond generic Go linters and race detection.
- Keep the current analyzer package and thin driver adapters. Introduce shared
  models or abstractions only for concrete consumers and measured needs.
- Preserve diagnostic IDs; never recycle them. Classification and default
  enablement are separate decisions. New optional policies start disabled.
- Follow the fixture-first workflow and review requirements in
  [AGENTS.md](AGENTS.md). No release date or automatic publication is implied.

## v0.9.0 — Hardening, architecture and compatibility

**Status: planned.** Stabilize the existing 27 rules and gather the evidence
needed for v1.0.0. New rules, broader discovery and a new analysis framework are
outside the required scope. Execute the following stages in order; documentation
and compatibility work may proceed alongside corpus investigation once the
baseline is recorded.

### 1. Capture a reproducible baseline

- Compare against v0.8.2 and record the reviewed commit, working-tree changes,
  hardware, Go, Testo, Allure and golangci-lint versions.
- Run the verification checklist below, including a freshly built custom plugin.
  Investigate environment failures and record retries separately from defects.
- Capture diagnostics for existing fixtures: IDs, messages, primary and related
  locations, ranges, fixes, deduplication and rule selections.
- Record benchmark time, memory and allocations for the existing small/large,
  no-Testo, valid, diagnostic-heavy and embedded-suite scenarios.
- Measure full CLI execution separately: analyzer benchmarks exclude package
  loading and type checking. Include opt-in style/no-sleep configurations and
  disabled rule groups when measuring pipeline costs.

**Deliverable:** reproducible commands, results and a comparison baseline.
Choose performance budgets after measurement, without invented throughput targets.

### 2. Audit real Testo projects and fix confirmed defects

- Audit at least two available real repositories with different parameterization,
  embedding, custom context/plugin, parallel lifecycle and Allure usage. Record
  repository revisions and exact configurations; identify uncovered scenarios.
- Run defaults and opt-in rules separately. Classify each finding as a confirmed
  error, useful policy diagnostic, false positive or unresolved case. Investigate
  representative missed cases; a clean run alone does not establish coverage.
- Review aliases, pointer/value method sets, promoted/imported methods, shadowing,
  ambiguity, invalid signatures, generic declarations versus instantiated calls,
  and cascading diagnostics where relevant to the findings.
- Reproduce semantic issues with actual pinned Testo source and runtime cases,
  distinguishing failures, framework warnings and linter policy.
- For every rule/discovery fix, first add a type-checking e2e fixture containing
  failing and valid examples, then observe the expected diagnostic failure.
  Implement only after red; rerun the targeted test and full checks after green.
  Add corresponding CLI/plugin coverage first when driver behavior changes.
- Keep broader unsupported patterns in the backlog unless real usage justifies
  a precise, maintainable extension.

**Deliverable:** triaged corpus results and regression fixtures for fixes, with
observed red/green outcomes. Resolve confirmed release-blocking defects. If corpus
access is unavailable, record the missing evidence rather than claiming completion.

### 3. Simplify the pipeline while preserving behavior

- Reuse suites already discovered by the pipeline in TESTO027 instead of calling
  `discoverSuites` again inside `checkSleep`.
- Add a fast path when no rules are selected.
- Use profiling to decide whether call facts and run discovery should be skipped
  for rule selections that do not consume them. Record a decision to retain the
  existing implementation when added complexity is not justified.
- Document execution groups, the purpose of `afterParallel`, shared prerequisites
  and diagnostic ownership. Distinguish semantic dependencies from output order.
- Preserve single execution of multi-ID checkers, individual ID selection,
  suppression, structured diagnostic fields and deterministic deduplication.
- Keep behavior-preserving refactors under green regression tests. Any discovered
  rule/discovery defect follows the red/green procedure in stage 2.

**Deliverable:** focused changes with before/after diagnostics and performance
comparisons. Do not introduce a universal dependency graph, SSA, a public rules
framework or package reorganization solely for hypothetical future needs.

### 4. Define the user-facing contract

- Document guarantees for IDs, names, classification, configuration validation,
  precedence, defaults, CLI exit codes, Category, URLs and SuggestedFixes.
- Explain which message, range and related-location refinements may ship in patch
  releases. Detection fixes may change diagnostics within supported scope, but
  must be described as fixes with regression coverage.
- Move detailed rule descriptions, invalid/valid examples, remedies and
  rule-specific limitations out of README into one document per rule under
  `docs/rules/`, using stable ID-based filenames such as `TESTO001.md`.
  Keep README focused on installation, usage, configuration and a concise rule
  table linking to those documents. Turn `docs/rules.md` into a linked index;
  avoid duplicating descriptions and examples across these locations.
- Give every registry entry a rule-specific documentation URL, replacing the
  shared README URL. Preserve the structured diagnostic `URL` and include the
  same link in the diagnostic message text so it is available in plain CLI and
  golangci-lint output.
- Define and test a consistent message-link format retaining exactly one ID
  prefix. Update diagnostic expectations and driver e2e tests before changing
  reporting, observing the intended red phase first. Verify links appear once
  and match the rule ID, while preserving Category, ranges, related locations,
  suggested fixes and deduplication behavior.
- Test agreement between registry URLs, documentation files, the README table,
  the rule index and both drivers. Update internal links and retain navigation
  from previously published documentation entry points. Describe the intentional
  diagnostic message change in migration notes.
- Verify style fixes preserve valid Go, comments and LF/CRLF behavior, and that
  applying fixes again produces no further edits.
- Resolve TESTO027 default enablement after the corpus audit. The previous plan
  proposed enabling it in v0.9.0; this plan makes that change conditional on
  evidence and review. Either document the enablement migration and opt-out in
  both drivers, or explicitly retain opt-in behavior and record the decision.
  No other default changes are planned.

**Deliverable:** a documented contract, tested CLI/plugin parity and an explicit
TESTO027 decision. Local source suppression is optional and requires demonstrated
adoption problems that cannot be addressed reasonably with existing selection.

### 5. Verify compatibility and maintain CI evidence

- Define the minimum supported toolchain and any additional supported Go/Testo/
  golangci-lint combinations. Cover the declared minimum and the selected current
  supported toolchain; these may initially be the same verified configuration.
- Publish only combinations exercised by runtime and integration checks. Include
  Testo Allure compatibility for style behavior and update CI alongside the table.
- Retain fresh custom golangci-lint builds in CI and run all integration variants,
  including configuration and real-Allure fixes.
- Run real-Testo clean fixtures without cached success. Keep runtime regression
  cases for affected semantics in addition to simplified analyzer fixtures.
- Record performance comparisons and investigate material regressions before
  accepting a change. Avoid noisy timing thresholds without measured stability.

**Deliverable:** an exercised compatibility matrix and repeatable CI checks.
A narrow, explicit support policy is acceptable; unverified claims are not.

### 6. Complete documentation and release review

- Update README limitations, remedies, rule counts, defaults, compatibility and
  installation examples to match the reviewed implementation.
- Record all changes since v0.8.2 under the unreleased v0.9.0 changelog entry,
  without inventing a release date or describing the target tag as installable.
- Follow the full pre-release review in AGENTS.md. Read pinned Testo `suite.go`
  and `collector.go`, plus affected lifecycle/construction/plugin source. Include
  versioned source links and separate confirmed bugs from coverage gaps and
  maintenance suggestions.
- Save the local review report outside the repository, for example at
  `/tmp/testolint-v0.9.0-review.md`. Include scope, versions, findings, corpus
  evidence, benchmarks, exact check outcomes, failures/retries/skips, remaining
  conditions and a readiness recommendation. Do not commit release-review reports.
- Present changes and evidence for review. Commit, tagging and publication require
  explicit authorization; updating this roadmap does not authorize them.

### Verification checklist

Run from the reviewed checkout:

```sh
make test
go test ./... -count=1
go vet ./...
make lint
go mod tidy -diff
go mod verify
```

Run `make lint` when golangci-lint is available; report its absence explicitly.
Rebuild the custom binary using the repository's `.custom-gcl.yml` and the
selected supported bootstrap version before running:

```sh
TESTOLINT_GOLANGCI_BINARY="$PWD/bin/golangci-lint" go test ./plugin/golangci -run '^TestIntegration' -count=1 -v
```

From `plugin/golangci/testdata/integration`, also run:

```sh
go test ./clean -count=1
```

An old custom binary, cached test success or skipped external integration does
not verify the current implementation. Resolve missing essential verification
before recommending publication.

### v0.9.0 acceptance criteria

- No confirmed release-blocking defect remains unresolved.
- Every rule/discovery fix has a regression fixture and recorded red/green results.
- The selected real-project corpus has been audited and findings triaged.
- CLI and the freshly rebuilt plugin agree on diagnostics and configuration.
- Supported compatibility combinations have recorded successful checks.
- Performance changes are measured; material regressions are explained and reviewed.
- Pipeline changes preserve diagnostic ownership, suppression and structured fields.
- Defaults, limitations, migration guidance and rule documentation are consistent.
- Every rule has a dedicated documentation file with examples, linked from the
  registry, README/index and plain-text diagnostics in both drivers.
- The release review records exact outcomes and remaining conditions for approval.

## v1.0.0 — Stable supported behavior

Use the v0.9.0 evidence to finalize the stability contract. No mandatory new rules
or intermediate research release are required. Before recommending v1.0.0:

- Revalidate the supported semantics and toolchain matrix against the target
  checkout using the same full release-review process.
- Publish compatibility and migration policies for configuration, defaults,
  diagnostic identity and structured output. Breaking supported contracts in
  v1.x requires a major version; document bug fixes and additive coverage.
- Preserve published IDs, document retirement without reuse, and maintain working
  rule documentation links with a migration strategy when pages move.
- Keep unknown flow, generics, imported bodies, runtime plugins, filtering and
  concurrency limitations explicit. Do not promise complete race detection,
  universal Testo coverage or immutable diagnostic prose and byte offsets.
- Resolve blocking findings and present the final changelog and local review
  report for approval. Deferred research does not block stability.

## Deferred work

These items have no assigned release or diagnostic ID. Require real examples,
precision evidence and a maintenance justification before implementation.

| Area | Evidence needed before expanding scope |
| --- | --- |
| Generic suite discovery and broader recursion detection | Concrete missed cases, instantiation/reachability requirements and clear diagnostic ownership. |
| Broader empty-provider and provider-initialization checks | Proof patterns that preserve valid zero values, BeforeAll initialization and unknown instance history. |
| Plugin initialization gaps | A demonstrated failure at use, accounting for intentional nil fields, manual initialization and nil-safe methods. |
| Shared collection/global-state mutation | Alias, ownership, overlap and restoration evidence beyond generic race/style checks. |
| Polling, timing policies and case-count limits | Real costs and Testo-specific value; unknown bounds or large combinations alone do not establish an error. |
| Local source suppression | Repeated legitimate exceptions that make package-wide rule selection insufficient. |
| Interprocedural analysis and custom rules API | High-value downstream cases that simpler local analysis or built-in rules cannot address. |

Extend existing rules when possible, including TESTO027 for justified sleep-related
coverage. Optional new policies remain opt-in. Research may conclude that no
implementation is warranted and must not delay the stable release.
