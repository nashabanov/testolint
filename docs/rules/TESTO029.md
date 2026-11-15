# TESTO029: unreported-action

style; disabled by default.

Wrap configured actions in an Allure Step or Setup. Enable TESTO029 and declare
actions in `style.actions`; no actions are classified by default.

Examples use `Suite` embedding `testo.Suite[T]`, a custom context `T` embedding
`*testo.T` and `*allure.PluginAllure`, and import `allure "github.com/ozontech/testo-allure"`.
Local `action() Client` is configured as an action; `Client.Read()` is unclassified:

```yaml
style:
  actions:
    - {package: example.org/orders, function: action}
```

Violation (in package `example.org/orders`):

```go
func (*Suite) TestOrder(t T) {
    action().Read() // TESTO029: action executes while evaluating the receiver
}
```

Valid:

```go
func (*Suite) TestOrder(t T) {
    allure.Step(t, "Order", func(t T) { action().Read() })
}
```

## Analysis boundaries

Valid local Allure tests are checked. One local helper transition is followed;
helper violations point to the test's callsite with the operation as a related
location. Cover the call in the test or the operation inside the helper.
Computed callees and receivers are inspected before arguments.

Configured symbols work without source bodies. `style.ignore` excludes the
matched operation and its helper body, but evaluated receivers and arguments
are still checked. Unknown operations are not assumed pure. Arbitrary closures,
async/deferred calls and TearDown bodies are excluded. No automatic fix is offered.

[Rule index](../rules.md) · [Configuration](../configuration.md#operation-symbols) · [Shared scope](../analysis-scope.md#optional-test-structure-checks)
