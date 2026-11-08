# TESTO012: nil-suite

policy; enabled by default.

Report a statically nil suite argument to `RunSuite` or `RunSubSuite`.

Examples use a suite embedding `testo.Suite[*testo.T]`; `t` is the
current Testo context and `s` is the suite receiver. Other fields and helpers
are omitted where they do not affect the illustrated rule.

Violation: `testo.RunSuite(t, (*S)(nil))`.
Valid: `testo.RunSuite(t, &S{})`.

For suite execution, `TESTO012` reports at the suite argument:

```go
var suite *Suite = nil
testo.RunSuite(t, suite) // TESTO012: suite argument is statically nil
```

Disable the rule when nil-safe suite execution is intentional.
Pass an initialized instance such as `&Suite{}`. This rule reports a known nil
value, not a guaranteed panic. In Testo v1.8.0, inherited value-receiver hooks
panic when invoked through a nil suite pointer, but explicitly nil-safe pointer
hooks and tests can run successfully. Reporting those intentional nil suites is
also part of this rule's policy; method bodies are not analyzed for nil safety.

## Analysis boundaries

- `TESTO012` checks explicit nil, pointer conversions of nil and local pointer
  declarations with a nil initializer or zero value (`var`/`:=`). Any assignment
  or address-taking disables local inference, even after the call or in a closure.
  Other values, interface variables and control-flow-dependent nilness are skipped.

[Rule index](../rules.md) · [Configuration](../configuration.md) · [Shared scope](../analysis-scope.md)
