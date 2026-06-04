# JSON Patch Go

[![Go Reference](https://pkg.go.dev/badge/github.com/kaptinlin/jsonpatch.svg)](https://pkg.go.dev/github.com/kaptinlin/jsonpatch)
[![Go Report Card](https://goreportcard.com/badge/github.com/kaptinlin/jsonpatch)](https://goreportcard.com/report/github.com/kaptinlin/jsonpatch)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A Go-native JSON Patch+ library with compiled patches, type-preserving results, and immutable-by-default application

## Features

- **Compiled patch path**: Compile operation vocabulary once, then reuse the resulting `Patch`.
- **Type-preserving results**: `Apply` converts patched documents back to the caller's original Go type when safe.
- **Immutable by default**: `Apply` clones before execution; `ApplyInPlace` makes mutation explicit.
- **Capability-gated vocabulary**: RFC 6902 is the default; predicates, regex predicates, and extended operations opt in at compile time.
- **Explicit document shapes**: Plain `string` is scalar text; use `JSONText` or `[]byte` for JSON text.
- **Structured errors**: Match stable failure classes with `errors.Is` and inspect operation context with `errors.As`.
- **Wire-format codecs**: Use JSON, compact array, or MessagePack codecs without moving execution semantics out of operations.

## Installation

```bash
go get github.com/kaptinlin/jsonpatch
```

Requires **Go 1.26.3+**.

## Quick Start

```go
package main

import (
    "fmt"
    "log"

    "github.com/kaptinlin/jsonpatch"
    "github.com/kaptinlin/jsonpatch/op"
)

func main() {
    doc := map[string]any{"name": "John", "tags": []any{"go"}}

    patch, err := jsonpatch.Compile(
        op.NewTest([]string{"name"}, "John"),
        op.NewReplace([]string{"name"}, "Jane"),
        op.NewAdd([]string{"email"}, "jane@example.com"),
    )
    if err != nil {
        log.Fatal(err)
    }

    result, err := jsonpatch.Apply(patch, doc)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(result.Doc["name"])
    fmt.Println(doc["name"])
}
```

## Core API

| API | Use when |
| --- | --- |
| `Compile` | You have Go-built operations from `op/` and want the default RFC 6902 vocabulary. |
| `CompileOps` | You have an operation slice or need compile options. |
| `CompileOperations` | You have JSON-shaped `codec/json.Operation` values. |
| `CompileJSON` | You have a JSON patch document as bytes. |
| `Apply` | You want immutable, type-preserving patch application. |
| `ApplyInPlace` | You intentionally want to write the patched result back to the input variable. |
| `JSONText` | You want a string document parsed as JSON text. |

## Capabilities

Default compilation accepts only RFC 6902 operations. Enable additional operation families explicitly.

```go
patch, err := jsonpatch.CompileJSON(data,
    jsonpatch.WithCapabilities(
        jsonpatch.RFC6902,
        jsonpatch.Predicate,
        jsonpatch.RegexPredicate,
        jsonpatch.Extended,
    ),
)
if err != nil {
    return err
}
```

Use `jsonpatch.AllCapabilities` when your boundary intentionally accepts every operation implemented by the package.

## Document Shapes

| Input | Processing model | Output |
| --- | --- | --- |
| `map[string]any` | Apply directly | `map[string]any` |
| `[]byte` | Decode JSON, apply, encode JSON | `[]byte` |
| `JSONText` | Decode JSON, apply, encode JSON | `JSONText` |
| `string` | Treat as scalar text | `string` |
| Structs and concrete types | Marshal to JSON, apply, unmarshal back | Original Go type |
| Primitives and `[]any` | Apply directly when assignable | Original Go type |

```go
patch, err := jsonpatch.CompileJSON([]byte(`[{"op":"replace","path":"/name","value":"Jane"}]`))
if err != nil {
    return err
}

result, err := jsonpatch.Apply(patch, jsonpatch.JSONText(`{"name":"John"}`))
if err != nil {
    return err
}

fmt.Println(result.Doc)
```

## In-Place Application

Use `ApplyInPlace` when mutation is intentional and visible at the call site.

```go
doc := map[string]any{"name": "John"}
patch, err := jsonpatch.Compile(op.NewReplace([]string{"name"}, "Jane"))
if err != nil {
    return err
}

if err := jsonpatch.ApplyInPlace(patch, &doc); err != nil {
    return err
}

fmt.Println(doc["name"])
```

## Structured Errors

Compile and apply failures wrap stable sentinel errors and expose operation context.

```go
result, err := jsonpatch.Apply(patch, doc)
if err != nil {
    if errors.Is(err, jsonpatch.ErrTestFailed) {
        return err
    }

    var patchErr *jsonpatch.Error
    if errors.As(err, &patchErr) {
        fmt.Printf("operation %d %s %s failed: %v\n",
            patchErr.Index(),
            patchErr.Op(),
            patchErr.Path(),
            patchErr.Cause(),
        )
    }
    return err
}

_ = result.Doc
```

## Codecs

Codec packages translate wire formats. Operation behavior stays in `op/`.

### JSON Codec

Use `codec/json` when you want to encode or decode JSON-shaped operation values directly.

```go
operations := []jsoncodec.Operation{
    {Op: "replace", Path: "/name", Value: "Jane"},
}

patch, err := jsonpatch.CompileOperations(operations)
if err != nil {
    return err
}
```

### Compact Codec

Use `codec/compact` for compact array-form operations with segment-array paths.

```go
ops := []jsonpatch.Op{
    op.NewAdd([]string{"name"}, "Jane"),
    op.NewInc([]string{"version"}, 1),
}

encoded, err := compact.EncodeJSON(ops)
if err != nil {
    return err
}

decoded, err := compact.DecodeJSON(encoded)
if err != nil {
    return err
}

fmt.Println(len(decoded))
```

### Binary Codec

Use `codec/binary` for MessagePack encoding.

```go
codec := binary.New()
ops := []jsonpatch.Op{
    op.NewAdd([]string{"name"}, "Jane"),
    op.NewInc([]string{"version"}, 1),
}

data, err := codec.Encode(ops)
if err != nil {
    return err
}

decoded, err := codec.Decode(data)
if err != nil {
    return err
}

fmt.Println(len(decoded))
```

## Examples

Explore runnable examples in [`examples/`](examples/):

- [`examples/basic-operations/`](examples/basic-operations/)
- [`examples/array-operations/`](examples/array-operations/)
- [`examples/conditional-operations/`](examples/conditional-operations/)
- [`examples/copy-move-operations/`](examples/copy-move-operations/)
- [`examples/string-operations/`](examples/string-operations/)
- [`examples/struct-patch/`](examples/struct-patch/)
- [`examples/map-patch/`](examples/map-patch/)
- [`examples/json-bytes-patch/`](examples/json-bytes-patch/)
- [`examples/json-string-patch/`](examples/json-string-patch/)
- [`examples/apply-in-place/`](examples/apply-in-place/)
- [`examples/compact-codec/`](examples/compact-codec/)
- [`examples/binary-codec/`](examples/binary-codec/)
- [`examples/error-handling/`](examples/error-handling/)
- [`examples/batch-update/`](examples/batch-update/)

## Development

```bash
task test           # Run all tests with race detection
task golangci-lint  # Run golangci-lint
task lint           # Run golangci-lint and tidy checks
task vet            # Run go vet
task bench          # Run benchmarks
```

For development guidelines, see [AGENTS.md](AGENTS.md).
For recorded behavior contracts, see [SPECS/](SPECS/).

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
