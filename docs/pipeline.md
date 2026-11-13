# Analyzer pipeline

Rule selection happens once per analyzer invocation. A registration executes once
when any of its IDs is selected; the reporter filters individual IDs. An empty
selection returns before suite discovery or call collection.

Execution order is stable:

1. For each discovered suite: suite checks, parallel checks, then `afterParallel`.
2. Run checks, in discovered call order.
3. Package checks, in registration order.

`afterParallel` preserves the historical reporting position of TESTO020. It does
not consume diagnostics or mutable results of parallel checks. This is an output
ordering constraint, not a semantic dependency.

Suite checks consume the discovered suite model. Parallel checks share one
parallel model per suite, built only when that group is selected. Run checks
consume call facts and resolved run calls, linked to local suite models. Package
checks consume call facts and can reuse the same discovered suites. TESTO027
checks suite bodies and then callbacks in its original package-stage position;
it does not rediscover suites.

Call facts are collected only if run or package checks are selected. Run discovery
is performed only for run checks. Suite discovery remains shared and unconditional
for nonempty selections: package checks can need suites even without suite checks.
No dependency graph or per-ID prerequisite metadata is required.

Checkers own semantic suppression (for example, avoiding cascades after invalid
signatures). Selection filters emitted IDs, independently of checker execution.
The invocation-wide reporter deduplicates diagnostics across groups and repeated
promoted methods. Refactoring execution must preserve full diagnostic fields,
including category, URL, positions, related information and suggested fixes, as
well as emission order. `run_pipeline`, configuration tests, diagnostic tests,
style tests and `no_sleep` cover these contracts.

## Measuring prerequisites

Run on preloaded packages (package loading and type checking are excluded):

```sh
go test ./analyzer -run '^$' -bench '^BenchmarkPipelineSelections$' -benchmem -benchtime=300ms -count=3
go test ./analyzer -run '^$' -bench '^BenchmarkPipelinePrerequisites$' -benchmem -benchtime=300ms -count=3
```

These benchmarks use synthetic preloaded packages, so they do not measure CLI
latency. The embedded-suite fixture has no run calls; run discovery measurements
on it do not establish performance on call-heavy packages. Record the checkout
revisions, Go version, hardware and sampling parameters when comparing results.

See [Contributing](../CONTRIBUTING.md#benchmarks) for full analyzer benchmarks.

## Structure-rule benchmark

```sh
go test ./analyzer -run '^$' -bench '^BenchmarkStructureSharedHelper$' -benchmem -benchtime=100ms -count=2
```

This compares registry defaults with TESTO029–031 enabled, using the same helper
from 1, 10 and 100 tests. Loading and type checking are excluded. Structure checks share one lazy walker across Allure suites in an analyzer
invocation. The walker collects call facts and local declarations once, then
caches context-independent helper operations and phase events. Each use rebases
the callsite and outer coverage without mutating the cached summary. No helper
summary follows another helper. State is not shared between analyzer runs.

When TESTO030 is disabled, the shared coverage walker does not collect phase
events, including in cached helper summaries.
