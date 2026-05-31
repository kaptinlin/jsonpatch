# JSON Patch Go

[![Go Reference](https://pkg.go.dev/badge/github.com/kaptinlin/jsonpatch.svg)](https://pkg.go.dev/github.com/kaptinlin/jsonpatch)
[![Go Report Card](https://goreportcard.com/badge/github.com/kaptinlin/jsonpatch)](https://goreportcard.com/report/github.com/kaptinlin/jsonpatch)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A Go-native JSON Patch+ library with type-preserving results and immutable-by-default patch application

Requires Go 1.26+.

## Features

- **Type-preserving results**: `ApplyPatch`, `ApplyOp`, and `ApplyOps` preserve the input shape whenever the patched value can be converted back safely.
- **Immutable by default**: Patch application clones first unless you pass `jsonpatch.WithMutate(true)`.
- **Three entry points**: Use `ApplyPatch` for JSON-shaped operations, `ApplyOp` for one executable operation, and `ApplyOps` for executable operation slices.
- **Patch vocabulary**: Use RFC 6902 operations, predicate operations, and extended operations in one package.
- **Multiple document shapes**: Patch `map[string]any`, structs, `[]byte`, JSON strings, plain strings, primitives, and `[]any`.
- **Codec control**: Use `codec/json`, `codec/compact`, and `codec/binary` when you need a specific wire format.

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
| `codec/json` | You need explicit JSON operation encoding or decoding. |
| `codec/compact` | You need compact array-form operations with segment-array paths. |
| `codec/binary` | You need MessagePack encoding for executable operations. |
| `ValidateOperation` / `ValidateOperations` | You want to preflight Go-built `Operation` values before applying them. |

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
Binary supports the same operation tree as the compact codec, including second-order predicates (`and`, `or`, `not`).

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
task golangci-lint
task lint
```

For development guidelines, see [AGENTS.md](AGENTS.md).
For recorded behavior contracts, see [SPECS/](SPECS/).

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
