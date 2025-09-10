# Compact Codec for JSON Patch

The `compact` codec provides a space-efficient array-based encoding format for JSON Patch operations. It uses arrays instead of objects to represent operations, significantly reducing the physical space required while maintaining readability.

## Format Comparison

### Standard JSON format:
```json
{"op": "add", "path": "/foo/bar", "value": 123}
```

### Compact format:
```json
[0, "/foo/bar", 123]
```

Or with string opcodes:
```json
["add", "/foo/bar", 123]
```

## Usage

### Basic Operations

```go
package main

import (
    "fmt"
    "github.com/kaptinlin/jsonpatch/codec/compact"
    "github.com/kaptinlin/jsonpatch/op"
)

func main() {
    // Create operations
    ops := []internal.Op{
        op.NewAdd([]string{"foo"}, "bar"),
        op.NewReplace([]string{"baz"}, 42),
    }
    
    // Encode to compact format
    encoded, err := compact.Encode(ops)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Encoded: %v\n", encoded)
    // Output: [[0 "/foo" "bar"] [2 "/baz" 42]]
    
    // Decode back to operations
    decoded, err := compact.Decode(encoded)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Decoded: %d operations\n", len(decoded))
}
```

### String Opcodes

```go
// Use string opcodes instead of numeric ones
encoder := compact.NewEncoder(compact.WithStringOpcode(true))
encoded, err := encoder.Encode(ops)
// Result: [["add" "/foo" "bar"] ["replace" "/baz" 42]]
```

### JSON Marshaling

```go
// Encode to JSON bytes
jsonData, err := compact.EncodeJSON(ops)
if err != nil {
    panic(err)
}

// Decode from JSON bytes
decoded, err := compact.DecodeJSON(jsonData)
if err != nil {
    panic(err)
}
```

## Operation Mapping

### Standard JSON Patch Operations

| Operation | Numeric Code | String Code | Compact Format |
|-----------|--------------|-------------|----------------|
| add       | 0            | "add"       | `[0, path, value]` |
| remove    | 1            | "remove"    | `[1, path]` |
| replace   | 2            | "replace"   | `[2, path, value]` |
| move      | 3            | "move"      | `[3, path, from]` |
| copy      | 4            | "copy"      | `[4, path, from]` |
| test      | 5            | "test"      | `[5, path, value]` |

### Extended Operations

| Operation | Numeric Code | String Code | Compact Format |
|-----------|--------------|-------------|----------------|
| flip      | 10           | "flip"      | `[10, path]` |
| inc       | 11           | "inc"       | `[11, path, delta]` |
| merge     | 12           | "merge"     | `[12, path, value]` |
| extend    | 13           | "extend"    | `[13, path, value]` |

### Predicate Operations

| Operation | Numeric Code | String Code | Compact Format |
|-----------|--------------|-------------|----------------|
| defined   | 20           | "defined"   | `[20, path]` |
| undefined | 21           | "undefined" | `[21, path]` |
| contains  | 22           | "contains"  | `[22, path, value, ignoreCase?]` |
| starts    | 23           | "starts"    | `[23, path, value, ignoreCase?]` |
| ends      | 24           | "ends"      | `[24, path, value, ignoreCase?]` |

## API Reference

### Encoder

```go
type Encoder struct { ... }

// Create a new encoder
func NewEncoder(opts ...EncoderOption) *Encoder

// Encode a single operation
func (e *Encoder) Encode(op internal.Op) (CompactOp, error)

// Encode multiple operations
func (e *Encoder) EncodeSlice(ops []internal.Op) ([]CompactOp, error)
```

### Decoder

```go
type Decoder struct { ... }

// Create a new decoder
func NewDecoder(opts ...DecoderOption) *Decoder

// Decode a single compact operation
func (d *Decoder) Decode(compactOp CompactOp) (internal.Op, error)

// Decode multiple compact operations
func (d *Decoder) DecodeSlice(compactOps []CompactOp) ([]internal.Op, error)
```

### Standalone Functions

```go
// Encode operations using default options
func Encode(ops []internal.Op, opts ...EncoderOption) ([]CompactOp, error)

// Encode operations to JSON bytes
func EncodeJSON(ops []internal.Op, opts ...EncoderOption) ([]byte, error)

// Decode compact operations using default options
func Decode(compactOps []CompactOp, opts ...DecoderOption) ([]internal.Op, error)

// Decode compact operations from JSON bytes
func DecodeJSON(data []byte, opts ...DecoderOption) ([]internal.Op, error)
```

### Options

```go
// Use string opcodes instead of numeric codes
func WithStringOpcode(useString bool) EncoderOption
```

## Features

- **Space Efficient**: Significantly smaller than standard JSON format
- **Fast**: Optimized encoding and decoding performance
- **Flexible**: Supports both numeric and string opcodes
- **Compatible**: Works with all existing operation types
- **Type Safe**: Full Go type safety and error handling

## Supported Operations

Currently supports all standard JSON Patch operations and basic extended/predicate operations:

✅ **Standard**: add, remove, replace, move, copy, test  
✅ **Extended**: flip, inc  
✅ **Predicates**: defined, undefined, contains, starts, ends  
⏳ **Coming Soon**: merge, extend, string operations, composite predicates 
