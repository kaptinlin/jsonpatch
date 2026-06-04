# Predicate Specs

## Overview

This spec defines all predicate operations and second-order predicate composition rules.

Covered operations: `test`, `contains`, `defined`, `undefined`, `starts`, `ends`, `in`, `less`, `more`, `matches`, `type`, `test_type`, `test_string`, `test_string_len`, `and`, `or`, `not`.

## Negation and Composition Model

- Direct `not` is supported only by `test`, `test_string`, and `test_string_len`.
- All other negated predicates must be wrapped in the second-order `not` operation.
- `and` and `or` accept a non-empty predicate list.
- `not` accepts exactly one predicate. Multiple negated predicates are expressed with explicit structure such as `not(or(...))`.
- Child paths inside `apply` are relative to the containing predicate path. Use `path: ""` on the container when children should be root-scoped absolute paths.
- Compact and binary codecs preserve this model directly: child paths are encoded relative to the parent and decoded into executable absolute paths.

> **Why**: This keeps simple unary negation where it exists, makes structural negation unambiguous, and lets composite predicates name a common parent path once.
>
> **Rejected**: Allowing direct `not` on every predicate would blur the operation vocabulary. Letting `not` carry an implicit list would hide whether the caller meant `not(and(...))` or `not(or(...))`.

## Predicate Contracts

| Operation | Payload | Contract |
|-----------|---------|----------|
| `test` | `value`, optional `not` | Compare the target value with deep equality. When `not` is true, success is inverted. |
| `defined` | none | Succeeds when the target path exists. |
| `undefined` | none | Succeeds when the target path does not exist. |
| `type` | `value` string | Check one JSON type name: `string`, `number`, `boolean`, `object`, `array`, `null`, or `integer`. |
| `test_type` | `type` string or list | Check one or more JSON type names. |
| `contains` | `value` string, optional `ignore_case` | Check whether a string or `[]byte` target contains the given substring. |
| `starts` | `value` string, optional `ignore_case` | Check string or `[]byte` prefix match. |
| `ends` | `value` string, optional `ignore_case` | Check string or `[]byte` suffix match. |
| `in` | `value` array | Check whether the target value is equal to one of the array entries. |
| `less` | `value` number | Compare after coercing the target through JavaScript-like `Number()` semantics. |
| `more` | `value` number | Compare after coercing the target through JavaScript-like `Number()` semantics. |
| `matches` | `value` regex pattern, optional `ignore_case` | Match a string or `[]byte` target against a regex. The default matcher uses Go's `regexp` package; `WithCompileMatcher` overrides it during compilation. |
| `test_string` | `str`, `pos`, optional `not`, optional `ignore_case` | Compare the substring starting at `pos` against `str`. |
| `test_string_len` | `len`, optional `not` | Check whether the target string length is at least `len`; `not` inverts that result. |

## Second-Order Predicates

| Operation | Payload | Contract |
|-----------|---------|----------|
| `and` | `apply` predicate list | Succeeds only when every child predicate succeeds. |
| `or` | `apply` predicate list | Succeeds when any child predicate succeeds. |
| `not` | `apply` with one predicate | Negates its single child predicate. |

## Type Vocabulary

`type` and `test_type` use the JSON type names defined in `SPECS/30-data-model-specs.md`.

## Boundary Rules

- Use `contains` only for string containment.
- Use `in` for membership in a provided array of acceptable values.
- Use `WithCompileMatcher` only when the default Go regex implementation is not sufficient.

## Forbidden

- Do not add a direct `not` field to predicates other than `test`, `test_string`, and `test_string_len`.
- Do not put mutation operations inside `and`, `or`, or `not`; every child must be a predicate.
- Do not put multiple direct children inside `not`; wrap them in `and` or `or` first.
- Do not use `contains` for array membership; use `in`.
- Do not document `matches` as requiring a custom matcher; the library ships with a default matcher.
- Do not describe `test_string_len` as exact-length matching; its contract is minimum length (`>= len`).

## Acceptance Criteria

- [ ] Each predicate operation has one documented contract.
- [ ] Negation rules are explicit and non-duplicated.
- [ ] Composite predicate path rules are documented.
- [ ] String, numeric, and type predicates describe the actual payload fields they accept.
