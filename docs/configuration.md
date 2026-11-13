# Rule configuration

Rule selection is available since 0.8.0 in both the CLI and the
[golangci-lint plugin](golangci-lint.md). See the [rule catalog](rules.md) for
individual defaults.

Both integrations use the same selection model:

1. Start with each rule's registry default. Correctness and policy rules are enabled; style rules are opt-in.
2. With `disable-all`, start with no rules enabled instead.
3. Apply explicit `enable` and `disable` lists. An ID in both lists is an error,
   regardless of flag or YAML key order. Duplicate and unknown IDs are errors.

IDs are case-sensitive (`TESTO019`). No class changes a rule's default. Empty
lists preserve the starting selection; `disable-all` with no enabled IDs produces
no diagnostics. A disabled prerequisite may still be analyzed to preserve
suppression and diagnostic ownership. For example, disabling TESTO022 does not
make TESTO021 report fields in a construction graph that contains a cycle.
Selecting TESTO018 does not enable TESTO019, even though they share a checker.
Disabling an owning diagnostic does not transfer its reports to another rule.

## CLI

CLI flags take comma-separated lists; whitespace around each ID is trimmed.
Each list flag may appear once. A trailing comma or an empty list element is an
error; `-enable=` and `-disable=` mean empty lists. Invalid settings fail during
flag parsing, before package analysis.

```sh
# Preserve defaults except for two policy checks.
testolint -disable=TESTO006,TESTO012 ./...
# Run only the unused-provider policy.
testolint -disable-all -enable=TESTO006 ./...
# Run only nil-callback and construction-cycle checks.
testolint -disable-all -enable=TESTO019,TESTO022 ./...
```

## golangci-lint

The equivalent plugin configuration goes under the custom linter's `settings`:

```yaml
version: "2"
linters:
  default: standard
  enable: [testolint]
  settings:
    custom:
      testolint:
        type: module
        settings:
          disable-all: true
          enable: [TESTO019, TESTO022]
          disable: []
```

Plugin lists must be YAML sequences of exact IDs, and `disable-all` must be a
boolean. Unknown setting names, invalid types, unknown/duplicate IDs and conflicts
are rejected when the plugin is constructed. Each driver uses a separate analyzer
instance. For Go API usage, see [Contributing](../CONTRIBUTING.md#go-analysis-drivers).

See the [diagnostic and configuration contract](contract.md) for stability,
exit codes and suppression.

## Migrating to 0.9.0

TESTO027 is now enabled by default following the maintainer's project audit.
To keep the previous selection, use `-disable=TESTO027` in the CLI, or:

```yaml
settings:
  disable: [TESTO027]
```

under `linters.settings.custom.testolint` in the plugin configuration. TESTO025
and TESTO026 remain opt-in. Diagnostic messages now end with a rule-specific
link: `TESTOxxx: message (https://github.com/nashabanov/testolint/blob/main/docs/rules/TESTOxxx.md)`.
Update exact-message snapshots; the ID prefix remains unchanged and `Category`
remains the ID. The structured `URL` now points to the same rule page.
TESTO027 skips malformed test names: TESTO007 owns those diagnostics, even when
TESTO007 is disabled. Valid tests and lifecycle hooks remain checked for sleep.

## Unreleased header configuration

The shared Go/plugin `style.header` settings are:

```yaml
style:
  header:
    methods: [Title, Titlef, Tags, Owner]
    required: [[Title, Titlef], [Tags]]
    include-before-each: true
```

`methods` defaults to the existing TESTO025 metadata set. `required` defaults
empty; groups combine with AND, alternatives within a group with OR.
Unknown methods, duplicates and empty groups are errors. Required methods must
be recognized. Include-before-each defaults true. These settings do not enable
rules: select TESTO025 and/or TESTO028 independently. CLI accepts the same style object as JSON through `-style-config`.

### Operation symbols (implementation in progress)

`style.actions` and `style.assertions` are lists of exact symbols:

```yaml
style:
  actions:
    - package: example.org/client
      receiver: Client
      method: Send
  assertions:
    - package: example.org/check
      function: Equal
```

Choose package/function or package/receiver/method. Receiver is the declared
type name, without a pointer prefix. Import aliases do not change identity.
Duplicate/conflicting classifications and malformed identifiers are rejected.
These matchers currently have no built-in defaults; testify defaults are pending. Rules remain independently opt-in.

`style.phases.strict` defaults to true; false permits Assert → Action for
TESTO030. Returning to Setup remains invalid.

```sh
testolint -disable-all -enable=TESTO028 \
  -style-config='{"header":{"required":[["Title","Titlef"],["Tags"]]}}' ./...
testolint -enable=TESTO030 -phase-order-strict=false ./...
```

The explicit `-phase-order-strict` flag overrides JSON regardless of argument
order. Both style flags may appear only once. Unknown JSON fields, wrong types,
invalid symbol identities and invalid required groups are rejected. The phase
flag affects only TESTO030; it does not enable it or change other rules.
