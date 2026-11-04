# TESTO008: cases-name

correctness; enabled by default.

After `Cases`, the suffix is empty or begins with a non-lowercase Unicode rune.

Violation: `func (S) Casesage() []int { return []int{18} }`.
Valid: `func (S) CasesAge() []int { return []int{18} }` with a matching parameter.

Apply Testo's non-lowercase suffix rule after `Cases`. Unicode is supported;
the method named exactly `Cases` is ignored by Testo and by this check. A
well-named provider can still have an invalid signature or be unused.

[Rule index](../rules.md) · [Configuration](../../README.md#rule-configuration) · [Shared scope](../analysis-scope.md)
