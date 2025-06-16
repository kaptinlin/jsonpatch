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
    {"op": "str_ins", "path": "/content", "pos": 6, "str": "Beautiful "},
    {"op": "inc", "path": "/counter", "inc": 5},
    {"op": "flip", "path": "/active"},
    {"op": "extend", "path": "/config", "props": map[string]interface{}{
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
{"op": "str_ins", "path": "/text", "pos": 0, "str": "Hi! "}

// Insert in middle
{"op": "str_ins", "path": "/text", "pos": 6, "str": "Beautiful "}

// Insert at end
{"op": "str_ins", "path": "/text", "pos": 11, "str": "!"}
```

### String Deletion (str_del)

Delete text from strings using position and length or substring matching.

```go
// Delete by position and length
{"op": "str_del", "path": "/text", "pos": 6, "len": 10}

// Delete by substring
{"op": "str_del", "path": "/text", "pos": 6, "str": "Beautiful "}

// Delete from position to end
{"op": "str_del", "path": "/text", "pos": 5, "len": 100}
```

## Numeric Operations

### Increment (inc)

Increment numeric values by a specified amount.

```go
// Increment positive
{"op": "inc", "path": "/counter", "inc": 5}

// Increment negative (decrement)
{"op": "inc", "path": "/score", "inc": -10}

// Increment float
{"op": "inc", "path": "/price", "inc": 2.5}
```

## Boolean Operations

### Flip

Toggle boolean values.

```go
// Flip boolean
{"op": "flip", "path": "/active"}

// Works with any boolean field
{"op": "flip", "path": "/settings/notifications"}
```

## Object Operations

### Extend

Add properties to an object without replacing existing ones.

```go
{"op": "extend", "path": "/config", "props": {
    "version": "1.0",
    "debug": true
}}
```

### Merge

Merge objects, replacing existing properties.

```go
{"op": "merge", "path": "/settings", "props": {
    "theme": "light",
    "language": "zh"
}}
```

### Split

Split an object into multiple properties.

```go
// Split object properties to parent level
{"op": "split", "path": "/user/address", "props": ["street", "city", "zip"]}
```

## Common Patterns

### Text Editing

```go
patch := []jsonpatch.Operation{
    // Insert prefix
    {"op": "str_ins", "path": "/title", "pos": 0, "str": "[DRAFT] "},
    // Append suffix
    {"op": "str_ins", "path": "/title", "pos": -1, "str": " - Updated"},
    // Remove unwanted text
    {"op": "str_del", "path": "/content", "str": "TODO: "},
}
```

### Counter Operations

```go
patch := []jsonpatch.Operation{
    // Increment view count
    {"op": "inc", "path": "/stats/views", "inc": 1},
    // Decrement inventory
    {"op": "inc", "path": "/inventory/count", "inc": -1},
    // Update score
    {"op": "inc", "path": "/user/score", "inc": 100},
}
```

### Configuration Updates

```go
patch := []jsonpatch.Operation{
    // Toggle feature flag
    {"op": "flip", "path": "/features/newUI"},
    // Add new settings
    {"op": "extend", "path": "/settings", "props": {
        "autoSave": true,
        "timeout": 30
    }},
    // Merge user preferences
    {"op": "merge", "path": "/preferences", "props": {
        "theme": "dark",
        "language": "en"
    }},
}
```

### Data Transformation

```go
patch := []jsonpatch.Operation{
    // Split address into separate fields
    {"op": "split", "path": "/user/fullAddress", "props": ["street", "city", "country"]},
    // Merge contact info
    {"op": "merge", "path": "/user", "props": {
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
    {"op": "str_ins", "path": "/title", "pos": 0, "str": "[Updated] "},
    // Add content
    {"op": "str_ins", "path": "/content", "pos": -1, "str": " More content added."},
    // Increment views
    {"op": "inc", "path": "/stats/views", "inc": 1},
    // Publish
    {"op": "flip", "path": "/draft"},
    // Add metadata
    {"op": "extend", "path": "/", "props": {
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
    {"op": "str_ins", "path": "/user/name", "pos": -1, "str": " Doe"},
    // Increment login count
    {"op": "inc", "path": "/stats/loginCount", "inc": 1},
    // Merge new settings
    {"op": "merge", "path": "/user/settings", "props": {
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
    {"op": "inc", "path": "/product/price", "inc": -100.0},
    // Reduce stock
    {"op": "inc", "path": "/product/stock", "inc": -1},
    // Enable sale
    {"op": "flip", "path": "/sale"},
    // Add sale info
    {"op": "extend", "path": "/product", "props": {
        "salePrice": 899.99,
        "discount": "10%"
    }},
}
```
