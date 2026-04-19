# Predicate Specs

## Overview

This spec defines all predicate operations and second-order predicate composition rules.

Covered operations: `test`, `contains`, `defined`, `undefined`, `starts`, `ends`, `in`, `less`, `more`, `matches`, `type`, `test_type`, `test_string`, `test_string_len`, `and`, `or`, `not`.

## Negation and Composition Model

- Direct `not` is supported only by `test`, `test_string`, and `test_string_len`.
- All other negated predicates must be wrapped in the second-order `not` operation.
- `and`, `or`, and `not` use `path: ""` and absolute child paths inside `apply`.

> **Why**: This matches the json-joy predicate model and keeps simple unary negation where it exists while reserving structural composition for second-order predicates.
>
> **Rejected**: Allowing direct `not` on every predicate would diverge from the compatibility target. Relative child paths inside `apply` would make nested predicate trees harder to read and reason about.

## Predicate Contracts

| Operation | Payload | Contract |
|-----------|---------|----------|
| `test` | `value`, optional `not` | Compare the target value with deep equality. When `not` is true, success is inverted. |
| `defined` | none | Succeeds when the target path exists. |
| `undefined` | none | Succeeds when the target path does not exist. |
| `type` | `value` string | Check one JSON type name: `string`, `number`, `boolean`, `object`, `array`, `null`, or `integer`. |
| `test_type` | `type` string or list | Check one or more JSON type names. The decoder also accepts `value` for compatibility, but `type` is canonical. |
| `contains` | `value` string, optional `ignore_case` | Check whether a string or `[]byte` target contains the given substring. |
| `starts` | `value` string, optional `ignore_case` | Check string or `[]byte` prefix match. |
| `ends` | `value` string, optional `ignore_case` | Check string or `[]byte` suffix match. |
| `in` | `value` array | Check whether the target value is equal to one of the array entries. |
| `less` | `value` number | Compare after coercing the target through JavaScript-like `Number()` semantics. |
| `more` | `value` number | Compare after coercing the target through JavaScript-like `Number()` semantics. |
| `matches` | `value` regex pattern, optional `ignore_case` | Match a string or `[]byte` target against a regex. The default matcher uses Go's `regexp` package; `WithMatcher` overrides it. |
| `test_string` | `str`, `pos`, optional `not`, optional `ignore_case` | Compare the substring starting at `pos` against `str`. |
| `test_string_len` | `len`, optional `not` | Check whether the target string length is at least `len`; `not` inverts that result. |

## Second-Order Predicates

| Operation | Payload | Contract |
|-----------|---------|----------|
| `and` | `apply` predicate list | Succeeds only when every child predicate succeeds. |
| `or` | `apply` predicate list | Succeeds when any child predicate succeeds. |
| `not` | `apply` predicate list | Negates the child predicate sequence. |

## Type Vocabulary

`type` and `test_type` use the JSON type names defined in `SPECS/30-data-model-specs.md`.

## Boundary Rules

- Use `contains` only for string containment.
- Use `in` for membership in a provided array of acceptable values.
- Use `WithMatcher` only when the default Go regex implementation is not sufficient.

## Forbidden

- Do not add a direct `not` field to predicates other than `test`, `test_string`, and `test_string_len`.
- Do not use relative paths inside `and`, `or`, or `not`; child predicates must carry absolute paths.
- Do not use `contains` for array membership; use `in`.
- Do not document `matches` as requiring a custom matcher; the library ships with a default matcher.
- Do not describe `test_string_len` as exact-length matching; its contract is minimum length (`>= len`).

## Acceptance Criteria

- [ ] Each predicate operation has one canonical contract.
- [ ] Negation rules are explicit and non-duplicated.
- [ ] Composite predicate path rules are documented.
- [ ] String, numeric, and type predicates describe the actual payload fields they accept.

**Origin:** former docs file `json-predicate.md`.
