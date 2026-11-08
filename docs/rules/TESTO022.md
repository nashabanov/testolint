# TESTO022: construction-cycle

correctness; enabled by default.

Report a type cycle traversed during custom T/plugin initialization.

Violation: an executed context embeds `*Node` where
`type Node struct { Next *Node }`.
Valid: remove the initialization cycle, for example
`type Node struct { Value *int }`.

TESTO022 reports a cycle along the active initialization path, such as an
exported `Next *Node` field on `Node` reached from a custom T. Testo cannot
finish constructing this graph. Reusing a type on independent branches is safe.
Keep intentional recursive data behind unexported fields and initialize it
manually, or remove the recursive construction dependency.

These are construction failures under Testo v1.8.0, rather than struct-style
policies. Fill precedes plugin methods and suite hooks, so initializing an
exported invalid/cyclic graph in those methods cannot fix the construction
failure. A cycle takes precedence over TESTO021 for the same context. Detecting
uninitialized plugin fields remains deferred and is not included in this release.

Inaccessible fields, container contents and unknown traversed type parameters
are not treated as definite cycles.

See [shared construction boundaries](../analysis-scope.md) and [construction roots](TESTO021.md).

[Rule index](../rules.md) · [Configuration](../configuration.md) · [Shared scope](../analysis-scope.md)
