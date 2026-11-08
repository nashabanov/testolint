# TESTO027: no-sleep

policy; enabled by default.

Avoid time.Sleep in Testo tests and hooks; use polling or condition-based waiting.

Violation: `time.Sleep(time.Second)` in a Testo test.

Valid: wait for an explicit completion signal with a timeout. In this fragment,
`done` is a channel closed by the operation when it finishes; `t` is the Testo context.

```go
select {
case <-done:
    // The operation completed; check its result.
case <-time.After(time.Second):
    t.Fatal("operation did not complete within one second")
}
```

When no completion signal is available, poll the expected condition with a
bounded timeout instead of assuming a fixed sleep is long enough.

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

[Rule index](../rules.md) · [Configuration](../configuration.md) · [Shared scope](../analysis-scope.md)
