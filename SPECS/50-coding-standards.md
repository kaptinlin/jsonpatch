# Coding Standards

## Overview

This spec defines the implementation, error-handling, testing, and documentation rules that keep `jsonpatch` consistent.

## Core Rules

- Prefer type-preserving APIs over `any`-only convenience layers.
- Use reference implementations as evidence when they clarify the patch vocabulary; prefer the Go implementation's tested behavior when a reference shape conflicts with this library's model.
- Default to immutable application; require `WithMutate(true)` for in-place updates.
- Keep the library pure: return errors and leave logging to callers.

> **Why**: The library is valuable because it combines Go's type system with a precise patch vocabulary. Every rule here protects one of those two properties.
>
> **Rejected**: Hidden mutation for speed would make patches harder to reason about. Logging inside operations would leak policy into a library package.

## Behavior Rules

- Only `test`, `test_string`, and `test_string_len` support direct `not`.
- Composite predicates merge child paths with their own `path`; use `path: ""` for root-scoped children.
- `not` has exactly one direct child predicate; use explicit `and` or `or` to negate a sequence.
- Compact and binary encoders emit only one wire shape: segment-array paths, parent-relative composite children, and no false optional booleans.
- Numeric operations use shared JavaScript-like `Number()` coercion for target values.
- `test_type` uses `type` for one-or-many JSON type names.
- When a reference implementation cannot be represented cleanly in Go, document the Go behavior in the spec that owns it.

## Error Rules

- Use sentinel errors for stable failure classes.
- Add dynamic context with `fmt.Errorf("%w: ...", err)`.
- Match failures with `errors.Is`; do not inspect error message text.
- Validation and execution errors are separate concerns: validation checks payload shape, execution checks runtime document state.

## Testing Rules

- Use table-driven tests for behavior with more than one case.
- Use `t.Parallel()` for top-level tests and subtests when the cases are independent.
- Include success and failure cases for every new behavior.
- Protect codec wire contracts with focused golden tests: compact arrays, binary bytes, optional-field omission, parent-relative predicate paths, and JSON field presence.
- Use `testing.B.Loop()` for new benchmarks.
- Do not add `_test.go` files whose only purpose is to enforce `SPECS/` layout or markdown link structure.

## Documentation Rules

- `SPECS/` records technical behavior protected by code and tests.
- `README.md` stays user-facing and example-oriented.
- `CLAUDE.md` stays operational and should point to `SPECS/` instead of restating contracts.
- `SPECS/**/*.md` must stay covered by markdownlint and pre-commit checks.

## Forbidden

- Do not compare error strings.
- Do not add direct `not` support to unsupported predicates.
- Do not exclude `SPECS/**/*.md` from markdownlint.
- Do not put untested API or architecture promises into `README.md` or `CLAUDE.md`.
- Do not add layout-guard tests for `SPECS/`; enforce documentation quality through markdownlint and review.

## Acceptance Criteria

- [ ] Behavior-sensitive rules are explicit.
- [ ] Error matching rules use sentinels and wrapping.
- [ ] New tests follow the repo's table-driven and parallel-safe style.
- [ ] Documentation ownership between `SPECS/`, `README.md`, and `CLAUDE.md` is clear.
