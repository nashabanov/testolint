# TESTO024: plugin-parent-nil

correctness; enabled by default.

A Plugin parent interface compared with nil cannot identify the suite root.

Violation inside a reached Plugin implementation: `if parent == nil { ... }`.
Valid for the plugin's own concrete type `MyPlugin`:

```go
previous, _ := parent.(*MyPlugin)
if previous == nil {
    // Root initialization, or a different parent plugin type.
}
```

Testo supplies a typed nil at the root, so the parent interface itself is non-nil.
Assert the expected concrete type before comparing its pointer. The rule checks
supported local implementations reached from valid construction roots; it skips
reassigned/address-taken parameters, explicit local method references and unknown
graphs. Manual Plugin calls may legitimately pass a nil interface. A diagnostic
does not prove that the branch will fail a test.

TESTO024 checks plugin initialization. Testo supplies a typed nil parent at the
suite root: the interface itself is non-nil. A direct `if parent == nil` or
`if parent != nil` cannot distinguish that root from a child. Assert the plugin's
concrete pointer first, then compare that pointer with nil:

```go
func (p *Plugin) Plugin(parent testoplugin.Plugin, options ...testoplugin.Option) testoplugin.Spec {
    previous := parent.(*Plugin)
    if previous == nil {
        // Initialize the suite root.
    }
    return testoplugin.Spec{}
}
```

This checks the framework's parent contract, rather than proving that the branch
will fail a test. Manual Plugin calls can legitimately pass a nil interface;
explicit references to the implementation method in the analyzed package
suppress this check. External manual calls are not inferred.

### Analysis boundaries

- TESTO024 uses the same concrete execution roots as construction checks and
  follows exported pointer fields through valid acyclic construction graphs.
  It checks direct local implementations of the actual testoplugin.Plugin
  interface, including concrete generic instances and pointer/value receivers.
  The context type itself is not a plugin. Invalid/unknown construction graphs,
  imported or promoted implementation bodies, inaccessible fields, container
  contents and unexecuted adapters do not establish checked implementations.
  Only a direct if-condition comparison of the unchanged parent parameter with
  the actual nil identifier is supported, with either operand order and
  parentheses. Reassignment, address-taking or explicit local method references
  suppress inference. Parent aliases, compound conditions, nested closures,
  wrappers and external manual invocations are not analyzed.

See [shared construction boundaries](../analysis-scope.md) and [construction roots](TESTO021.md).

[Rule index](../rules.md) · [Configuration](../../README.md#rule-configuration) · [Shared scope](../analysis-scope.md)
