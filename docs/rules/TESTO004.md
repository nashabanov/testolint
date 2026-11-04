# TESTO004: case-type-mismatch

correctness; enabled by default.

Provider slice elements must be assignable to the parameter field type.

Violation: an `Age int` parameter field with `CasesAge() []string`.
Valid: `func (S) CasesAge() []int { return []int{18} }`.

Return elements assignable to the field type under Go's assignability rules.
For example, `[]int` can supply an `any` field. Invalid provider signatures and
unexported parameter fields suppress this diagnostic.

```go
package integration

import "github.com/ozontech/testo"

type Suite struct{ testo.Suite[*testo.T] }

func (Suite) TestAge(t *testo.T, p struct{ Age int }) {}
func (Suite) CasesAge() []string { return []string{"18"} }
```

```text
suite_test.go:8:14: TESTO004: CasesAge provides string, but parameter "Age" expects int (https://github.com/nashabanov/testolint/blob/main/docs/rules/TESTO004.md)
```

Return `[]int{18}` from `CasesAge` to satisfy the `Age int` parameter.

[Rule index](../rules.md) · [Configuration](../../README.md#rule-configuration) · [Shared scope](../analysis-scope.md)
