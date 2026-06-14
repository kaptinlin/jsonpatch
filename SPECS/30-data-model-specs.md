# Data Model Specs

## Overview

This spec defines the data structures shared by the public API, codecs, and operation implementations.

## `codec/json.Operation` Payload

`codec/json.Operation` is the JSON-shaped input format for `CompileOperations` and JSON codec conversion.

Because `codec/json.Operation` is a Go struct, it cannot represent raw JSON field presence. Empty `path` and `from` are valid root pointers, and `nil` `value` is a valid JSON `null` payload. Raw JSON/map decoding owns missing-field validation when presence matters.

| Field | Used by | Contract |
|-------|---------|----------|
| `op` | all operations | Operation name. |
| `path` | most operations | JSON Pointer target path. |
| `value` | `add`, `replace`, `test`, `type`, `contains`, `starts`, `ends`, `in`, `less`, `more`, `matches` | Primary payload field for operations that consume one value. |
| `from` | `move`, `copy` | Source JSON Pointer. |
| `inc` | `inc` | Numeric delta. `0` is meaningful and therefore not omitted. |
| `pos` | `str_ins`, `str_del`, `split`, `merge`, `test_string` | Position field. `0` is meaningful and therefore not omitted. |
| `str` | `str_ins`, `str_del`, `test_string` | String operand. Empty string is meaningful and therefore not omitted. |
| `len` | `str_del`, `test_string_len` | Length operand. `0` is meaningful and therefore not omitted. |
| `not` | `test`, `test_string`, `test_string_len` | Direct negation flag. |
| `type` | `test_type` | One JSON type name or a list of type names. |
| `ignore_case` | string and regex predicates | Case-insensitive matching flag when supported. |
| `apply` | `and`, `or`, `not` | Nested predicate operations. |
| `props` | `extend`, `split`, `merge` | Object properties used by structural extended operations. |
| `deleteNull` | `extend` | Delete keys whose incoming property value is `nil` instead of storing them. |
| `oldValue` | `remove`, `replace`, encoded prior-value payloads | Optional prior value metadata. |

## Compact Operation Payload

`CompactOperation` is an array DTO for the compact and binary codecs, not a separate semantic model. It uses numeric operation codes, segment-array paths, and the minimal payload needed to reconstruct an executable operation.

| Shape | Contract |
|-------|----------|
| `[code, path]` | Path-only operations. |
| `[code, path, value]` | Value operations and required scalar payloads. |
| `[code, path, from]` | `move` and `copy`, where both paths are segment arrays. |
| `[code, path, pos, str, not?]` | `test_string`; `not` appears only when true. |
| `[code, path, apply]` | `and`, `or`, and `not`; child paths inside `apply` are parent-relative. |

## Supported JSON Type Names

| Name | Meaning |
|------|---------|
| `string` | JSON string |
| `number` | Any numeric value |
| `integer` | Whole-number numeric value |
| `boolean` | JSON boolean |
| `object` | JSON object |
| `array` | JSON array |
| `null` | JSON null |

## Document and Result Types

### `Document`

The generic API accepts values matching the `Document` constraint and dispatches by runtime shape so the result can be converted back to the caller's type.

Document dispatch is a closed classifier:

- `JSONText`, `[]byte`, and byte-slice aliases are JSON text and must parse as JSON.
- Plain `string` and string aliases are scalar string documents.
- `map[string]any`, `[]any`, interface values, numbers, and booleans apply directly.
- Struct-like values are marshaled to JSON-shaped data, patched, and unmarshaled back to the original type.

### `JSONText`

`JSONText` is a string wrapper that marks a document as JSON text for the compiled patch path. Plain `string` values are scalar string documents; `JSONText` values are decoded as JSON, patched, and encoded back to `JSONText`.

### `Patch`

`Patch` is a compiled operation sequence. It stores operations accepted by compile-time capability policy and can be reused with `Apply` or `ApplyInPlace`.

Compiled operations are independent executable clones, so later mutation of the caller-provided operation value or payload does not change the compiled patch. `Compile` and `CompileOps` freeze Go-built operations without JSON projection. Executable operations that cannot freeze themselves for compiled storage are rejected at compile time. Regex matcher behavior for `matches` operations decoded from JSON-shaped input is bound through `WithCompileMatcher` when callers need a custom matcher.

### `Capability`

| Capability | Contract |
|------------|----------|
| `RFC6902` | Core JSON Patch operations. |
| `Predicate` | Non-regex predicate operations. |
| `RegexPredicate` | `matches` predicate operations. |
| `Extended` | JSON Patch Extended operations. |
| `AllCapabilities` | All operation vocabularies implemented by the package. |

Capabilities describe operation vocabulary only; codecs remain wire-format translators.

Operation family, capability, and compact/binary code share one internal vocabulary spine. Payload field presence, nullability, and constructor rules stay with the JSON codec and executable operations instead of moving into a global manifest.

### `Result[T]`

| Field | Contract |
|-------|----------|
| `Doc` | The final patched document converted back to `T`. |
| `Steps` | Per-operation facts for successfully applied operations. |

### `Step`

`Step` exposes accessor methods for operation facts: `Index`, `Op`, `Path`, `From`, `Old`, and `Applied`. It does not expose a typed per-step document.

### `Error`

`Error` carries a stable failure kind plus optional operation index, op, path, from, codec, and cause context. It supports `errors.Is` through its kind and cause, and supports `errors.As` for callers that need structured context.

> **Why**: One shared payload type keeps JSON decoding, examples, and public API usage aligned while a separate executable `Op` layer keeps runtime behavior explicit.
>
> **Rejected**: A different payload struct per operation would make JSON patch assembly much harder for callers. Returning only the final document would remove per-operation `Old` values that callers sometimes need for auditing or post-processing.

## Forbidden

- Do not introduce a second JSON-shaped patch payload type outside `codec/json.Operation`.
- Do not use `value` as the encoded multi-type field for `test_type`; use `type`.
- Do not treat `Old` as universally populated; only some operations return it.
- Do not treat `Step` as a typed document history; it records operation facts only.
- Do not remove zero-value fields such as `inc`, `pos`, `str`, or `len` from the model contract; zero is meaningful for these fields.

## Acceptance Criteria

- [ ] The payload fields used by every operation family are documented once.
- [ ] The JSON type vocabulary is explicit.
- [ ] Result shapes and ordering semantics are documented.
- [ ] The difference between JSON-shaped `codec/json.Operation` values and executable operations remains clear.
