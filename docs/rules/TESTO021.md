# TESTO021: construction-field

correctness; enabled by default.

Exported fields traversed during custom T/plugin wiring must be pointers.

Violation in an executed custom context: `type C struct { *testo.T; Count int }`.
Valid: `type C struct { *testo.T; Count *int }`.

Testo constructs a fresh custom T and fills its accessible pointer fields before
wiring plugins. TESTO021 reports an exported non-pointer field reached during
this wiring, including fields in nested structs that are not themselves plugins.
Use an exported pointer for automatic initialization, or an unexported field
with manual initialization when the field should remain outside automatic wiring.

```go
type CustomT struct {
    *testo.T
    Count int // TESTO021: exported fields must be pointers
}
func TestCustom(t *testing.T) {
    testo.RunTest(t, func(ctx *CustomT) {})
}
```

Cycles and unknown type parameters suppress dependent field diagnostics even
when TESTO022 is disabled.

See [construction cycle precedence and runtime behavior](TESTO022.md).

## Analysis boundaries

- TESTO021/022 start from typed RunSuite/RunSubSuite calls, RunTest, Run and
  immediately executed Test adapters (direct invocation or direct testing.T.Run
  callback). Import aliases, explicit/inferred generics and unchanged direct
  local function aliases are supported. Unused suite/context declarations and
  stored Test adapters do not establish construction roots. A call need not be
  proven reachable: diagnostics describe its construction graph if executed.

See [shared construction boundaries](../analysis-scope.md).

[Rule index](../rules.md) · [Configuration](../configuration.md) · [Shared scope](../analysis-scope.md)
