# TESTO020: premature-teardown

correctness; enabled by default.

AfterEach closes an `*os.File` before a parallel child uses the same file.

Violation: AfterEach does `_ = s.file.Close()` before a parallel child calls
`s.file.Stat()` on the same file.
Valid replacement in AfterEach:

```go
file := s.file
t.Cleanup(func() { _ = file.Close() })
```

Defer resource release until children finish, or keep the child serial. The
proof is restricted to os.File Read/Write/Stat, one regular test and supported
direct lifecycle bodies; arbitrary Close-like methods are not modeled. Ignoring
an I/O error does not necessarily fail the test, but the file is already closed.

TESTO020 detects premature teardown of an `*os.File`: AfterEach calls Close
before a parallel child resumes and calls Read, Write or Stat on the same direct
suite field. These operations return an error after Close; ignored errors do not
necessarily fail the test. The rule combines Testo lifecycle order with the
concrete os.File contract, rather than interpreting arbitrary Close-like names.
Move the Close into `t.Cleanup(func() { _ = s.file.Close() })` in AfterEach,
which runs after the children complete, or keep the child serial.

```go
func (s *Suite) TestFile(t *testo.T) {
    testo.Run(t, "child", func(child *testo.T) {
        child.Parallel()
        _, _ = s.file.Read(nil)
    })
}
func (s *Suite) AfterEach(t *testo.T) {
    _ = s.file.Close() // TESTO020: the parallel child has not resumed yet
}
```

### Analysis boundaries

- TESTO020 requires one regular test, no case providers, standard `*testo.T`,
  default BeforeAll/BeforeEach/AfterAll, and local non-promoted test/AfterEach
  methods. The test body must consist of one direct `testo.Run` without options,
  on its context parameter and with a constant name. Its literal callback must
  contain exactly the child's Parallel call followed by Read, Write or Stat on
  the same direct `*os.File` field closed by the sole AfterEach statement.
  Calls may discard results as expressions or blank assignments; named-result
  handling, branches, joins, aliases, deferred calls and helpers are skipped.
  File operation arguments must be nil or constant-only literals. File/context
  type aliases and pointer/value suite receivers are supported. Custom resource
  types (including wrappers around os.File), multiple tests, parameterized tests,
  field replacement and custom surrounding hooks remain unknown. The model
  assumes no external mutation of the suite field or custom runtime scheduling.
  It is a narrow file-lifetime check, not general resource or race analysis.

[Rule index](../rules.md) · [Configuration](../../README.md#rule-configuration) · [Shared scope](../analysis-scope.md)
