# Overview

## Scope

`jsonpatch` applies RFC 6902 operations, predicate operations, and extended operations to Go values while preserving the caller's document type whenever the result can be converted back safely.

Recorded technical contracts live in `SPECS/`. `README.md` is user-facing, and `CLAUDE.md` is operational guidance for agents.

## Library Contract

- The public module is `github.com/kaptinlin/jsonpatch`.
- The library targets Go 1.26.
- The operation vocabulary follows JSON Patch, predicate, and extended patch behavior proven by the Go implementation and tests.
- Mutation is opt-in. Use `Apply` for immutable application and `ApplyInPlace` for explicit mutation.
- The compiled path defaults to the RFC 6902 operation vocabulary. Predicate, regex predicate, and extended operations require compile-time capabilities.
- The library stays pure: operations return errors and never log.

> **Why**: The library exists to give Go callers a rich JSON Patch+ vocabulary without giving up static typing or forcing every caller into `map[string]any`.
>
> **Rejected**: A map-only API would discard struct typing and make the library less useful in Go applications. Implicit mutation would make patch application surprising and harder to reason about in shared state.

## Supported Document Shapes

| Input shape | Processing model | Result shape |
|-------------|------------------|--------------|
| `map[string]any` | Applied directly | `map[string]any` |
| `[]byte` | Decoded as JSON, patched, re-encoded | `[]byte` |
| `JSONText` | Parsed as JSON, patched, re-encoded | `JSONText` |
| `string` | Treated as a plain scalar string | `string` |
| Structs and other concrete types | Marshaled to JSON, patched as untyped data, unmarshaled back | Original Go type |
| Primitive values and `[]any` | Applied directly when the result remains assignable | Original Go type |

`JSONText` and `[]byte` documents re-encode root `null` results as JSON `null`;
plain scalar strings fail conversion when an operation yields JSON null because the
result can no longer be represented as a string value.

## Spec Map

- `SPECS/20-api-specs.md` — public entry points and RFC 6902 mutating operations
- `SPECS/25-predicate-specs.md` — predicate and second-order predicate contracts
- `SPECS/26-extended-operation-specs.md` — extended operation contracts
- `SPECS/30-data-model-specs.md` — operation payload, type vocabulary, and result shapes
- `SPECS/40-architecture-specs.md` — package boundaries and execution pipeline
- `SPECS/50-coding-standards.md` — implementation, errors, tests, and documentation rules

## Forbidden

- Do not treat `README.md` examples as technical contracts; use code, tests, and `SPECS/` for behavior.
- Do not assume string inputs are JSON documents in the compiled path; use `JSONText` or `[]byte` when JSON text is intended.
- Do not depend on implicit mutation; use `ApplyInPlace` when in-place updates are required.
- Do not add new technical design rules to `CLAUDE.md`; add or update a spec in `SPECS/` instead.

## Acceptance Criteria

- [ ] `SPECS/` records behavior protected by code and tests.
- [ ] Each supported document shape has a documented processing rule.
- [ ] Mutation and typing defaults are explicit.
- [ ] `README.md` and `CLAUDE.md` stay concise because this file records the overview contract.
