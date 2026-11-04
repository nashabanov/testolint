# TESTO005: cases-signature

correctness; enabled by default.

Providers accept no arguments and return exactly one slice.

Violation: `func (S) CasesAge(n int) []int { return []int{n} }`.
Valid: `func (S) CasesAge() []int { return []int{18} }`.

Providers take no arguments and return exactly one slice; named slices are
supported. Methods not recognized as providers are handled by naming rules.
Fix the signature before relying on provider type or emptiness diagnostics.

[Rule index](../rules.md) · [Configuration](../../README.md#rule-configuration) · [Shared scope](../analysis-scope.md)
