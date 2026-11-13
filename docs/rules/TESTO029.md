# TESTO029 — unreported-action

Opt-in style check for explicitly configured actions in valid local Allure tests.

Examples assume `Suite` embeds `testo.Suite[T]` and `T` contains the Testo and
Allure plugins. `allure` denotes `github.com/ozontech/testo-allure`; native assertions use `t.Require()` and `t.Assert()`. Domain fields and methods are
illustrative. Each Valid example is an alternative to its Violation example.

## Violation

Assume `CreateOrder` is configured as an action:

```yaml
style:
  actions:
    - package: example.org/orders
      receiver: Client
      method: CreateOrder
```

```go
func (s *Suite) TestCreate(t T) {
    s.createOrder(t) // TESTO029: uncovered action inside helper
}

func (s *Suite) createOrder(t T) {
    s.client.CreateOrder()
}
```

A configured operation called directly outside Step or Setup is diagnosed.
A local helper with an uncovered operation is diagnosed at its callsite, with
related information pointing to that operation.

## Valid

Either wrap the helper call in the test, or cover its actions internally:

```go
func (s *Suite) TestCreate(t T) {
    allure.Step(t, "Create order", func(t T) {
        s.createOrder(t)
    })
}
```

Alternatively, the test can call this helper directly:

```go
func (s *Suite) createOrder(t T) {
    allure.Step(t, "Create order", func(t T) {
        s.client.CreateOrder()
    })
}
```

A helper containing both a Step and another uncovered action still violates
the rule. Setup also covers preparation actions.

Wrap the operation or helper call in Step/Setup, or cover all configured
operations inside the helper. Data-only helpers need no reporting container.

Only one local helper transition is followed. Nested Step callbacks do not add
helper depth. Arbitrary function literals, asynchronous/deferred calls and
TearDown bodies are excluded. Unknown operations are not classified as pure.
No automatic fix is offered.

## Classifying a client

With `actions: [{package: example.org/orders, receiver: Client, method: "*"}]`
and `ignore: [{package: example.org/orders, receiver: Client, method: Close}]`,
all Client methods are actions except Close. Imports and suite assumptions are
the same as above.

Violation:
```go
func (s *Suite) TestCreate(t T) {
    s.client.CreateOrder() // TESTO029
}
```

Valid:
```go
func (s *Suite) TestCreate(t T) {
    allure.Step(t, "Create order", func(t T) { s.client.CreateOrder() })
    s.client.Close() // explicitly excluded
}
```

Use groups for domain clients or dedicated action packages. A classified helper
is treated as an operation even if its implementation already creates a Step.

## Ignored helpers and arguments

Assume local `action` is configured as an action, returns an integer, and
`ignored` is configured in `style.ignore` for this package. The suite/context
and `allure` import assumptions above apply.

Violation:
```go
func (s *Suite) TestArgument(t T) {
    ignored(action()) // TESTO029: the argument executes before ignored
}
```

Valid:
```go
func ignored(int) { action() }

func (s *Suite) TestIgnored(t T) {
    ignored(0) // explicit ignore also excludes the helper body
    allure.Step(t, "Prepare argument", func(t T) { ignored(action()) })
}
```

Ignore is a trusted exclusion, not verification that the helper produces a
report. Unclassified local helpers still have their bodies inspected.
