# Data Model Specs

## Overview

This spec defines the data structures shared by the public API, codecs, and operation implementations.

## `Operation` Payload

`Operation` is the JSON-shaped input format for `ApplyPatch`, validation, and codec conversion.

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
| `type` | `test_type` and nested encoded operations | One JSON type name or a list of type names. |
| `ignore_case` | string and regex predicates | Case-insensitive matching flag when supported. |
| `apply` | `and`, `or`, `not` | Nested predicate operations. |
| `props` | `extend`, `split`, `merge` | Object properties used by structural extended operations. |
| `deleteNull` | `extend` | Delete keys whose incoming property value is `nil` instead of storing them. |
| `oldValue` | `remove`, `replace`, encoded compatibility payloads | Optional prior value metadata. |

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

### `OpResult[T]`

| Field | Contract |
|-------|----------|
| `Doc` | The document after one operation completes. |
| `Old` | The previous value when the operation reports one. |

### `PatchResult[T]`

| Field | Contract |
|-------|----------|
| `Doc` | The final patched document. |
| `Res` | Per-operation results in application order. |

> **Why**: One shared payload type keeps JSON decoding, validation, examples, and public API usage aligned while a separate executable `Op` layer keeps runtime behavior explicit.
>
> **Rejected**: A different payload struct per operation would make JSON patch assembly much harder for callers. Returning only the final document would remove per-operation `Old` values that callers sometimes need for auditing or post-processing.

## Forbidden

- Do not introduce a second JSON-shaped patch payload type outside `Operation`.
- Do not use `value` as the canonical multi-type field for `test_type`; use `type`.
- Do not treat `Old` as universally populated; only some operations return it.
- Do not remove zero-value fields such as `inc`, `pos`, `str`, or `len` from the model contract; zero is meaningful for these fields.

## Acceptance Criteria

- [ ] The payload fields used by every operation family are documented once.
- [ ] The JSON type vocabulary is explicit.
- [ ] Result shapes and ordering semantics are documented.
- [ ] The difference between JSON-shaped `Operation` values and executable operations remains clear.

**Origin:** `CLAUDE.md` (Key Types and Interfaces), `internal/types.go`, `internal/jsonpatch_type.go`.
