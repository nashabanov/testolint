# Rule examples and remedies

[The README registry table](../README.md#rules) lists every rule's name, class
and default. All rules are enabled by default. Selection settings are documented
in [Rule configuration](../README.md#rule-configuration-080-development).

Examples are replacement fragments, not complete packages. Assume `T` aliases
`*testo.T`, `S` embeds `testo.Suite[T]`, and fields, imports and surrounding tests
exist where used. “Valid” means valid for the illustrated rule. Runtime plugins,
filtering and code outside the supported static patterns can change execution;
these checks do not execute tests. The README's
[limitations](../README.md#scope-and-limitations) describe the complete analysis boundaries.

## TESTO001

Violation: `func (S) TestValue() {}`.
Valid: `func (S) TestValue(t T) {}`.

Match the suite's exact T, optionally add a struct parameter, and return no
values. Named structs and aliases are supported; unrelated T types are rejected.
Invalid signatures suppress dependent parameter diagnostics.

## TESTO002

Violation: `func (S) BeforeEach() {}`.
Valid: `func (S) BeforeEach(t T) {}`.

Use exactly one T parameter and no results for each lifecycle hook. Promoted
methods are checked through the suite method set. Signature validity does not
establish that an arbitrary hook body can be analyzed by lifecycle rules.

## TESTO003

Violation: `TestValue(t T, p struct{ Age int })` with no `CasesAge` method.
Valid: add `func (S) CasesAge() []int { return []int{18} }`.

Provide cases for each exported parameter field, matching the field name rather
than the test name. Invalid test signatures and unexported fields suppress this
check; a malformed existing provider is handled by its own signature check.

## TESTO004

Violation: an `Age int` parameter field with `CasesAge() []string`.
Valid: `func (S) CasesAge() []int { return []int{18} }`.

Return elements assignable to the field type under Go's assignability rules.
For example, `[]int` can supply an `any` field. Invalid provider signatures and
unexported parameter fields suppress this diagnostic.

## TESTO005

Violation: `func (S) CasesAge(n int) []int { return []int{n} }`.
Valid: `func (S) CasesAge() []int { return []int{18} }`.

Providers take no arguments and return exactly one slice; named slices are
supported. Methods not recognized as providers are handled by naming rules.
Fix the signature before relying on provider type or emptiness diagnostics.

## TESTO006

Violation: `CasesUnused() []int` when no test parameter has an `Unused` field.
Valid: remove the provider, or use `TestValue(t T, p struct{ Unused int })`.

This is policy: Testo permits unused providers. Keep intentional providers by
disabling TESTO006. Matching follows recognized test parameters, not arbitrary
references to the provider method; invalid test signatures are not valid uses.

## TESTO007

Violation: `func (S) Testlower(t T) {}`.
Valid: `func (S) TestLower(t T) {}`.

Rename the method so the suffix after `Test` is empty or starts with a
non-lowercase Unicode rune. This mirrors Testo's name validation, not an
ASCII-only or uppercase-only convention. Unrelated methods are ignored.

## TESTO008

Violation: `func (S) Casesage() []int { return []int{18} }`.
Valid: `func (S) CasesAge() []int { return []int{18} }` with a matching parameter.

Apply Testo's non-lowercase suffix rule after `Cases`. Unicode is supported;
the method named exactly `Cases` is ignored by Testo and by this check. A
well-named provider can still have an invalid signature or be unused.

## TESTO009

Violation: `TestValue(t T, p struct{ age int })`.
Valid: `TestValue(t T, p struct{ Age int })` with a matching `CasesAge` provider.

Export parameter fields so reflection can set them. Named structs are supported.
The same unexported field does not also receive missing-provider or provider-type
diagnostics. Invalid test signatures are handled first.

## TESTO010

Violation: `func (S) CasesAge() []int { return nil }` for a used provider.
Valid: `func (S) CasesAge() []int { return []int{18} }`.

Supply at least one case when the parameterized test should execute. Testo warns
and skips the affected test for an empty case set; this is not a guaranteed
runtime failure. Only supported single-return bodies prove emptiness: nil,
empty slice literals and built-in make with proven zero length and valid capacity.
Imported or more complex bodies are not inferred.

## TESTO011

Violation: `type S struct{ testo.Suite[T] }` with no recognized tests.
Valid: add `func (S) TestValue(t T) {}`.

Add a test or remove the unused suite. Hooks and providers alone do not run tests;
Testo warns for an empty suite. Promoted tests count, and naming recognition is
separate from signature validity: TESTO001 handles a recognized but invalid test.

## TESTO012

Violation: `testo.RunSuite(t, (*S)(nil))`.
Valid: `testo.RunSuite(t, &S{})`.

Pass an initialized suite. This policy includes intentionally nil-safe pointer
suites, so a diagnostic is not a claim of an inevitable panic. Explicit and
unchanged local nil pointers are supported; interface boxing and
control-flow-dependent nilness are not inferred. Disable the rule when nil-safe
suite execution is intentional.

## TESTO013

Violation in `BeforeAll`: `testo.RunSubSuite(t, s)` on the current receiver.
Valid: `testo.RunSubSuite(t, &ChildSuite{})` for a separate child suite.

Avoid directly collecting the same receiver recursively. Analysis covers
BeforeAll and a sole regular test under documented default-hook conditions;
arbitrary recursion graphs and different instances are not reported. Runtime
filtering can stop execution, so the check does not guarantee infinite recursion.

## TESTO014

With `BeforeEach` assigning `s.value` and another possible execution overwriting
it, a Violation is:

```go
t.Parallel()
t.Log(s.value)
```

Valid replacement:

```go
value := s.value
t.Parallel()
t.Log(value)
```

Capture the per-execution value before suspension. This ownership policy requires
a proven local overwrite path, not merely a Parallel call. Unknown accesses can
suppress the proof. Cleanup and stored-context hazards retain TESTO015/016
ownership even when those diagnostics are disabled.

## TESTO015

Violation for a replaceable BeforeEach field:
`t.Cleanup(func() { t.Log(s.value) })` after parallel suspension.

Valid: capture `value := s.value` before suspension, then register
`t.Cleanup(func() { t.Log(value) })`.

Capture the resource itself rather than re-reading shared receiver storage.
This conservative ownership policy uses understood direct lifecycle operations;
it does not analyze arbitrary closures or prove deep ownership. Stored T fields
belong to TESTO016, and this check suppresses related TESTO014 reports.

## TESTO016

Violation: BeforeEach stores `s.ctx = t`, then a parallel test uses
`s.ctx.Log("value")` after another execution can overwrite `s.ctx`.
Valid: use the current test's `t.Log("value")` directly.

Pass the current context explicitly, or capture it locally before suspension.
This ownership policy checks direct shared fields with a proven local overwrite;
it does not infer arbitrary context aliases, asynchronous ownership or
plugin-defined scheduling.

## TESTO017

Violation: `t.Parallel(); t.Parallel()` on a known lifecycle/test path.
Valid: call `t.Parallel()` once across BeforeEach, the test and AfterEach.

Remove the duplicate call at its owning stage. The rule considers understood
paths and context identity under default Testo behavior; unsupported operations
can block later proofs without erasing duplicates already established. Parallel
calls on distinct contexts are not duplicates.

## TESTO018

Violation: `t.Cleanup(func() { testo.Run(t, "child", nil) })`.
Valid: call `testo.Run(t, "child", nil)` before registering cleanup.

Start subtests before cleanup; use Cleanup for resource release. Nil is a valid
callback for Run and is not the error here. Only supported direct calls on the
unchanged cleanup context are checked; overridden Cleanup implementations and
unknown registration graphs are outside the proof.

## TESTO019

Violation: `testo.RunTest[T](t, nil)`.
Valid: `testo.RunTest(t, func(t T) {})`.

Supply a non-nil function, including for an intentionally empty test. The rule
also checks immediately executed Test adapters and direct testing.T.Run use.
Stored adapters, unknown callbacks and reassigned or escaped locals are not
inferred. `testo.Run(t, "child", nil)` remains valid.

## TESTO020

Violation: AfterEach does `_ = s.file.Close()` before a parallel child calls
`s.file.Stat()` on the same file.
Valid replacement in AfterEach:

```go
file := s.file
t.Cleanup(func() { _ = file.Close() })
```

Defer resource release until children finish, or keep the child serial. The
proof is restricted to os.File Read/Write/Stat, one regular test and supported
direct lifecycle bodies; arbitrary Close-like methods are not modeled. Ignoring
an I/O error does not necessarily fail the test, but the file is already closed.

## TESTO021

Violation in an executed custom context: `type C struct { *testo.T; Count int }`.
Valid: `type C struct { *testo.T; Count *int }`.

Use pointer fields for exported values reached by Testo initialization. Traversal
starts at recognized execution roots and includes supported imported and generic
types; unused contexts and container contents are not arbitrary roots. Cycles
and unknown type parameters suppress dependent field diagnostics even when
TESTO022 is disabled.

## TESTO022

Violation: an executed context embeds `*Node` where
`type Node struct { Next *Node }`.
Valid: remove the initialization cycle, for example
`type Node struct { Value *int }`.

Use an acyclic structure for Testo's automatic initialization. The check follows
active-path types; the same leaf type in independent branches is valid.
Inaccessible fields, container contents and unknown traversed type parameters
are not treated as definite cycles. Cycles take precedence over TESTO021.

## TESTO023

Violation: `BeforeEach` initializes `s.value`, but collection first runs
`func (s *S) CasesAge() []int { return []int{*s.value} }` on a fresh suite.
Valid: move the necessary allocation into a supported `BeforeAll` path, or
initialize the field in the suite passed to RunSuite.

Initialize provider dependencies before collection, not only before each test.
The rule requires a proven fresh nil field, a supported direct BeforeEach
allocation, one provider and valid collection signatures. Constructors, earlier
provider setup, BeforeAll setup, unknown histories and dynamic providers are
preserved rather than guessed.

## TESTO024

Violation inside a reached Plugin implementation: `if parent == nil { ... }`.
Valid for the plugin's own concrete type `MyPlugin`:

```go
previous, _ := parent.(*MyPlugin)
if previous == nil {
    // Root initialization, or a different parent plugin type.
}
```

Testo supplies a typed nil at the root, so the parent interface itself is non-nil.
Assert the expected concrete type before comparing its pointer. The rule checks
supported local implementations reached from valid construction roots; it skips
reassigned/address-taken parameters, explicit local method references and unknown
graphs. Manual Plugin calls may legitimately pass a nil interface. A diagnostic
does not prove that the branch will fail a test.
