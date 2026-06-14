# API Specs

## Overview

This spec defines the public entry points and the RFC 6902 mutating operation contracts exposed by `github.com/kaptinlin/jsonpatch`.

The `test` operation is part of RFC 6902, but its behavioral contract lives in `SPECS/25-predicate-specs.md` so predicate rules are defined in exactly one place.

## Public Entry Points

| API | Input model | Contract |
|-----|-------------|----------|
| `Compile(ops ...Op)` | Go-built operation values | Compiles operations with the default RFC 6902 capability and returns a reusable `Patch`. Operations must be able to freeze themselves for compiled patch storage. |
| `CompileOps(ops []Op, opts ...CompileOption)` | Go-built operation values | Compiles operations with explicit compile options such as capabilities. Operations must be able to freeze themselves for compiled patch storage. |
| `CompileOperations(ops []codec/json.Operation, opts ...CompileOption)` | JSON-shaped `codec/json.Operation` values | Decodes through the JSON codec and compiles the resulting operations. This is a migration boundary for the field-bag shape. |
| `CompileJSON(data []byte, opts ...CompileOption)` | JSON patch document bytes | Decodes a JSON patch document and compiles it with operation-family policy. |
| `Apply[T Document](patch *Patch, doc T)` | Compiled patch and one document | Applies the patch immutably and returns `Result[T]`. |
| `ApplyInPlace[T Document](patch *Patch, doc *T)` | Compiled patch and document pointer | Applies the patch with mutation enabled and writes the final result back to `doc`. |

## Compile Options

| Compile option | Contract |
|----------------|----------|
| `WithCapabilities(caps...)` | Sets the allowed operation families. Default compilation accepts only RFC 6902 operations. |
| `WithCompileMatcher(factory)` | Binds the regex matcher factory used when compiling `matches` operations from JSON-shaped input. |

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

## Compile Boundary Contract

- `Compile`, `CompileOps`, `CompileOperations`, and `CompileJSON` reject invalid operation shape before any document is touched.
- Capability policy is enforced at compile time. `matches` requires `RegexPredicate`; non-regex predicates require `Predicate`; extended operations require `Extended`.
- Operation family, required capability, and compact/binary code come from the internal operation vocabulary spine. Codec payload fields and operation constructors remain owned by the codec and operation packages.
- Empty `path` and `from` values are valid JSON Pointers that target the root document. Missing field presence is a raw JSON/map concern and is enforced by the JSON codec, not by zero-value `codec/json.Operation` structs.
- `nil` `value` in a `codec/json.Operation` means JSON `null` for `add`, `replace`, and `test`; raw JSON decoding still rejects omitted required `value` fields.
- Go-built operations are cloned through the executable operation layer during compilation. `Compile` and `CompileOps` do not use JSON projection or JSON codec decoding to freeze operations.
- Patch payloads that become document values are cloned during compilation and insertion so later mutation of caller-provided operation fields cannot mutate the compiled patch or result document.

## Error Contract

- `Compile`, `CompileOps`, `CompileOperations`, and `CompileJSON` return structured `*Error` values for invalid payloads and unsupported capabilities.
- `Compile` and `CompileOps` reject executable operations that cannot be cloned for compilation, because compiled patches must be isolated from later caller mutation. The package does not promise a public plugin runtime for arbitrary external operation implementations.
- `Apply` and `ApplyInPlace` return structured `*Error` values for runtime conflicts, failed predicates, type mismatches, and conversion failures.
- `*Error` supports `errors.Is` for stable failure classes and `errors.As` for operation index, op, path, from, codec, and cause context.
- Execution errors are wrapped with operation index context when they happen during a sequence.
- Compile and execution errors are intended to be matched with `errors.Is` against sentinel errors.
- JSON, compact, and binary codecs expose codec-local sentinels for codec encode/decode failures. Root compile entry points wrap codec failures in the root `*Error` surface with codec and operation context.

## Forbidden

- Do not add a second validation path parallel to compilation.
- Do not add runtime options for mutation; use `ApplyInPlace`.
- Do not duplicate the `test` operation contract here; `SPECS/25-predicate-specs.md` records predicate behavior.
- Do not model whole-object replacement with `extend`; use `replace` when the target value must be replaced atomically.

## Acceptance Criteria

- [ ] Every public entry point has one documented contract.
- [ ] The RFC 6902 mutating operations are defined once and only once.
- [ ] Compile-time validation and capability behavior are explicit.
- [ ] Callers can tell when to compile, apply immutably, apply in place, and use a codec.
