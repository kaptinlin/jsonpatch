# JSON Patch Go

A comprehensive Go implementation of JSON Patch (RFC 6902), JSON Predicate, and extended operations for JSON document manipulation.

> **Note**: This is a Golang port of the powerful [json-joy/json-patch](https://github.com/streamich/json-joy/tree/master/src/json-patch), bringing JSON Patch extended operations to the Go ecosystem.

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

    // Apply patch
    options := jsonpatch.ApplyPatchOptions{Mutate: false}
    result, err := jsonpatch.ApplyPatch(doc, patch, options)
    if err != nil {
        log.Fatalf("Failed to apply patch: %v", err)
    }

    // Output result
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

## 📖 Examples

Explore comprehensive examples in the [`examples/`](examples/) directory:

- **[Basic Operations](examples/basic-operations/)** - Fundamental operations: `add`, `replace`, `remove`, `test`
- **[Array Operations](examples/array-operations/)** - Working with arrays and collections
- **[Conditional Operations](examples/conditional-operations/)** - Safe updates using test conditions
- **[Batch Updates](examples/batch-update/)** - Efficient batch operations
- **[Copy & Move Operations](examples/copy-move-operations/)** - Data restructuring and migration
- **[String Operations](examples/string-operations/)** - Text editing with string insertion
- **[Error Handling](examples/error-handling/)** - Robust error handling patterns

### Running Examples

```bash
# Run a specific example
cd examples/basic-operations && go run main.go
```

## 🎯 Features

### ✅ RFC 6902 Standard Operations

- **`add`** - Add new values to the document
- **`remove`** - Remove existing values
- **`replace`** - Replace existing values
- **`move`** - Move values to different locations
- **`copy`** - Copy values to new locations
- **`test`** - Test values for conditional operations

### 🔍 JSON Predicate Operations

Advanced querying and validation capabilities:

- **`contains`** - Check if strings contain substrings
- **`defined`** - Test if paths exist
- **`ends`** - Test string endings
- **`in`** - Check membership in arrays
- **`matches`** - Regular expression matching
- **`starts`** - Test string beginnings
- **`type`** - Type validation
- **`less`/`more`** - Numeric comparisons

### 🔧 Extended Operations

Beyond RFC 6902 for enhanced functionality:

- **`str_ins`** - String insertion for text editing
- **`str_del`** - String deletion
- **`inc`** - Increment numeric values
- **`flip`** - Boolean toggling
- **`split`** - Object splitting
- **`merge`** - Object merging

## 📚 API Reference

### Core Functions

```go
// Apply a single operation
func ApplyOp(doc interface{}, operation Op, mutate bool) (*OpResult, error)

// Apply multiple operations
func ApplyOps(doc interface{}, operations []Op, mutate bool) (*PatchResult, error)

// Apply JSON Patch
func ApplyPatch(doc interface{}, patch []Operation, options ApplyPatchOptions) (*PatchResult, error)
```

### Configuration Options

```go
type ApplyPatchOptions struct {
    Mutate           bool                // Modify the original document
    JsonPatchOptions JsonPatchOptions   // Additional options
}

type JsonPatchOptions struct {
    CreateMatcher CreateRegexMatcher     // Custom regex matcher
}
```

## 🔨 Common Patterns

### 1. Safe Updates with Test Operations

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
```

### 2. Batch Operations

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

result, err := jsonpatch.ApplyPatch(doc, patch, options)
```

### 3. Array Manipulation

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
```

### 4. String Operations

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
```

### 5. Error Handling

```go
result, err := jsonpatch.ApplyPatch(doc, patch, options)
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
```

## 🚀 Advanced Features

### Custom Regex Matchers

```go
customMatcher := func(pattern string, ignoreCase bool) jsonpatch.RegexMatcher {
    flags := ""
    if ignoreCase {
        flags = "(?i)"
    }
    
    re, _ := regexp.Compile(flags + pattern)
    return func(value string) bool {
        return re.MatchString(value)
    }
}

options := jsonpatch.ApplyPatchOptions{
    JsonPatchOptions: jsonpatch.JsonPatchOptions{
        CreateMatcher: customMatcher,
    },
}
```

### Performance Optimization

```go
// For large patches, use mutate mode to avoid deep copying
if len(patch) > 100 {
    options := jsonpatch.ApplyPatchOptions{Mutate: true}
    result, err := jsonpatch.ApplyPatch(doc, patch, options)
}

// For concurrent access, use immutable mode
options := jsonpatch.ApplyPatchOptions{Mutate: false}
result, err := jsonpatch.ApplyPatch(doc, patch, options)
```

## 📈 Best Practices

### 1. Always Use Test Operations for Critical Updates

```go
// Good: Test before critical changes
patch := []jsonpatch.Operation{
    {"op": "test", "path": "/balance", "value": 1000.0},
    {"op": "replace", "path": "/balance", "value": 500.0},
}
```

### 2. Choose Mutate Mode Carefully

```go
// Use mutate: false for concurrent access
// Use mutate: true for large patches when performance matters
options := jsonpatch.ApplyPatchOptions{
    Mutate: len(patch) > 100, // Performance optimization
}
```

### 3. Handle Errors Gracefully

```go
if result, err := jsonpatch.ApplyPatch(doc, patch, options); err != nil {
    // Log error with context
    log.Printf("Patch failed for document %s: %v", docID, err)
    return originalDoc, err
}
```

### 4. Use Batch Operations for Multiple Changes

```go
// Efficient: Single patch with multiple operations
patch := []jsonpatch.Operation{
    {"op": "replace", "path": "/status", "value": "active"},
    {"op": "inc", "path": "/version", "inc": 1},
    {"op": "add", "path": "/lastModified", "value": time.Now()},
}
```

## 📖 Documentation

For detailed documentation on specific operation types:

- **[JSON Patch Operations](docs/json-patch.md)** - Complete guide to standard RFC 6902 operations
- **[JSON Predicate Operations](docs/json-predicate.md)** - Advanced querying and conditional operations
- **[JSON Patch Extended Operations](docs/json-patch-extended.md)** - Extended operations for specialized use cases

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
