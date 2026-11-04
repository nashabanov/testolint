# TESTO007: test-name

correctness; enabled by default.

After `Test`, the suffix is empty or begins with a non-lowercase Unicode rune.

Violation: `func (S) Testlower(t T) {}`.
Valid: `func (S) TestLower(t T) {}`.

Rename the method so the suffix after `Test` is empty or starts with a
non-lowercase Unicode rune. This mirrors Testo's name validation, not an
ASCII-only or uppercase-only convention. Unrelated methods are ignored.

[Rule index](../rules.md) · [Configuration](../../README.md#rule-configuration) · [Shared scope](../analysis-scope.md)
