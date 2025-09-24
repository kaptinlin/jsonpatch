# JSON Patch Operations

This document covers all [JSON Patch (RFC 6902)][json-patch] operations:

- `add` - Add new values to the document
- `remove` - Remove existing values
- `replace` - Replace existing values
- `move` - Move values to different locations
- `copy` - Copy values to new locations
- `test` - Test values for conditional operations

## Basic Usage

```go
import "github.com/kaptinlin/jsonpatch"

doc := map[string]interface{}{
    "name": "Alice",
    "age":  25,
}

patch := []jsonpatch.Operation{
    {Op: "add", Path: "/email", Value: "alice@example.com"},
    {Op: "replace", Path: "/age", Value: 26},
}

result, err := jsonpatch.ApplyPatch(doc, patch)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Result: %+v\n", result.Doc)
// Output: map[age:26 email:alice@example.com name:Alice]
```

## Operations

### Add Operation

Add new values to objects or arrays.

```go
// Add to object
{Op: "add", Path: "/newField", Value: "newValue"}

// Add to end of array
{Op: "add", Path: "/items/-", Value: "newItem"}

// Insert at specific array position
{Op: "add", Path: "/items/0", Value: "firstItem"}
```

### Remove Operation

Remove values from objects or arrays.

```go
// Remove object property
{Op: "remove", Path: "/fieldToRemove"}

// Remove array element
{Op: "remove", Path: "/items/0"}
```

### Replace Operation

Replace existing values.

```go
{Op: "replace", Path: "/existingField", Value: "newValue"}
```

### Move Operation

Move values to new locations.

```go
{Op: "move", From: "/oldPath", Path: "/newPath"}
```

### Copy Operation

Copy values to new locations.

```go
{Op: "copy", From: "/sourcePath", Path: "/targetPath"}
```

### Test Operation

Test if values match expectations. Supports optional `not` flag for inverted tests.

```go
// Test equality
{Op: "test", Path: "/status", Value: "active"}

// Test inequality (inverted)
{Op: "test", Path: "/status", Value: "inactive", Not: true}
```

## Common Patterns

### Conditional Updates

```go
patch := []jsonpatch.Operation{
    // Test current value first
    {Op: "test", Path: "/version", Value: 1},
    // Then make changes
    {Op: "replace", Path: "/status", Value: "updated"},
    {Op: "replace", Path: "/version", Value: 2},
}
```

### Array Manipulation

```go
patch := []jsonpatch.Operation{
    // Add to end
    {Op: "add", Path: "/items/-", Value: "new item"},
    // Insert at beginning
    {Op: "add", Path: "/items/0", Value: "first item"},
    // Remove specific element
    {Op: "remove", Path: "/items/1"},
}
```

### Batch Operations

```go
var patch []jsonpatch.Operation

// Update multiple fields
for i := 0; i < 3; i++ {
    patch = append(patch, jsonpatch.Operation{
        Op:    "replace",
        Path:  fmt.Sprintf("/items/%d/status", i),
        Value: "processed",
    })
}

result, err := jsonpatch.ApplyPatch(doc, patch)
```

## Options

```go
// Default: creates a copy
result, err := jsonpatch.ApplyPatch(doc, patch)

// Mutate original document for better performance
result, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.WithMutate(true))
```

[json-patch]: https://tools.ietf.org/html/rfc6902
