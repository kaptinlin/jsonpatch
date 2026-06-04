# Architecture Specs

## Overview

This spec defines the package boundaries and execution pipeline of `jsonpatch`.

## Package Roles

| Package | Responsibility |
|---------|----------------|
| root package (`patch.go`, `errors.go`, `index.go`, `util.go`) | Compiled patch API, structured errors, operation constants, dispatch by document shape, and compile-time capability policy |
| `op` | Executable operation implementations and shared apply helpers |
| `internal` | Shared interfaces, constants, apply options, and codec payload types |
| `codec/json` | Decode `codec/json.Operation` payloads into executable operations and encode operations back to JSON form |
| `codec/compact` | Compact array codec |
| `codec/binary` | Binary codec |

## Interface Hierarchy

| Interface | Contract |
|-----------|----------|
| `internal.Op` | Executable operation with `Op`, `Path`, `Apply`, and `Validate`. |
| `internal.JSONOp` | `Op` plus `ToJSON` for JSON projection and compile-time freezing. |
| `internal.CompactOp` | `Op` plus `Code` and `ToCompact` for compact-array projection. |
| `internal.PredicateOp` | `Op` plus `Test` and `Not`. |
| `internal.SecondOrderPredicateOp` | `PredicateOp` plus child predicate access through `Ops`. |
| `internal.Codec` | Encode and decode operations between wire formats and executable operations. |

## Compiled Execution Pipeline

1. `Compile`, `CompileOps`, `CompileOperations`, or `CompileJSON` creates a `Patch`.
2. JSON-shaped inputs decode through `codec/json` before compile policy is applied.
3. Compile policy validates operation shape and rejects operation families outside enabled capabilities.
4. `Apply` dispatches by runtime document shape and clones the working document.
5. `ApplyInPlace` dispatches by runtime document shape with mutation enabled and writes the final result back to the caller's variable.
6. Operations run sequentially, and each operation's output document becomes the next operation's input.
7. The final document is converted back to the caller's original type, and successful operation facts become `Step` values.

## Document-Shape Dispatch

| Input shape | Path |
|-------------|------|
| `map[string]any` | Direct apply without JSON round-trip |
| `[]byte` | JSON decode → apply → JSON encode |
| `JSONText` | JSON decode → apply → JSON encode |
| `string` | Scalar-string apply |
| Structs and other concrete types | JSON marshal → apply → JSON unmarshal |
| Primitives and `[]any` | Direct apply |

> **Why**: The root package owns shape dispatch so operation implementations can stay focused on patch behavior instead of type conversion and codec concerns.
>
> **Rejected**: Pushing document conversion into each operation would duplicate conversion rules across the library. Collapsing codec logic into the root package would make alternative encodings harder to support and test.

## Codec Wire Contract

- Compact and binary codecs use path segment arrays as their only path representation.
- Compact and binary operation arrays are `[code, path, ...payload]`; `move` and `copy` use `[code, path, from]`.
- `test_string` uses `[code, path, pos, str, not?]`.
- Optional boolean fields are emitted only when true: `test.not`, `test_string.not`, `test_string_len.not`, and `ignore_case` for `matches`, `contains`, `starts`, and `ends`. `extend.deleteNull` is likewise omitted when false.
- Optional structural payloads such as `split.props` and `merge.props` are omitted when absent.
- Composite predicates encode child predicate paths relative to the containing predicate path. Decoding merges those paths into executable absolute paths.
- Binary supports the same operation tree as compact, including `and`, `or`, and unary `not`.

## Dependency Rules

- The root package is the public entry point.
- `op` depends on `internal` contracts and helpers, not on the root package.
- Codec packages translate between wire formats and `internal.Op`; they do not own patch execution.
- JSON and compact encode paths require the decoded operation value to implement the matching projection interface and fail when a custom executable operation cannot represent itself in that wire format.
- Compile capability policy belongs to the root package, after codec decoding and before operation application.
- `internal` defines contracts and shared constants only.

## Forbidden

- Do not bypass the codec layer when the public input is `[]Operation`; JSON-shaped operations are decoded through `codec/json`.
- Do not create package cycles between root, `op`, `codec/*`, and `internal`.
- Do not put architectural promises in `CLAUDE.md`; keep CLAUDE operational and record tested architecture here.
- Do not let format-specific codecs define behavioral semantics that conflict with executable operations.
- Do not put codec names in capability policy; codecs are wire formats, not operation vocabularies.

## Acceptance Criteria

- [ ] Each package has one documented responsibility.
- [ ] The path from compilation to operation execution is explicit.
- [ ] Document-shape conversion rules are defined once.
- [ ] Codec responsibilities stay separate from operation behavior.
