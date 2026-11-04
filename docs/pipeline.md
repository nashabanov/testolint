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

On Apple M4 / Go 1.27.1, the initial 100 embedded-suite scenario measured
suite discovery at 0.516–0.604 ms and call collection at 0.527–0.630 ms.
This justified skipping call collection for suite-only selections. Run discovery
on this fixture has no calls to process, so its ~2 ns result does not establish
performance on call-heavy packages. Its guard expresses the existing group
prerequisite without introducing additional rule metadata. Defaults continue to
use both prerequisites. These are synthetic measurements, not CLI latency claims
or performance budgets.

Median results over three 300 ms samples for this change:

| Selection | Before, ns/op | After, ns/op | Before, B/op | After, B/op |
| --- | ---: | ---: | ---: | ---: |
| Defaults | 1,249,568 | 1,245,925 | 2,331,669 | 2,330,096 |
| None | 1,084,507 | 248 | 2,173,698 | 248 |
| TESTO027 only | 1,646,966 | 1,094,197 | 3,060,926 | 2,197,691 |
| Style only | 1,100,817 | 553,275 | 2,175,830 | 873,447 |

Existing analyzer fixtures passed before and after the refactor, including
`no_sleep` and `run_pipeline`; no diagnostic expectations were changed. The
initial sandboxed baseline failed on missing Go cache entries and passed on a
retry with cache access. The CPU profile was captured, but the installed Go
has no `pprof` tool; the prerequisite microbenchmarks above informed the decision.
`make test`, `go vet ./...` and `make lint` passed on the final implementation.
External custom golangci-lint integration was not run for this internal refactor;
its skipped test is not evidence of external integration verification.
