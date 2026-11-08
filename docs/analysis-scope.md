# Shared analysis scope

These checks inspect supported source patterns without running tests. No SSA or
interprocedural analysis is used. Rule pages describe narrower boundaries;
a clean result is not proof of race freedom.

## Suite discovery

- Suite discovery uses Testo v1.8.0's private marker and recognizes package-level
  named structs embedding `testo.Suite[T]` directly or indirectly, including
  aliases and pointer embedding. Aliases do not duplicate suites. Generic suite
  declarations are skipped by suite rules; instantiated calls support run checks.

- Suite rules use the pointer method set for declared, promoted and imported
  methods, respecting Go shadowing and ambiguity. Local diagnostics point to
  methods or fields; imported-method diagnostics point to the local suite type.

## Run discovery

- Run discovery uses Go types and supports import aliases, inferred/explicit
  generics and unchanged direct local function aliases. Globals, alias chains
  and wrappers are not resolved.

## Context and plugin construction

- Construction follows pointers (including named pointers) and exported struct
  fields. Fill also traverses exported value structs, so their cycles can take
  precedence over a non-pointer-field error. Wiring stops at the actual testo.T;
  user types with a similar name receive no special treatment. Unexported fields
  and contents of arrays, slices, maps, channels, functions and interfaces are
  not traversed. Pointer fields leading to those containers are valid.
  Repeated acyclic types are cached independently of active-path cycle detection.

- Known concrete generic instantiations are checked; unknown traversed type
  parameters suppress field checks. Interface contexts with unknown concrete
  runtime types, wrappers, indirect adapters and alias chains remain unsupported.
  Only the first encountered cycle per context is reported; shared cycles and
  identical local-field diagnostics are deduplicated. Local diagnostics point
  to the offending field, imported fields to the first local construction call
  for that context. Fixing a reported obstruction may reveal another one.
  Ordinary struct cycles produce Testo's recursion panic; pointer-only cycles
  also prevent initialization from terminating, without promising that specific
  panic. Plugin method bodies, their overrides, manually initialized hidden
  dependencies and general nilness are not analyzed.

[Rule-specific boundaries](rules.md)

## Collection and suppression

For `testo.Suite[T]`, tests have the form `TestX(t T)` or
`TestX(t T, p struct{ Age int })`. Named structs and named slice returns are
accepted, as are aliases. Providers match **field names**, not test names, and
can be shared. Assignability follows Go semantics: `[]int` can supply an `any`
field, but `[]any` cannot supply an `int` field.

`Test`, `Test1`, `Test_Foo`, and `TestÉ` are valid names; `Testfoo` is not.
Testo ignores a method named exactly `Cases`, including its signature.
Methods without the exact `Test` or `Cases` prefix are ignored, except hooks.
`TESTO011` counts tests with names recognized by Testo, including promoted and
imported tests. Malformed names do not count; signature errors are reported
separately by `TESTO001`. Empty-suite diagnostics point to the suite type name.

Unexported fields receive only `TESTO009`, without missing/type provider
messages. Invalid test signatures suppress dependent parameter checks; invalid
provider signatures suppress type mismatch checks. Identical diagnostics at the
same position are emitted once. Independent problems may still produce multiple
messages. Rule selection is available in both integrations; source-level
suppression is not implemented.

## Global options

`testo.Options` intentionally appends process-global configuration for later
suite initializations. It does not update an existing context's option snapshot.
Registration from a running test can deliberately configure a later RunSubSuite
or RunTest; it is not diagnosed.
