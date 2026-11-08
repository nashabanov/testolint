# TESTO011: empty-suite

correctness; enabled by default.

Suite must contain at least one runnable Testo test.

Examples use `T = *testo.T` and a suite `S` embedding `testo.Suite[T]`.

Violation: `type S struct{ testo.Suite[T] }` with no recognized tests.
Valid: add `func (S) TestValue(t T) {}`.

Add a test or remove the unused suite. Hooks and providers alone do not run tests;
Testo warns for an empty suite. Promoted tests count, and naming recognition is
separate from signature validity: TESTO001 handles a recognized but invalid test.

[Rule index](../rules.md) · [Configuration](../configuration.md) · [Shared scope](../analysis-scope.md)
