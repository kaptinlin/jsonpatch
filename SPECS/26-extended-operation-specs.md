# Extended Operation Specs

## Overview

This spec defines the non-RFC operations implemented by `jsonpatch`: `flip`, `inc`, `str_ins`, `str_del`, `split`, `merge`, and `extend`.

These operations exist to match the json-joy extended patch vocabulary while still fitting Go's type-preserving API.

## Operation Contracts

| Operation | Payload | Contract |
|-----------|---------|----------|
| `flip` | `path` | Toggle a boolean target value. |
| `inc` | `path`, `inc` | Increment a numeric target after JavaScript-like `Number()` coercion. A missing final path is treated as `0` and the target is created with the incremented value. |
| `str_ins` | `path`, `pos`, `str` | Insert `str` into a string target at rune position `pos`. Negative positions count from the end. |
| `str_del` | `path`, `pos`, and either `len` or `str` | Delete runes starting at `pos`. When `str` is supplied it takes precedence over `len`. Negative positions count from the end. |
| `split` | `path`, `pos`, optional `props` | Split a target value. Strings split by rune index, numbers split into `[pos, original-pos]`, Slate-like nodes split into two nodes, and array-element splits replace one element with two. |
| `merge` | `path`, `pos`, optional `props` | Merge array elements at `pos-1` and `pos`. Strings concatenate, numbers add, Slate-like nodes merge structurally, and unsupported pairs fall back to `[one, two]`. |
| `extend` | `path`, `props`, optional `deleteNull` | Shallow-extend an object target. When `deleteNull` is true, incoming `nil` values delete keys instead of storing `null`. |

## Behavioral Notes

### `inc`

- Root-level increment is allowed.
- `nil`, booleans, numeric strings, and numeric Go types are coerced through the shared `ToFloat64` helper.
- Missing final targets are treated as zero because Go does not use JavaScript `NaN` as a document value.

### `split`

- Empty path splits the root document and returns the split result as the new document.
- When splitting a string with `props`, the result is wrapped into text-node-like maps so both halves receive the supplied properties.
- Slate-like text and element nodes are handled specially; other unsupported values split into a two-element tuple of the same value.

### `merge`

- `merge` operates on arrays and interprets `pos` as the index of the second element in the pair.
- `pos <= 0` is invalid.
- `props` are applied only to merged Slate-like node results.

### `extend`

- Empty path extends the root object.
- `extend` is shallow, not recursive.
- `__proto__` keys are ignored during extension.

> **Why**: Extended operations are intentionally narrow. They encode the behaviors the compatibility target already defines instead of inventing a generic transformation language.
>
> **Rejected**: Replacing `extend` with recursive merge semantics would hide too much behavior inside one operation. Expanding `merge` beyond arrays would break the operation's adjacency-based contract.

## Forbidden

- Do not use `inc` as a generic conversion operator for non-numeric targets.
- Do not describe `str_del` length mode as substring verification; only `str` mode binds deletion to specific text.
- Do not use `merge` on non-array targets.
- Do not rely on `extend` to preserve `nil` properties when `deleteNull` is true.
- Do not document `split` as string-only; numbers and Slate-like nodes are part of the contract.

## Acceptance Criteria

- [ ] Each extended operation has one canonical payload definition.
- [ ] Numeric coercion and missing-path behavior for `inc` are explicit.
- [ ] `split`, `merge`, and `extend` document their non-obvious structural behavior.
- [ ] Security-relevant `extend` behavior (`__proto__` skipping) is preserved in the spec.

**Origin:** former docs file `json-patch-extended.md`.
