# v0.3.0 pre-release review

Reviewed on 2026-10-10. Target v0.3.0 is unpublished; no tag, release or push
was created. Implementation commit: `719182a938bd9b9f91dfb7789d0f1eef689792ef`.
Compared against v0.2.0 (`d2cfb94f03d559e77ae24437f1fbc16823537d96`).
The pending documentation stage changes README.md, CHANGELOG.md and this report.

Environment: Go 1.27.1 darwin/arm64; stock golangci-lint 2.14.0 (built with
Go 1.27.0, commit 114493f9); integration dependency Testo v1.8.0.
The custom binary was rebuilt from the reviewed checkout, using local module
replacement, with Go 1.27.1:
`v2.14.0-custom-gcl-LsDq73qFkQIRzr0YXXCWCgjVJnC4Oa7dosjsvp6RTzc`.

## Scope and findings

No blocking defect was found in the reviewed scope. Existing TESTO001–TESTO011
implementations remain unchanged. Suite and run rules share position/message
deduplication, with per-pass state and deterministic sequential execution.
Run discovery uses go/types function identity, one AST traversal, and a shared
index of unchanged local initializers; no SSA, call graph or new dependency.

The following are intentional limitations or coverage gaps, not confirmed bugs:

- `analyzer/rule_nil_suite.go:14`: TESTO012 reports known nil pointers, including
  intentionally nil-safe suites. It does not promise a panic. Interface values,
  globals, assignment-based inference and control-flow nilness are excluded.
- `analyzer/rule_recursive_subsuite.go:17`: TESTO013 requires a sole direct call
  in BeforeAll or the sole regular test with inherited preceding hooks. It checks
  receiver identity, the parent parameter and the actual argument method set.
  Options, providers, guarded/stateful bodies, imported bodies and other paths
  are excluded. Plugins and filters can interrupt this structural recurrence.
- `analyzer/run_discovery.go:13`: function aliases resolve one unchanged local
  initializer only. Any assignment or address-taking disables inference, even
  after the call. This sacrifices coverage to avoid execution-order assumptions.
- Partial type information has nil/empty-info regression coverage; exhaustive
  sparse-info permutations and large-package performance benchmarks are absent.
  Add these if analyzer consumers or measured performance require them.

Review covered aliases, generic instantiation, pointer/value method sets,
promoted/imported methods, shadowing, ambiguity, malformed signatures,
diagnostic positions, cascades, deduplication, CLI, plugin registration and CI.
Suite rules retain their documented pointer-method-set scope; run recursion
uses the argument's actual method set. No release-only refactoring was introduced.

## Comparison with pinned Testo

The actual v1.8.0 module source was read locally, including
[suite.go](https://github.com/ozontech/testo/blob/v1.8.0/suite.go),
[collector.go](https://github.com/ozontech/testo/blob/v1.8.0/collector.go) and
[runner.go](https://github.com/ozontech/testo/blob/v1.8.0/runner.go).

- The private suite marker and T-bearing methods support typed discovery.
  Runtime reflection uses the supplied type's method set. Naming, hook/test
  signatures, field-name provider matching, slice element assignability and
  field visibility were compared with the existing rules. Exactly Cases is
  ignored. Malformed collected signatures fail at runtime; unexported parameter
  fields cannot be set by reflection. TESTO006 is policy: unused valid providers
  are permitted. Empty case sets and empty suites produce runtime warnings;
  TESTO010 and TESTO011 expose those patterns statically.
- Both run functions take the suite as argument two. There is no universal nil
  rejection. A nil outer pointer with inherited value-receiver hooks panicked
  in a local runtime reproduction; a suite with nil-safe pointer hooks/tests
  completed successfully. TESTO012 therefore describes nilness rather than an
  unconditional runtime failure.
- Collection precedes hooks and tests. BeforeAll runs before tests; custom
  preceding hooks can bound recursion. A runtime reproduction observed repeated
  BeforeAll entry, stopped by a plugin guard. The committed clean integration
  fixture executes bounded receiver recursion successfully, without TESTO013.
  The rule excludes custom preceding hooks for test-method recurrence.

## Tests and observed red/green phases

Fixtures were written before implementation and loaded/type-checked successfully
for the counted red phases. Loading/cache failures were not counted as TDD red.

| Fixture or test | Red evidence | Green coverage |
| --- | --- | --- |
| run_discovery | Missing discovery diagnostics | Real function identity, aliases, generics, unrelated names, changed/escaped aliases; nil/empty TypesInfo unit checks |
| run_pipeline | Four missing run diagnostics | Existing eleven rules alongside run rules, shared deduplication, repeated-pass determinism |
| nil_suite | Sixteen missing TESTO012 diagnostics | Explicit/conversion/local-zero nil, aliases, generic calls, source positions; reassignment, closures, address-taking, globals, interfaces and unrelated calls excluded |
| recursive_subsuite | Six missing TESTO013 diagnostics | Pointer/value/addressed receivers, aliases, generics and BeforeAll; different instances, guards, hooks, providers, method-set differences, options and unrelated calls excluded |
| CLI and custom plugin integration | Missing TESTO012/TESTO013 with the original valid fixture | Real Testo: three nil and two recursive diagnostics alongside TESTO001; exact CLI locations/counts and exit codes; clean child/value/generic/alias/bounded runs |

The pre-existing typed-nil discovery example gained its new TESTO012 expectation;
old discovery expectations were preserved. Integration lint initially found an
unused initial nil assignment in a new clean fixture; the fixture now checks that
initial value before reassignment. No existing rule was weakened.

## Final validation

All commands below succeeded on the reviewed checkout:

| Command | Result |
| --- | --- |
| make test | Passed; external binary test is skipped without its environment variable |
| go test ./... -count=1 | Passed uncached, including CLI tests |
| go vet ./... | Passed |
| go test -race ./... -count=1 | Passed |
| make lint | 0 issues, golangci-lint 2.14.0 |
| go mod tidy -diff | Passed, no diff |
| go mod verify | All modules verified |
| ./bin/bootstrap/golangci-lint custom -v | Rebuilt successfully from current sources |
| TESTOLINT_GOLANGCI_BINARY="$PWD/bin/golangci-lint" go test ./plugin/golangci -run '^TestIntegration$' -count=1 -v | Passed with the freshly rebuilt binary; not skipped |
| go test ./clean -count=1 | Passed in the real-Testo integration module |
| go test -race ./clean -count=1 | Passed in the real-Testo integration module |
| git diff --check | Passed |

Earlier sandbox-only Go package/cache loading failures were retried with approved
cache access and passed; these failures were not counted as analyzer failures or
successful checks. The final Go checks ran with that access. Ordinary test-suite
integration skips are not evidence of external integration success; the explicit
fresh-binary check above supplies that evidence. No required final check remains
unexecuted. CI pins the same golangci-lint version and builds/runs the custom
plugin and clean Testo module.

## Readiness

Ready for documentation review and approval of `docs: prepare v0.3.0 release`.
README installation examples use latest, meaning the latest published release;
local checkout builds include unpublished changes. The changelog records v0.3.0
without a publication date. After documentation approval and commit, verify the
working tree is clean. Tagging and publication require separate authorization.
