# TESTO009: unexported-param-field

correctness; enabled by default.

Parameter fields must be exported so reflection can set them.

Violation: `TestValue(t T, p struct{ age int })`.
Valid: `TestValue(t T, p struct{ Age int })` with a matching `CasesAge` provider.

Export parameter fields so reflection can set them. Named structs are supported.
The same unexported field does not also receive missing-provider or provider-type
diagnostics. Invalid test signatures are handled first.

[Rule index](../rules.md) · [Configuration](../../README.md#rule-configuration) · [Shared scope](../analysis-scope.md)
