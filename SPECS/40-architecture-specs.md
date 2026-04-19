# Architecture Specs

## Overview

This spec defines the package boundaries and execution pipeline of `jsonpatch`.

## Package Roles

| Package | Responsibility |
|---------|----------------|
| root package (`index.go`, `jsonpatch.go`, `validate.go`) | Public entry points, type aliases, dispatch by document shape, JSON-shaped validation helpers |
| `op` | Executable operation implementations and shared apply helpers |
| `internal` | Shared interfaces, constants, options, and payload/result types |
| `codec/json` | Decode `Operation` payloads into executable operations and encode operations back to JSON form |
| `codec/compact` | Compact array codec |
| `codec/binary` | Binary codec |

## Interface Hierarchy

| Interface | Contract |
|-----------|----------|
| `internal.Op` | Executable operation with `Op`, `Code`, `Path`, `Apply`, `ToJSON`, `ToCompact`, and `Validate`. |
| `internal.PredicateOp` | `Op` plus `Test` and `Not`. |
| `internal.SecondOrderPredicateOp` | `PredicateOp` plus child predicate access through `Ops`. |
| `internal.Codec` | Encode and decode operations between wire formats and executable operations. |

## Execution Pipeline

1. `ApplyPatch` dispatches by runtime document shape.
2. The JSON codec decodes `[]Operation` into executable operations, honoring `WithMatcher` when regex predicates are present.
3. `applyInternalOps` clones the working document unless `WithMutate(true)` is set.
4. Operations run sequentially, and each operation's output document becomes the next operation's input.
5. The final document is converted back to the caller's original type.

## Document-Shape Dispatch

| Input shape | Path |
|-------------|------|
| `map[string]any` | Direct apply without JSON round-trip |
| `[]byte` | JSON decode → apply → JSON encode |
| `string` | JSON decode for object/array strings, otherwise scalar-string apply |
| Structs and other concrete types | JSON marshal → apply → JSON unmarshal |
| Primitives and `[]any` | Direct apply |

> **Why**: The root package owns shape dispatch so operation implementations can stay focused on patch behavior instead of type conversion and codec concerns.
>
> **Rejected**: Pushing document conversion into each operation would duplicate conversion rules across the library. Collapsing codec logic into the root package would make alternative encodings harder to support and test.

## Dependency Rules

- The root package is the public entry point.
- `op` depends on `internal` contracts and helpers, not on the root package.
- Codec packages translate between wire formats and `internal.Op`; they do not own patch execution.
- `internal` defines contracts and shared constants only.

## Forbidden

- Do not bypass the codec layer when the public input is `[]Operation`; JSON-shaped operations are decoded through `codec/json`.
- Do not create package cycles between root, `op`, `codec/*`, and `internal`.
- Do not move technical contracts back into `CLAUDE.md`; keep architectural contracts here.
- Do not let format-specific codecs define behavioral semantics that conflict with executable operations.

## Acceptance Criteria

- [ ] Each package has one documented responsibility.
- [ ] The path from `ApplyPatch` to operation execution is explicit.
- [ ] Document-shape conversion rules are defined once.
- [ ] Codec responsibilities stay separate from operation behavior.

**Origin:** `CLAUDE.md` (Architecture), `jsonpatch.go`, `index.go`, `internal/interfaces.go`.
