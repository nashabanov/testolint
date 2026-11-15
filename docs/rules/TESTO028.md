# TESTO028: required-test-header

style; disabled by default.

Require configured Allure metadata in tests, optionally supplied by BeforeEach.
Enable TESTO028 and set `style.header.required`, for example
`[[Title, Titlef], [Tags]]`: every group is required, with alternatives inside it.
No metadata is required by default; `include-before-each` defaults to true.

Examples use `Suite` embedding `testo.Suite[T]` and a custom context `T` embedding
`*testo.T` and `*allure.PluginAllure` from `github.com/ozontech/testo-allure`.

Violation:

```go
func (*Suite) TestOrder(t T) { // TESTO028: missing Tags
    t.Title("Order")
}
```

Valid (alternative suite definition):

```go
func (*Suite) BeforeEach(t T) { t.Tags("smoke") }
func (*Suite) TestOrder(t T) { t.Titlef("%s", "Order") }
```

## Analysis boundaries

Valid local Allure tests and their effective BeforeEach are checked, including
local generic promoted methods. Shadowing follows Go method sets. Set
`include-before-each: false` to require metadata in each test itself.

Only direct calls on the original context prove presence. Conditional metadata
makes the corresponding methods uncertain; independently missing groups are
still reported. Unavailable hooks, context aliases, stored method values and
opaque calls receiving the context suppress unproven absence. Uninvoked
closures contribute nothing. No automatic fix is offered.

**Known limitation:** passing the context to `allure.Step` also suppresses this
check, even without the required metadata. Child-step metadata does not satisfy
the parent header, but callbacks can capture the outer context; the checker does
not yet distinguish these cases. A clean result is not proof of complete metadata.

[Rule index](../rules.md) · [Configuration](../configuration.md#unreleased-header-configuration) · [Shared scope](../analysis-scope.md#optional-test-structure-checks)
