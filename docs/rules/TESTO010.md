# TESTO010: empty-cases

correctness; enabled by default.

Providers returning `nil`, an empty slice literal or built-in `make` with constant zero length always return an empty case set.

Examples use `T = *testo.T` and a suite `S` embedding `testo.Suite[T]`.

Violation: `func (S) CasesAge() []int { return nil }` for a used provider.
Valid: `func (S) CasesAge() []int { return []int{18} }`.

Supply at least one case when the parameterized test should execute. Testo warns
and skips the affected test for an empty case set; this is not a guaranteed
runtime failure. Only supported single-return bodies prove emptiness: nil,
empty slice literals and built-in make with proven zero length and valid capacity.
Imported or more complex bodies are not inferred.

## Analysis boundaries

- `TESTO010` checks valid local providers whose sole statement returns nil or an
  empty slice literal, or built-in `make` with constant zero length and an
  omitted or constant nonnegative capacity, including named slices and aliases.
  Shadowed make functions, unknown length/capacity, local slice variables, imported
  bodies and other expressions or control flow are skipped. Invalid signatures/names
  suppress this check; malformed parameterized tests make orphan checks conservative.

[Rule index](../rules.md) · [Configuration](../configuration.md) · [Shared scope](../analysis-scope.md)
