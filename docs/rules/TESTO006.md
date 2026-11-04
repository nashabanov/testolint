# TESTO006: orphan-cases-provider

policy; enabled by default.

Policy: providers should be referenced by a test parameter; Testo permits unused providers.

Violation: `CasesUnused() []int` when no test parameter has an `Unused` field.
Valid: remove the provider, or use `TestValue(t T, p struct{ Unused int })`.

This is policy: Testo permits unused providers. Keep intentional providers by
disabling TESTO006. Matching follows recognized test parameters, not arbitrary
references to the provider method; invalid test signatures are not valid uses.

[Rule index](../rules.md) · [Configuration](../../README.md#rule-configuration) · [Shared scope](../analysis-scope.md)
