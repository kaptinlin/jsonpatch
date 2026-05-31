# Binary Codec for JSON Patch Operations

The **binary** codec serializes JSON Patch operations to MessagePack using the same operation tree as the compact codec.

The binary wire contract is protected by golden tests for MessagePack bytes, optional fields, and parent-relative predicate paths.

## Format Overview

The binary codec uses **MessagePack arrays** with the same structure as the compact codec:

```text
[opcode, path, ...args]
```

The payload is encoded as MessagePack bytes instead of JSON arrays.

## Usage

```go
import (
    "fmt"
    "log"

    "github.com/kaptinlin/jsonpatch"
    "github.com/kaptinlin/jsonpatch/codec/binary"
    "github.com/kaptinlin/jsonpatch/op"
)

func main() {
    ops := []jsonpatch.Op{
        op.NewAdd([]string{"user", "name"}, "Alice"),
        op.NewInc([]string{"user", "score"}, 100),
        op.NewStrIns([]string{"user", "bio"}, 0, "Hello! "),
        op.NewFlip([]string{"user", "active"}),
    }

    codec := binary.New()

    data, err := codec.Encode(ops)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Binary size: %d bytes\n", len(data))

    decoded, err := codec.Decode(data)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(len(decoded))
}
```

## Operation Support

The binary codec uses the same operation codes and path segment arrays as the compact codec:

### Standard JSON Patch Operations (RFC 6902)

| Operation | Code | Binary Format | Struct Example |
|-----------|------|---------------|----------------|
| **add**     | 0    | `[0, path_array, value]` | `{Op: "add", Path: "/foo", Value: 123}` |
| **remove**  | 1    | `[1, path_array]` | `{Op: "remove", Path: "/foo"}` |
| **replace** | 2    | `[2, path_array, value]` | `{Op: "replace", Path: "/foo", Value: 456}` |
| **move**    | 4    | `[4, path_array, from_array]` | `{Op: "move", Path: "/bar", From: "/foo"}` |
| **copy**    | 3    | `[3, path_array, from_array]` | `{Op: "copy", Path: "/bar", From: "/foo"}` |
| **test**    | 5    | `[5, path_array, value, not?]` | `{Op: "test", Path: "/foo", Value: 123}` |

### Extended Operations

| Operation | Code | Binary Format | Struct Example |
|-----------|------|---------------|----------------|
| **flip**     | 8    | `[8, path_array]` | `{Op: "flip", Path: "/active"}` |
| **inc**      | 9    | `[9, path_array, delta]` | `{Op: "inc", Path: "/count", Inc: 5}` |
| **str_ins**  | 6    | `[6, path_array, pos, str]` | `{Op: "str_ins", Path: "/text", Pos: 0, Str: "Hi"}` |
| **str_del**  | 7    | `[7, path_array, pos, len]` | `{Op: "str_del", Path: "/text", Pos: 5, Len: 3}` |
| **split**    | 10   | `[10, path_array, pos, props?]` | `{Op: "split", Path: "/obj", Pos: 2, Props: {...}}` |
| **extend**   | 12   | `[12, path_array, props, deleteNull?]` | `{Op: "extend", Path: "/config", Props: {...}}` |
| **merge**    | 11   | `[11, path_array, pos, props?]` | `{Op: "merge", Path: "/obj", Pos: 0, Props: {...}}` |

### JSON Predicate Operations

| Operation | Code | Binary Format | Struct Example |
|-----------|------|---------------|----------------|
| **defined**   | 31   | `[31, path_array]` | `{Op: "defined", Path: "/field"}` |
| **undefined** | 38   | `[38, path_array]` | `{Op: "undefined", Path: "/field"}` |
| **contains**  | 30   | `[30, path_array, value, ignoreCase?]` | `{Op: "contains", Path: "/text", Value: "hello"}` |
| **starts**    | 37   | `[37, path_array, value, ignoreCase?]` | `{Op: "starts", Path: "/text", Value: "Hello"}` |
| **ends**      | 32   | `[32, path_array, value, ignoreCase?]` | `{Op: "ends", Path: "/text", Value: ".com"}` |
| **matches**   | 35   | `[35, path_array, pattern, ignoreCase?]` | `{Op: "matches", Path: "/email", Value: ".*@.*"}` |
| **in**        | 33   | `[33, path_array, values]` | `{Op: "in", Path: "/role", Value: ["admin"]}` |
| **less**      | 34   | `[34, path_array, value]` | `{Op: "less", Path: "/age", Value: 30}` |
| **more**      | 36   | `[36, path_array, value]` | `{Op: "more", Path: "/score", Value: 100}` |
| **type**      | 42   | `[42, path_array, type]` | `{Op: "type", Path: "/data", Value: "string"}` |
| **test_type** | 39   | `[39, path_array, types]` | `{Op: "test_type", Path: "/data", Type: ["string"]}` |
| **test_string** | 40 | `[40, path_array, pos, str, not?]` | `{Op: "test_string", Path: "/text", Str: "test", Pos: 5}` |
| **test_string_len** | 41 | `[41, path_array, len, not?]` | `{Op: "test_string_len", Path: "/text", Len: 10}` |

### Second-Order Predicates

Child paths are encoded relative to the containing predicate path.

| Operation | Code | Binary Format | Struct Example |
|-----------|------|---------------|----------------|
| **and** | 43 | `[43, path_array, ops[]]` | `{Op: "and", Path: "/profile", Apply: [...]}` |
| **not** | 44 | `[44, path_array, ops[]]` | `{Op: "not", Path: "/profile", Apply: [...]}` |
| **or**  | 45 | `[45, path_array, ops[]]` | `{Op: "or", Path: "/profile", Apply: [...]}` |

## MessagePack Technical Details

### Path Encoding

Paths are encoded as **variable-length arrays**:

```text
[path_length, segment1, segment2, ...]
```

For example, path `"/user/name"` becomes:

```text
[2, "user", "name"]  // Binary MessagePack representation
```

### Value Encoding  

Values use MessagePack's native type preservation:

- **Strings**: Direct UTF-8 binary encoding
- **Numbers**: Compact binary number representation
- **Objects**: Recursive MessagePack object encoding
- **Arrays**: Recursive MessagePack array encoding

## API Reference

### Codec Structure

```go
type Codec struct{}

func New() *Codec
func (c *Codec) Encode(ops []jsonpatch.Op) ([]byte, error)
func (c *Codec) Decode(data []byte) ([]jsonpatch.Op, error)
```

## Testing Contract

The codec has golden coverage for MessagePack bytes, optional-field omission, and parent-relative composite predicate paths.
