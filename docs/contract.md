# Diagnostic and configuration contract

## Identity and selection

Diagnostic IDs are stable and are never recycled for another meaning. Rule names
are the human-readable registry names; configuration selects exact, case-sensitive
IDs, not names or classes. Renaming or reclassifying a rule must be documented.
`correctness` describes a supported Testo contract or execution hazard; `policy`
can reject intentional code; `style` concerns optional formatting and test structure. Classification does not
implicitly enable or disable a rule.

In 0.9.0, TESTO001–TESTO024 and TESTO027 are enabled by default. TESTO025 and
TESTO026 are disabled. New optional policies start disabled. See [rule configuration](configuration.md) 
for selection semantics, validation and migration instructions.

A multi-ID checker executes once if any of its IDs is selected. Only selected
IDs are reported. Disabling an owning rule does not transfer diagnostics to a
related rule or remove semantic suppression of cascading errors. Only identical
complete diagnostics are deduplicated, retaining the first occurrence. There is
no analyzer-level source suppression syntax; driver exclusions are separate.

## Reports and fixes

Every diagnostic has `Category` equal to its ID, `URL` equal to its registry rule
page, and this message format:

```text
TESTO004: message (https://github.com/nashabanov/testolint/blob/main/docs/rules/TESTO004.md)
```

The ID prefix and documentation suffix each occur once. The suffix makes links
available in plain CLI and golangci-lint output. URLs use stable ID-based filenames;
the linked main-branch documentation can evolve. Existing README section anchors
and `docs/rules.md#testoxxx` remain navigation entry points.

Message wording, source ranges and related locations may be refined in patch
releases. Consumers should use Category/ID rather than match complete prose.
Detection fixes can change diagnostics within documented scope and must include
regression coverage and changelog notes. Ranges, related information and
SuggestedFixes are preserved by reporting and deduplication, but their presence
is rule-specific; not every diagnostic has a fix or related location.

TESTO025/026 offer blank-line insertions where safe, preserving comments and
LF/CRLF. Same-line statements may be diagnosed without a fix. Applying the fixes
again makes no further changes. No automatic time.Sleep replacement is offered.
Drivers decide how fixes are exposed or applied (`-fix` in the CLI, `--fix` in
golangci-lint).

## Driver results

The standalone CLI uses the Go singlechecker driver: a successful normal run
without findings exits 0, diagnostics exit 3, and analysis/package-loading errors
exit 1. Invalid flags/configuration exit 2; help exits 0. Fix and JSON modes follow
the underlying driver behavior and should not be interpreted as normal text mode.
The plugin uses golangci-lint's exit policy (normally 1 for findings, configurable
by that driver); it does not promise the standalone CLI's exit codes.

In unreleased 0.10.0, TESTO028–031 are also disabled by default and offer no
automatic fixes. TESTO030 consumes the same operation classification as coverage
checks without enabling their diagnostics.
