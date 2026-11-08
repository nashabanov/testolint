# TESTO015: shared-fixture-cleanup

policy; enabled by default.

Cleanup re-reads a shared BeforeEach field that another execution can replace.

Examples use a suite embedding `testo.Suite[*testo.T]`; `t` is the
current Testo context and `s` is the suite receiver. Other fields and helpers
are omitted where they do not affect the illustrated rule.

Violation for a replaceable BeforeEach field:
`t.Cleanup(func() { t.Log(s.value) })` after parallel suspension.

Valid: capture `value := s.value` before suspension, then register
`t.Cleanup(func() { t.Log(value) })`.

Capture the resource itself rather than re-reading shared receiver storage.
This conservative ownership policy uses understood direct lifecycle operations;
it does not analyze arbitrary closures or prove deep ownership. Stored T fields
belong to TESTO016, and this check suppresses related TESTO014 reports.

See [shared lifecycle semantics and diagnostic precedence](TESTO014.md).

[Rule index](../rules.md) · [Configuration](../configuration.md) · [Shared scope](../analysis-scope.md)
