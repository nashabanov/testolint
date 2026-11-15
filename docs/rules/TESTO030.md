# TESTO030: test-phase-order

style; disabled by default.

Keep Setup before actions and assertions. Enable TESTO030; set
`style.phases.strict: true` to also reject Action after Assert. The default
`false` allows repeated Action/Assert sequences. Phases are optional.

Examples use `Suite` embedding `testo.Suite[T]`, a custom context `T` embedding
`*testo.T` and `*allure.PluginAllure`, and import `allure "github.com/ozontech/testo-allure"`.
Local `action() Client` is configured in `style.actions` by package/function;
`Client.Read()` is unclassified. These examples isolate TESTO030.

Violation with `strict: true`:

```go
func (*Suite) TestOrder(t T) {
    t.Require().Equal(1, 1)
    action().Read() // TESTO030: Action after Assert
}
```

Valid in both modes:

```go
func (*Suite) TestOrder(t T) {
    allure.Setup(t, "Prepare", func(t T) { t.Require().Equal(1, 1) })
    action().Read()
    t.Require().Equal(1, 1)
}
```

The first example is also valid with `strict: false`. Setup after either Action
or Assert is invalid in both modes. If TESTO029/031 are enabled, also wrap the
operations in Step/Setup.

## Analysis boundaries

Uses the operation identities and traversal of TESTO029/031, including one local
helper transition and computed receivers before arguments. Step preserves phase
order; Setup is a single preparation event, including its assertions.
Only explicit Setup establishes preparation: literals, assignments and
unclassified factories are neutral. Ignored operations and helper bodies add no
phase, but their evaluated receivers and arguments remain checked.

Branches, loops and short-circuit expressions interrupt inference. Arbitrary
closures, async/deferred calls and TearDown bodies are excluded. Unknown calls
receive no assumed phase. No automatic reorder fix is offered.

[Rule index](../rules.md) · [Configuration](../configuration.md#operation-symbols) · [Shared scope](../analysis-scope.md#optional-test-structure-checks)
