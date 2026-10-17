# testolint

A small Go linter for tests written with [Testo](https://github.com/ozontech/testo).
It checks suite signatures, parameter providers, naming, suite execution, and
conservative parallel lifecycle hazards without running tests.
Both the standalone CLI and golangci-lint integration run the same seventeen rules.
Runtime semantics were checked against Testo v1.8.0.

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
| `TESTO014` | BeforeEach writes a shared receiver field read after Parallel suspension. |
| `TESTO015` | Cleanup re-reads a shared BeforeEach field that another execution can replace. |
| `TESTO016` | A shared BeforeEach field can expose another execution's test context. |
| `TESTO017` | Repeated Parallel calls on the same context across BeforeEach, test and AfterEach. |

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

## Parallel lifecycle checks

Testo reuses a pointer suite instance. BeforeEach, the test and AfterEach share
one per-test T; BeforeAll and AfterAll use a separate suite context. Parallel
suspends the current test until its parent has scheduled the sibling tests.
A later serial or parallel execution can run BeforeEach and replace the suite
field while an earlier parallel consumer is suspended. Testo collects regular
methods first, then parameterized methods, in method-name order within each
group.
A serial execution completed before the consumer does not establish an overwrite.
Cleanup runs later and can re-read the other execution's resource. These are
lifecycle hazards; TESTO014–016 are
conservative linter policy, not a claim that every run fails. TESTO017 identifies
a repeated Parallel call, which fails at runtime with default Testo behavior.
Plugins and runtime filtering can change these conclusions.

```go
func (s *Suite) BeforeEach(t *testo.T) { s.resource = new(int) }
func (s *Suite) TestA(t *testo.T) {
    t.Parallel()
    t.Log(s.resource) // TESTO014: TestB runs BeforeEach while TestA is suspended
}
func (s *Suite) TestB(t *testo.T) { t.Parallel() }
```

Capture the resource in a local **before** Parallel, or keep fixture ownership
local to the test. In BeforeEach, use `resource := s.resource` and close over
`resource` in Cleanup rather than re-reading `s.resource`. Keep T in the test's
parameter or a local. Choose one Parallel call per test context across its
hooks and body. TESTO015 takes precedence over TESTO014 for the same field in
the same execution; context fields use TESTO016. Shared hook diagnostics are
reported once at each source position.

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
- TESTO014–017 support the actual `*testo.T` type (including aliases), local
  straight-line bodies, resolved direct calls on the context parameter, and
  direct fields of the suite's own pointer receiver. Supported statements are
  single assignments, local captures, `new`, `Log`, `Logf`, `Helper`, `Parallel`,
  and direct `Cleanup(func() { ... })` closures containing local captures and Log.
  A terminal method call on a direct pointer field (such as
  `s.resource.Close()`) is also
  supported inside Cleanup: evaluating its receiver is the hazardous read;
  the method body and anything it might synchronize are not inferred.
  Value receivers can participate in TESTO017; their fields are not analyzed.
- TESTO014–016 require a BeforeEach write before the first Parallel suspension
  and an established later execution on the same receiver, serial or parallel.
  A single parameterized method qualifies when every used provider returns a
  non-empty, unkeyed slice literal of constants in its sole return statement.
  Other case counts are unknown.
  Fields written by tests or AfterEach are skipped. Storing T in a test body,
  shared referenced storage behind value receivers, and nested fields are outside
  this model. Unknown case counts do not suppress a proof established by other
  known executions, but cannot establish an interfering execution themselves.
  Proven empty providers are distinguished from unknown counts: an empty
  path cannot interfere, even if its body is unsupported. Provider sizes are
  inspected once per suite.
- TESTO017 requires an established execution of the checked test: a regular
  method or a parameterized method with known non-empty constant cases. It does
  not require sibling executions or known sibling case counts. Collection errors
  (including missing or incompatible providers) still suppress the check.
- Invalid lifecycle/test signatures and collection errors suppress dependent
  checks. BeforeAll/AfterAll signatures are validated; their bodies do not enter
  per-test timelines. Suite-context Parallel does not reorder the inner siblings.
  An unknown test body (calls, branches, loops, early exits, aliases or unavailable
  bodies) blocks inference about its own events and later stages, without
  suppressing known sibling paths.
  TESTO017 uses complete lifecycle stages up to the first unknown body. An
  unsupported AfterEach does not hide a duplicate already established by
  BeforeEach and the test; an unknown BeforeEach or test blocks reasoning
  about later stages. Unknown hooks are never treated as empty hooks.
  TESTO014–016 require understood BeforeEach/AfterEach bodies and a complete
  consumer path. A later sibling can establish the BeforeEach overwrite even
  when its test body is unsupported, provided the local body never references
  its suite receiver or exposes T beyond the resolved standard operations.
  Such a body supplies no reads, writes, or Parallel facts.
  Any receiver reference in unknown code (including aliases, field addresses
  and helper arguments), context escapes (T can expose the suite through
  `testo.Reflect`), or an unavailable/promoted body suppresses shared-field
  proofs. Known test/AfterEach writes still exclude the affected field.
  No partial events from an unsupported body are used. Receiver access through
  globals or other hidden aliases and helper side effects are not inferred;
  the supported ownership model assumes no such external receiver access.
  Helper is accepted only on the actual context parameter with its resolved
  standard method; arbitrary helpers and context aliases remain unknown.
  Other expressions, including arithmetic, indexing and conversions,
  remain unsupported. Identical constant field assignments retain existing policy.
  Synchronization, helper internals, subtests,
  custom T types, plugin overrides and runtime filters are not modeled. A clean
  result is not proof of race freedom. No general concurrency analysis is used.
- Other lifecycle checks, plugin hooks, standalone `Run`/`RunTest`
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
