# JSON Patch Extended Operations

This document covers the extended operations beyond the standard JSON Patch (RFC 6902) specification:

- `str_ins` - String insertion for text editing
- `str_del` - String deletion with position and length/substring
- `inc` - Increment numeric values
- `flip` - Boolean toggling
- `split` - Object splitting operations
- `merge` - Object merging operations
- `extend` - Object extension with property merging

These operations provide enhanced functionality for specialized use cases like text editing, mathematical operations, and advanced object manipulation.

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
        "content": "Hello World",
        "counter": 10,
        "active":  false,
        "config": map[string]interface{}{
            "theme": "dark",
            "lang":  "en",
        },
    }
    
    // Extended operations
    patch := []jsonpatch.Operation{
        // String insertion
        {
            "op":   "str_ins",
            "path": "/content",
            "pos":  6,
            "str":  "Beautiful ",
        },
        // Increment number
        {
            "op":   "inc",
            "path": "/counter",
            "inc":  5,
        },
        // Flip boolean
        {
            "op":   "flip",
            "path": "/active",
        },
        // Extend object
        {
            "op":   "extend",
            "path": "/config",
            "props": map[string]interface{}{
                "version": "1.0",
                "debug":   true,
            },
        },
    }
    
    // Apply patch
    options := jsonpatch.ApplyPatchOptions{Mutate: false}
    result, err := jsonpatch.ApplyPatch(doc, patch, options)
    if err != nil {
        log.Fatalf("Failed to apply extended patch: %v", err)
    }
    
    // Output result
    output, _ := json.MarshalIndent(result.Doc, "", "  ")
    fmt.Println(string(output))
    // Output:
    // {
    //   "active": true,
    //   "config": {
    //     "debug": true,
    //     "lang": "en",
    //     "theme": "dark",
    //     "version": "1.0"
    //   },
    //   "content": "Hello Beautiful World",
    //   "counter": 15
    // }
}
```

## String Operations

### String Insertion (str_ins)

Insert text at a specific position within a string.

```go
doc := map[string]interface{}{
    "text": "Hello World",
}

// Insert at beginning
patch := []jsonpatch.Operation{
    {
        "op":   "str_ins",
        "path": "/text",
        "pos":  0,
        "str":  "Hi! ",
    },
}
// Result: "Hi! Hello World"

// Insert in middle
patch = []jsonpatch.Operation{
    {
        "op":   "str_ins",
        "path": "/text",
        "pos":  6,
        "str":  "Beautiful ",
    },
}
// Result: "Hello Beautiful World"

// Insert at end (position can be string length)
patch = []jsonpatch.Operation{
    {
        "op":   "str_ins",
        "path": "/text",
        "pos":  11,
        "str":  "!",
    },
}
// Result: "Hello World!"
```

### String Deletion (str_del)

Delete text from strings using position and length or substring matching.

```go
doc := map[string]interface{}{
    "text": "Hello Beautiful World",
}

// Delete by position and length
patch := []jsonpatch.Operation{
    {
        "op":   "str_del",
        "path": "/text",
        "pos":  6,
        "len":  10, // Delete "Beautiful "
    },
}
// Result: "Hello World"

// Delete by substring
patch = []jsonpatch.Operation{
    {
        "op":   "str_del",
        "path": "/text",
        "pos":  6,
        "str":  "Beautiful ", // Delete specific substring
    },
}
// Result: "Hello World"

// Delete from position to end
patch = []jsonpatch.Operation{
    {
        "op":   "str_del",
        "path": "/text",
        "pos":  5,
        "len":  100, // Large length deletes to end
    },
}
// Result: "Hello"
```

## Numeric Operations

### Increment (inc)

Increment numeric values by a specified amount.

```go
doc := map[string]interface{}{
    "counter":     10,
    "score":       95.5,
    "negative":    -5,
    "nested": map[string]interface{}{
        "value": 100,
    },
}

// Increment positive
patch := []jsonpatch.Operation{
    {
        "op":   "inc",
        "path": "/counter",
        "inc":  5,
    },
}
// Result: counter becomes 15

// Increment with decimal
patch = []jsonpatch.Operation{
    {
        "op":   "inc",
        "path": "/score",
        "inc":  2.5,
    },
}
// Result: score becomes 98.0

// Decrement (negative increment)
patch = []jsonpatch.Operation{
    {
        "op":   "inc",
        "path": "/counter",
        "inc":  -3,
    },
}
// Result: counter becomes 7

// Increment nested value
patch = []jsonpatch.Operation{
    {
        "op":   "inc",
        "path": "/nested/value",
        "inc":  50,
    },
}
// Result: nested.value becomes 150
```

## Boolean Operations

### Flip (flip)

Toggle boolean values.

```go
doc := map[string]interface{}{
    "active":    true,
    "enabled":   false,
    "settings": map[string]interface{}{
        "debug": false,
    },
}

// Flip boolean values
patch := []jsonpatch.Operation{
    {
        "op":   "flip",
        "path": "/active",
    },
    {
        "op":   "flip",
        "path": "/enabled",
    },
    {
        "op":   "flip",
        "path": "/settings/debug",
    },
}

options := jsonpatch.ApplyPatchOptions{Mutate: false}
result, err := jsonpatch.ApplyPatch(doc, patch, options)
// Result: active=false, enabled=true, settings.debug=true
```

## Object Operations

### Extend (extend)

Extend objects with new properties, optionally deleting null values.

```go
doc := map[string]interface{}{
    "user": map[string]interface{}{
        "name":  "Alice",
        "email": "alice@example.com",
        "temp":  nil,
    },
}

// Basic extend
patch := []jsonpatch.Operation{
    {
        "op":   "extend",
        "path": "/user",
        "props": map[string]interface{}{
            "age":    25,
            "active": true,
            "role":   "admin",
        },
    },
}

// Extend with null deletion
patch = []jsonpatch.Operation{
    {
        "op":   "extend",
        "path": "/user",
        "props": map[string]interface{}{
            "age":     25,
            "temp":    nil,  // This will be deleted
            "updated": true,
        },
        "deleteNull": true,
    },
}
```

### Merge (merge)

Merge objects at specified positions in arrays or merge object properties.

```go
doc := map[string]interface{}{
    "items": []interface{}{
        map[string]interface{}{
            "id":   1,
            "name": "Item 1",
        },
        map[string]interface{}{
            "id":   2,
            "name": "Item 2",
        },
    },
}

// Merge objects in array
patch := []jsonpatch.Operation{
    {
        "op":   "merge",
        "path": "/items",
        "pos":  1, // Position in array (1-based)
        "props": map[string]interface{}{
            "description": "Updated item",
            "active":      true,
        },
    },
}
```

### Split (split)

Split objects by extracting properties to new locations.

```go
doc := map[string]interface{}{
    "user": map[string]interface{}{
        "name":    "Alice",
        "email":   "alice@example.com",
        "age":     25,
        "address": "123 Main St",
        "phone":   "555-1234",
    },
}

// Split user object
patch := []jsonpatch.Operation{
    {
        "op":   "split",
        "path": "/user",
        "pos":  2, // Split position
        // Properties before pos go to one part, after pos go to another
    },
}
```

## Advanced Usage Examples

### Text Editor Operations

```go
func textEditorOperations() {
    doc := map[string]interface{}{
        "document": map[string]interface{}{
            "title":   "My Document",
            "content": "The quick brown fox jumps over the lazy dog.",
            "version": 1,
        },
    }
    
    // Simulate text editing operations
    patch := []jsonpatch.Operation{
        // Add prefix to title
        {
            "op":   "str_ins",
            "path": "/document/title",
            "pos":  0,
            "str":  "Draft: ",
        },
        // Replace word in content
        {
            "op":   "str_del",
            "path": "/document/content",
            "pos":  16,
            "str":  "fox",
        },
        {
            "op":   "str_ins",
            "path": "/document/content",
            "pos":  16,
            "str":  "cat",
        },
        // Add exclamation
        {
            "op":   "str_ins",
            "path": "/document/content",
            "pos":  43, // Adjusted for previous changes
            "str":  "!",
        },
        // Increment version
        {
            "op":   "inc",
            "path": "/document/version",
            "inc":  1,
        },
    }
    
    options := jsonpatch.ApplyPatchOptions{Mutate: false}
    result, err := jsonpatch.ApplyPatch(doc, patch, options)
    if err != nil {
        log.Printf("Text editing failed: %v", err)
        return
    }
    
    fmt.Printf("Edited document: %+v\n", result.Doc)
}
```

### Gaming Score System

```go
func gamingScoreSystem() {
    doc := map[string]interface{}{
        "player": map[string]interface{}{
            "name":    "Player1",
            "score":   1000,
            "level":   5,
            "active":  true,
            "achievements": []interface{}{
                "first_win", "level_5",
            },
        },
        "game": map[string]interface{}{
            "round":     10,
            "bonus":     false,
            "completed": false,
        },
    }
    
    // Game actions: win round with bonus
    patch := []jsonpatch.Operation{
        // Increment score with bonus
        {
            "op":   "inc",
            "path": "/player/score",
            "inc":  250, // Base score + bonus
        },
        // Increment level
        {
            "op":   "inc",
            "path": "/player/level",
            "inc":  1,
        },
        // Next round
        {
            "op":   "inc",
            "path": "/game/round",
            "inc":  1,
        },
        // Activate bonus
        {
            "op":   "flip",
            "path": "/game/bonus",
        },
        // Add new achievement
        {
            "op":    "add",
            "path":  "/player/achievements/-",
            "value": "bonus_round",
        },
        // Extend player with new stats
        {
            "op":   "extend",
            "path": "/player",
            "props": map[string]interface{}{
                "lastPlayed":  "2024-01-15",
                "winStreak":   3,
                "bonusPoints": 50,
            },
        },
    }
    
    options := jsonpatch.ApplyPatchOptions{Mutate: false}
    result, err := jsonpatch.ApplyPatch(doc, patch, options)
    if err != nil {
        log.Printf("Game update failed: %v", err)
        return
    }
    
    fmt.Printf("Updated game state: %+v\n", result.Doc)
}
```

### Configuration Management

```go
func configurationManagement() {
    doc := map[string]interface{}{
        "app": map[string]interface{}{
            "name":    "MyApp",
            "version": "1.0.0",
            "debug":   false,
            "config": map[string]interface{}{
                "theme":    "light",
                "language": "en",
                "timeout":  30,
            },
        },
    }
    
    // Update configuration
    patch := []jsonpatch.Operation{
        // Update version string
        {
            "op":   "str_del",
            "path": "/app/version",
            "pos":  0,
            "len":  5, // Remove "1.0.0"
        },
        {
            "op":   "str_ins",
            "path": "/app/version",
            "pos":  0,
            "str":  "1.1.0",
        },
        // Enable debug mode
        {
            "op":   "flip",
            "path": "/app/debug",
        },
        // Extend configuration
        {
            "op":   "extend",
            "path": "/app/config",
            "props": map[string]interface{}{
                "maxRetries":    3,
                "cacheEnabled":  true,
                "logLevel":      "info",
                "apiVersion":    "v2",
                "deprecated":    nil, // Will be ignored unless deleteNull is true
            },
            "deleteNull": false,
        },
        // Increase timeout
        {
            "op":   "inc",
            "path": "/app/config/timeout",
            "inc":  10,
        },
    }
    
    options := jsonpatch.ApplyPatchOptions{Mutate: false}
    result, err := jsonpatch.ApplyPatch(doc, patch, options)
    if err != nil {
        log.Printf("Configuration update failed: %v", err)
        return
    }
    
    fmt.Printf("Updated configuration: %+v\n", result.Doc)
}
```

### Batch Operations

```go
func batchExtendedOperations() {
    doc := map[string]interface{}{
        "counters": map[string]interface{}{
            "views":      100,
            "downloads":  50,
            "shares":     25,
        },
        "flags": map[string]interface{}{
            "featured":   false,
            "published":  true,
            "deprecated": false,
        },
        "text": map[string]interface{}{
            "title":       "Sample Title",
            "description": "This is a sample description.",
        },
    }
    
    // Batch update multiple values
    patch := []jsonpatch.Operation{
        // Increment all counters
        {
            "op":   "inc",
            "path": "/counters/views",
            "inc":  10,
        },
        {
            "op":   "inc",
            "path": "/counters/downloads",
            "inc":  5,
        },
        {
            "op":   "inc",
            "path": "/counters/shares",
            "inc":  2,
        },
        // Flip flags
        {
            "op":   "flip",
            "path": "/flags/featured",
        },
        {
            "op":   "flip",
            "path": "/flags/deprecated",
        },
        // Update text content
        {
            "op":   "str_ins",
            "path": "/text/title",
            "pos":  0,
            "str":  "Updated: ",
        },
        {
            "op":   "str_ins",
            "path": "/text/description",
            "pos":  46, // End of sentence
            "str":  " Updated with new information.",
        },
        // Extend with metadata
        {
            "op":   "extend",
            "path": "/",
            "props": map[string]interface{}{
                "lastModified": "2024-01-15T10:30:00Z",
                "modifiedBy":   "system",
                "changeLog":    []interface{}{"incremented counters", "flipped flags", "updated text"},
            },
        },
    }
    
    options := jsonpatch.ApplyPatchOptions{Mutate: false}
    result, err := jsonpatch.ApplyPatch(doc, patch, options)
    if err != nil {
        log.Printf("Batch operations failed: %v", err)
        return
    }
    
    fmt.Printf("Batch update result: %+v\n", result.Doc)
}
```

## Error Handling

```go
func extendedOperationErrors() {
    doc := map[string]interface{}{
        "text":    "Hello",
        "number":  "not a number", // Wrong type
        "boolean": "not a boolean", // Wrong type
    }
    
    // Operations that will fail
    patch := []jsonpatch.Operation{
        {
            "op":   "str_ins",
            "path": "/text",
            "pos":  100, // Position out of range
            "str":  "test",
        },
        {
            "op":   "inc",
            "path": "/number",
            "inc":  5, // Cannot increment non-numeric value
        },
        {
            "op":   "flip",
            "path": "/boolean", // Cannot flip non-boolean value
        },
    }
    
    options := jsonpatch.ApplyPatchOptions{Mutate: false}
    result, err := jsonpatch.ApplyPatch(doc, patch, options)
    if err != nil {
        // Handle specific error types
        switch {
        case strings.Contains(err.Error(), "position out of range"):
            log.Printf("String operation position error: %v", err)
        case strings.Contains(err.Error(), "not a number"):
            log.Printf("Numeric operation type error: %v", err)
        case strings.Contains(err.Error(), "not a boolean"):
            log.Printf("Boolean operation type error: %v", err)
        default:
            log.Printf("Extended operation error: %v", err)
        }
        return
    }
    
    fmt.Printf("Operations succeeded: %+v\n", result.Doc)
}
```

## Performance Considerations

Extended operations are generally more efficient than equivalent combinations of standard operations:

```go
// Instead of multiple operations for text editing:
patch := []jsonpatch.Operation{
    {"op": "test", "path": "/text", "value": "Hello World"},
    {"op": "replace", "path": "/text", "value": "Hello Beautiful World"},
}

// Use single str_ins operation:
patch = []jsonpatch.Operation{
    {
        "op":   "str_ins",
        "path": "/text",
        "pos":  6,
        "str":  "Beautiful ",
    },
}

// Instead of test + replace for increment:
patch = []jsonpatch.Operation{
    {"op": "test", "path": "/counter", "value": 10},
    {"op": "replace", "path": "/counter", "value": 15},
}

// Use single inc operation:
patch = []jsonpatch.Operation{
    {
        "op":   "inc",
        "path": "/counter",
        "inc":  5,
    },
}
```
