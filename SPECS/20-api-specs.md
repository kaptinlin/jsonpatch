# API Specs

## Overview

This spec defines the public entry points and the RFC 6902 mutating operation contracts exposed by `github.com/kaptinlin/jsonpatch`.

The `test` operation is part of RFC 6902, but its behavioral contract lives in `SPECS/25-predicate-specs.md` so predicate rules are defined in exactly one place.

## Public Entry Points

| API | Input model | Contract |
|-----|-------------|----------|
| `ApplyPatch[T Document](doc T, patch []Operation, opts ...Option)` | JSON-shaped `Operation` values | Decodes operations through the JSON codec, applies them in order, and returns `PatchResult[T]`. |
| `ApplyOp[T Document](doc T, op Op, opts ...Option)` | One executable `Op` instance | Applies one already-decoded operation and returns `OpResult[T]`. |
| `ApplyOps[T Document](doc T, ops []Op, opts ...Option)` | Executable `Op` instances | Applies decoded operations without a JSON decode round-trip. |
| `ValidateOperations(ops []Operation, allowMatchesOp bool)` | JSON-shaped `Operation` values | Performs preflight validation and returns the first validation error with operation index context. |
| `ValidateOperation(op Operation, allowMatchesOp bool)` | One JSON-shaped `Operation` value | Performs preflight validation for one operation. |

## Options

| Option | Contract |
|--------|----------|
| `WithMutate(true)` | Apply changes to the original working document instead of cloning first. |
| `WithMatcher(factory)` | Override the regex matcher factory used by the `matches` predicate. |

> **Why**: The library exposes both JSON-shaped and executable-operation entry points so callers can choose between ergonomic patch payloads and lower-level direct operation reuse without losing the same result model.
>
> **Rejected**: A single untyped entry point returning `any` would throw away the type-preserving API. Making validation implicit in `ApplyPatch` would prevent callers from using preflight validation as a separate step.

## RFC 6902 Mutating Operations

| Operation | Required fields | Contract |
|-----------|-----------------|----------|
| `add` | `path`, `value` | Insert or replace at the target path. `/-` appends to arrays. |
| `remove` | `path` | Remove an existing value at the target path. Missing targets fail. |
| `replace` | `path`, `value` | Replace an existing value. Empty path replaces the entire document. |
| `move` | `path`, `from` | Move a value from `from` to `path`. Validation rejects moving into a descendant of `from`. |
| `copy` | `path`, `from` | Copy a value from `from` to `path`. |

## Validation Contract

- `ValidateOperations` rejects `nil` patches, empty patches, missing required fields, invalid JSON Pointer values, and operation-specific shape errors.
- Validation uses the `allowMatchesOp` flag to permit or reject the `matches` predicate when the caller needs a restricted feature set.
- `ApplyPatch` decodes and applies operations directly. Call `ValidateOperations` yourself when a preflight validation step is required before execution.

## Error Contract

- `ApplyPatch`, `ApplyOp`, and `ApplyOps` stop at the first decode or apply failure.
- Execution errors are wrapped with operation index context when they happen during a sequence.
- Validation and execution errors are intended to be matched with `errors.Is` against sentinel errors.

## Forbidden

- Do not assume `ApplyPatch` performs an explicit `ValidateOperations` pass before execution.
- Do not use `ApplyPatch` when you already have `Op` instances and need to avoid JSON decode overhead; use `ApplyOp` or `ApplyOps`.
- Do not redefine the `test` operation contract here; `SPECS/25-predicate-specs.md` is the single source of truth for predicate behavior.
- Do not model whole-object replacement with `extend`; use `replace` when the target value must be replaced atomically.

## Acceptance Criteria

- [ ] Every public entry point has one documented contract.
- [ ] The RFC 6902 mutating operations are defined once and only once.
- [ ] Validation behavior is explicit, including the `allowMatchesOp` gate.
- [ ] Callers can tell when to use JSON-shaped operations versus executable `Op` instances.

**Origin:** former docs file `json-patch.md`.
