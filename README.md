# testolint

[![CI](https://github.com/nashabanov/testolint/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/nashabanov/testolint/actions/workflows/ci.yml)
[![Codecov](https://codecov.io/github/nashabanov/testolint/graph/badge.svg?branch=main)](https://codecov.io/github/nashabanov/testolint)
[![Latest release](https://img.shields.io/github/v/release/nashabanov/testolint)](https://github.com/nashabanov/testolint/releases/latest)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

A Go linter for tests written with [Testo](https://github.com/ozontech/testo).
It finds incorrect suite signatures, parameter providers and lifecycle usage
without running tests, with optional formatting and test structure checks for Testo Allure.

## Installation

Requires Go 1.27.1 or newer. Install the latest published release:

### Standalone

```sh
go install github.com/nashabanov/testolint/cmd/testolint@latest
testolint ./...
```

Run from your Go project directory. Add `go env GOBIN` (or `$(go env GOPATH)/bin`
when GOBIN is empty) to `PATH` if needed.

### golangci-lint

Build a custom golangci-lint v2 binary with testolint using the
[plugin setup guide](docs/golangci-lint.md). YAML alone cannot add the plugin
to a stock binary.

## Example

A provider must return values assignable to its parameter field:

```go
type Suite struct{ testo.Suite[*testo.T] }

func (Suite) TestAge(t *testo.T, p struct{ Age int }) {}
func (Suite) CasesAge() []string { return []string{"18"} } // TESTO004
```

Fix the provider by returning integers:

```go
func (Suite) CasesAge() []int { return []int{18} }
```

See [TESTO004](docs/rules/TESTO004.md) for the diagnostic and a complete source example.

## Usage

```sh
testolint ./...
testolint ./internal/integration
```

Test files are included by default; packages must load and type-check first.
Normal CLI runs exit 0 for clean code and 3 for diagnostics; see
[all exit codes](docs/contract.md#driver-results). `testolint -h` lists flags.
The same commands work in CI.

## Rule configuration

```sh
# Disable a policy check while keeping the other defaults.
testolint -disable=TESTO006 ./...
# Enable Allure formatting checks and apply available fixes.
testolint -enable=TESTO025,TESTO026 -fix ./...
```

See [configuration](docs/configuration.md) for CLI and YAML settings.

<a id="migrating-to-090"></a>
For upgrades, see [migration to 0.9.0](docs/configuration.md#migrating-to-090)
and the [changelog](CHANGELOG.md).

## Rules

The [rule catalog](docs/rules.md) lists all 31 checks, examples and remedies.
The 25 correctness/policy checks are enabled by default; six style checks are opt-in.
Policy findings can flag intentional code and do not necessarily indicate a runtime failure.

<a id="parallel-lifecycle-checks"></a>
<a id="execution-and-cleanup-safety"></a>
<a id="custom-context-and-plugin-construction"></a>
<a id="provider-and-plugin-initialization"></a>
Browse [parallel lifecycle](docs/rules/TESTO014.md),
[execution and cleanup](docs/rules/TESTO018.md),
[context construction](docs/rules/TESTO021.md) and
[provider initialization](docs/rules/TESTO023.md) checks in the catalog.

In the unreleased 0.10.0 work, [TESTO028](docs/rules/TESTO028.md) checks configured
required Allure metadata in valid local tests and their effective BeforeEach.
Unknown sources suppress unproven absence; this is not a proof of runtime metadata.

### Allure operation configuration

Optional TESTO029/031 checks support groups of typed client methods and package
functions, with exact exceptions that exclude local helper bodies but still
analyze evaluated arguments; see [configuration](docs/configuration.md#grouped-operations).
TESTO030 defaults to relaxed Action/Assert sequences; strict AAA is opt-in.
Only explicit Setup establishes preparation order; literal and factory-based
input construction is neutral. Coverage follows one local helper transition,
so unclassified calls in deeper or imported helpers may remain unchecked.

## Scope and limitations

Static checks cover supported source patterns and do not prove race freedom.
Source-level suppression is not implemented. See [shared scope](docs/analysis-scope.md),
individual rule boundaries and the [diagnostic contract](docs/contract.md).

### Verified compatibility

The unreleased 0.10.0 candidate is verified with:

| Go | Testo | Testo Allure | Custom golangci-lint |
| --- | --- | --- | --- |
| 1.27.1 | 1.8.0 | 1.5.0 | 2.14.0 |

CI uses this combination for CLI and freshly built plugin integration checks.
The minimum Go requirement does not imply verification of every newer toolchain
or dependency version. Other combinations are unverified.

## Development

See [Contributing](CONTRIBUTING.md) for local builds, tests and benchmarks,
and the [analyzer pipeline](docs/pipeline.md) for implementation details.

## License

[MIT](LICENSE).
