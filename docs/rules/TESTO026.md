# TESTO026: blank-between-steps

style; disabled by default.

Separate adjacent Allure Step statements with a blank line.

Enable with `-enable=TESTO026`; use `-fix` to insert separators.

Examples use a custom Testo context `T` embedding `*allure.PluginAllure`
from `github.com/ozontech/testo-allure`.

Violation:

```go
func (*Suite) TestOrder(t T) {
    allure.Step(t, "create", func(T) {})
    allure.Step(t, "check", func(T) {})
}
```

Valid:

```go
func (*Suite) TestOrder(t T) {
    allure.Step(t, "create", func(T) {})

    allure.Step(t, "check", func(T) {})
}
```

Nested blocks and callbacks are checked. Other statements break adjacency.
Comments do not substitute for a blank line; fixes preserve their placement.

## Analysis boundaries

- TESTO026 recognizes direct `allure.Step` calls in AST block statement lists,
  including nested callbacks, aliases/dot imports and generic arguments. It skips
  wrappers, function variables, Setup/TearDown and switch/select clause lists.
  Another statement breaks adjacency. Comments do not count as blank lines.
  Fixes preserve comments and LF/CRLF; same-line statements have no auto-fix.

See [Allure metadata scope](TESTO025.md).

[Rule index](../rules.md) · [Configuration](../configuration.md) · [Shared scope](../analysis-scope.md)
