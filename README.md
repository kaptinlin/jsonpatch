# JSON Patch Go

A comprehensive Go implementation of JSON Patch (RFC 6902), JSON Predicate, and extended operations for JSON document manipulation with **full type safety** and **generic support**.

> **json-joy Compatible**: This is a Go port of [json-joy/json-patch](https://github.com/streamich/json-joy/tree/master/src/json-patch) with 95%+ behavioral compatibility, bringing all JSON Patch extended operations to the Go ecosystem with modern Go generics.

[![Go Reference](https://pkg.go.dev/badge/github.com/kaptinlin/jsonpatch.svg)](https://pkg.go.dev/github.com/kaptinlin/jsonpatch)
[![Go Report Card](https://goreportcard.com/badge/github.com/kaptinlin/jsonpatch)](https://goreportcard.com/report/github.com/kaptinlin/jsonpatch)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Quick Start

### Installation

```bash
go get github.com/kaptinlin/jsonpatch
```

### Basic Usage

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"

    "github.com/kaptinlin/jsonpatch"
)

func main() {
    // Original document
    doc := map[string]any{
        "name": "John",
        "age":  30,
    }

    // Create patch operations using struct syntax
    patch := []jsonpatch.Operation{
        {Op: "replace", Path: "/name", Value: "Jane"},
        {Op: "add", Path: "/email", Value: "jane@example.com"},
    }

    // Apply patch with type-safe generic API
    result, err := jsonpatch.ApplyPatch(doc, patch)
    if err != nil {
        log.Fatalf("Failed to apply patch: %v", err)
    }

    // result.Doc is automatically typed as map[string]any
    fmt.Printf("Name: %s\n", result.Doc["name"])
    fmt.Printf("Email: %s\n", result.Doc["email"])

    output, _ := json.MarshalIndent(result.Doc, "", "  ")
    fmt.Println(string(output))
    // Output:
    // {
    //   "age": 30,
    //   "email": "jane@example.com",
    //   "name": "Jane"
    // }
}
```

### Type-Safe Struct Usage

```go
type User struct {
    Name  string `json:"name"`
    Email string `json:"email,omitempty"`
    Age   int    `json:"age"`
}

func main() {
    user := User{Name: "John", Age: 30}

    patch := []jsonpatch.Operation{
        {Op: "replace", Path: "/name", Value: "Jane"},
        {Op: "add", Path: "/email", Value: "jane@example.com"},
    }

    // Type-safe: result.Doc is automatically typed as User
    result, err := jsonpatch.ApplyPatch(user, patch)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Updated user: %+v\n", result.Doc)
    // Output: Updated user: {Name:Jane Email:jane@example.com Age:30}
}
```

## Features

### Type-Safe Generic API

- **Full Generic Support** - No `interface{}` or type assertions needed
- **Compile-Time Type Safety** - Catch type errors at compile time
- **Automatic Type Inference** - Result types match input types
- **Multiple Document Types** - `map[string]any`, structs, `[]byte`, `string`

### RFC 6902 Standard Operations ([docs](docs/json-patch.md))

| Operation | Description |
|-----------|-------------|
| `add` | Add new values to objects or arrays |
| `remove` | Remove existing values |
| `replace` | Replace existing values |
| `move` | Move values to different locations |
| `copy` | Copy values to new locations |
| `test` | Test values for conditional operations |

### JSON Predicate Operations ([docs](docs/json-predicate.md))

| Operation | Description |
|-----------|-------------|
| `contains` | Check string/array containment |
| `defined` / `undefined` | Test path existence |
| `starts` / `ends` | Test string prefix/suffix |
| `in` | Check membership in arrays |
| `less` / `more` | Numeric comparisons |
| `matches` | Regular expression matching |
| `type` / `test_type` | Type validation |
| `test_string` | Position-based string testing |
| `test_string_len` | String length validation |
| `and` / `or` / `not` | Logical predicate combinations |

### Extended Operations ([docs](docs/json-patch-extended.md))

| Operation | Description |
|-----------|-------------|
| `str_ins` | String insertion at position |
| `str_del` | String deletion by position/substring |
| `inc` | Increment/decrement numeric values |
| `flip` | Toggle boolean values |
| `split` | Split values at position |
| `merge` | Merge adjacent array elements |
| `extend` | Extend objects with properties |

### Codecs

Three codec formats for different use cases:

| Codec | Package | Description |
|-------|---------|-------------|
| **JSON** | `codec/json` | Standard RFC 6902 JSON format |
| **Compact** | `codec/compact` | Array-based format with ~35% space savings |
| **Binary** | `codec/binary` | MessagePack binary format for maximum efficiency |

## json-joy Compatibility

This implementation provides **95%+ behavioral compatibility** with the TypeScript [json-joy/json-patch](https://github.com/streamich/json-joy/tree/master/src/json-patch) reference implementation.

### Predicate Negation Pattern

```go
// Direct negation (only test, test_string, test_string_len)
{Op: "test", Path: "/value", Value: 42, Not: true}
{Op: "test_string", Path: "/name", Pos: 0, Str: "test", Not: true}
{Op: "test_string_len", Path: "/name", Len: 5, Not: true}

// Second-order predicate negation (for all other predicates)
{
    Op:   "not",
    Path: "",
    Apply: []jsonpatch.Operation{
        {Op: "starts", Path: "/name", Value: "John"},
    },
}
```

### Complex Predicate Logic

```go
// Logical AND - all conditions must pass
{
    Op:   "and",
    Path: "",
    Apply: []jsonpatch.Operation{
        {Op: "starts", Path: "/name", Value: "John"},
        {Op: "ends", Path: "/name", Value: "Doe"},
    },
}

// Logical OR - any condition can pass
{
    Op:   "or",
    Path: "",
    Apply: []jsonpatch.Operation{
        {Op: "contains", Path: "/email", Value: "@gmail.com"},
        {Op: "contains", Path: "/email", Value: "@yahoo.com"},
    },
}
```

## API Reference

### Core Functions

```go
// Apply a JSON Patch document to any supported document type
func ApplyPatch[T Document](doc T, patch []Operation, opts ...Option) (*PatchResult[T], error)

// Apply a single operation
func ApplyOp[T Document](doc T, operation Op, opts ...Option) (*OpResult[T], error)

// Apply multiple operations
func ApplyOps[T Document](doc T, operations []Op, opts ...Option) (*PatchResult[T], error)
```

### Validation Functions

```go
// Validate an array of operations
func ValidateOperations(ops []Operation, allowMatchesOp bool) error

// Validate a single operation
func ValidateOperation(operation Operation, allowMatchesOp bool) error
```

### Functional Options

```go
// Configure mutation behavior (default: false, creates deep copy)
jsonpatch.WithMutate(true)

// Configure custom regex matcher for "matches" operations
jsonpatch.WithMatcher(func(pattern string, ignoreCase bool) jsonpatch.RegexMatcher {
    // return custom matcher
})
```

### Result Types

```go
// Result of a single operation
type OpResult[T Document] struct {
    Doc T   // Result document with preserved type
    Old any // Previous value at the path
}

// Result of applying a patch
type PatchResult[T Document] struct {
    Doc T             // Result document with preserved type
    Res []OpResult[T] // Results for each operation
}
```

## Codec Usage

### Compact Codec

```go
import "github.com/kaptinlin/jsonpatch/codec/compact"

// Standard operations
ops := []jsonpatch.Operation{
    {Op: "add", Path: "/name", Value: "John"},
    {Op: "replace", Path: "/age", Value: 30},
    {Op: "remove", Path: "/temp"},
}

// Encode to compact format with numeric opcodes (~35% space savings)
encoder := compact.NewEncoder(compact.WithStringOpcode(false))
compactData, err := encoder.EncodeJSON(ops)
// Result: [[0,"/name","John"],[2,"/age",30],[1,"/temp"]]

// Decode back to operations
decoder := compact.NewDecoder()
decoded, err := decoder.DecodeJSON(compactData)
```

### Binary Codec

```go
import "github.com/kaptinlin/jsonpatch/codec/binary"

ops := []jsonpatch.Operation{
    {Op: "add", Path: "/user/name", Value: "Alice"},
    {Op: "inc", Path: "/user/score", Inc: 100},
}

codec := &binary.Codec{}

// Encode to MessagePack binary (40-60% smaller than JSON)
data, err := codec.Encode(ops)

// Decode back
decoded, err := codec.Decode(data)
```

## Common Patterns

### Type-Safe Operations

```go
type Config struct {
    Version int    `json:"version"`
    Status  string `json:"status"`
    Enabled bool   `json:"enabled"`
}

config := Config{Version: 1, Status: "active", Enabled: true}

patch := []jsonpatch.Operation{
    {Op: "inc", Path: "/version", Inc: 1},
    {Op: "replace", Path: "/status", Value: "updated"},
    {Op: "flip", Path: "/enabled"},
}

// result.Doc is automatically typed as Config
result, err := jsonpatch.ApplyPatch(config, patch)
```

### Safe Updates with Test Operations

```go
patch := []jsonpatch.Operation{
    {Op: "test", Path: "/version", Value: 1},
    {Op: "replace", Path: "/status", Value: "updated"},
    {Op: "inc", Path: "/version", Inc: 1},
}

result, err := jsonpatch.ApplyPatch(doc, patch)
```

### Mutation Control

```go
// Preserve original document (default)
result, err := jsonpatch.ApplyPatch(doc, patch)

// Mutate original document for performance
result, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.WithMutate(true))
```

### Batch Operations

```go
var patch []jsonpatch.Operation

for i := range itemCount {
    patch = append(patch, jsonpatch.Operation{
        Op:    "replace",
        Path:  fmt.Sprintf("/items/%d/status", i),
        Value: "processed",
    })
}

result, err := jsonpatch.ApplyPatch(doc, patch)
```

### Array Manipulation

```go
patch := []jsonpatch.Operation{
    {Op: "add", Path: "/users/-", Value: map[string]any{"name": "New User"}},
    {Op: "add", Path: "/tags/0", Value: "important"},
    {Op: "remove", Path: "/items/2"},
}

result, err := jsonpatch.ApplyPatch(doc, patch)
```

### String Operations

```go
patch := []jsonpatch.Operation{
    {Op: "str_ins", Path: "/content", Pos: 0, Str: "Prefix: "},
    {Op: "str_del", Path: "/content", Pos: 6, Len: 10},
}

result, err := jsonpatch.ApplyPatch(doc, patch)
```

### Error Handling

```go
import "errors"

result, err := jsonpatch.ApplyPatch(doc, patch)
if err != nil {
    log.Printf("Patch application failed: %v", err)
    return err
}
```

### Compact Codec for Storage/Network

```go
import "github.com/kaptinlin/jsonpatch/codec/compact"

func storeOperations(ops []jsonpatch.Operation) error {
    encoder := compact.NewEncoder(compact.WithStringOpcode(false))
    compactData, err := encoder.EncodeJSON(ops)
    if err != nil {
        return err
    }
    return database.Store(compactData)
}
```

## Examples

Explore comprehensive examples in the [`examples/`](examples/) directory (see [`examples/README.md`](examples/README.md) for the complete guide):

**Core Operations:**
[Basic Operations](examples/basic-operations/) |
[Array Operations](examples/array-operations/) |
[Conditional Operations](examples/conditional-operations/) |
[Copy & Move](examples/copy-move-operations/) |
[String Operations](examples/string-operations/) |
[Extended Operations](examples/extended/)

**Document Types:**
[Struct Patch](examples/struct-patch/) |
[Map Patch](examples/map-patch/) |
[JSON Bytes](examples/json-bytes-patch/) |
[JSON String](examples/json-string-patch/)

**Codecs:**
[Compact Codec](examples/compact-codec/) |
[Binary Codec](examples/binary-codec/)

**Advanced:**
[Batch Updates](examples/batch-update/) |
[Error Handling](examples/error-handling/) |
[Mutate Option](examples/mutate-option/)

```bash
# Run any example
cd examples/<example-name> && go run main.go
```

## Related Specifications

- [RFC 6902 - JSON Patch](https://tools.ietf.org/html/rfc6902)
- [RFC 6901 - JSON Pointer](https://tools.ietf.org/html/rfc6901)
- [JSON Predicate Draft](https://tools.ietf.org/id/draft-snell-json-test-01.html)

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## Credits

This project is a Go port of [json-joy/json-patch](https://github.com/streamich/json-joy/tree/master/src/json-patch). Thanks to the original authors for their excellent work.

Original project: [streamich/json-joy](https://github.com/streamich/json-joy)

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) file for details.
