# JSON Patch Go

A comprehensive Go implementation of JSON Patch (RFC 6902), JSON Predicate, and extended operations for JSON document manipulation with **full type safety** and **generic support**.

> **json-joy Compatible**: This is a Go port of [json-joy/json-patch](https://github.com/streamich/json-joy/tree/master/src/json-patch) with 95%+ behavioral compatibility, bringing all JSON Patch extended operations to the Go ecosystem with modern Go generics.

[![Go Reference](https://pkg.go.dev/badge/github.com/kaptinlin/jsonpatch.svg)](https://pkg.go.dev/github.com/kaptinlin/jsonpatch)
[![Go Report Card](https://goreportcard.com/badge/github.com/kaptinlin/jsonpatch)](https://goreportcard.com/report/github.com/kaptinlin/jsonpatch)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## 🚀 Quick Start

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
    doc := map[string]interface{}{
        "name": "John",
        "age":  30,
    }

    // Create patch operations using struct syntax
    patch := []jsonpatch.Operation{
        {
            Op:    "replace",
            Path:  "/name",
            Value: "Jane",
        },
        {
            Op:    "add",
            Path:  "/email",
            Value: "jane@example.com",
        },
    }

    // Apply patch with type-safe generic API
    result, err := jsonpatch.ApplyPatch(doc, patch)
    if err != nil {
        log.Fatalf("Failed to apply patch: %v", err)
    }

    // result.Doc is automatically typed as map[string]interface{}
    // No type assertions needed!
    fmt.Printf("Name: %s\n", result.Doc["name"])
    fmt.Printf("Email: %s\n", result.Doc["email"])

    // Output result as JSON
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

### Type-Safe Generic Usage

```go
// Define your own types for complete type safety
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

    // No type assertions needed - compile-time type safety!
    fmt.Printf("Updated user: %+v\n", result.Doc)
    // Output: Updated user: {Name:Jane Email:jane@example.com Age:30}
}
```

### Compact Codec Usage

```go
import "github.com/kaptinlin/jsonpatch/codec/compact"

func main() {
    // Standard JSON Patch operations
    standardOps := []jsonpatch.Operation{
        {Op: "add", Path: "/name", Value: "John"},
        {Op: "replace", Path: "/age", Value: 30},
        {Op: "remove", Path: "/temp"},
    }

    // Encode to compact format with numeric opcodes (35.9% space savings)
    encoder := compact.NewEncoder(compact.WithNumericOpcodes(true))
    compactData, err := encoder.Encode(standardOps)
    if err != nil {
        log.Fatal(err)
    }

    // Compact format: [[0,"/name","John"],[2,"/age",30],[1,"/temp"]]
    // vs Standard: [{"op":"add","path":"/name","value":"John"},...]

    // Decode back to standard operations
    decoder := compact.NewDecoder()
    decodedOps, err := decoder.Decode(compactData)
    if err != nil {
        log.Fatal(err)
    }

    // Perfect round-trip compatibility
    fmt.Printf("Original ops count: %d\n", len(standardOps))
    fmt.Printf("Decoded ops count: %d\n", len(decodedOps))
    // Output: Original ops count: 3
    //         Decoded ops count: 3
}
```

## 🎪 json-joy Compatibility

This implementation provides **95%+ behavioral compatibility** with the TypeScript [json-joy/json-patch](https://github.com/streamich/json-joy/tree/master/src/json-patch) reference implementation:

### ✅ **Predicate Negation Pattern**
```go
// Direct negation (only for specific operations)
{Op: "test", Path: "/value", Value: 42, Not: true}
{Op: "test_string", Path: "/name", Pos: 0, Str: "test", Not: true}
{Op: "test_string_len", Path: "/name", Len: 5, Not: true}

// Second-order predicate negation (for all other predicates)
{
    Op: "not",
    Path: "",
    Apply: []Operation{
        {Op: "starts", Path: "/name", Value: "John"},
    },
}
```

### ✅ **Complex Predicate Logic**
```go
// Logical AND - all conditions must pass
{
    Op: "and", 
    Path: "",
    Apply: []Operation{
        {Op: "starts", Path: "/name", Value: "John"},
        {Op: "ends", Path: "/name", Value: "Doe"},
    },
}

// Logical OR - any condition can pass
{
    Op: "or",
    Path: "",
    Apply: []Operation{
        {Op: "contains", Path: "/email", Value: "@gmail.com"},
        {Op: "contains", Path: "/email", Value: "@yahoo.com"},
    },
}
```

## 🎯 Features

### ✨ **Type-Safe Generic API**

- **Full Generic Support** - No more `interface{}` or type assertions
- **Compile-Time Type Safety** - Catch type errors at compile time
- **Automatic Type Inference** - Result types are automatically inferred
- **Zero-Value Usability** - Works without configuration

### ✅ RFC 6902 Standard Operations ([docs](docs/json-patch.md))

- **`add`** - Add new values to the document
- **`remove`** - Remove existing values
- **`replace`** - Replace existing values
- **`move`** - Move values to different locations
- **`copy`** - Copy values to new locations
- **`test`** - Test values for conditional operations

### 🔍 JSON Predicate Operations ([docs](docs/json-predicate.md))

- **`contains`** - Check if strings contain substrings
- **`defined`** - Test if paths exist
- **`ends`** - Test string endings
- **`in`** - Check membership in arrays
- **`matches`** - Regular expression matching
- **`starts`** - Test string beginnings
- **`type`** - Type validation
- **`less`/`more`** - Numeric comparisons
- **`test_string`** - Position-based string testing with negation
- **`test_string_len`** - String length validation with negation
- **`and`/`or`/`not`** - Second-order logical predicate combinations

### 🔧 Extended Operations ([docs](docs/json-patch-extended.md))

- **`str_ins`** - String insertion for text editing
- **`str_del`** - String deletion
- **`inc`** - Increment numeric values
- **`flip`** - Boolean toggling
- **`split`** - Object splitting
- **`merge`** - Object merging

### 🗜️ **Compact Codec** ([codec/compact/](codec/compact/))

- **Space Efficient** - 35.9% size reduction with numeric opcodes
- **Array Format** - `[opcode, path, ...args]` instead of verbose objects
- **Both Formats** - Support for numeric (0,1,2...) and string ("add","remove"...) opcodes
- **Round-trip Compatible** - Perfect conversion between standard and compact formats
- **High Performance** - O(1) opcode lookups and optimized encoding/decoding

## 📖 Examples

Explore comprehensive examples in the [`examples/`](examples/) directory (see [`examples/README.md`](examples/README.md) for complete guide):

### 🏗️ **Core Operations**
- **[Basic Operations](examples/basic-operations/)** - `add`, `replace`, `remove`, `test`
- **[Array Operations](examples/array-operations/)** - Array manipulation
- **[Conditional Operations](examples/conditional-operations/)** - Safe updates with tests
- **[Copy & Move](examples/copy-move-operations/)** - Data restructuring
- **[String Operations](examples/string-operations/)** - Text editing

### 📄 **Document Types**
- **[Struct Patch](examples/struct-patch/)** - Type-safe Go structs
- **[Map Patch](examples/map-patch/)** - Dynamic `map[string]any`
- **[JSON Bytes](examples/json-bytes-patch/)** - Raw JSON byte data
- **[JSON String](examples/json-string-patch/)** - JSON string data

### 🗜️ **Codecs**
- **[Compact Codec](examples/compact-codec/)** - Space-efficient array format

### 🚀 **Advanced**
- **[Batch Updates](examples/batch-update/)** - Bulk operations
- **[Error Handling](examples/error-handling/)** - Error patterns
- **[Mutate Option](examples/mutate-option/)** - Performance optimization

```bash
# Run any example
cd examples/<example-name> && go run main.go

# Try the struct patching demo
cd examples/struct-patch && go run main.go
```

## 📚 API Reference

### Core Generic Functions

```go
// Apply a single operation with type safety
func ApplyOp[T Document](doc T, operation Op, opts ...Option) (*OpResult[T], error)

// Apply multiple operations with type safety
func ApplyOps[T Document](doc T, operations []Op, opts ...Option) (*PatchResult[T], error)

// Apply JSON Patch with type safety
func ApplyPatch[T Document](doc T, patch []Operation, opts ...Option) (*PatchResult[T], error)
```

### Compact Codec Functions

```go
import "github.com/kaptinlin/jsonpatch/codec/compact"

// Create encoder with options
func NewEncoder(opts ...EncoderOption) *Encoder
func WithNumericOpcodes(numeric bool) EncoderOption

// Create decoder with options  
func NewDecoder(opts ...DecoderOption) *Decoder

// Encode operations to compact format
func (e *Encoder) Encode(ops []jsonpatch.Operation) ([]byte, error)

// Decode from compact format to operations
func (d *Decoder) Decode(data []byte) ([]jsonpatch.Operation, error)

// Convenience functions
func Encode(ops []jsonpatch.Operation, opts ...EncoderOption) ([]byte, error)
func Decode(data []byte, opts ...DecoderOption) ([]jsonpatch.Operation, error)
```

### Functional Options

```go
// Configure mutation behavior
func WithMutate(mutate bool) Option

// Configure custom regex matcher
func WithMatcher(createMatcher func(string, bool) RegexMatcher) Option
```

### Result Types

```go
// Generic result for single operation
type OpResult[T Document] struct {
    Doc T    // Result document with preserved type
    Old any  // Previous value at the path
}

// Generic result for multiple operations
type PatchResult[T Document] struct {
    Doc T                // Result document with preserved type
    Res []OpResult[any]  // Results for each operation
}
```

## 🔨 Common Patterns

### 1. Type-Safe Operations

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
if err != nil {
    log.Fatal(err)
}

// No type assertions needed!
fmt.Printf("Version: %d, Status: %s, Enabled: %t\n", 
    result.Doc.Version, result.Doc.Status, result.Doc.Enabled)
```

### 2. Safe Updates with Test Operations

```go
patch := []jsonpatch.Operation{
    // Test current value before modifying
    {
        Op:    "test",
        Path:  "/version",
        Value: 1,
    },
    // Make the change
    {
        Op:    "replace",
        Path:  "/status",
        Value: "updated",
    },
    // Increment version
    {
        Op:   "inc",
        Path: "/version",
        Inc:  1,
    },
}

result, err := jsonpatch.ApplyPatch(doc, patch)
```

### 3. Mutation Control

```go
// Preserve original document (default)
result, err := jsonpatch.ApplyPatch(doc, patch)

// Mutate original document for performance
result, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.WithMutate(true))
```

### 4. Batch Operations

```go
var patch []jsonpatch.Operation

// Update multiple items efficiently
for i := 0; i < itemCount; i++ {
    patch = append(patch, jsonpatch.Operation{
        Op:    "replace",
        Path:  fmt.Sprintf("/items/%d/status", i),
        Value: "processed",
    })
}

result, err := jsonpatch.ApplyPatch(doc, patch)
```

### 5. Array Manipulation

```go
patch := []jsonpatch.Operation{
    // Add to end of array
    {
        Op:   "add",
        Path: "/users/-",
        Value: map[string]interface{}{"name": "New User"},
    },
    // Insert at specific position
    {
        Op:   "add",
        Path: "/tags/0",
        Value: "important",
    },
    // Remove array element
    {
        Op:   "remove",
        Path: "/items/2",
    },
}

result, err := jsonpatch.ApplyPatch(doc, patch)
```

### 6. String Operations

```go
patch := []jsonpatch.Operation{
    // Insert text at position
    {
        Op:   "str_ins",
        Path: "/content",
        Pos:  0,
        Str:  "Prefix: ",
    },
    // Insert at end
    {
        Op:   "str_ins",
        Path: "/content",
        Pos:  20,
        Str:  " (Updated)",
    },
}

result, err := jsonpatch.ApplyPatch(doc, patch)
```

### 8. Error Handling

```go
result, err := jsonpatch.ApplyPatch(doc, patch)
if err != nil {
    switch {
    case strings.Contains(err.Error(), "path not found"):
        log.Printf("Invalid path in patch: %v", err)
    case strings.Contains(err.Error(), "test operation failed"):
        log.Printf("Condition not met: %v", err)
    default:
        log.Printf("Patch application failed: %v", err)
    }
    return err
}

// Type-safe access to result
processedDoc := result.Doc
```

### 9. Compact Codec Usage

```go
import "github.com/kaptinlin/jsonpatch/codec/compact"

// Standard operations
ops := []jsonpatch.Operation{
    {Op: "add", Path: "/users/-", Value: map[string]string{"name": "Alice"}},
    {Op: "inc", Path: "/counter", Inc: 1},
    {Op: "flip", Path: "/enabled"},
}

// Choose encoding format
numericEncoder := compact.NewEncoder(compact.WithNumericOpcodes(true))
stringEncoder := compact.NewEncoder(compact.WithNumericOpcodes(false))

// Encode (numeric opcodes for maximum space savings)
compactData, err := numericEncoder.Encode(ops)
if err != nil {
    log.Fatal(err)
}

// Decode
decoder := compact.NewDecoder()
decoded, err := decoder.Decode(compactData)
if err != nil {
    log.Fatal(err)
}

// Perfect round-trip compatibility
fmt.Printf("Space savings: %d%%\n", (len(standardJSON)-len(compactData))*100/len(standardJSON))
// Output: Space savings: 35%
```

## 📈 Best Practices

### 1. Leverage Type Safety

```go
// Define specific types for your data
type UserProfile struct {
    ID       string   `json:"id"`
    Name     string   `json:"name"`
    Email    string   `json:"email"`
    Tags     []string `json:"tags"`
    Settings map[string]interface{} `json:"settings"`
}

// Get compile-time type safety
result, err := jsonpatch.ApplyPatch(userProfile, patch)
// result.Doc is automatically typed as UserProfile
```

### 2. Use Functional Options

```go
// Default behavior (no mutation)
result, err := jsonpatch.ApplyPatch(doc, patch)

// Custom configuration
result, err := jsonpatch.ApplyPatch(doc, patch, 
    jsonpatch.WithMutate(true),
    jsonpatch.WithMatcher(customRegexMatcher),
)
```

### 3. Always Use Test Operations for Critical Updates

```go
// Good: Test before critical changes
patch := []jsonpatch.Operation{
    {Op: "test", Path: "/balance", Value: 1000.0},
    {Op: "replace", Path: "/balance", Value: 500.0},
}
```

### 4. Choose Mutate Mode Carefully

```go
// Use mutate: false for concurrent access (default)
// Use mutate: true for large patches when performance matters
result, err := jsonpatch.ApplyPatch(doc, patch, 
    jsonpatch.WithMutate(len(patch) > 100), // Performance optimization
)
```

### 5. Handle Errors Gracefully

```go
if result, err := jsonpatch.ApplyPatch(doc, patch); err != nil {
    // Log error with context
    log.Printf("Patch failed for document %s: %v", docID, err)
    return originalDoc, err
}
```

### 6. Use Compact Codec for Storage/Network Efficiency

```go
import "github.com/kaptinlin/jsonpatch/codec/compact"

// For maximum space savings (35.9% reduction)
compactData, err := compact.Encode(patch, compact.WithNumericOpcodes(true))
if err == nil {
    // Store or transmit compactData instead of standard JSON
    saveToDatabase(compactData) // 35.9% smaller
    sendOverNetwork(compactData) // 35.9% less bandwidth
}
```

### 7. Use Batch Operations for Multiple Changes

```go
// Efficient: Single patch with multiple operations
patch := []jsonpatch.Operation{
    {Op: "replace", Path: "/status", Value: "active"},
    {Op: "inc", Path: "/version", Inc: 1},
    {Op: "add", Path: "/lastModified", Value: time.Now()},
}
```

### 10. Optimize with Compact Codec for Storage/Network

```go
import "github.com/kaptinlin/jsonpatch/codec/compact"

// For storage or network transmission
func storeOperations(ops []jsonpatch.Operation) error {
    // Use compact format for 35.9% space savings
    compactData, err := compact.Encode(ops, compact.WithNumericOpcodes(true))
    if err != nil {
        return err
    }
    
    return database.Store(compactData) // Much smaller than JSON
}

func loadOperations() ([]jsonpatch.Operation, error) {
    compactData, err := database.Load()
    if err != nil {
        return nil, err
    }
    
    return compact.Decode(compactData)
}
```

## 🔗 Related Specifications

- [RFC 6902 - JSON Patch](https://tools.ietf.org/html/rfc6902)
- [RFC 6901 - JSON Pointer](https://tools.ietf.org/html/rfc6901)
- [JSON Predicate Draft](https://tools.ietf.org/id/draft-snell-json-test-01.html)

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## 🎯 Credits

This project is a Golang port of [json-joy/json-patch](https://github.com/streamich/json-joy/tree/master/src/json-patch). Thanks to the original authors for their excellent work.

Original project: [streamich/json-joy](https://github.com/streamich/json-joy)

## 📄 License

This project is licensed under the MIT License. See [LICENSE](LICENSE) file for details.
