# Contributing

See [AGENTS.md](AGENTS.md) for repository guidelines, review and release procedures.
The [analyzer pipeline](docs/pipeline.md) describes execution and diagnostic ownership.

## Build locally

Use the Go version required by [go.mod](go.mod):

```sh
go build -o bin/testolint ./cmd/testolint
./bin/testolint ./...
```

## Rule and discovery changes

1. Add an end-to-end fixture in `analyzer/testdata/src/`, with `// want "TESTOxxx..."`
   diagnostics and a valid example without a diagnostic. Register new fixture
   packages in `analyzer/analyzer_test.go`.
2. Run the fixture before implementation and observe the expected diagnostic
   failure. Loading or compilation errors do not count as the red phase.
3. Make the smallest implementation change, then rerun the fixture and full suite.

```sh
go test ./analyzer -run '^TestAnalyzer/<fixture>$' -count=1
```

CLI and plugin changes likewise require end-to-end tests first in
`cmd/testolint/main_test.go` or `plugin/golangci/integration_test.go`.
Keep rule IDs stable and update the [catalog](docs/rules.md) and affected rule
pages when behavior changes. Documentation checks compare the catalog with the registry.

## Checks

Format modified Go files with `gofmt`, then run:

```sh
make test
go vet ./...
make lint
go mod tidy -diff
```

`make lint` requires golangci-lint v2. Tests cover analyzer diagnostics and CLI
exit codes using real Testo fixtures. The external plugin test is skipped unless
`TESTOLINT_GOLANGCI_BINARY` is set; a skip does not verify integration.

## Coverage

```sh
make test-cov
go tool cover -html=coverage.out
```

`make test-cov` runs the root module's tests without caching, writes `coverage.out`
and a per-function summary in `coverage.txt`, and requires at least **80% total
statement coverage**. The threshold uses statement counts across all packages,
not an average of package percentages or a rounded display value. Both reports
are ignored by Git.

CI runs the same target and uploads `coverage.out` using
[Codecov's GitHub Action](https://github.com/codecov/codecov-action/tree/v5), as in
conformly. To enable uploads, activate `nashabanov/testolint` in Codecov and add
its upload token as the GitHub Actions secret `CODECOV_TOKEN`. Fork pull requests
do not receive repository secrets; tokenless uploads depend on Codecov's settings.
The local 80% gate runs independently of Codecov. Upload errors use the action's
default non-blocking behavior.

CLI tests build and execute a separate binary without coverage instrumentation,
so `cmd/testolint` appears as 0% despite its end-to-end tests. Execution inside
the external custom golangci-lint binary is also outside this profile. Those
integration checks still run separately in CI; this report does not measure them.

## Plugin integration

From the checkout root, install the bootstrap and build local sources using the
checked-in `.custom-gcl.yml`:

```sh
curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b ./bin/bootstrap v2.14.0
./bin/bootstrap/golangci-lint custom -v
TESTOLINT_GOLANGCI_BINARY="$PWD/bin/golangci-lint" go test ./plugin/golangci -run '^TestIntegration$' -count=1 -v
(cd plugin/golangci/testdata/integration && go test ./clean -count=1)
```

Rebuild after source changes. CI builds the custom binary and runs integration
checks; an older binary does not verify the current checkout.

## Benchmarks

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

See [pipeline benchmarks](docs/pipeline.md#measuring-prerequisites) for selection
and prerequisite measurements.

## Go analysis drivers

The analyzer is available as `analyzer.Analyzer` for standard `go/analysis` drivers.
Use `analyzer.New(analyzer.Config{...})` for an independently configured instance;
`analyzer.Analyzer` retains default behavior and has no mutable selection flags.
