# JSON Patch Go

[![Go Reference](https://pkg.go.dev/badge/github.com/kaptinlin/jsonpatch.svg)](https://pkg.go.dev/github.com/kaptinlin/jsonpatch)
[![Go Report Card](https://goreportcard.com/badge/github.com/kaptinlin/jsonpatch)](https://goreportcard.com/report/github.com/kaptinlin/jsonpatch)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A Go library for RFC 6902 JSON Patch plus predicate and extended operations with type-preserving APIs

> json-joy compatible: this package follows the JSON Patch, predicate, and extended-operation behavior used by [streamich/json-joy](https://github.com/streamich/json-joy/tree/master/src/json-patch).

## Features

- **Type-preserving results**: `ApplyPatch`, `ApplyOp`, and `ApplyOps` preserve the input document shape whenever the result can be converted back safely.
- **Broad operation support**: Use RFC 6902 operations, predicate operations, and extended operations in one package.
- **Multiple document shapes**: Patch `map[string]any`, structs, `[]byte`, JSON strings, plain strings, primitives, and `[]any`.
- **Executable operations**: Build operations with the `op` package when you want typed operation values instead of JSON-shaped payloads.
- **Codec support**: Encode and decode operations through `codec/json`, `codec/compact`, and `codec/binary`.
- **Predictable defaults**: Patch application is immutable unless you pass `jsonpatch.WithMutate(true)`.

## Installation

```bash
go get github.com/kaptinlin/jsonpatch
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"

    "github.com/kaptinlin/jsonpatch"
)

func main() {
    doc := map[string]any{
        "name": "John",
        "tags": []any{"golang"},
    }

    patch := []jsonpatch.Operation{
        {Op: "test", Path: "/name", Value: "John"},
        {Op: "replace", Path: "/name", Value: "Jane"},
        {Op: "add", Path: "/email", Value: "jane@example.com"},
    }

    result, err := jsonpatch.ApplyPatch(doc, patch)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(result.Doc["name"])
    fmt.Println(result.Doc["email"])
    fmt.Println(doc["name"])
}
```

## API Overview

| API | Use when |
| --- | --- |
| `ApplyPatch` | You already have JSON-shaped `[]Operation` values. |
| `ApplyOp` | You want to execute one compiled operation from `op/`. |
| `ApplyOps` | You want to execute compiled operations directly. |
| `ValidateOperation` / `ValidateOperations` | You want to validate JSON-shaped operation payloads before applying them. |

## Executable Operations

Use `op/` when you want operation values with methods such as `Validate`, `ToJSON`, and `ToCompact`.

```go
package main

import (
    "fmt"
    "log"

    "github.com/kaptinlin/jsonpatch"
    "github.com/kaptinlin/jsonpatch/op"
)

type User struct {
    Name   string   `json:"name"`
    Active bool     `json:"active"`
    Roles  []string `json:"roles"`
}

func main() {
    user := User{Name: "John", Active: true, Roles: []string{"admin"}}

    ops := []jsonpatch.Op{
        op.NewTest([]string{"active"}, true),
        op.NewReplace([]string{"name"}, "Jane"),
        op.NewAdd([]string{"roles", "-"}, "owner"),
    }

    result, err := jsonpatch.ApplyOps(user, ops)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(result.Doc.Name)
    fmt.Println(result.Doc.Roles)
}
```

## Document Shapes

| Input | Processing model | Output |
| --- | --- | --- |
| `map[string]any` | Apply directly | `map[string]any` |
| `[]byte` | Decode JSON, apply, encode JSON | `[]byte` |
| `string` starting with `{` or `[` | Decode JSON, apply, encode JSON | `string` |
| Other `string` | Treat as a plain string value | `string` |
| Structs and concrete types | Marshal to JSON, apply, unmarshal back | Original Go type |
| Primitives and `[]any` | Apply directly when assignable | Original Go type |

## Codecs

### Compact Codec

Use `codec/compact` when you want the array-based wire format.

```go
ops := []jsonpatch.Op{
    op.NewAdd([]string{"name"}, "Jane"),
    op.NewInc([]string{"version"}, 1),
}

encoded, err := compact.EncodeJSON(ops)
if err != nil {
    log.Fatal(err)
}

decoded, err := compact.DecodeJSON(encoded)
if err != nil {
    log.Fatal(err)
}

fmt.Println(len(decoded))
```

### Binary Codec

Use `codec/binary` when you want MessagePack encoding for executable operations.
Second-order predicates (`and`, `or`, `not`) are not supported by the binary codec.

```go
codec := binary.New()

data, err := codec.Encode(ops)
if err != nil {
    log.Fatal(err)
}

decoded, err := codec.Decode(data)
if err != nil {
    log.Fatal(err)
}

fmt.Println(len(decoded))
```

## Examples

Explore the runnable examples in [`examples/`](examples/):

- [`examples/basic-operations/`](examples/basic-operations/)
- [`examples/array-operations/`](examples/array-operations/)
- [`examples/conditional-operations/`](examples/conditional-operations/)
- [`examples/copy-move-operations/`](examples/copy-move-operations/)
- [`examples/string-operations/`](examples/string-operations/)
- [`examples/struct-patch/`](examples/struct-patch/)
- [`examples/map-patch/`](examples/map-patch/)
- [`examples/json-bytes-patch/`](examples/json-bytes-patch/)
- [`examples/json-string-patch/`](examples/json-string-patch/)
- [`examples/compact-codec/`](examples/compact-codec/)
- [`examples/binary-codec/`](examples/binary-codec/)
- [`examples/error-handling/`](examples/error-handling/)
- [`examples/mutate-option/`](examples/mutate-option/)
- [`examples/batch-update/`](examples/batch-update/)

## Development

```bash
task test
task lint
task markdownlint
```

For development guidelines, see [AGENTS.md](AGENTS.md).
For technical contracts, see [SPECS/](SPECS/).

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
