# testolint

A small Go linter for tests written with [Testo](https://github.com/ozontech/testo).
It checks suite signatures, parameter providers, naming, and suite execution
without running tests.
Both the standalone CLI and golangci-lint integration run the same thirteen rules.
Runtime semantics were checked against Testo v1.8.0.

## Installation

Requires Go 1.27.1 or newer. The commands using `v0.2.0` below become available
**after the release tag is published**.

### Standalone

```sh
go install github.com/nashabanov/testolint/cmd/testolint@v0.2.0
```

Add your Go binary directory (`go env GOBIN`, or `$(go env GOPATH)/bin` when
GOBIN is empty) to `PATH`. Before the tag exists, build from this checkout:

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

After release, build in your own project using this `.custom-gcl.yml`:

```yaml
version: v2.14.0
name: golangci-lint
destination: ./bin
plugins:
  - module: github.com/nashabanov/testolint
    import: github.com/nashabanov/testolint/plugin/golangci
    version: v0.2.0
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

## Rules

| ID | Purpose |
| --- | --- |
| `TESTO001` | Tests accept exactly the suite's `T`, optionally a struct, and return no values. |
| `TESTO002` | `BeforeAll`, `BeforeEach`, `AfterEach`, `AfterAll` accept exactly `T` and return no values. |
| `TESTO003` | Exported parameter fields require a `Cases<Field>` method. |
| `TESTO004` | Provider slice elements must be assignable to the parameter field type. |
| `TESTO005` | Providers accept no arguments and return exactly one slice. |
| `TESTO006` | Policy: providers should be referenced by a test parameter; Testo permits unused providers. |
| `TESTO007` | After `Test`, the suffix is empty or begins with a non-lowercase Unicode rune. |
| `TESTO008` | After `Cases`, the suffix is empty or begins with a non-lowercase Unicode rune. |
| `TESTO009` | Parameter fields must be exported so reflection can set them. |
| `TESTO010` | Providers whose body is a single return of `nil` or an empty slice literal always return an empty case set. |
| `TESTO011` | Suite must contain at least one runnable Testo test. |
| `TESTO012` | Report a statically nil suite argument to `RunSuite` or `RunSubSuite`. |
| `TESTO013` | Report a directly recursive `RunSubSuite` invocation on the current receiver. |

For `testo.Suite[T]`, tests have the form `TestX(t T)` or
`TestX(t T, p struct{ Age int })`. Named structs and named slice returns are
accepted, as are aliases. Providers match **field names**, not test names, and
can be shared. Assignability follows Go semantics: `[]int` can supply an `any`
field, but `[]any` cannot supply an `int` field.

`Test`, `Test1`, `Test_Foo`, and `TestÉ` are valid names; `Testfoo` is not.
Testo ignores a method named exactly `Cases`, including its signature.
Methods without the exact `Test` or `Cases` prefix are ignored, except hooks.
`TESTO011` counts tests with names recognized by Testo, including promoted and
imported tests. Malformed names do not count; signature errors are reported
separately by `TESTO001`. Empty-suite diagnostics point to the suite type name.

Unexported fields receive only `TESTO009`, without missing/type provider
messages. Invalid test signatures suppress dependent parameter checks; invalid
provider signatures suppress type mismatch checks. Identical diagnostics at the
same position are emitted once. Independent problems may still produce multiple
messages. All rules are enabled; standalone rule selection and source-level
suppression are not implemented.

## Example

```go
package integration

import "github.com/ozontech/testo"

type Suite struct{ testo.Suite[*testo.T] }

func (Suite) TestAge(t *testo.T, p struct{ Age int }) {}
func (Suite) CasesAge() []string { return []string{"18"} }
```

```text
suite_test.go:8:14: TESTO004: CasesAge provides string, but parameter "Age" expects int
```

Return `[]int{18}` from `CasesAge` to satisfy the `Age int` parameter.

For suite execution, `TESTO012` reports at the suite argument:

```go
var suite *Suite = nil
testo.RunSuite(t, suite) // TESTO012: suite argument is statically nil
```

Pass an initialized instance such as `&Suite{}`. This rule reports a known nil
value, not a guaranteed panic. In Testo v1.8.0, inherited value-receiver hooks
panic when invoked through a nil suite pointer, but explicitly nil-safe pointer
hooks and tests can run successfully. Reporting those intentional nil suites is
also part of this rule's policy; method bodies are not analyzed for nil safety.

`TESTO013` detects direct receiver recursion in a narrowly defined execution
path, for example:

```go
type RecursiveSuite struct{ testo.Suite[*testo.T] }

func (s *RecursiveSuite) TestRecursive(t *testo.T) {
    testo.RunSubSuite(t, s) // TESTO013: recursive RunSubSuite invocation
}
```

The nested run selects the same method on the receiver again. Execute a child
suite with its own tests instead. The diagnostic points to the `RunSubSuite`
call. Different instances of the same type are not reported. Runtime plugins
and filtering may interrupt execution; the message identifies direct recursion
and does not claim that every run must loop forever.

## Scope and limitations

- Suite discovery uses Testo v1.8.0's private marker and recognizes package-level
  named structs embedding `testo.Suite[T]` directly or indirectly, including
  aliases and pointer embedding. Aliases do not duplicate suites. Generic suite
  declarations are skipped by suite rules; instantiated calls support run checks.
- Suite rules use the pointer method set for declared, promoted and imported
  methods, respecting Go shadowing and ambiguity. Local diagnostics point to
  methods or fields; imported-method diagnostics point to the local suite type.
- Run discovery uses Go types and supports import aliases, inferred/explicit
  generics and unchanged direct local function aliases. Globals, alias chains
  and wrappers are not resolved.
- `TESTO012` checks explicit nil, pointer conversions of nil and local pointer
  declarations with a nil initializer or zero value (`var`/`:=`). Any assignment
  or address-taking disables local inference, even after the call or in a closure.
  Other values, interface variables and control-flow-dependent nilness are skipped.
- `TESTO013` requires a sole `RunSubSuite` call without options, passing the
  method's own Testo parameter and receiver (or a value receiver's address).
  The argument's method set must select that same method. Suites must have valid
  regular tests and no providers: `BeforeAll` is checked, or the sole test when
  `BeforeAll`/`BeforeEach` are inherited defaults. Other execution paths, imported
  bodies, plugins and runtime filters are outside this check's scope.
- `TESTO010` checks valid local providers whose sole statement returns nil or an
  empty slice literal, including named slices and aliases. Imported bodies and
  other expressions or control flow are skipped. Invalid signatures/names
  suppress this check; malformed parameterized tests make orphan checks conservative.
- Other lifecycle/parallelism checks, plugin hooks, standalone `Run`/`RunTest`
  and broader provider logic are outside this release's scope. No SSA or
  interprocedural analysis is used.

## Development

```sh
make test
make lint
go vet ./...
go mod tidy
```

`make lint` requires golangci-lint v2. Tests cover analyzer diagnostics and CLI
exit codes using real Testo fixtures. After building the custom binary:

```sh
TESTOLINT_GOLANGCI_BINARY="$PWD/bin/golangci-lint" go test ./plugin/golangci -run '^TestIntegration$' -count=1 -v
(cd plugin/golangci/testdata/integration && go test ./clean)
```

The external binary test is otherwise skipped. CI builds the custom binary and
runs this check. The analyzer is also available as `analyzer.Analyzer` for
standard `go/analysis` drivers.

## License

[MIT](LICENSE).
