# Changelog

## 0.3.0

### Added

- Discover `RunSuite` and `RunSubSuite` using Go types, including generic calls,
  import aliases and unchanged direct local function aliases.
- `TESTO012`: report explicit nil and unchanged local nil suite pointers.
- `TESTO013`: report direct receiver recursion in `BeforeAll` or a sole regular
  test with default preceding hooks; skip execution paths requiring inference.
- Regression fixtures for discovery, both rules, false positives, diagnostic
  positions, deduplication and deterministic execution alongside existing rules.
- CLI and custom golangci-lint integration coverage for both new rules and
  clean execution with real Testo, including bounded receiver recursion.
- Document detection boundaries and Testo v1.8.0 runtime behavior.

### Changed

- Add separate run-level rules with shared discovery and diagnostic deduplication,
  keep analysis state local, and modernize loops and test error handling.
- Use the latest published release in installation examples and record pre-release checks.

## 0.2.0

### Added

- `TESTO010`: report providers whose single return statement returns `nil` or an
  empty slice literal, including named slices and aliases. Local promoted
  providers are checked; imported bodies and more complex provider logic are skipped.
- `TESTO011`: report suites with no tests recognized by Testo's naming rules,
  including suites with only hooks or providers.
- Discover suites that indirectly embed `testo.Suite[T]`.
- Check promoted and imported tests, providers, and hooks using Go's pointer
  method set, respecting shadowing and ambiguous methods. Imported method
  diagnostics point to the local suite declaration.
- Repository guidelines for fixture-first TDD, pre-release review, and changelog
  updates before requesting commit approval.

### Fixed

- Discover package-level suite types even when they declare no methods.
  Type aliases do not create duplicate suites.

### Changed

- Share provider signature validation between `TESTO004`, `TESTO005`, and
  `TESTO010`, preserving diagnostic messages and behavior.
- Update README installation examples to target v0.2.0 once its tag is published.

## 0.1.0

Initial release.

### Added

- Standalone `testolint` CLI and reusable `go/analysis` analyzer.
- golangci-lint v2 module plugin using the same analyzer and rules as the CLI.
- `TESTO001` and `TESTO002`: validate test and lifecycle hook signatures.
- `TESTO003`, `TESTO004`, and `TESTO005`: validate required field-name providers,
  assignable slice element types, and provider signatures.
- `TESTO006`: report providers unused by test parameters as a linter policy.
- `TESTO007` and `TESTO008`: validate Testo test and provider naming, including
  Unicode suffixes; ignore the method named exactly `Cases`.
- `TESTO009`: report unexported test parameter fields.
- Direct `testo.Suite[T]` discovery, including aliases and pointer embedding,
  with checks for declared pointer and value receiver methods.
- Support named parameter structs, named provider slices, and type aliases.
- Suppress dependent diagnostics for invalid signatures and unexported fields;
  deduplicate identical diagnostics at the same source position.
- Analyzer fixtures, CLI exit-code tests, real-Testo integration fixtures,
  custom golangci-lint integration tests, and CI checks.
- Installation, usage, compatibility, and limitations documentation, with runtime
  semantics checked against Testo v1.8.0.
