# JSON Patch Extended Operations

This document covers extended operations beyond the standard JSON Patch (RFC 6902) specification:

- `str_ins` - String insertion at a specific position
- `str_del` - String deletion by position and length or substring
- `inc` - Increment/decrement numeric values
- `flip` - Toggle boolean values
- `split` - Split values at a position
- `merge` - Merge adjacent array elements
- `extend` - Extend objects with additional properties

## Basic Usage

```go
import "github.com/kaptinlin/jsonpatch"

doc := map[string]any{
    "content": "Hello World",
    "counter": 10,
    "active":  false,
    "config": map[string]any{
        "theme": "dark",
        "lang":  "en",
    },
}

patch := []jsonpatch.Operation{
    {Op: "str_ins", Path: "/content", Pos: 6, Str: "Beautiful "},
    {Op: "inc", Path: "/counter", Inc: 5},
    {Op: "flip", Path: "/active"},
    {Op: "extend", Path: "/config", Props: map[string]any{
        "version": "1.0",
        "debug":   true,
    }},
}

result, err := jsonpatch.ApplyPatch(doc, patch)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Result: %+v\n", result.Doc)
// Output: map[active:true config:map[debug:true lang:en theme:dark version:1.0] content:Hello Beautiful World counter:15]
```

## String Operations

### String Insertion (str_ins)

Insert text at a specific position within a string.

```go
// Insert at beginning
{Op: "str_ins", Path: "/text", Pos: 0, Str: "Hi! "}

// Insert in middle
{Op: "str_ins", Path: "/text", Pos: 6, Str: "Beautiful "}

// Insert at end (use the string length as position)
{Op: "str_ins", Path: "/text", Pos: 11, Str: "!"}
```

### String Deletion (str_del)

Delete text from strings using position and length or substring matching.

```go
// Delete by position and length
{Op: "str_del", Path: "/text", Pos: 6, Len: 10}

// Delete by substring at position
{Op: "str_del", Path: "/text", Pos: 6, Str: "Beautiful "}

// Delete from position to end (use a large length value)
{Op: "str_del", Path: "/text", Pos: 5, Len: 100}
```

## Numeric Operations

### Increment (inc)

Increment numeric values by a specified amount.

```go
// Increment positive
{Op: "inc", Path: "/counter", Inc: 5}

// Increment negative (decrement)
{Op: "inc", Path: "/score", Inc: -10}

// Increment float
{Op: "inc", Path: "/price", Inc: 2.5}
```

## Boolean Operations

### Flip

Toggle boolean values.

```go
// Flip boolean
{Op: "flip", Path: "/active"}

// Works with any boolean field
{Op: "flip", Path: "/settings/notifications"}
```

## Object Operations

### Extend

Add properties to an object without replacing existing ones.

```go
{Op: "extend", Path: "/config", Props: map[string]any{
    "version": "1.0",
    "debug":   true,
}}
```

### Split

Split a value at a specified position. Works with strings, numbers, and Slate.js-style nodes.

```go
// Split at position
{Op: "split", Path: "/items/0", Pos: 5}
```

### Merge

Merge adjacent array elements. The `Pos` field specifies the number of elements to merge.

```go
// Merge adjacent elements
{Op: "merge", Path: "/items/1", Pos: 1}
```

## Common Patterns

### Text Editing

```go
patch := []jsonpatch.Operation{
    // Insert prefix
    {Op: "str_ins", Path: "/title", Pos: 0, Str: "[DRAFT] "},
    // Remove unwanted text by position and length
    {Op: "str_del", Path: "/content", Pos: 0, Len: 6},
}
```

### Counter Operations

```go
patch := []jsonpatch.Operation{
    // Increment view count
    {Op: "inc", Path: "/stats/views", Inc: 1},
    // Decrement inventory
    {Op: "inc", Path: "/inventory/count", Inc: -1},
    // Update score
    {Op: "inc", Path: "/user/score", Inc: 100},
}
```

### Configuration Updates

```go
patch := []jsonpatch.Operation{
    // Toggle feature flag
    {Op: "flip", Path: "/features/newUI"},
    // Add new settings without overwriting existing ones
    {Op: "extend", Path: "/settings", Props: map[string]any{
        "autoSave": true,
        "timeout":  30,
    }},
}
```

## Advanced Examples

### Blog Post Editor

```go
doc := map[string]any{
    "title":   "My Blog Post",
    "content": "This is the content.",
    "stats":   map[string]any{"views": 0, "likes": 0},
    "draft":   true,
}

patch := []jsonpatch.Operation{
    // Update title
    {Op: "str_ins", Path: "/title", Pos: 0, Str: "[Updated] "},
    // Append to content
    {Op: "str_ins", Path: "/content", Pos: 20, Str: " More content added."},
    // Increment views
    {Op: "inc", Path: "/stats/views", Inc: 1},
    // Publish (toggle draft)
    {Op: "flip", Path: "/draft"},
    // Add metadata
    {Op: "extend", Path: "/", Props: map[string]any{
        "publishedAt": "2024-01-01T00:00:00Z",
        "author":      "John Doe",
    }},
}
```

### User Profile Update

```go
doc := map[string]any{
    "user": map[string]any{
        "name": "John",
        "settings": map[string]any{
            "theme": "light",
        },
    },
    "stats": map[string]any{
        "loginCount": 5,
    },
}

patch := []jsonpatch.Operation{
    // Append to name
    {Op: "str_ins", Path: "/user/name", Pos: 4, Str: " Doe"},
    // Increment login count
    {Op: "inc", Path: "/stats/loginCount", Inc: 1},
    // Add new settings
    {Op: "extend", Path: "/user/settings", Props: map[string]any{
        "notifications": true,
        "language":      "en",
    }},
}
```

### E-commerce Operations

```go
doc := map[string]any{
    "product": map[string]any{
        "name":  "Laptop",
        "price": 999.99,
        "stock": 10,
    },
    "sale": false,
}

patch := []jsonpatch.Operation{
    // Apply discount
    {Op: "inc", Path: "/product/price", Inc: -100.0},
    // Reduce stock
    {Op: "inc", Path: "/product/stock", Inc: -1},
    // Enable sale
    {Op: "flip", Path: "/sale"},
    // Add sale info
    {Op: "extend", Path: "/product", Props: map[string]any{
        "salePrice": 899.99,
        "discount":  "10%",
    }},
}
```
