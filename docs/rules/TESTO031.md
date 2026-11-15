# TESTO031: assert-inside-step

style; disabled by default.

Wrap assertions in an outer Allure Step or Setup. Enable TESTO031; native
Testo Allure assertions and the testify v1.11.1 assert/require APIs are recognized
by default. Add custom assertion symbols through `style.assertions`.

Examples use `Suite` embedding `testo.Suite[T]`, a custom context `T` embedding
`*testo.T` and `*allure.PluginAllure`, and import `allure "github.com/ozontech/testo-allure"`.
`Client.Read()` is unclassified; the local helper contains an assertion.

Violation:

```go
func helper(t T) Client {
    t.Require().Equal(1, 1)
    return Client{}
}
func (*Suite) TestOrder(t T) {
    helper(t).Read() // TESTO031: uncovered assertion inside helper
}
```

Valid (using the same helper):

```go
func (*Suite) TestOrder(t T) {
    allure.Step(t, "Check order", func(t T) { helper(t).Read() })
}
```

## Analysis boundaries

Valid local Allure tests follow the same coverage boundaries as TESTO029: one
local helper transition, direct Step/Setup callbacks, and computed callees and
receivers before arguments. Cover a helper call in the test or its assertions
internally. `style.ignore` excludes matched operations and helper bodies, while
evaluated receivers and arguments remain checked.

Native assertion report entries do not replace the required outer Step/Setup.
`Require()`/`Assert()` obtain helpers; `Hidden()`/`Masked()` only configure them.
Stored assertion helpers retain classification. Testify constructors and
comparison utilities are excluded. Arbitrary closures, async/deferred calls and
TearDown bodies are skipped. No automatic fix is offered.

[Rule index](../rules.md) · [Configuration](../configuration.md#operation-symbols) · [Shared scope](../analysis-scope.md#optional-test-structure-checks)
