# TESTO027: no-sleep

policy; enabled by default.

Avoid time.Sleep in Testo tests and hooks; use polling or condition-based waiting.

Violation: `time.Sleep(time.Second)` in a Testo test.

Valid: poll the expected condition with a deadline.

TESTO027 checks local valid Testo tests, lifecycle hooks and literal callbacks to
`testo.Run`, `testo.RunTest` and `testo.Test`, including nested functions. It
recognizes direct standard-library `time.Sleep` calls, including renamed and dot
imports. Helpers, Cases providers, imported bodies and function-variable calls
are outside its scope. No automatic replacement is offered.

Malformed test names (such as `Testfoo`) are handled by TESTO007 and do not
produce TESTO027, even when TESTO007 is disabled. Invalid test/hook signatures
are also excluded from this check.

Enabled by default since 0.9.0; disable with `-disable=TESTO027` or plugin
`settings: {disable: [TESTO027]}`. This is policy, not proof of a runtime failure.

[Rule index](../rules.md) · [Configuration](../../README.md#rule-configuration) · [Shared scope](../analysis-scope.md)
