# TESTO031 — assert-inside-step

Opt-in style check for built-in and explicitly configured assertions in valid local Allure tests.

Examples assume `Suite` embeds `testo.Suite[T]` and `T` contains the Testo and
Allure plugins. `allure` denotes `github.com/ozontech/testo-allure`; native assertions use `t.Require()` and `t.Assert()`. Domain fields and methods are
illustrative. Each Valid example is an alternative to its Violation example.

## Violation

Native Testo Allure assertions are recognized without custom configuration:

```go
func (s *Suite) TestResponse(t T) {
    t.Require().Equal(200, s.response.StatusCode) // TESTO031
    t.Assert().NotNil(s.response.Body)           // TESTO031
}
```

A configured operation called directly outside Step or Setup is diagnosed.
A local helper with an uncovered operation is diagnosed at its callsite, with
related information pointing to that operation.

## Valid

An assertion Step and Setup preconditions are both supported:

```go
func (s *Suite) TestResponse(t T) {
    allure.Setup(t, "Check preconditions", func(t T) {
        t.Require().NotNil(s.response)
    })
    allure.Step(t, "Check response", func(t T) {
        t.Require().Equal(200, s.response.StatusCode)
        t.Assert().NotNil(s.response.Body)
    })
}
```

The rule does not require an assertion in every Step. `t.Require()` and `t.Assert()` only obtain assertion helpers;
`Hidden()` and `Masked()` only configure them. The assertion methods are checks.

Wrap the operation or helper call in Step/Setup, or cover all configured
operations inside the helper. Data-only helpers need no reporting container.

Only one local helper transition is followed. Nested Step callbacks do not add
helper depth. Arbitrary function literals, asynchronous/deferred calls and
TearDown bodies are excluded. Unknown operations are not classified as pure.
No automatic fix is offered. Testify v1.11.1 assert/require assertion functions and Assertions methods
are recognized by default; constructors and comparison utilities are excluded.

## Native reporting versus scenario structure

`t.Require().Equal(...)` and `t.Assert().Equal(...)` already produce assertion
entries in Allure. TESTO031 intentionally requires an **outer** Step or Setup
for scenario structure; those automatic assertion entries do not satisfy this
rule. Both APIs also produce Assert events for TESTO030, except inside Setup,
where they remain preconditions in Arrange. Stored helpers such as
`checks := t.Assert(); checks.Equal(...)` retain the same classification.

Direct testify assert/require functions and methods remain supported as
additional APIs. User-defined assertions can be configured explicitly.

Custom assertion symbols also accept the grouped patterns documented in
[configuration](../configuration.md#grouped-operations). For example, with
`assertions: [{package: example.org/check, function: "*"}]` and import
`check "example.org/check"`, the same suite assumptions as above apply.

Violation:
```go
func (s *Suite) TestResult(t T) {
    check.Equal(1, 1) // TESTO031
}
```

Valid:
```go
func (s *Suite) TestResult(t T) {
    allure.Step(t, "Check result", func(t T) { check.Equal(1, 1) })
}
```

## Excluding assertion helpers

Assume the local function `ignored` is configured in `style.ignore`, while
`check` is unclassified. The suite/context and native assertion assumptions
above apply.

```go
func check(t T) { t.Require().Equal(1, 1) }
func ignored(t T) { t.Require().Equal(1, 1) }

// Violation:
func (s *Suite) TestCheck(t T) {
    check(t) // TESTO031: uncovered assertion in a local helper
}

// Valid:
func (s *Suite) TestIgnored(t T) {
    ignored(t) // explicit exclusion prevents traversal of its body
}
```

The exclusion does not verify reporting. Assertions in evaluated arguments
remain subject to coverage checks.
