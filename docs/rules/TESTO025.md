# TESTO025: blank-after-header

style; disabled by default.

Separate initial Allure metadata from the test body with a blank line.

For the style examples, `T` is a custom Testo context embedding
`*allure.PluginAllure` from `github.com/ozontech/testo-allure`.

`blank-after-header` (style, disabled by default) requires a blank line after
initial Allure metadata. Enable with `-enable=TESTO025`; use `-fix` to insert it.

Violation:

```go
func (*Suite) TestOrder(t T) {
    t.Title("Order")
    t.Tags("smoke")
    allure.Step(t, "create", func(T) {})
}
```

Valid:

```go
func (*Suite) TestOrder(t T) {
    t.Title("Order")
    t.Tags("smoke")

    allure.Step(t, "create", func(T) {})
}
```

A metadata-only body needs no separator. Only metadata on the original test
context is recognized; see the analysis boundaries below for supported methods and limitations.

### Analysis boundaries

- TESTO025/026 are opt-in Allure formatting checks on valid local Testo test
  bodies, including local promoted methods. Hooks, helpers and imported bodies
  are skipped. Metadata is an initial sequence on the original context:
  `Title`, `Titlef`, `Description`, `Descriptionf`, `Tags`, `ID`, `Labels`, `Links`,
  `Owner`, `Severity`, `Epic`, `Feature`, `Story`, `Parameters`, `Flaky`, `Muted`,
  `Known`, resolving to `github.com/ozontech/testo-allure`.

Direct `t.Parallel()` on the original context also counts as header metadata.

[Rule index](../rules.md) · [Configuration](../../README.md#rule-configuration) · [Shared scope](../analysis-scope.md)
