# TESTO016: shared-test-context

policy; enabled by default.

A shared BeforeEach field can expose another execution's test context.

Violation: BeforeEach stores `s.ctx = t`, then a parallel test uses
`s.ctx.Log("value")` after another execution can overwrite `s.ctx`.
Valid: use the current test's `t.Log("value")` directly.

Pass the current context explicitly, or capture it locally before suspension.
This ownership policy checks direct shared fields with a proven local overwrite;
it does not infer arbitrary context aliases, asynchronous ownership or
plugin-defined scheduling.

See [shared lifecycle semantics and diagnostic precedence](TESTO014.md).

See [shared scope](../analysis-scope.md) and [parallel analysis boundaries](TESTO014.md).

[Rule index](../rules.md) · [Configuration](../../README.md#rule-configuration) · [Shared scope](../analysis-scope.md)
