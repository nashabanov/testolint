# TESTO018: run-during-cleanup

correctness; enabled by default.

Direct `testo.Run` on the cleanup context inside its Cleanup callback.

Violation: `t.Cleanup(func() { testo.Run(t, "child", nil) })`.
Valid: call `testo.Run(t, "child", nil)` before registering cleanup.

TESTO018 reports a direct `testo.Run` inside `t.Cleanup(func() { ... })` when
both operations use the same unchanged local `*testo.T` variable. Testo delegates
Run to the underlying testing context, where running a subtest during Cleanup
panics. Run subtests before cleanup; use Cleanup to release their resources.

```go
t.Cleanup(func() {
    testo.Run(t, "child", nil) // TESTO018, even though nil is valid for Run
})
```

Overridden Cleanup implementations and unknown registration graphs are outside
the proof.

## Analysis boundaries

- TESTO018 inspects direct synchronous expression statements in the straight-line
  prefix of a direct Cleanup closure. The context must be the same local variable
  or parameter with the actual `*testo.T` type (aliases accepted); any reassignment
  or address-taking anywhere disables identity inference. Different contexts,
  context alias chains, fields, custom contexts, Cleanup method values/expressions,
  stored callbacks, nested function bodies, branches, calls following control-flow
  statements, and go/defer Run calls are skipped. It does not analyze helper
  bodies or prove that registration/callback execution is reached.

[Rule index](../rules.md) · [Configuration](../configuration.md) · [Shared scope](../analysis-scope.md)
