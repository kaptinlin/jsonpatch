# API Specs

## Overview

This spec defines the public entry points and the RFC 6902 mutating operation contracts exposed by `github.com/kaptinlin/jsonpatch`.

The `test` operation is part of RFC 6902, but its behavioral contract lives in `SPECS/25-predicate-specs.md` so predicate rules are defined in exactly one place.

## Public Entry Points

| API | Input model | Contract |
|-----|-------------|----------|
| `Compile(ops ...Op)` | Go-built operation values | Compiles operations with the default RFC 6902 capability and returns a reusable `Patch`. Operations must be able to freeze themselves for compiled patch storage. |
| `CompileOps(ops []Op, opts ...CompileOption)` | Go-built operation values | Compiles operations with explicit compile options such as capabilities. Operations must be able to freeze themselves for compiled patch storage. |
| `CompileOperations(ops []codec/json.Operation, opts ...CompileOption)` | JSON-shaped `codec/json.Operation` values | Decodes the structured field-bag shape through the JSON codec and compiles the resulting operations. |
| `CompileJSON(data []byte, opts ...CompileOption)` | JSON patch document bytes | Decodes a JSON patch document and compiles it with operation-family policy. |
| `Apply[T Document](patch *Patch, doc T)` | Compiled patch and one document | Applies the patch immutably and returns `Result[T]`. |
| `ApplyInPlace[T Document](patch *Patch, doc *T)` | Compiled patch and document pointer | Applies the patch with mutation enabled and writes the final result back to `doc`. |

## Compile Options

| Compile option | Contract |
|----------------|----------|
| `WithCapabilities(caps...)` | Sets the allowed operation families. Default compilation accepts only RFC 6902 operations. |
| `WithCompileMatcher(factory)` | Binds the regex matcher factory used when compiling `matches` operations from JSON-shaped input. The factory returns `(RegexMatcher, error)` so invalid patterns fail compilation. |

> **Why**: The compiled patch path gives callers one stable lifecycle: compile operation vocabulary once, then apply it to documents. Mutation has its own entry point so destructive application is visible at the call site.
>
> **Rejected**: A single untyped entry point returning `any` would throw away the type-preserving API. A runtime mutation option would make a destructive action look like ordinary configuration.

## Capability Contract

- Default compilation enables `RFC6902` only.
- `Predicate` enables non-regex predicate operations.
- `RegexPredicate` enables `matches`; it is separate because regex matching has its own safety and semantic boundary.
- `Extended` enables JSON Patch Extended operations.
- Codec choice is not a capability. JSON, compact, and binary codecs translate wire formats; compile policy decides whether decoded operations may run.

## RFC 6902 Mutating Operations

| Operation | Required fields | Contract |
|-----------|-----------------|----------|
| `add` | `path`, `value` | Insert or replace at the target path. Empty path replaces the entire document. `/-` appends to arrays. |
| `remove` | `path` | Remove an existing value at the target path. Empty path removes the root and yields `nil`; missing targets fail. |
| `replace` | `path`, `value` | Replace an existing value. Empty path replaces the entire document. |
| `move` | `path`, `from` | Move a value from `from` to `path`. Empty `from` means the root document. Validation rejects moving into a descendant of `from`. |
| `copy` | `path`, `from` | Copy a value from `from` to `path` using `add` target semantics, including array insertion and `/-` append. Empty `from` means the root document. |

A failing `move` leaves the document as it was at the start of that operation. Under `ApplyInPlace`,
operations that completed before the failing move remain applied; the failed move does not leave its
source removed.

## Apply Failure and Conversion Contract

- `Apply` executes on a clone. The first operation failure returns a nil result and leaves the caller's
  document unchanged.
- `ApplyInPlace` executes operations sequentially without promising rollback. When an operation fails,
  the caller's document is left at the complete successful prefix when that prefix can be represented by
  the original document type. When the first operation fails, that prefix is empty and the original
  carrier is not decoded, encoded, or otherwise rewritten.
- JSON text/bytes, typed or named containers, and structs execute through a logical document that is
  restored to the caller's type at the API boundary. An intermediate state need not be representable by
  that type when later operations restore a representable final state.
- If an operation fails after producing a successful prefix that cannot be restored to the caller's type,
  the adapter does not write back that prefix. The returned error matches both the original operation
  failure and `ErrConversionFailed`.
- If every operation succeeds but the final document cannot be restored to the caller's type, the call
  returns `ErrConversionFailed` and does not perform adapter writeback.
- A conversion failure that cannot be assigned to one operation carries no operation context and reports
  index `-1`. It is a document-boundary failure, not a second operation lifecycle.
- Direct maps and slices may already contain successful in-place mutations when a later failure occurs.
  `ApplyInPlace` does not clone them to manufacture rollback semantics.

> **Why**: Operations transform one logical document in order, while conversion answers a different
> question: whether the observable result can cross the caller's Go type boundary. Converting after every
> operation would reject patches whose temporary state is repaired before return and would make in-place
> application perform a whole-document conversion per operation.
>
> **Rejected**: Per-operation conversion checkpoints, implicit coercion, rollback options, and public
> partial-result APIs add competing execution models without a caller need.

## Compile Boundary Contract

- `Compile`, `CompileOps`, `CompileOperations`, and `CompileJSON` reject invalid operation shape before any document is touched.
- Capability policy is enforced at compile time. `matches` requires `RegexPredicate`; non-regex predicates require `Predicate`; extended operations require `Extended`.
- Invalid default regex patterns and custom matcher factory rejections are payload errors. The structured root error preserves the matcher construction cause.
- Operation family, required capability, and compact/binary code come from the internal operation vocabulary spine. Codec payload fields and operation constructors remain owned by the codec and operation packages.
- Empty `path` and `from` values are valid JSON Pointers that target the root document. Missing field presence is a raw JSON/map concern and is enforced by the JSON codec, not by zero-value `codec/json.Operation` structs.
- `nil` `value` in a `codec/json.Operation` means JSON `null` for `add`, `replace`, and `test`; raw JSON decoding still rejects omitted required `value` fields.
- Go-built operations are cloned through the executable operation layer during compilation. `Compile` and `CompileOps` do not use JSON projection or JSON codec decoding to freeze operations.
- Patch payloads that become document values are cloned during compilation and insertion so later mutation of caller-provided operation fields cannot mutate the compiled patch or result document.
- Composite predicate cloning rejects cycles in public `and`, `or`, and `not` operand graphs with `op.ErrPredicateCycle`; root compilation wraps that cause as `ErrPayloadInvalid`. Reusing one operand in multiple acyclic branches remains valid and each branch receives an independent clone.

## Error Contract

- `Compile`, `CompileOps`, `CompileOperations`, and `CompileJSON` return structured `*Error` values for invalid payloads and unsupported capabilities.
- The public `Op` method set includes `Clone`, so every operation accepted by `Compile` and `CompileOps` can freeze itself for compiled storage. Clone failures are structured payload errors. The package does not promise a plugin runtime for arbitrary external operation implementations.
- `Apply` and `ApplyInPlace` return structured `*Error` values for runtime conflicts, failed predicates, type mismatches, and conversion failures.
- `Apply` returns no partial `Result` when an operation fails; the structured error index identifies the first failed operation.
- `*Error` supports `errors.Is` for stable failure classes and `errors.As` for operation index, op, path, from, codec, and cause context.
- Predicate mismatches are classified from the operation-layer `op.ErrTestFailed` base. The returned root error matches `jsonpatch.ErrTestFailed` and preserves the concrete operation cause; predicate missing-path and wrong-target-type results do not become runtime-conflict or root type-mismatch classifications.
- Non-predicate runtime value type failures are classified from `op.ErrTypeMismatch`. The returned root error matches `jsonpatch.ErrTypeMismatch` and preserves the concrete operation cause. Missing paths, index failures, invalid positions, and other state conflicts remain `jsonpatch.ErrRuntimeConflict`.
- Execution errors are wrapped with operation index context when they happen during a sequence.
- Compile and execution errors are intended to be matched with `errors.Is` against sentinel errors.
- JSON, compact, and binary codecs expose codec-local sentinels for codec encode/decode failures. Root compile entry points wrap codec failures in the root `*Error` surface with codec and operation context.

## Forbidden

- Do not add a second validation path parallel to compilation.
- Do not add runtime options for mutation; use `ApplyInPlace`.
- Do not duplicate the `test` operation contract here; `SPECS/25-predicate-specs.md` records predicate behavior.
- Do not model whole-object replacement with `extend`; use `replace` when the target value must be replaced atomically.

## Acceptance Criteria

- Public entry points compile and apply through the single lifecycle documented above. Verification: `patch_test.go`, `readme_examples_test.go`, and `go test ./...`.
- Compile-time validation, clone isolation, matcher construction, and capability failures preserve structured error context. Verification: `compile_clone_test.go`, `patch_test.go`, and `codec_presence_test.go`.
- RFC 6902 mutation semantics are verified through root, operation, and scenario tests; predicate behavior remains owned by `SPECS/25-predicate-specs.md`. Verification: `op/*_test.go`, `tests/basic/`, and `tests/scenarios/`.
- `ApplyInPlace` success-prefix and conversion-boundary behavior is verified by the `TestApplyInPlace*` tests in `patch_test.go`.
