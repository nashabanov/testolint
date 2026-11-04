# testolint

A small Go linter for tests written with [Testo](https://github.com/ozontech/testo).
It checks suite signatures, parameter providers, naming, suite execution, and
conservative parallel lifecycle hazards, cleanup misuse, nil test callbacks,
custom context/plugin construction errors and provider/plugin initialization
misuse and sleep-based waits, plus optional Allure formatting, without running tests.
Both the standalone CLI and golangci-lint integration run the same 27 rules.
Runtime semantics were checked against Testo v1.8.0; style integration was checked
against Testo Allure v1.5.0.

## Installation

Requires Go 1.27.1 or newer. Installation examples use the latest published release.

### Standalone

```sh
go install github.com/nashabanov/testolint/cmd/testolint@latest
```

Add your Go binary directory (`go env GOBIN`, or `$(go env GOPATH)/bin` when
GOBIN is empty) to `PATH`. To use the current checkout, build locally:

```sh
go build -o bin/testolint ./cmd/testolint
./bin/testolint ./...
```

### golangci-lint

Uses golangci-lint v2's official [Module Plugin System](https://golangci-lint.run/docs/plugins/module-plugins/).
You need Go, Git, network access, curl and sh. From a testolint checkout:

```sh
curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b ./bin/bootstrap v2.14.0
./bin/bootstrap/golangci-lint custom -v
export PATH="$(pwd)/bin:$PATH"
```

The checked-in `.custom-gcl.yml` builds local sources into `bin/golangci-lint`.
This **custom binary** contains testolint; the bootstrap and your ordinary system
binary do not. Rebuild it after updating the plugin or golangci-lint.

Build in your own project using this `.custom-gcl.yml`:

```yaml
version: v2.14.0
name: golangci-lint
destination: ./bin
plugins:
  - module: github.com/nashabanov/testolint
    import: github.com/nashabanov/testolint/plugin/golangci
    version: latest
```

Install the bootstrap and run `custom` using the commands above. In the project
you want to analyze, create `.golangci.yml` (or merge into your v2 config):

```yaml
version: "2"
linters:
  default: standard
  enable:
    - testolint
  settings:
    custom:
      testolint:
        type: module
        description: Check Testo suite usage.
```

Run `golangci-lint run` with the custom binary on `PATH`. YAML alone cannot add
testolint to a stock binary. To verify the local build, from the checkout root:

```sh
cd plugin/golangci/testdata/integration
go mod download
golangci-lint run
```

This intentionally broken example exits with status 1 and includes:

```text
suite_test.go:8:14: TESTO001: invalid test signature: expected func(T) or func(T, struct{...}) (testolint)
```

## Usage

From your Go project directory:

```sh
testolint ./...
testolint ./internal/integration
```

Test files are included by default. Packages must load and type-check first.
The standalone CLI exits with 0 for clean code, 3 for diagnostics, and a nonzero
status for loading errors. `testolint -h` lists the standard go/analysis flags.
The same installation and invocation commands work in CI.

## Rule configuration

These settings are available since 0.8.0;
older releases do not support them. Build the CLI and custom plugin locally using
the instructions above.

Both integrations use the same selection model:

1. Start with each rule's registry default. Correctness and policy rules are enabled; style rules are opt-in.
2. With `disable-all`, start with no rules enabled instead.
3. Apply explicit `enable` and `disable` lists. An ID in both lists is an error,
   regardless of flag or YAML key order. Duplicate and unknown IDs are errors.

IDs are case-sensitive (`TESTO019`). No class changes a rule's default. Empty
lists preserve the starting selection; `disable-all` with no enabled IDs produces
no diagnostics. A disabled prerequisite may still be analyzed to preserve
suppression and diagnostic ownership. For example, disabling TESTO022 does not
make TESTO021 report fields in a construction graph that contains a cycle.
Selecting TESTO018 does not enable TESTO019, even though they share a checker.
Disabling an owning diagnostic does not transfer its reports to another rule.

CLI flags take comma-separated lists; whitespace around each ID is trimmed.
Each list flag may appear once. A trailing comma or an empty list element is an
error; `-enable=` and `-disable=` mean empty lists. Invalid settings fail during
flag parsing, before package analysis.

```sh
# Preserve defaults except for two policy checks.
./bin/testolint -disable=TESTO006,TESTO012 ./...
# Opt in to the unused-provider policy alone.
./bin/testolint -disable-all -enable=TESTO006 ./...
# Run only nil-callback and construction-cycle checks.
./bin/testolint -disable-all -enable=TESTO019,TESTO022 ./...
```

The equivalent plugin configuration goes under the custom linter's `settings`:

```yaml
version: "2"
linters:
  default: standard
  enable: [testolint]
  settings:
    custom:
      testolint:
        type: module
        settings:
          disable-all: true
          enable: [TESTO019, TESTO022]
          disable: []
```

Plugin lists must be YAML sequences of exact IDs, and `disable-all` must be a
boolean. Unknown setting names, invalid types, unknown/duplicate IDs and conflicts
are rejected when the plugin is constructed. Each driver uses a separate analyzer
instance. Go callers can use `analyzer.New(analyzer.Config{...})`; the existing
`analyzer.Analyzer` retains default behavior and has no mutable selection flags.

See the [diagnostic and configuration contract](docs/contract.md) for stability,
exit codes, suppression and migration details.

### Migrating to 0.9.0

TESTO027 is now enabled by default following the maintainer's project audit.
To keep the previous selection, use `-disable=TESTO027` in the CLI, or:

```yaml
settings:
  disable: [TESTO027]
```

under `linters.settings.custom.testolint` in the plugin configuration. TESTO025
and TESTO026 remain opt-in. Diagnostic messages now end with a rule-specific
link: `TESTOxxx: message (https://github.com/nashabanov/testolint/blob/main/docs/rules/TESTOxxx.md)`.
Update exact-message snapshots; the ID prefix remains unchanged and `Category`
remains the ID. The structured `URL` now points to the same rule page.
Version 0.9.0 is unreleased; build this checkout to use these changes.

## Rules

`correctness` covers Testo contract violations, missing execution and concrete
lifecycle hazards. `policy` covers unused providers, possibly intentional nil
suites and conservative shared-field ownership checks; it does not assert that
every reported program fails. All 25 correctness/policy rules are enabled by default.
`style` covers source formatting; its two rules are disabled by default.
Enable them with `-enable=TESTO025,TESTO026` and apply fixes with `-fix`.

The table is checked against the authoritative registry. See the
[per-rule examples and remedies](docs/rules.md) and the detailed limitations below.

| ID | Name | Class | Default | Purpose |
| --- | --- | --- | --- | --- |
| [`TESTO001`](docs/rules/TESTO001.md) | `test-signature` | correctness | enabled | Tests accept exactly the suite's `T`, optionally a struct, and return no values. |
| [`TESTO002`](docs/rules/TESTO002.md) | `hook-signature` | correctness | enabled | `BeforeAll`, `BeforeEach`, `AfterEach`, `AfterAll` accept exactly `T` and return no values. |
| [`TESTO003`](docs/rules/TESTO003.md) | `missing-case-provider` | correctness | enabled | Exported parameter fields require a `Cases<Field>` method. |
| [`TESTO004`](docs/rules/TESTO004.md) | `case-type-mismatch` | correctness | enabled | Provider slice elements must be assignable to the parameter field type. |
| [`TESTO005`](docs/rules/TESTO005.md) | `cases-signature` | correctness | enabled | Providers accept no arguments and return exactly one slice. |
| [`TESTO006`](docs/rules/TESTO006.md) | `orphan-cases-provider` | policy | enabled | Policy: providers should be referenced by a test parameter; Testo permits unused providers. |
| [`TESTO007`](docs/rules/TESTO007.md) | `test-name` | correctness | enabled | After `Test`, the suffix is empty or begins with a non-lowercase Unicode rune. |
| [`TESTO008`](docs/rules/TESTO008.md) | `cases-name` | correctness | enabled | After `Cases`, the suffix is empty or begins with a non-lowercase Unicode rune. |
| [`TESTO009`](docs/rules/TESTO009.md) | `unexported-param-field` | correctness | enabled | Parameter fields must be exported so reflection can set them. |
| [`TESTO010`](docs/rules/TESTO010.md) | `empty-cases` | correctness | enabled | Providers returning `nil`, an empty slice literal or built-in `make` with constant zero length always return an empty case set. |
| [`TESTO011`](docs/rules/TESTO011.md) | `empty-suite` | correctness | enabled | Suite must contain at least one runnable Testo test. |
| [`TESTO012`](docs/rules/TESTO012.md) | `nil-suite` | policy | enabled | Report a statically nil suite argument to `RunSuite` or `RunSubSuite`. |
| [`TESTO013`](docs/rules/TESTO013.md) | `recursive-subsuite` | correctness | enabled | Report a directly recursive `RunSubSuite` invocation on the current receiver. |
| [`TESTO014`](docs/rules/TESTO014.md) | `shared-before-each-write` | policy | enabled | BeforeEach writes a shared receiver field read after Parallel suspension. |
| [`TESTO015`](docs/rules/TESTO015.md) | `shared-fixture-cleanup` | policy | enabled | Cleanup re-reads a shared BeforeEach field that another execution can replace. |
| [`TESTO016`](docs/rules/TESTO016.md) | `shared-test-context` | policy | enabled | A shared BeforeEach field can expose another execution's test context. |
| [`TESTO017`](docs/rules/TESTO017.md) | `duplicate-parallel-lifecycle` | correctness | enabled | Repeated Parallel calls on the same context across BeforeEach, test and AfterEach. |
| [`TESTO018`](docs/rules/TESTO018.md) | `run-during-cleanup` | correctness | enabled | Direct `testo.Run` on the cleanup context inside its Cleanup callback. |
| [`TESTO019`](docs/rules/TESTO019.md) | `nil-callback` | correctness | enabled | Statically nil callbacks to `RunTest` or an immediately executed `Test` adapter. |
| [`TESTO020`](docs/rules/TESTO020.md) | `premature-teardown` | correctness | enabled | AfterEach closes an `*os.File` before a parallel child uses the same file. |
| [`TESTO021`](docs/rules/TESTO021.md) | `construction-field` | correctness | enabled | Exported fields traversed during custom T/plugin wiring must be pointers. |
| [`TESTO022`](docs/rules/TESTO022.md) | `construction-cycle` | correctness | enabled | Report a type cycle traversed during custom T/plugin initialization. |
| [`TESTO023`](docs/rules/TESTO023.md) | `provider-before-each` | correctness | enabled | A provider dereferences a fresh suite’s nil field before BeforeEach initializes it. |
| [`TESTO024`](docs/rules/TESTO024.md) | `plugin-parent-nil` | correctness | enabled | A Plugin parent interface compared with nil cannot identify the suite root. |
| [`TESTO025`](docs/rules/TESTO025.md) | `blank-after-header` | style | disabled | Separate initial Allure metadata from the test body with a blank line. |
| [`TESTO026`](docs/rules/TESTO026.md) | `blank-between-steps` | style | disabled | Separate adjacent Allure Step statements with a blank line. |
| [`TESTO027`](docs/rules/TESTO027.md) | `no-sleep` | policy | enabled | Avoid time.Sleep in Testo tests and hooks; use polling or condition-based waiting. |

## Example

See [TESTO004](docs/rules/TESTO004.md), [TESTO012](docs/rules/TESTO012.md), [TESTO013](docs/rules/TESTO013.md).

## Parallel lifecycle checks

See [TESTO014](docs/rules/TESTO014.md), [TESTO015](docs/rules/TESTO015.md), [TESTO016](docs/rules/TESTO016.md), [TESTO017](docs/rules/TESTO017.md).

## Execution and cleanup safety

See [TESTO018](docs/rules/TESTO018.md), [TESTO019](docs/rules/TESTO019.md), [TESTO020](docs/rules/TESTO020.md).

## Custom context and plugin construction

See [TESTO021](docs/rules/TESTO021.md), [TESTO022](docs/rules/TESTO022.md).

## Provider and plugin initialization

See [TESTO023](docs/rules/TESTO023.md), [TESTO024](docs/rules/TESTO024.md).

## Scope and limitations

See [shared analysis scope](docs/analysis-scope.md) and the boundaries on each
[rule page](docs/rules.md). These static checks do not execute tests or prove
race freedom. Source-level suppression is not implemented.

## Development

```sh
make test
make lint
go vet ./...
go mod tidy
```

Run analyzer benchmarks with `make bench`. They report time and allocations for
packages without Testo, valid suites, invalid signatures, and embedded suites with
promoted parameterized tests and case providers, each at sizes 1, 10 and 100.
Complex workloads have one base and one derived suite per size unit. Inputs use
the local Testo stub; source generation, package loading, type checking and an
initial diagnostic-count check are excluded. Each measured iteration runs the
complete default analyzer; diagnostic formatting and deduplication are included,
while the report callback only counts diagnostics. CLI and golangci-lint startup,
loading, and output rendering are not measured.

For repeated samples or comparisons with `benchstat`:

```sh
make bench BENCH_ARGS='-count=5 -benchtime=1s' > before.txt
# Repeat after the change, saving to after.txt.
benchstat before.txt after.txt
```

`make lint` requires golangci-lint v2. Tests cover analyzer diagnostics and CLI
exit codes using real Testo fixtures. After building the custom binary:

```sh
TESTOLINT_GOLANGCI_BINARY="$PWD/bin/golangci-lint" go test ./plugin/golangci -run '^TestIntegration' -count=1 -v
(cd plugin/golangci/testdata/integration && go test ./clean)
```

The external binary test is otherwise skipped. CI builds the custom binary and
runs this check. The analyzer is also available as `analyzer.Analyzer` for
standard `go/analysis` drivers.

## License

[MIT](LICENSE).
