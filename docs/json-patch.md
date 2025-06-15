# JSON Patch Operations

This document covers all [JSON Patch (RFC 6902)][json-patch] operations implemented in this library:

- `add` - Add new values to the document
- `remove` - Remove existing values
- `replace` - Replace existing values
- `move` - Move values to different locations
- `copy` - Copy values to new locations
- `test` - Test values for conditional operations

The `test` operation is extended with an optional `not` property. When `not` is set to `true`, the result of the `test` operation is inverted.

## Basic Usage

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
        "name": "Alice",
        "age":  25,
    }
    
    // Create patch operations
    patch := []jsonpatch.Operation{
        {
            "op":    "add",
            "path":  "/email",
            "value": "alice@example.com",
        },
        {
            "op":    "replace",
            "path":  "/age",
            "value": 26,
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
    //   "age": 26,
    //   "email": "alice@example.com",
    //   "name": "Alice"
    // }
}
```

## Operation Details

### Add Operation

Add new values to the document.

```go
// Add to object
patch := []jsonpatch.Operation{
    {
        "op":    "add",
        "path":  "/newField",
        "value": "newValue",
    },
}

// Add to end of array
patch = []jsonpatch.Operation{
    {
        "op":    "add",
        "path":  "/items/-",
        "value": "newItem",
    },
}

// Add to specific array position
patch = []jsonpatch.Operation{
    {
        "op":    "add",
        "path":  "/items/0",
        "value": "firstItem",
    },
}
```

### Remove Operation

Remove values from the document.

```go
patch := []jsonpatch.Operation{
    {
        "op":   "remove",
        "path": "/fieldToRemove",
    },
}

// Remove array element
patch = []jsonpatch.Operation{
    {
        "op":   "remove",
        "path": "/items/0",
    },
}
```

### Replace Operation

Replace existing values.

```go
patch := []jsonpatch.Operation{
    {
        "op":    "replace",
        "path":  "/existingField",
        "value": "newValue",
    },
}
```

### Move Operation

Move values to new locations.

```go
patch := []jsonpatch.Operation{
    {
        "op":   "move",
        "from": "/oldPath",
        "path": "/newPath",
    },
}
```

### Copy Operation

Copy values to new locations.

```go
patch := []jsonpatch.Operation{
    {
        "op":   "copy",
        "from": "/sourcePath",
        "path": "/targetPath",
    },
}
```

### Test Operation

Test if values match expectations.

```go
patch := []jsonpatch.Operation{
    {
        "op":    "test",
        "path":  "/status",
        "value": "active",
    },
}

// Use not flag for inverted test
patch = []jsonpatch.Operation{
    {
        "op":    "test",
        "path":  "/status",
        "value": "inactive",
        "not":   true, // Test that value is NOT "inactive"
    },
}
```

## Advanced Usage

### Batch Operations

```go
func batchOperations() {
    doc := map[string]interface{}{
        "users": []interface{}{
            map[string]interface{}{"id": 1, "active": false},
            map[string]interface{}{"id": 2, "active": false},
        },
    }
    
    // Batch activate users
    patch := []jsonpatch.Operation{
        {
            "op":    "replace",
            "path":  "/users/0/active",
            "value": true,
        },
        {
            "op":    "replace",
            "path":  "/users/1/active",
            "value": true,
        },
    }
    
    options := jsonpatch.ApplyPatchOptions{Mutate: false}
    result, err := jsonpatch.ApplyPatch(doc, patch, options)
    if err != nil {
        log.Printf("Batch operation failed: %v", err)
        return
    }
    
    fmt.Printf("Updated document: %+v\n", result.Doc)
}
```

### Conditional Updates

```go
func conditionalUpdate() {
    doc := map[string]interface{}{
        "counter": 5,
        "status":  "pending",
    }
    
    // Only update if current value matches expectation
    patch := []jsonpatch.Operation{
        {
            "op":    "test",
            "path":  "/status",
            "value": "pending",
        },
        {
            "op":    "replace",
            "path":  "/status",
            "value": "processing",
        },
        {
            "op":    "replace",
            "path":  "/counter",
            "value": 6,
        },
    }
    
    options := jsonpatch.ApplyPatchOptions{Mutate: false}
    result, err := jsonpatch.ApplyPatch(doc, patch, options)
    if err != nil {
        log.Printf("Conditional update failed: %v", err)
        return
    }
    
    fmt.Printf("Updated document: %+v\n", result.Doc)
}
```

### Array Manipulation

```go
func arrayManipulation() {
    doc := map[string]interface{}{
        "items": []interface{}{"apple", "banana", "cherry"},
        "tags":  []interface{}{"fruit", "food"},
    }
    
    patch := []jsonpatch.Operation{
        // Insert at beginning
        {
            "op":    "add",
            "path":  "/items/0",
            "value": "orange",
        },
        // Append to end
        {
            "op":    "add",
            "path":  "/tags/-",
            "value": "healthy",
        },
        // Remove middle element
        {
            "op":   "remove",
            "path": "/items/2", // Note: indices shift after insertion
        },
        // Move element
        {
            "op":   "move",
            "from": "/items/0",
            "path": "/items/-",
        },
    }
    
    options := jsonpatch.ApplyPatchOptions{Mutate: false}
    result, err := jsonpatch.ApplyPatch(doc, patch, options)
    if err != nil {
        log.Printf("Array manipulation failed: %v", err)
        return
    }
    
    fmt.Printf("Result: %+v\n", result.Doc)
}
```

## Error Handling

```go
func errorHandling() {
    doc := map[string]interface{}{
        "name": "Alice",
    }
    
    // This patch will fail because path doesn't exist
    patch := []jsonpatch.Operation{
        {
            "op":   "remove",
            "path": "/nonexistent",
        },
    }
    
    options := jsonpatch.ApplyPatchOptions{Mutate: false}
    result, err := jsonpatch.ApplyPatch(doc, patch, options)
    if err != nil {
        // Handle specific error types
        switch {
        case strings.Contains(err.Error(), "path not found"):
            log.Printf("Invalid path in patch: %v", err)
        case strings.Contains(err.Error(), "test operation failed"):
            log.Printf("Test condition not met: %v", err)
        default:
            log.Printf("Patch application failed: %v", err)
        }
        return
    }
    
    fmt.Printf("Success: %+v\n", result.Doc)
}
```

## Performance Considerations

### Mutate Option

```go
// For better performance with large documents, use mutate: true
// This modifies the original document in-place
largeDoc := createLargeDocument()
patch := []jsonpatch.Operation{
    {
        "op":    "replace",
        "path":  "/status",
        "value": "updated",
    },
}

// Fast: modifies original document
options := jsonpatch.ApplyPatchOptions{Mutate: true}
result, err := jsonpatch.ApplyPatch(largeDoc, patch, options)

// Safe: creates a copy (slower but safer for concurrent access)
options = jsonpatch.ApplyPatchOptions{Mutate: false}
result, err = jsonpatch.ApplyPatch(largeDoc, patch, options)
```

[json-patch]: https://tools.ietf.org/html/rfc6902
