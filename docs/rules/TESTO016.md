# TESTO016: shared-test-context

policy; enabled by default.

A shared BeforeEach field can expose another execution's test context.

Examples use a suite embedding `testo.Suite[*testo.T]`; `t` is the
current Testo context and `s` is the suite receiver. Other fields and helpers
are omitted where they do not affect the illustrated rule.

Violation: BeforeEach stores `s.ctx = t`, then a parallel test uses
`s.ctx.Log("value")` after another execution can overwrite `s.ctx`.
Valid: use the current test's `t.Log("value")` directly.

Pass the current context explicitly, or capture it locally before suspension.
This ownership policy checks direct shared fields with a proven local overwrite;
it does not infer arbitrary context aliases, asynchronous ownership or
plugin-defined scheduling.

See [shared lifecycle semantics and diagnostic precedence](TESTO014.md).

[Rule index](../rules.md) · [Configuration](../configuration.md) · [Shared scope](../analysis-scope.md)
