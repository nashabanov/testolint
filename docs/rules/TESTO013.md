# TESTO013: recursive-subsuite

correctness; enabled by default.

Report a directly recursive `RunSubSuite` invocation on the current receiver.

Examples use a suite embedding `testo.Suite[*testo.T]`; `t` is the
current Testo context and `s` is the suite receiver. Other fields and helpers
are omitted where they do not affect the illustrated rule.

Violation in `BeforeAll`: `testo.RunSubSuite(t, s)` on the current receiver.
Valid: `testo.RunSubSuite(t, &ChildSuite{})` for a separate child suite.

`TESTO013` detects direct receiver recursion in a narrowly defined execution
path, for example:

```go
type RecursiveSuite struct{ testo.Suite[*testo.T] }

func (s *RecursiveSuite) TestRecursive(t *testo.T) {
    testo.RunSubSuite(t, s) // TESTO013: recursive RunSubSuite invocation
}
```

The nested run selects the same method on the receiver again. Execute a child
suite with its own tests instead. The diagnostic points to the `RunSubSuite`
call. Different instances of the same type are not reported. Runtime plugins
and filtering may interrupt execution; the message identifies direct recursion
and does not claim that every run must loop forever.

## Analysis boundaries

- `TESTO013` requires a sole `RunSubSuite` call without options, passing the
  method's own Testo parameter and receiver (or a value receiver's address).
  The argument's method set must select that same method. Suites must have valid
  regular tests and no providers: `BeforeAll` is checked, or the sole test when
  `BeforeAll`/`BeforeEach` are inherited defaults. Other execution paths, imported
  bodies, plugins and runtime filters are outside this check's scope.

[Rule index](../rules.md) · [Configuration](../configuration.md) · [Shared scope](../analysis-scope.md)
