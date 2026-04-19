# Coding Standards

## Overview

This spec defines the compatibility, error-handling, testing, and documentation rules that keep `jsonpatch` consistent.

## Core Rules

- Prefer type-preserving APIs over `any`-only convenience layers.
- Maintain json-joy-compatible behavior unless the Go implementation has a documented, deliberate adaptation.
- Default to immutable application; require `WithMutate(true)` for in-place updates.
- Keep the library pure: return errors and leave logging to callers.

> **Why**: The library is valuable because it combines Go's type system with a compatibility-oriented patch vocabulary. Every rule here protects one of those two properties.
>
> **Rejected**: Hidden mutation for speed would make patches harder to reason about. Logging inside operations would leak policy into a library package.

## Compatibility Rules

- Only `test`, `test_string`, and `test_string_len` support direct `not`.
- Composite predicates use `path: ""` and absolute paths inside `apply`.
- Numeric operations use shared JavaScript-like `Number()` coercion for target values.
- `test_type` uses `type` as the canonical field for one-or-many JSON type names.
- When Go cannot represent the compatibility target literally, the adaptation must be documented in the spec that owns that behavior.

## Error Rules

- Use sentinel errors for stable failure classes.
- Add dynamic context with `fmt.Errorf("%w: ...", err)`.
- Match failures with `errors.Is`; do not inspect error message text.
- Validation and execution errors are separate concerns: validation checks payload shape, execution checks runtime document state.

## Testing Rules

- Use table-driven tests for behavior with more than one case.
- Use `t.Parallel()` for top-level tests and subtests when the cases are independent.
- Include success and failure cases for every new behavior.
- Use `testing.B.Loop()` for new benchmarks.
- Do not add `_test.go` files whose only purpose is to enforce `SPECS/` layout or markdown link structure.

## Documentation Rules

- `SPECS/` is the canonical home for technical behavior.
- `README.md` stays user-facing and example-oriented.
- `CLAUDE.md` stays operational and should point to `SPECS/` instead of restating contracts.
- `SPECS/**/*.md` must stay covered by markdownlint and pre-commit checks.

## Forbidden

- Do not compare error strings.
- Do not add direct `not` support to unsupported predicates.
- Do not exclude `SPECS/**/*.md` from markdownlint.
- Do not put canonical API or architecture rules back into `README.md` or `CLAUDE.md`.
- Do not add layout-guard tests for `SPECS/`; enforce documentation quality through markdownlint and review.

## Acceptance Criteria

- [ ] Compatibility-sensitive rules are explicit.
- [ ] Error matching rules use sentinels and wrapping.
- [ ] New tests follow the repo's table-driven and parallel-safe style.
- [ ] Documentation ownership between `SPECS/`, `README.md`, and `CLAUDE.md` is clear.

**Origin:** `CLAUDE.md` (Coding Rules, Testing), `validate.go`, `lefthook.yml`.
