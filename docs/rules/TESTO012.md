# TESTO012: nil-suite

policy; enabled by default.

Report a statically nil suite argument to `RunSuite` or `RunSubSuite`.

Violation: `testo.RunSuite(t, (*S)(nil))`.
Valid: `testo.RunSuite(t, &S{})`.

Pass an initialized suite. This policy includes intentionally nil-safe pointer
suites, so a diagnostic is not a claim of an inevitable panic. Explicit and
unchanged local nil pointers are supported; interface boxing and
control-flow-dependent nilness are not inferred. Disable the rule when nil-safe
suite execution is intentional.

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

### Analysis boundaries

- `TESTO012` checks explicit nil, pointer conversions of nil and local pointer
  declarations with a nil initializer or zero value (`var`/`:=`). Any assignment
  or address-taking disables local inference, even after the call or in a closure.
  Other values, interface variables and control-flow-dependent nilness are skipped.

[Rule index](../rules.md) · [Configuration](../../README.md#rule-configuration) · [Shared scope](../analysis-scope.md)
