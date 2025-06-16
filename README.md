# JSON Patch Go

A comprehensive Go implementation of JSON Patch (RFC 6902), JSON Predicate, and extended operations for JSON document manipulation with **full type safety** and **generic support**.

> **Note**: This is a Golang port of the powerful [json-joy/json-patch](https://github.com/streamich/json-joy/tree/master/src/json-patch), bringing JSON Patch extended operations to the Go ecosystem with modern Go generics.

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

    // Create patch operations
    patch := []jsonpatch.Operation{
        {
            "op":    "replace",
            "path":  "/name",
            "value": "Jane",
        },
        {
            "op":    "add",
            "path":  "/email",
            "value": "jane@example.com",
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
        {"op": "replace", "path": "/name", "value": "Jane"},
        {"op": "add", "path": "/email", "value": "jane@example.com"},
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

### 🔧 Extended Operations ([docs](docs/json-patch-extended.md))

- **`str_ins`** - String insertion for text editing
- **`str_del`** - String deletion
- **`inc`** - Increment numeric values
- **`flip`** - Boolean toggling
- **`split`** - Object splitting
- **`merge`** - Object merging

## 📖 Examples

Explore comprehensive examples in the [`examples/`](examples/) directory:

### Core Operations
- **[Basic Operations](examples/basic-operations/)** - `add`, `replace`, `remove`, `test`
- **[Array Operations](examples/array-operations/)** - Array manipulation
- **[Conditional Operations](examples/conditional-operations/)** - Safe updates with tests
- **[Copy & Move](examples/copy-move-operations/)** - Data restructuring
- **[String Operations](examples/string-operations/)** - Text editing

### Document Types
- **[Struct Patch](examples/struct-patch/)** - Type-safe Go structs
- **[Map Patch](examples/map-patch/)** - Dynamic `map[string]any`
- **[JSON Bytes](examples/json-bytes-patch/)** - Raw JSON byte data
- **[JSON String](examples/json-string-patch/)** - JSON string data

### Advanced
- **[Batch Updates](examples/batch-update/)** - Bulk operations
- **[Error Handling](examples/error-handling/)** - Error patterns
- **[Mutate Option](examples/mutate-option/)** - Performance optimization

```bash
# Run any example
cd examples/<example-name> && go run main.go
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
    {"op": "inc", "path": "/version", "inc": 1},
    {"op": "replace", "path": "/status", "value": "updated"},
    {"op": "flip", "path": "/enabled"},
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
        "op":    "test",
        "path":  "/version",
        "value": 1,
    },
    // Make the change
    {
        "op":    "replace",
        "path":  "/status",
        "value": "updated",
    },
    // Increment version
    {
        "op":   "inc",
        "path": "/version",
        "inc":  1,
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
        "op":    "replace",
        "path":  fmt.Sprintf("/items/%d/status", i),
        "value": "processed",
    })
}

result, err := jsonpatch.ApplyPatch(doc, patch)
```

### 5. Array Manipulation

```go
patch := []jsonpatch.Operation{
    // Add to end of array
    {
        "op":   "add",
        "path": "/users/-",
        "value": map[string]interface{}{"name": "New User"},
    },
    // Insert at specific position
    {
        "op":   "add",
        "path": "/tags/0",
        "value": "important",
    },
    // Remove array element
    {
        "op":   "remove",
        "path": "/items/2",
    },
}

result, err := jsonpatch.ApplyPatch(doc, patch)
```

### 6. String Operations

```go
patch := []jsonpatch.Operation{
    // Insert text at position
    {
        "op":   "str_ins",
        "path": "/content",
        "pos":  0,
        "str":  "Prefix: ",
    },
    // Insert at end
    {
        "op":   "str_ins",
        "path": "/content",
        "pos":  20,
        "str":  " (Updated)",
    },
}

result, err := jsonpatch.ApplyPatch(doc, patch)
```

### 7. Error Handling

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
    {"op": "test", "path": "/balance", "value": 1000.0},
    {"op": "replace", "path": "/balance", "value": 500.0},
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

### 6. Use Batch Operations for Multiple Changes

```go
// Efficient: Single patch with multiple operations
patch := []jsonpatch.Operation{
    {"op": "replace", "path": "/status", "value": "active"},
    {"op": "inc", "path": "/version", "inc": 1},
    {"op": "add", "path": "/lastModified", "value": time.Now()},
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
