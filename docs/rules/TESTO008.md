# TESTO008: cases-name

correctness; enabled by default.

After `Cases`, the suffix is empty or begins with a non-lowercase Unicode rune.

Examples use `T = *testo.T` and a suite `S` embedding `testo.Suite[T]`.

Violation: `func (S) Casesage() []int { return []int{18} }`.
Valid: `func (S) CasesAge() []int { return []int{18} }` with a matching parameter.

Apply Testo's non-lowercase suffix rule after `Cases`. Unicode is supported;
the method named exactly `Cases` is ignored by Testo and by this check. A
well-named provider can still have an invalid signature or be unused.

[Rule index](../rules.md) · [Configuration](../configuration.md) · [Shared scope](../analysis-scope.md)
