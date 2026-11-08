# TESTO003: missing-case-provider

correctness; enabled by default.

Exported parameter fields require a `Cases<Field>` method.

Examples use `T = *testo.T` and a suite `S` embedding `testo.Suite[T]`.

Violation: `TestValue(t T, p struct{ Age int })` with no `CasesAge` method.
Valid: add `func (S) CasesAge() []int { return []int{18} }`.

Provide cases for each exported parameter field, matching the field name rather
than the test name. Invalid test signatures and unexported fields suppress this
check; a malformed existing provider is handled by its own signature check.

[Rule index](../rules.md) · [Configuration](../configuration.md) · [Shared scope](../analysis-scope.md)
