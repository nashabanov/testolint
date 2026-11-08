# TESTO019: nil-callback

correctness; enabled by default.

Statically nil callbacks to `RunTest` or an immediately executed `Test` adapter.

Examples use `T = *testo.T`; `t` is the parent testing context.

Violation: `testo.RunTest[T](t, nil)`.
Valid: `testo.RunTest(t, func(t T) {})`.

TESTO019 reports a statically nil function supplied to `testo.RunTest`. It also
checks `testo.Test` when the returned adapter is immediately called or passed
directly to the resolved `testing.T.Run` method. Test creates an adapter without
invoking the callback; the failure occurs when that adapter runs. Supply a
non-nil function, such as `func(t *testo.T) {}`, for an intentionally empty test.
**`testo.Run(t, name, nil)` is valid** and uses a no-op callback.

```go
testo.RunTest[*testo.T](t, nil)       // TESTO019
t.Run("child", testo.Test[*testo.T](nil)) // TESTO019 when the adapter runs
```

These are runtime-failure checks under the default Testo behavior, verified
against v1.8.0. Plugins can override Cleanup or filter singleton execution;
such runtime changes are not inferred.

## Analysis boundaries

- TESTO019 accepts explicit nil, function-type conversions of nil (including
  named function types and aliases), and unchanged local zero/nil initializers.
  Reassignment or address-taking anywhere disables local inference. Unknown
  functions, global callback variables and alias chains are skipped. Stored,
  returned, discarded or indirectly passed Test adapters are not tracked.
  The rule supports custom Testo context types for RunTest/Test and does not
  require a suite declaration. It never rejects nil callbacks to Run.

[Rule index](../rules.md) · [Configuration](../configuration.md) · [Shared scope](../analysis-scope.md)
