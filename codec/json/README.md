# JSON Codec for JSON Patch Operations

The `json` codec converts object-shaped JSON Patch operations to executable operations and back.

Use this package when you need direct control over JSON operation decoding or encoding. Most callers should use `jsonpatch.ApplyPatch`, which uses this codec internally.

## Wire Shape

JSON operations use JSON Pointer strings at the JSON boundary:

```json
{"op": "add", "path": "/profile/name", "value": "Ada"}
```

Raw JSON/map decoding owns field-presence checks. Missing required fields are rejected, while present `null`, `0`, and empty-string values remain real payload values.

## Usage

```go
package main

import (
    "fmt"
    "log"

    jsoncodec "github.com/kaptinlin/jsonpatch/codec/json"
)

func main() {
    raw := []map[string]any{
        {"op": "test", "path": "/name", "value": "Ada"},
        {"op": "replace", "path": "/name", "value": "Grace"},
    }

    ops, err := jsoncodec.Decode(raw, jsoncodec.PatchOptions{})
    if err != nil {
        log.Fatal(err)
    }

    encoded, err := jsoncodec.EncodeJSON(ops)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(len(encoded) > 0)
}
```

## API

```go
func Decode(operations []map[string]any, opts PatchOptions) ([]jsonpatch.Op, error)
func DecodeOperations(operations []jsonpatch.Operation, opts PatchOptions) ([]jsonpatch.Op, error)
func DecodeJSON(data []byte, opts PatchOptions) ([]jsonpatch.Op, error)

func Encode(ops []jsonpatch.Op) ([]jsonpatch.Operation, error)
func EncodeJSON(ops []jsonpatch.Op) ([]byte, error)
```

`PatchOptions` configures JSON decoding. Its main use is providing the matcher factory for `matches` predicates.

## Operation Families

The JSON codec decodes the same operation language used by the root package:

- RFC 6902 operations: `add`, `remove`, `replace`, `move`, `copy`, `test`
- Predicate operations: `defined`, `undefined`, `contains`, `starts`, `ends`, `matches`, `type`, `test_type`, `test_string`, `test_string_len`, `in`, `less`, `more`
- Composite predicates: `and`, `or`, unary `not`
- Extended operations: `flip`, `inc`, `str_ins`, `str_del`, `split`, `merge`, `extend`

Composite predicate child paths are decoded relative to the containing predicate path. Use `path: ""` on the containing predicate when children should be root-scoped.

## Testing Contract

The codec has fixture coverage for field presence: missing required fields, present `null`, root path, zero numeric fields, and empty strings.
