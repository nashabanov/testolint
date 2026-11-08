# TESTO014: shared-before-each-write

policy; enabled by default.

BeforeEach writes a shared receiver field read after Parallel suspension.

Examples use a suite embedding `testo.Suite[*testo.T]`; `t` is the
current Testo context and `s` is the suite receiver. Other fields and helpers
are omitted where they do not affect the illustrated rule.

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

## Analysis boundaries

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

See also [shared analysis boundaries](TESTO017.md).

[Rule index](../rules.md) · [Configuration](../configuration.md) · [Shared scope](../analysis-scope.md)
