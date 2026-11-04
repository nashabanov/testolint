# TESTO001: test-signature

correctness; enabled by default.

Tests accept exactly the suite's `T`, optionally a struct, and return no values.

Violation: `func (S) TestValue() {}`.
Valid: `func (S) TestValue(t T) {}`.

Match the suite's exact T, optionally add a struct parameter, and return no
values. Named structs and aliases are supported; unrelated T types are rejected.
Invalid signatures suppress dependent parameter diagnostics.

[Rule index](../rules.md) · [Configuration](../../README.md#rule-configuration) · [Shared scope](../analysis-scope.md)
