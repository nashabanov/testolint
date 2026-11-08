# TESTO017: duplicate-parallel-lifecycle

correctness; enabled by default.

Repeated Parallel calls on the same context across BeforeEach, test and AfterEach.

Examples use a suite embedding `testo.Suite[*testo.T]`; `t` is the
current Testo context and `s` is the suite receiver. Other fields and helpers
are omitted where they do not affect the illustrated rule.

Violation: `t.Parallel(); t.Parallel()` on a known lifecycle/test path.
Valid: call `t.Parallel()` once across BeforeEach, the test and AfterEach.

Remove the duplicate call at its owning stage. The rule considers understood
paths and context identity under default Testo behavior; unsupported operations
can block later proofs without erasing duplicates already established. Parallel
calls on distinct contexts are not duplicates.

See [shared lifecycle semantics and diagnostic precedence](TESTO014.md).

See also [shared analysis boundaries](TESTO014.md).

## Analysis boundaries

- TESTO017 requires an established execution of the checked test: a regular
  method or a parameterized method with known non-empty constant cases. It does
  not require sibling executions or known sibling case counts. Collection errors
  (including missing or incompatible providers) still suppress the check.

- Invalid lifecycle/test signatures and collection errors suppress dependent
  checks. BeforeAll/AfterAll signatures are validated; their bodies do not enter
  per-test timelines. Suite-context Parallel does not reorder the inner siblings.
  An unknown test body (calls, branches, loops, early exits, aliases or unavailable
  bodies) blocks inference about its own events and later stages, without
  suppressing known sibling paths.
  TESTO017 uses complete lifecycle stages up to the first unknown body. An
  unsupported AfterEach does not hide a duplicate already established by
  BeforeEach and the test; an unknown BeforeEach or test blocks reasoning
  about later stages. Unknown hooks are never treated as empty hooks.
  TESTO014–016 require understood BeforeEach/AfterEach bodies and a complete
  consumer path. A later sibling can establish the BeforeEach overwrite even
  when its test body is unsupported, provided the local body never references
  its suite receiver or exposes T beyond the resolved standard operations.
  Such a body supplies no reads, writes, or Parallel facts.
  Any receiver reference in unknown code (including aliases, field addresses
  and helper arguments), context escapes (T can expose the suite through
  `testo.Reflect`), or an unavailable/promoted body suppresses shared-field
  proofs. Known test/AfterEach writes still exclude the affected field.
  No partial events from an unsupported body are used. Receiver access through
  globals or other hidden aliases and helper side effects are not inferred;
  the supported ownership model assumes no such external receiver access.
  Helper is accepted only on the actual context parameter with its resolved
  standard method; arbitrary helpers and context aliases remain unknown.
  Other expressions, including arithmetic, indexing and conversions,
  remain unsupported. Identical constant field assignments retain existing policy.
  Synchronization, helper internals, subtests,
  custom T types, plugin overrides and runtime filters are not modeled. A clean
  result is not proof of race freedom. No general concurrency analysis is used.

See [shared scope](../analysis-scope.md) and [parallel analysis boundaries](TESTO014.md).

[Rule index](../rules.md) · [Configuration](../configuration.md) · [Shared scope](../analysis-scope.md)
