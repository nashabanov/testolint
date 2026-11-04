# TESTO023: provider-before-each

correctness; enabled by default.

A provider dereferences a fresh suite’s nil field before BeforeEach initializes it.

Violation: `BeforeEach` initializes `s.value`, but collection first runs
`func (s *S) CasesAge() []int { return []int{*s.value} }` on a fresh suite.
Valid: move the necessary allocation into a supported `BeforeAll` path, or
initialize the field in the suite passed to RunSuite.

Initialize provider dependencies before collection, not only before each test.
The rule requires a proven fresh nil field, a supported direct BeforeEach
allocation, one provider and valid collection signatures. Constructors, earlier
provider setup, BeforeAll setup, unknown histories and dynamic providers are
preserved rather than guessed.

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

### Analysis boundaries

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

[Rule index](../rules.md) · [Configuration](../../README.md#rule-configuration) · [Shared scope](../analysis-scope.md)
