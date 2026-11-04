# Rule examples and remedies

[The README registry table](../README.md#rules) lists every rule's name, class
and default. Correctness/policy rules are enabled by default; style rules are
opt-in. Selection settings are documented
in [Rule configuration](../README.md#rule-configuration).

Examples are replacement fragments, not complete packages. Assume `T` aliases
`*testo.T`, `S` embeds `testo.Suite[T]`, and fields, imports and surrounding tests
exist where used. “Valid” means valid for the illustrated rule. Runtime plugins,
filtering and code outside the supported static patterns can change execution;
these checks do not execute tests. The README's
[limitations](analysis-scope.md) describe the complete analysis boundaries.

## TESTO001

[test-signature — description, examples, remedies and limitations](rules/TESTO001.md)

## TESTO002

[hook-signature — description, examples, remedies and limitations](rules/TESTO002.md)

## TESTO003

[missing-case-provider — description, examples, remedies and limitations](rules/TESTO003.md)

## TESTO004

[case-type-mismatch — description, examples, remedies and limitations](rules/TESTO004.md)

## TESTO005

[cases-signature — description, examples, remedies and limitations](rules/TESTO005.md)

## TESTO006

[orphan-cases-provider — description, examples, remedies and limitations](rules/TESTO006.md)

## TESTO007

[test-name — description, examples, remedies and limitations](rules/TESTO007.md)

## TESTO008

[cases-name — description, examples, remedies and limitations](rules/TESTO008.md)

## TESTO009

[unexported-param-field — description, examples, remedies and limitations](rules/TESTO009.md)

## TESTO010

[empty-cases — description, examples, remedies and limitations](rules/TESTO010.md)

## TESTO011

[empty-suite — description, examples, remedies and limitations](rules/TESTO011.md)

## TESTO012

[nil-suite — description, examples, remedies and limitations](rules/TESTO012.md)

## TESTO013

[recursive-subsuite — description, examples, remedies and limitations](rules/TESTO013.md)

## TESTO014

[shared-before-each-write — description, examples, remedies and limitations](rules/TESTO014.md)

## TESTO015

[shared-fixture-cleanup — description, examples, remedies and limitations](rules/TESTO015.md)

## TESTO016

[shared-test-context — description, examples, remedies and limitations](rules/TESTO016.md)

## TESTO017

[duplicate-parallel-lifecycle — description, examples, remedies and limitations](rules/TESTO017.md)

## TESTO018

[run-during-cleanup — description, examples, remedies and limitations](rules/TESTO018.md)

## TESTO019

[nil-callback — description, examples, remedies and limitations](rules/TESTO019.md)

## TESTO020

[premature-teardown — description, examples, remedies and limitations](rules/TESTO020.md)

## TESTO021

[construction-field — description, examples, remedies and limitations](rules/TESTO021.md)

## TESTO022

[construction-cycle — description, examples, remedies and limitations](rules/TESTO022.md)

## TESTO023

[provider-before-each — description, examples, remedies and limitations](rules/TESTO023.md)

## TESTO024

[plugin-parent-nil — description, examples, remedies and limitations](rules/TESTO024.md)

## TESTO025

[blank-after-header — description, examples, remedies and limitations](rules/TESTO025.md)

## TESTO026

[blank-between-steps — description, examples, remedies and limitations](rules/TESTO026.md)

## TESTO027

[no-sleep — description, examples, remedies and limitations](rules/TESTO027.md)
