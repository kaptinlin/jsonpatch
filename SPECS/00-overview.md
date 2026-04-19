# Overview

## Scope

`jsonpatch` applies RFC 6902 operations, predicate operations, and extended operations to Go values while preserving the caller's document type whenever the result can be converted back safely.

Canonical technical contracts live in `SPECS/`. `README.md` is user-facing, and `CLAUDE.md` is operational guidance for agents.

## Library Contract

- The public module is `github.com/kaptinlin/jsonpatch`.
- The library targets Go 1.26.
- The compatibility baseline is json-joy's JSON Patch behavior, including predicate and extended operations.
- Mutation is opt-in. `ApplyPatch`, `ApplyOp`, and `ApplyOps` clone the input unless `WithMutate(true)` is supplied.
- The library stays pure: operations return errors and never log.

> **Why**: The library exists to give Go callers the same operation vocabulary as json-joy without giving up static typing or forcing every caller into `map[string]any`.
>
> **Rejected**: A map-only API would discard struct typing and make the library less useful in Go applications. Implicit mutation would make patch application surprising and harder to reason about in shared state.

## Supported Document Shapes

| Input shape | Processing model | Result shape |
|-------------|------------------|--------------|
| `map[string]any` | Applied directly | `map[string]any` |
| `[]byte` | Decoded as JSON, patched, re-encoded | `[]byte` |
| `string` starting with `{` or `[` | Parsed as JSON, patched, re-encoded | `string` |
| Other `string` | Treated as a plain scalar string | `string` |
| Structs and other concrete types | Marshaled to JSON, patched as untyped data, unmarshaled back | Original Go type |
| Primitive values and `[]any` | Applied directly when the result remains assignable | Original Go type |

## Spec Map

- `SPECS/20-api-specs.md` — public entry points and RFC 6902 mutating operations
- `SPECS/25-predicate-specs.md` — predicate and second-order predicate contracts
- `SPECS/26-extended-operation-specs.md` — extended operation contracts
- `SPECS/30-data-model-specs.md` — operation payload, type vocabulary, and result shapes
- `SPECS/40-architecture-specs.md` — package boundaries and execution pipeline
- `SPECS/50-coding-standards.md` — compatibility, errors, tests, and documentation rules

## Forbidden

- Do not treat `README.md` examples as the canonical contract when a `SPECS/*.md` file says otherwise.
- Do not assume string inputs are always JSON documents; only strings beginning with `{` or `[` are parsed as JSON.
- Do not depend on implicit mutation; use `WithMutate(true)` when in-place updates are required.
- Do not add new technical design rules to `CLAUDE.md`; add or update a spec in `SPECS/` instead.

## Acceptance Criteria

- [ ] The canonical behavior of the library is defined only in `SPECS/`.
- [ ] Each supported document shape has a documented processing rule.
- [ ] Mutation, typing, and compatibility defaults are explicit.
- [ ] `README.md` and `CLAUDE.md` can stay concise because this file owns the overview contract.

**Origin:** `CLAUDE.md` (Project Overview, Design Philosophy).
