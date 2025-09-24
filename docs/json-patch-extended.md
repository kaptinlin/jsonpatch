# JSON Patch Extended Operations

This document covers extended operations beyond the standard JSON Patch (RFC 6902) specification:

- `str_ins` - String insertion for text editing
- `str_del` - String deletion with position and length/substring
- `inc` - Increment numeric values
- `flip` - Boolean toggling
- `split` - Object splitting operations
- `merge` - Object merging operations
- `extend` - Object extension with property merging

## Basic Usage

```go
import "github.com/kaptinlin/jsonpatch"

doc := map[string]interface{}{
    "content": "Hello World",
    "counter": 10,
    "active":  false,
    "config": map[string]interface{}{
        "theme": "dark",
        "lang":  "en",
    },
}

patch := []jsonpatch.Operation{
    {Op: "str_ins", Path: "/content", Pos: 6, Str: "Beautiful "},
    {Op: "inc", Path: "/counter", Inc: 5},
    {Op: "flip", Path: "/active"},
    {Op: "extend", Path: "/config", Props: map[string]interface{}{
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

// Insert at end
{Op: "str_ins", Path: "/text", Pos: 11, Str: "!"}
```

### String Deletion (str_del)

Delete text from strings using position and length or substring matching.

```go
// Delete by position and length
{Op: "str_del", Path: "/text", Pos: 6, Len: 10}

// Delete by substring
{Op: "str_del", Path: "/text", Pos: 6, Str: "Beautiful "}

// Delete from position to end
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
{Op: "extend", Path: "/config", Props: {
    "version": "1.0",
    "debug": true
}}
```

### Merge

Merge objects, replacing existing properties.

```go
{Op: "merge", Path: "/settings", Props: {
    "theme": "light",
    "language": "zh"
}}
```

### Split

Split an object into multiple properties.

```go
// Split object properties to parent level
{Op: "split", Path: "/user/address", Props: ["street", "city", "zip"]}
```

## Common Patterns

### Text Editing

```go
patch := []jsonpatch.Operation{
    // Insert prefix
    {Op: "str_ins", Path: "/title", Pos: 0, Str: "[DRAFT] "},
    // Append suffix
    {Op: "str_ins", Path: "/title", Pos: -1, Str: " - Updated"},
    // Remove unwanted text
    {Op: "str_del", Path: "/content", Str: "TODO: "},
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
    // Add new settings
    {Op: "extend", Path: "/settings", Props: {
        "autoSave": true,
        "timeout": 30
    }},
    // Merge user preferences
    {Op: "merge", Path: "/preferences", Props: {
        "theme": "dark",
        "language": "en"
    }},
}
```

### Data Transformation

```go
patch := []jsonpatch.Operation{
    // Split address into separate fields
    {Op: "split", Path: "/user/fullAddress", Props: ["street", "city", "country"]},
    // Merge contact info
    {Op: "merge", Path: "/user", Props: {
        "email": "user@example.com",
        "phone": "+1234567890"
    }},
}
```

## Advanced Examples

### Blog Post Editor

```go
doc := map[string]interface{}{
    "title":   "My Blog Post",
    "content": "This is the content.",
    "stats":   map[string]interface{}{"views": 0, "likes": 0},
    "draft":   true,
}

patch := []jsonpatch.Operation{
    // Update title
    {Op: "str_ins", Path: "/title", Pos: 0, Str: "[Updated] "},
    // Add content
    {Op: "str_ins", Path: "/content", Pos: -1, Str: " More content added."},
    // Increment views
    {Op: "inc", Path: "/stats/views", Inc: 1},
    // Publish
    {Op: "flip", Path: "/draft"},
    // Add metadata
    {Op: "extend", Path: "/", Props: {
        "publishedAt": "2024-01-01T00:00:00Z",
        "author": "John Doe"
    }},
}
```

### User Profile Update

```go
doc := map[string]interface{}{
    "user": map[string]interface{}{
        "name": "John",
        "settings": map[string]interface{}{
            "theme": "light",
        },
    },
    "stats": map[string]interface{}{
        "loginCount": 5,
    },
}

patch := []jsonpatch.Operation{
    // Update name
    {Op: "str_ins", Path: "/user/name", Pos: -1, Str: " Doe"},
    // Increment login count
    {Op: "inc", Path: "/stats/loginCount", Inc: 1},
    // Merge new settings
    {Op: "merge", Path: "/user/settings", Props: {
        "notifications": true,
        "language": "en"
    }},
}
```

### E-commerce Operations

```go
doc := map[string]interface{}{
    "product": map[string]interface{}{
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
    {Op: "extend", Path: "/product", Props: {
        "salePrice": 899.99,
        "discount": "10%"
    }},
}
```
