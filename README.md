# testolint

A small Go linter for tests written with [Testo](https://github.com/ozontech/testo).
It checks suite signatures, parameter providers, naming, suite execution, and
conservative parallel lifecycle hazards, cleanup misuse, nil test callbacks,
custom context/plugin construction errors and provider/plugin initialization
misuse without running tests.
Both the standalone CLI and golangci-lint integration run the same twenty-four rules.
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

## Rule configuration

These settings are available in the current checkout for the unpublished 0.8.0;
older releases do not support them. Build the CLI and custom plugin locally using
the instructions above.

Both integrations use the same selection model:

1. Start with each rule's registry default. All current rules are enabled.
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

Migration from 0.7.0 requires no selection settings: diagnostics and all defaults
are preserved. Tooling consuming analysis diagnostics now receives `Category`
equal to the rule ID and a documentation URL from the registry. Messages have one
`TESTO###: ` prefix. All URLs currently use the existing
[Rules section](https://github.com/nashabanov/testolint#rules). Full reports preserve
ranges, related locations and suggested fixes; only identical complete reports
are deduplicated.

## Rules

`correctness` covers Testo contract violations, missing execution and concrete
lifecycle hazards. `policy` covers unused providers, possibly intentional nil
suites and conservative shared-field ownership checks; it does not assert that
every reported program fails. All 24 rules remain enabled by default.

The table is checked against the authoritative registry. See the
[per-rule examples and remedies](docs/rules.md) and the detailed limitations below.

| ID | Name | Class | Default | Purpose |
| --- | --- | --- | --- | --- |
| `TESTO001` | `test-signature` | correctness | enabled | Tests accept exactly the suite's `T`, optionally a struct, and return no values. |
| `TESTO002` | `hook-signature` | correctness | enabled | `BeforeAll`, `BeforeEach`, `AfterEach`, `AfterAll` accept exactly `T` and return no values. |
| `TESTO003` | `missing-case-provider` | correctness | enabled | Exported parameter fields require a `Cases<Field>` method. |
| `TESTO004` | `case-type-mismatch` | correctness | enabled | Provider slice elements must be assignable to the parameter field type. |
| `TESTO005` | `cases-signature` | correctness | enabled | Providers accept no arguments and return exactly one slice. |
| `TESTO006` | `orphan-cases-provider` | policy | enabled | Policy: providers should be referenced by a test parameter; Testo permits unused providers. |
| `TESTO007` | `test-name` | correctness | enabled | After `Test`, the suffix is empty or begins with a non-lowercase Unicode rune. |
| `TESTO008` | `cases-name` | correctness | enabled | After `Cases`, the suffix is empty or begins with a non-lowercase Unicode rune. |
| `TESTO009` | `unexported-param-field` | correctness | enabled | Parameter fields must be exported so reflection can set them. |
| `TESTO010` | `empty-cases` | correctness | enabled | Providers returning `nil`, an empty slice literal or built-in `make` with constant zero length always return an empty case set. |
| `TESTO011` | `empty-suite` | correctness | enabled | Suite must contain at least one runnable Testo test. |
| `TESTO012` | `nil-suite` | policy | enabled | Report a statically nil suite argument to `RunSuite` or `RunSubSuite`. |
| `TESTO013` | `recursive-subsuite` | correctness | enabled | Report a directly recursive `RunSubSuite` invocation on the current receiver. |
| `TESTO014` | `shared-before-each-write` | policy | enabled | BeforeEach writes a shared receiver field read after Parallel suspension. |
| `TESTO015` | `shared-fixture-cleanup` | policy | enabled | Cleanup re-reads a shared BeforeEach field that another execution can replace. |
| `TESTO016` | `shared-test-context` | policy | enabled | A shared BeforeEach field can expose another execution's test context. |
| `TESTO017` | `duplicate-parallel-lifecycle` | correctness | enabled | Repeated Parallel calls on the same context across BeforeEach, test and AfterEach. |
| `TESTO018` | `run-during-cleanup` | correctness | enabled | Direct `testo.Run` on the cleanup context inside its Cleanup callback. |
| `TESTO019` | `nil-callback` | correctness | enabled | Statically nil callbacks to `RunTest` or an immediately executed `Test` adapter. |
| `TESTO020` | `premature-teardown` | correctness | enabled | AfterEach closes an `*os.File` before a parallel child uses the same file. |
| `TESTO021` | `construction-field` | correctness | enabled | Exported fields traversed during custom T/plugin wiring must be pointers. |
| `TESTO022` | `construction-cycle` | correctness | enabled | Report a type cycle traversed during custom T/plugin initialization. |
| `TESTO023` | `provider-before-each` | correctness | enabled | A provider dereferences a fresh suite’s nil field before BeforeEach initializes it. |
| `TESTO024` | `plugin-parent-nil` | correctness | enabled | A Plugin parent interface compared with nil cannot identify the suite root. |

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

## Execution and cleanup safety

TESTO018 reports a direct `testo.Run` inside `t.Cleanup(func() { ... })` when
both operations use the same unchanged local `*testo.T` variable. Testo delegates
Run to the underlying testing context, where running a subtest during Cleanup
panics. Run subtests before cleanup; use Cleanup to release their resources.

```go
t.Cleanup(func() {
    testo.Run(t, "child", nil) // TESTO018, even though nil is valid for Run
})
```

TESTO019 reports a statically nil function supplied to `testo.RunTest`. It also
checks `testo.Test` when the returned adapter is immediately called or passed
directly to the resolved `testing.T.Run` method. Test creates an adapter without
invoking the callback; the failure occurs when that adapter runs. Supply a
non-nil function, such as `func(t *testo.T) {}`, for an intentionally empty test.
**`testo.Run(t, name, nil)` is valid** and uses a no-op callback.

```go
testo.RunTest[*testo.T](t, nil)       // TESTO019
t.Run("child", testo.Test[*testo.T](nil)) // TESTO019 when the adapter runs
```

These are runtime-failure checks under the default Testo behavior, verified
against v1.8.0. Plugins can override Cleanup or filter singleton execution;
such runtime changes are not inferred.

TESTO020 detects premature teardown of an `*os.File`: AfterEach calls Close
before a parallel child resumes and calls Read, Write or Stat on the same direct
suite field. These operations return an error after Close; ignored errors do not
necessarily fail the test. The rule combines Testo lifecycle order with the
concrete os.File contract, rather than interpreting arbitrary Close-like names.
Move the Close into `t.Cleanup(func() { _ = s.file.Close() })` in AfterEach,
which runs after the children complete, or keep the child serial.

```go
func (s *Suite) TestFile(t *testo.T) {
    testo.Run(t, "child", func(child *testo.T) {
        child.Parallel()
        _, _ = s.file.Read(nil)
    })
}
func (s *Suite) AfterEach(t *testo.T) {
    _ = s.file.Close() // TESTO020: the parallel child has not resumed yet
}
```

## Custom context and plugin construction

Testo constructs a fresh custom T and fills its accessible pointer fields before
wiring plugins. TESTO021 reports an exported non-pointer field reached during
this wiring, including fields in nested structs that are not themselves plugins.
Use an exported pointer for automatic initialization, or an unexported field
with manual initialization when the field should remain outside automatic wiring.

```go
type CustomT struct {
    *testo.T
    Count int // TESTO021: exported fields must be pointers
}
func TestCustom(t *testing.T) {
    testo.RunTest(t, func(ctx *CustomT) {})
}
```

TESTO022 reports a cycle along the active initialization path, such as an
exported `Next *Node` field on `Node` reached from a custom T. Testo cannot
finish constructing this graph. Reusing a type on independent branches is safe.
Keep intentional recursive data behind unexported fields and initialize it
manually, or remove the recursive construction dependency.

These are construction failures under Testo v1.8.0, rather than struct-style
policies. Fill precedes plugin methods and suite hooks, so initializing an
exported invalid/cyclic graph in those methods cannot fix the construction
failure. A cycle takes precedence over TESTO021 for the same context. Detecting
uninitialized plugin fields remains deferred and is not included in this release.

## Provider and plugin initialization

Testo evaluates parameter providers after BeforeAll and before BeforeEach.
Dynamic providers and state prepared in BeforeAll are valid. TESTO023 reports a
narrow proven failure: a direct `RunSuite(t, &Suite{})` or `RunSubSuite(t, &Suite{})`
uses a fresh zero-initialized suite with exactly one provider, BeforeEach
assigns `new(T)` to a pointer
field, but a provider first returns a slice containing `*s.field`.
Move the initialization into BeforeAll or initialize the suite before running it.
The diagnostic points to the provider's dereference, once even across repeated runs.

```go
func (s *Suite) BeforeEach(t *testo.T) { s.value = new(int) }
func (s *Suite) CasesValue() []int {
    return []int{*s.value} // TESTO023 when RunSuite uses &Suite{}
}
```

TESTO024 checks plugin initialization. Testo supplies a typed nil parent at the
suite root: the interface itself is non-nil. A direct `if parent == nil` or
`if parent != nil` cannot distinguish that root from a child. Assert the plugin's
concrete pointer first, then compare that pointer with nil:

```go
func (p *Plugin) Plugin(parent testoplugin.Plugin, options ...testoplugin.Option) testoplugin.Spec {
    previous := parent.(*Plugin)
    if previous == nil {
        // Initialize the suite root.
    }
    return testoplugin.Spec{}
}
```

This checks the framework's parent contract, rather than proving that the branch
will fail a test. Manual Plugin calls can legitimately pass a nil interface;
explicit references to the implementation method in the analyzed package
suppress this check. External manual calls are not inferred.

`testo.Options` intentionally appends process-global configuration for later
suite initializations. It does not update an existing context's option snapshot.
Registration from a running test can deliberately configure a later RunSubSuite
or RunTest; it is not diagnosed. A global-options policy remains deferred and
is not included in the rule count.

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
  empty slice literal, or built-in `make` with constant zero length and an
  omitted or constant nonnegative capacity, including named slices and aliases.
  Shadowed make functions, unknown length/capacity, local slice variables, imported
  bodies and other expressions or control flow are skipped. Invalid signatures/names
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
- TESTO018 inspects direct synchronous expression statements in the straight-line
  prefix of a direct Cleanup closure. The context must be the same local variable
  or parameter with the actual `*testo.T` type (aliases accepted); any reassignment
  or address-taking anywhere disables identity inference. Different contexts,
  context alias chains, fields, custom contexts, Cleanup method values/expressions,
  stored callbacks, nested function bodies, branches, calls following control-flow
  statements, and go/defer Run calls are skipped. It does not analyze helper
  bodies or prove that registration/callback execution is reached.
- TESTO019 accepts explicit nil, function-type conversions of nil (including
  named function types and aliases), and unchanged local zero/nil initializers.
  Reassignment or address-taking anywhere disables local inference. Unknown
  functions, global callback variables and alias chains are skipped. Stored,
  returned, discarded or indirectly passed Test adapters are not tracked.
  The rule supports custom Testo context types for RunTest/Test and does not
  require a suite declaration. It never rejects nil callbacks to Run.
- TESTO020 requires one regular test, no case providers, standard `*testo.T`,
  default BeforeAll/BeforeEach/AfterAll, and local non-promoted test/AfterEach
  methods. The test body must consist of one direct `testo.Run` without options,
  on its context parameter and with a constant name. Its literal callback must
  contain exactly the child's Parallel call followed by Read, Write or Stat on
  the same direct `*os.File` field closed by the sole AfterEach statement.
  Calls may discard results as expressions or blank assignments; named-result
  handling, branches, joins, aliases, deferred calls and helpers are skipped.
  File operation arguments must be nil or constant-only literals. File/context
  type aliases and pointer/value suite receivers are supported. Custom resource
  types (including wrappers around os.File), multiple tests, parameterized tests,
  field replacement and custom surrounding hooks remain unknown. The model
  assumes no external mutation of the suite field or custom runtime scheduling.
  It is a narrow file-lifetime check, not general resource or race analysis.
- TESTO023 requires a local non-generic suite with standard `*testo.T`, a direct
  empty address-of suite literal, no run options and a default BeforeAll reached
  through value embeddings. Exactly one provider is allowed, so another provider
  cannot initialize the field during earlier collection. All test/hook/provider
  signatures and parameter bindings must be valid. BeforeEach must be a local pointer method whose sole
  statement assigns built-in `new(T)` to a direct pointer field; the provider
  must be a local pointer method whose sole return is a one-element slice literal
  containing a direct dereference of that field. Aliases are accepted. Unknown
  instance history, nonempty literals, custom BeforeAll, pointer embeddings on the
  BeforeAll path, promoted bodies, helper calls, other expressions and flow are
  skipped. Default runtime behavior is assumed; filters/plugins may prevent
  execution. A declaration or a valid zero-value read alone is never an error.
- TESTO024 uses the same concrete execution roots as construction checks and
  follows exported pointer fields through valid acyclic construction graphs.
  It checks direct local implementations of the actual testoplugin.Plugin
  interface, including concrete generic instances and pointer/value receivers.
  The context type itself is not a plugin. Invalid/unknown construction graphs,
  imported or promoted implementation bodies, inaccessible fields, container
  contents and unexecuted adapters do not establish checked implementations.
  Only a direct if-condition comparison of the unchanged parent parameter with
  the actual nil identifier is supported, with either operand order and
  parentheses. Reassignment, address-taking or explicit local method references
  suppress inference. Parent aliases, compound conditions, nested closures,
  wrappers and external manual invocations are not analyzed.
- Other lifecycle checks, plugin hooks and broader provider logic are outside
  this release's scope. No SSA or
  interprocedural analysis is used.
- TESTO021/022 start from typed RunSuite/RunSubSuite calls, RunTest, Run and
  immediately executed Test adapters (direct invocation or direct testing.T.Run
  callback). Import aliases, explicit/inferred generics and unchanged direct
  local function aliases are supported. Unused suite/context declarations and
  stored Test adapters do not establish construction roots. A call need not be
  proven reachable: diagnostics describe its construction graph if executed.
- Construction follows pointers (including named pointers) and exported struct
  fields. Fill also traverses exported value structs, so their cycles can take
  precedence over a non-pointer-field error. Wiring stops at the actual testo.T;
  user types with a similar name receive no special treatment. Unexported fields
  and contents of arrays, slices, maps, channels, functions and interfaces are
  not traversed. Pointer fields leading to those containers are valid.
  Repeated acyclic types are cached independently of active-path cycle detection.
- Known concrete generic instantiations are checked; unknown traversed type
  parameters suppress field checks. Interface contexts with unknown concrete
  runtime types, wrappers, indirect adapters and alias chains remain unsupported.
  Only the first encountered cycle per context is reported; shared cycles and
  identical local-field diagnostics are deduplicated. Local diagnostics point
  to the offending field, imported fields to the first local construction call
  for that context. Fixing a reported obstruction may reveal another one.
  Ordinary struct cycles produce Testo's recursion panic; pointer-only cycles
  also prevent initialization from terminating, without promising that specific
  panic. Plugin method bodies, their overrides, manually initialized hidden
  dependencies and general nilness are not analyzed.

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
TESTOLINT_GOLANGCI_BINARY="$PWD/bin/golangci-lint" go test ./plugin/golangci -run '^TestIntegration' -count=1 -v
(cd plugin/golangci/testdata/integration && go test ./clean)
```

The external binary test is otherwise skipped. CI builds the custom binary and
runs this check. The analyzer is also available as `analyzer.Analyzer` for
standard `go/analysis` drivers.

## License

[MIT](LICENSE).
