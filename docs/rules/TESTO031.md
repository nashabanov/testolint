# TESTO031 — assert-inside-step

Opt-in style check for explicitly configured assertions in valid local Allure tests.

## Violation

A configured operation called directly outside Step or Setup is diagnosed.
A local helper with an uncovered operation is diagnosed at its callsite, with
related information pointing to that operation.

## Valid

Wrap the operation or helper call in Step/Setup, or cover all configured
operations inside the helper. Data-only helpers need no reporting container.

Only one local helper transition is followed. Nested Step callbacks do not add
helper depth. Arbitrary function literals, asynchronous/deferred calls and
TearDown bodies are excluded. Unknown operations are not classified as pure.
No automatic fix is offered. Built-in assertion defaults are still pending.
