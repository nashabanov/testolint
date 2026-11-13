# TESTO030 — test-phase-order

Opt-in style check using the same configured operation identities and traversal
as TESTO029/031. `style.phases.strict` defaults to false; explicit false allows
Assert → Action. Setup after Action or Assert is invalid in either mode.

Examples assume `Suite` embeds `testo.Suite[T]` and `T` contains the Testo and
Allure plugins. `allure` denotes `github.com/ozontech/testo-allure`; native assertions use `t.Require()` and `t.Assert()`. Domain fields and methods are
illustrative. Each Valid example is an alternative to its Violation example.

## Violation

Assume `CreateOrder` is configured as an action as in TESTO029.
With explicit `strict: true`:

```go
func (s *Suite) TestCreate(t T) {
    s.client.CreateOrder()
    t.Require().Equal(1, s.orderCount())
    s.client.CreateOrder() // TESTO030: Action after Assert
}
```

A late Setup is invalid in both modes:

```go
func (s *Suite) TestCreate(t T) {
    s.client.CreateOrder()
    allure.Setup(t, "Prepare environment", func(t T) {}) // TESTO030
}
```

Action → Assert → Setup violates either mode. Assert → Action violates strict.

## Valid

Multiple actions and assertions in a phase are allowed:

```go
func (s *Suite) TestCreate(t T) {
    allure.Setup(t, "Prepare environment", func(t T) {
        t.Require().NotNil(s.client) // precondition remains Arrange
    })
    s.client.CreateOrder()
    s.client.CreateOrder()
    t.Require().Equal(2, s.orderCount())
    t.Require().NotNil(s.client)
}
```

The first violation example is valid with `style.phases.strict: false` (or
`-phase-order-strict=false`). These examples isolate TESTO030; when coverage
rules are also enabled, put actions and assertions in Step/Setup containers.

Setup → Action → Action → Assert → Assert is valid. No phase is mandatory.
Action → Assert → Action → Assert is valid with strict false.

Setup is one Arrange event, including any precondition assertions inside it.
Step retains the order of operations inside its direct callback. Local helpers
are followed once; their event order is retained and diagnostics point to the
callsite with related information. No automatic reorder fix is offered.

Branches, loops and short-circuit boolean expressions interrupt phase inference; their paths are not concatenated.
TearDown, deferred/async bodies and uninvoked closures are excluded. Ordinary
assignments and result calculations are neutral. Input preparation, including literals and factory calls, remains neutral;
only explicit Setup calls establish Arrange events. Unknown operations receive no assumed phase.

## Ignored helpers

With `strict: true`, local integer-returning `action` configured as an action
and `ignored` configured in `style.ignore`, the following is a violation.
Suite/context assumptions and the native assertion API are as above; these
examples isolate TESTO030.

```go
func (s *Suite) TestArgument(t T) {
    t.Require().Equal(1, 1)
    ignored(action()) // TESTO030: argument Action follows Assert
}
```

Valid:
```go
func ignored(int) { action() }

func (s *Suite) TestIgnored(t T) {
    t.Require().Equal(1, 1)
    ignored(0) // the excluded helper body contributes no phase
}
```

With `strict: false`, the first example is also valid for TESTO030.
Explicit ignore excludes local bodies just as it excludes imported helpers.
