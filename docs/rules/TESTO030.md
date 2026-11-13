# TESTO030 — test-phase-order

Opt-in style check using the same configured operation identities and traversal
as TESTO029/031. `style.phases.strict` defaults to true; explicit false allows
Assert → Action. Setup after Action or Assert is invalid in either mode.

## Violation

Action → Assert → Setup violates either mode. Assert → Action violates strict.

## Valid

Setup → Action → Action → Assert → Assert is valid. No phase is mandatory.
Action → Assert → Action → Assert is valid with strict false.

Setup is one Arrange event, including any precondition assertions inside it.
Step retains the order of operations inside its direct callback. Local helpers
are followed once; their event order is retained and diagnostics point to the
callsite with related information. No automatic reorder fix is offered.

Branches, loops and short-circuit boolean expressions interrupt phase inference; their paths are not concatenated.
TearDown, deferred/async bodies and uninvoked closures are excluded. Ordinary
assignments and result calculations are neutral. A local literal directly passed to a configured action without reassignment
is recognized as input Data preparation. Other assignments and opaque generators
remain neutral because their purpose is not proven. Unknown operations receive no assumed phase.
