# Extended Operation Specs

## Overview

This spec defines the non-RFC operations implemented by `jsonpatch`: `flip`, `inc`, `str_ins`, `str_del`, `split`, `merge`, and `extend`.

These operations keep the extended patch vocabulary useful in Go while fitting the library's type-preserving API.

## Operation Contracts

| Operation | Payload | Contract |
|-----------|---------|----------|
| `flip` | `path` | Toggle a boolean target value. |
| `inc` | `path`, `inc` | Increment a numeric target after JavaScript-like `Number()` coercion. A missing final path is treated as `0` and the target is created with the incremented value. |
| `str_ins` | `path`, integer `pos`, `str` | Insert `str` into a string target at Unicode code-point position `pos`. Negative positions count from the end. |
| `str_del` | `path`, integer `pos`, and either integer `len` or `str` | Delete Unicode code points starting at `pos`. When `str` is supplied it takes precedence over `len`. Negative positions count from the end. |
| `split` | `path`, integer `pos`, optional `props` | Split a target value. Strings split by Unicode code-point index, numbers split into `[pos, original-pos]`, Slate-like nodes split into two nodes, and array-element splits replace one element with two. |
| `merge` | `path`, integer `pos`, optional `props` | Merge the array element at `path` with its immediately preceding sibling. `pos` and `props` are reverse-split metadata and do not select the target. Strings concatenate, numbers add, Slate-like nodes merge structurally, and unsupported pairs fall back to `[one, two]`. |
| `extend` | `path`, `props`, optional `deleteNull` | Shallow-extend an object target. When `deleteNull` is true, incoming `nil` values delete keys instead of storing `null`. |

## Behavioral Notes

String positions and lengths count Unicode code points. They are not UTF-8 byte offsets, JavaScript UTF-16 code-unit offsets, or grapheme-cluster offsets; combining marks count separately.

`str_ins`, `str_del`, `split`, and `merge` constructors use Go `int` coordinates. Their JSON, compact, and binary decoders reject fractional, non-finite, and out-of-range values instead of truncating them. Binary encoding uses MessagePack integers for these fields so every Go `int` value round-trips exactly. Number-target splitting converts the integer boundary to `float64` only for the numeric calculation.

Wrong target types for `str_ins`, `str_del`, `inc`, and `extend` match `op.ErrTypeMismatch` and retain their operation-specific `ErrNotString`, `ErrNotNumber`, or `ErrNotObject` cause. Path, index, position, and other state failures are not type mismatches.

### `inc`

- Root-level increment is allowed.
- `nil`, booleans, numeric strings, and numeric Go types are coerced through the shared `ToFloat64` helper.
- Missing final targets are treated as zero because Go does not use JavaScript `NaN` as a document value.

### `split`

- Empty path splits the root document and returns the split result as the new document.
- Negative and out-of-range string positions retain the same clamping behavior as other integer split positions.
- When splitting a string with `props`, the result is wrapped into text-node-like maps so both halves receive the supplied properties.
- Slate-like text and element nodes are handled specially; other unsupported values split into a two-element tuple of the same value.
- `props` may override ordinary node metadata, but the computed `text` or `children` field is always written last.

### `merge`

- `path` identifies the second element of an adjacent pair; `/1` targets the second element of a root array and `/nodes/1` targets the second element of a nested array.
- The target index must be canonical, present, and greater than zero. The containing value must be an array.
- Negative, zero, and positive `pos` values are preserved as reverse-split metadata and do not participate in target resolution or bounds checks.
- `props` are reverse-split metadata and do not modify the merged value.

### `extend`

- Empty path extends the root object.
- `extend` is shallow, not recursive.
- `__proto__` is an ordinary `map[string]any` key in Go and is preserved during extension. JavaScript prototype-pollution concerns belong in a JavaScript adapter or codec boundary, not in Go document mutation.

> **Why**: Extended operations are intentionally narrow. They encode concrete patch behaviors instead of inventing a generic transformation language.
>
> **Rejected**: Replacing `extend` with recursive merge semantics would hide too much behavior inside one operation. Expanding `merge` beyond arrays would break the operation's adjacency-based contract.

## Forbidden

- Do not use `inc` as a generic conversion operator for non-numeric targets.
- Do not describe `str_del` length mode as substring verification; only `str` mode binds deletion to specific text.
- Do not point `merge.path` at the containing array; point it at the second element of the adjacent pair.
- Do not let split or merge metadata override computed `text` or `children` fields.
- Do not rely on `extend` to preserve `nil` properties when `deleteNull` is true.
- Do not document `split` as string-only; numbers and Slate-like nodes are part of the contract.

## Acceptance Criteria

- Extended operations follow the payload, Unicode-coordinate, and runtime type rules in this spec. Verification: `op/*_test.go`, `patch_test.go`, and `tests/extended/`.
- `inc` preserves numeric coercion and missing-final-target behavior. Verification: `op/inc_test.go` and extended scenario tests.
- `split`, `merge`, and `extend` preserve their structural, metadata, and ordinary-map-key semantics. Split and merge positions round-trip losslessly through every codec. Verification: `op/split_test.go`, `op/merge_test.go`, `op/extend_test.go`, `codec_coordinate_test.go`, and `tests/extended/merge_test.go`.
