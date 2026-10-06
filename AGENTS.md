# Repository Guidelines

These instructions apply to the entire repository.

## Structure

- `analyzer/discovery.go` — discovers Testo suites and builds the model used by rules.
- `analyzer/rule_*.go` and `analyzer/rules.go` — rules and their registration.
- `analyzer/testdata/src/` — fixtures for end-to-end analyzer tests.
- `analyzer/analyzer_test.go` — runs fixtures through `analysistest.Run`.
- `cmd/testolint/` — CLI; `plugin/golangci/` — golangci-lint integration.

## Mandatory TDD for Rules and Discovery

Every new rule, rule fix, or discovery change must follow this sequence:

1. **Create or extend an e2e fixture first** in `analyzer/testdata/src/` to exercise the intended behavior through a full analyzer run. For a new fixture package, add it to the `TestAnalyzer` list in `analyzer/analyzer_test.go`.
2. Specify expected diagnostics using `// want "TESTOxxx..."` comments on the relevant lines. Add a valid example without `want` to check for false positives. The fixture must load successfully and pass Go type checking.
3. **Run the test before changing the implementation and observe the expected failure**: a missing, unexpected, or incorrect diagnostic. Compilation, package loading, or environment errors do not count as the TDD red phase. For discovery changes, verify the result through analyzer diagnostics, including cases that should be ignored.
4. **Only after the red phase, implement or modify the rule or discovery code.** Make the minimum changes needed to pass the fixture.
5. Rerun the targeted test, then the full test suite. Refactor only while tests are green.

Do not implement changes before writing the fixture or add tests after the fact. Unit tests may supplement end-to-end fixtures but cannot replace them. Every bug fix must leave a regression fixture.

Cover boundaries relevant to the change: valid and invalid signatures, type aliases, pointer/value receivers, unrelated `Suite` types, and the absence of unexpected or duplicate diagnostics. Choose cases based on the change rather than mechanically adding every variant.

## Verification Commands

Run a targeted fixture (replace `<fixture>` with the package name):

```sh
go test ./analyzer -run '^TestAnalyzer/<fixture>$' -count=1
```

After changing Go code, format the modified files with `gofmt`, then run:

```sh
make test
go vet ./...
```

`make lint` runs golangci-lint v2 and should be run when the binary is available.

For CLI or integration changes, also add or update the corresponding end-to-end tests first in `cmd/testolint/main_test.go` or `plugin/golangci/integration_test.go` and their fixtures in `plugin/golangci/testdata/integration/`.

To verify the plugin using a previously built custom binary that includes the current changes:

```sh
TESTOLINT_GOLANGCI_BINARY="$PWD/bin/golangci-lint" go test ./plugin/golangci -run '^TestIntegration$' -count=1 -v
```

Without this variable, `TestIntegration` is skipped; a skip does not count as successful verification of the external integration. If the environment prevents a check from running, state that explicitly in the report.

## Change Consistency

Keep diagnostic identifiers (`TESTOxxx`) stable. When public behavior changes, update the rule descriptions and limitations in `README.md`. In the work report, identify the added fixture, the red and green phase results, and the checks performed.

## Review and Commit

After completing changes and verification, present the changes and check results to the user for review. Wait for explicit user approval before creating a commit. After approval, commit the reviewed changes using the project's existing commit convention (Conventional Commits, such as `fix: ...`, `feat: ...`, or `docs: ...`). If further changes are made after approval, present them for review again before committing.
