# TESTO002: hook-signature

correctness; enabled by default.

`BeforeAll`, `BeforeEach`, `AfterEach`, `AfterAll` accept exactly `T` and return no values.

Violation: `func (S) BeforeEach() {}`.
Valid: `func (S) BeforeEach(t T) {}`.

Use exactly one T parameter and no results for each lifecycle hook. Promoted
methods are checked through the suite method set. Signature validity does not
establish that an arbitrary hook body can be analyzed by lifecycle rules.

[Rule index](../rules.md) · [Configuration](../../README.md#rule-configuration) · [Shared scope](../analysis-scope.md)
