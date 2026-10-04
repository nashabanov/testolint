# testolint

`testolint` is a static analyzer for Go tests written with the
[Testo](https://github.com/ozontech/testo) framework. It checks suite method
signatures, parameterized test providers, and method naming without running tests.

The analyzer uses Go type information to recognize Testo suites and compare test
parameters with their case providers.

## Installation

### Standalone

The current module requires Go 1.27.1 or newer. Until a release tag is available,
build from this checkout:

```sh
mkdir -p bin
go build -o bin/testolint ./cmd/testolint
./bin/testolint ./...
```

After a release, install a specific published tag (replace `vX.Y.Z` with that tag):

```sh
go install github.com/nashabanov/testolint/cmd/testolint@vX.Y.Z
```

Ensure the Go binary directory is on your `PATH`.

### golangci-lint

This integration uses the official [Module Plugin System](https://golangci-lint.run/docs/plugins/module-plugins/),
which upstream recommends over the Go plugin system. The official `custom`
command builds a separate golangci-lint binary containing testolint. Both
frontends use the same `analyzer.Analyzer` and all nine rules.

Prerequisites: Go 1.27.1 or newer, Git, and network access for the custom build.
The commands below also use curl and sh to install the pinned upstream bootstrap
binary, following the [official installation instructions](https://golangci-lint.run/docs/welcome/install/local/).

From the root of a testolint checkout, install golangci-lint v2.14.0 and build:

```sh
curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b ./bin/bootstrap v2.14.0
./bin/bootstrap/golangci-lint custom -v
```

The checked-in [.custom-gcl.yml](.custom-gcl.yml) uses local sources:

```yaml
version: v2.14.0
name: golangci-lint
destination: ./bin
plugins:
  - module: github.com/nashabanov/testolint
    import: github.com/nashabanov/testolint/plugin/golangci
    path: .
```

Run the build command from this repository's root. It creates `bin/golangci-lint`;
the bootstrap and your system golangci-lint remain separate binaries. Put this
custom binary's directory first on `PATH` to use the ordinary command:

```sh
export PATH="$(pwd)/bin:$PATH"
```

In the project you want to analyze, create `.golangci.yml` (or merge these entries
into an existing v2 configuration):

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

Then run `golangci-lint run` in that project. Standard linters and testolint run
together. A system binary built without the plugin cannot load it just from YAML.

To verify with the included [example project](plugin/golangci/testdata/integration),
run from this repository's root after exporting `PATH`:

```sh
cd plugin/golangci/testdata/integration
go mod download
golangci-lint run
```

The intentionally broken `TestBroken` method produces a diagnostic like:

```text
suite_test.go:8:14: TESTO001: invalid test signature: expected func(T) or func(T, struct{...}) (testolint)
```

The lint command exits with status 1 because the example contains a violation.

For a released version, replace the local `path` entry in `.custom-gcl.yml` with
`version` and run `golangci-lint custom` in your own project:

```yaml
version: v2.14.0
name: golangci-lint
destination: ./bin
plugins:
  - module: github.com/nashabanov/testolint
    import: github.com/nashabanov/testolint/plugin/golangci
    version: vX.Y.Z # Replace with an actual published testolint tag.
```

No testolint tag is required for the local setup above; no release is created by
this integration. Rebuild the custom binary when changing the pinned versions
or updating local plugin sources.

## Quick start

Run the analyzer from your Go project directory:

```sh
testolint ./...
```

You can also check a specific package:

```sh
testolint ./internal/integration
```

Test files are included by default. The analyzer reports diagnostics with a file
location and a rule ID:

```text
suite_test.go:18:3: TESTO003: parameter "Role" requires CasesRole
```

Packages must load and type-check successfully before they can be analyzed.

## Example

A parameterized test binds each field of its second argument to a provider named
`Cases<FieldName>`. For example, the `UserID` field uses `CasesUserID`:

```go
package integration

import (
    "testing"

    "github.com/ozontech/testo"
)

type UserSuite struct {
    testo.Suite[*testo.T]
}

func (UserSuite) CasesUserID() []int {
    return []int{1, 2, 3}
}

func (UserSuite) TestUser(t *testo.T, p struct{ UserID int }) {
    if p.UserID <= 0 {
        t.Error("user ID must be positive")
    }
}

func TestUsers(t *testing.T) {
    testo.RunSuite(t, new(UserSuite))
}
```

This suite satisfies all currently implemented rules.

## Rules

All nine rules run on every discovered suite. Individual rule selection and
source-level suppression are not currently implemented.

| ID | Rule | Checks |
| --- | --- | --- |
| `TESTO001` | `invalid-test-signature` | Tests accept the suite's `T`, optionally followed by a struct, and return no values. |
| `TESTO002` | `invalid-hook-signature` | Suite hooks accept exactly the suite's `T` and return no values. |
| `TESTO003` | `missing-cases-provider` | Every field in a test's parameter struct has a matching `Cases<Field>` method. |
| `TESTO004` | `cases-type-mismatch` | The provider's slice element type is assignable to the parameter field type. |
| `TESTO005` | `invalid-cases-signature` | Case providers accept no parameters and return exactly one slice. |
| `TESTO006` | `orphan-cases-provider` | A provider with a valid name is referenced by at least one test parameter field in the suite. |
| `TESTO007` | `malformed-test-name` | Methods starting with `Test` have an empty suffix or a first suffix rune that is not lowercase. |
| `TESTO008` | `malformed-cases-name` | Methods starting with `Cases` have an empty suffix or a first suffix rune that is not lowercase. |
| `TESTO009` | `unexported-param-field` | Every field in a test's parameter struct is exported so Testo can set it through reflection. |

### Test and hook signatures

For a suite embedding `testo.Suite[T]`, the accepted test forms are:

```go
func (Suite) TestSimple(t T) {}
func (Suite) TestParameterized(t T, p struct{ UserID int }) {}
```

The first argument must match the suite's type argument exactly. A named type
whose underlying type is a struct is also accepted as the second argument.
Returning a value, accepting extra arguments, or using a non-struct second
argument produces `TESTO001`.

`TESTO002` checks `BeforeAll`, `BeforeEach`, `AfterEach`, and `AfterAll`:

```go
func (Suite) BeforeEach(t T) {}
func (Suite) AfterEach(t T) {}
```

### Parameterized tests

The analyzer binds providers by field name, rather than by test method name.
A provider can be shared by multiple tests in the same suite.

Parameter fields must be exported according to Go identifier semantics.
Unexported fields receive `TESTO009` at the field and are excluded from
`TESTO003` and `TESTO004` checks to avoid secondary provider diagnostics.

```go
func (Suite) TestUser(t T, p struct {
    UserID int
    Role   string // TESTO003 if CasesRole is missing
}) {}

func (Suite) CasesUserID() []string { // TESTO004: expected []int
    return []string{"admin"}
}
```

Type comparison uses assignability from the provider element to the parameter
field, matching Testo runtime behavior. For example, an `int` element can populate
an `any` field, but an `any` element cannot populate an `int` field.
Named slice types are accepted as provider return types.

```go
func (Suite) CasesRole(limit int) []string { // TESTO005: parameters are not allowed
    return nil
}

func (Suite) CasesUnused() []bool { // TESTO006 if no test has an Unused field
    return nil
}
```

Providers with malformed names receive `TESTO008` and are excluded from the
orphan-provider check. Other independent checks can report multiple diagnostics
on the same method, such as `TESTO005` and `TESTO006`.

### Naming

The naming rules match Testo: the suffix may be empty; otherwise, its first
Unicode rune must not be lowercase:

```go
func (Suite) TestUser(t T) {}       // valid name
func (Suite) Testuser(t T) {}       // TESTO007
func (Suite) Test(t T) {}           // valid name
func (Suite) Casesuser() []int {    // TESTO008
    return nil
}
```

Digits and underscores immediately after the prefix are accepted, as are Unicode
runes without lowercase status. `Cases` alone is also a valid provider name.
Methods without the exact `Test` or `Cases` prefix, such as `testUser`,
are not analyzed as tests or providers.

## Scope and limitations

- Suites are recognized by directly embedding `Suite[T]` from
  `github.com/ozontech/testo`, including pointer embedding. Recognition uses the
  package path and type name, so an unrelated type named `Suite` is ignored.
- Only methods declared on the suite type in the analyzed package are collected.
  Promoted methods and suites embedding another user-defined suite are not
  included in the discovery model.
- The analyzer does not verify that a suite is passed to `testo.RunSuite`.
- Empty suites, parallel execution, shared state, plugin hooks, and standalone
  `testo.Run` or `testo.RunTest` callbacks are not currently checked.
- Provider bodies and returned values are not evaluated. For example, a nil or
  empty case slice does not produce a diagnostic.

## CI

Add installation and analysis to your existing Go CI job:

```sh
go install github.com/nashabanov/testolint/cmd/testolint@vX.Y.Z
testolint ./...
```

Replace `vX.Y.Z` with a specific published testolint tag.
The command exits unsuccessfully when it reports diagnostics or cannot analyze
the requested packages.


## Development

From a checkout of this repository, run the analyzer with:

```sh
go run ./cmd/testolint ./...
```

Run tests and lint checks:

```sh
make test
make lint
```

After building the custom binary, run the opt-in end-to-end test from the
repository root. It executes golangci-lint against the example project and checks
both the diagnostic and the exit status:

```sh
TESTOLINT_GOLANGCI_BINARY="$PWD/bin/golangci-lint" go test ./plugin/golangci -run '^TestIntegration$' -count=1 -v
```

The default test run skips this external binary test; the adapter contract test
always runs and verifies that it returns the existing analyzer with type loading.

`make test` runs `go test ./...`. `make lint` requires `golangci-lint` v2 and uses
the repository's [.golangci.yml](.golangci.yml) configuration.

The analyzer is built on `golang.org/x/tools/go/analysis` and is exposed as
`analyzer.Analyzer` for integration into custom analysis drivers.

Rule implementations live in [analyzer](analyzer). Test fixtures live in
[analyzer/testdata/src](analyzer/testdata/src) and use `analysistest` comments to
declare expected diagnostics:

```go
func (Suite) Testinvalid(t T) {} // want "TESTO007"
```

When adding a rule, register it in [analyzer/rules.go](analyzer/rules.go), add
fixtures for invalid and valid usage, and include any new fixture package in
[analyzer/analyzer_test.go](analyzer/analyzer_test.go).

## License

[MIT](LICENSE).
