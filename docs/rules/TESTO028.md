# TESTO028 — required-test-header

Opt-in style check. Configure `style.header.required` as groups of recognized
header methods: AND between groups, OR within each group. No groups are required
by default. `include-before-each` defaults to true.

Examples assume `Suite` embeds `testo.Suite[T]` and `T` contains the Testo and
Allure plugins. `allure` denotes `github.com/ozontech/testo-allure`; native assertions use `t.Require()` and `t.Assert()`. Domain fields and methods are
illustrative. Each Valid example is an alternative to its Violation example.

## Violation

With `required: [[Title, Titlef], [Tags]]`:

```go
func (s *Suite) TestCreate(t T) { // TESTO028: missing Tags
    t.Title("Create order")
}
```

With required `[[Title, Titlef], [Tags]]`, an empty Allure test is diagnosed.

## Valid

The effective BeforeEach can supply Tags; Titlef satisfies the title alternative.

```go
func (s *Suite) BeforeEach(t T) {
    t.Tags("smoke")
}

func (s *Suite) TestCreate(t T) {
    t.Titlef("Create %s", "order")
}
```

A test calling `t.Title("title")` and `t.Tags("smoke")` meets those groups.
Tags may instead come from the effective local BeforeEach hook.

Only valid local Testo tests with the typed Allure marker are checked.
Promoted and shadowed hooks follow discovery. Unavailable hooks and opaque calls receiving the context suppress unproven
absence diagnostics. Conditional known metadata makes only those methods
uncertain, preserving independently proven missing groups.
Uninvoked function literals contribute no metadata. No automatic fix is offered.
