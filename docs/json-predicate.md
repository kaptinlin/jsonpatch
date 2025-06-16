# JSON Predicate Operations

This document covers [JSON Predicate][json-predicate] operations for conditional testing and validation:

- `test` - Test equality with optional negation
- `contains` - Check if arrays contain values or strings contain substrings
- `defined` - Test if paths exist and are not undefined
- `undefined` - Test if paths don't exist or are undefined
- `starts` - Test if strings start with specific prefixes
- `ends` - Test if strings end with specific suffixes
- `in` - Check membership in arrays
- `less` - Numeric less-than comparison
- `more` - Numeric greater-than comparison
- `matches` - Regular expression matching
- `type` - Type validation
- `and` - Logical AND operation
- `or` - Logical OR operation
- `not` - Logical NOT operation

## Basic Usage

```go
import "github.com/kaptinlin/jsonpatch"

doc := map[string]interface{}{
    "user": map[string]interface{}{
        "name":   "Alice",
        "email":  "alice@example.com",
        "age":    25,
        "active": true,
    },
    "tags": []interface{}{"admin", "user"},
}

patch := []jsonpatch.Operation{
    {"op": "defined", "path": "/user/name"},
    {"op": "type", "path": "/user/age", "value": "number"},
    {"op": "contains", "path": "/tags", "value": "admin"},
}

result, err := jsonpatch.ApplyPatch(doc, patch)
if err != nil {
    log.Fatal(err)
}
```

## Operations

### Test Operation

Test if a value equals the expected value.

```go
// Test equality
{"op": "test", "path": "/status", "value": "active"}

// Test inequality (inverted)
{"op": "test", "path": "/status", "value": "inactive", "not": true}
```

### Defined/Undefined Operations

Check if paths exist.

```go
// Check if path exists
{"op": "defined", "path": "/user/email"}

// Check if path doesn't exist
{"op": "undefined", "path": "/user/phone"}
```

### Type Operation

Check the type of a value.

```go
// Single type
{"op": "type", "path": "/user/age", "value": "number"}

// Multiple types
{"op": "type", "path": "/data", "value": ["string", "number"]}
```

Supported types: `"string"`, `"number"`, `"boolean"`, `"object"`, `"array"`, `"null"`, `"integer"`

### Contains Operation

Check if arrays contain values or strings contain substrings.

```go
// Array contains
{"op": "contains", "path": "/tags", "value": "admin"}

// String contains
{"op": "contains", "path": "/user/email", "value": "@example.com"}
```

### String Operations

Test string prefixes and suffixes.

```go
// Starts with
{"op": "starts", "path": "/user/email", "value": "alice"}

// Ends with
{"op": "ends", "path": "/user/email", "value": ".com"}

// Case-insensitive
{"op": "starts", "path": "/user/name", "value": "ALICE", "ignore_case": true}
```

### Numeric Comparisons

Compare numeric values.

```go
// Less than
{"op": "less", "path": "/user/age", "value": 30}

// Greater than
{"op": "more", "path": "/user/age", "value": 18}
```

### In Operation

Check if a value is in an array.

```go
{"op": "in", "path": "/user/role", "value": ["admin", "moderator", "user"]}
```

### Matches Operation

Regular expression matching. Requires custom matcher configuration.

```go
// Define custom regex matcher
customMatcher := func(pattern string, ignoreCase bool) jsonpatch.RegexMatcher {
    var flags string
    if ignoreCase {
        flags = "(?i)"
    }
    re := regexp.MustCompile(flags + pattern)
    return func(value string) bool {
        return re.MatchString(value)
    }
}

patch := []jsonpatch.Operation{
    {"op": "matches", "path": "/user/email", "value": `^[^@]+@[^@]+\.[^@]+$`},
}

result, err := jsonpatch.ApplyPatch(doc, patch, 
    jsonpatch.WithMatcher(customMatcher),
)
```

## Logical Operations

### And Operation

All conditions must be true.

```go
{
    "op": "and",
    "apply": [
        {"op": "defined", "path": "/user/email"},
        {"op": "type", "path": "/user/age", "value": "number"},
        {"op": "more", "path": "/user/age", "value": 18}
    ]
}
```

### Or Operation

At least one condition must be true.

```go
{
    "op": "or",
    "apply": [
        {"op": "contains", "path": "/tags", "value": "admin"},
        {"op": "contains", "path": "/tags", "value": "moderator"}
    ]
}
```

### Not Operation

Invert a condition.

```go
{
    "op": "not",
    "apply": [
        {"op": "contains", "path": "/tags", "value": "banned"}
    ]
}
```

## Common Patterns

### User Validation

```go
patch := []jsonpatch.Operation{
    {
        "op": "and",
        "apply": []jsonpatch.Operation{
            {"op": "defined", "path": "/user/name"},
            {"op": "defined", "path": "/user/email"},
            {"op": "type", "path": "/user/age", "value": "number"},
            {"op": "more", "path": "/user/age", "value": 17},
            {"op": "test", "path": "/user/active", "value": true},
        },
    },
}
```

### Permission Check

```go
patch := []jsonpatch.Operation{
    {
        "op": "or",
        "apply": []jsonpatch.Operation{
            {"op": "contains", "path": "/roles", "value": "admin"},
            {
                "op": "and",
                "apply": []jsonpatch.Operation{
                    {"op": "contains", "path": "/roles", "value": "user"},
                    {"op": "contains", "path": "/permissions", "value": "write"},
                },
            },
        },
    },
}
```

### Data Validation

```go
patch := []jsonpatch.Operation{
    {"op": "defined", "path": "/required_field"},
    {"op": "type", "path": "/required_field", "value": "string"},
    {
        "op": "not",
        "apply": []jsonpatch.Operation{
            {"op": "test", "path": "/required_field", "value": ""},
        },
    },
}
```

[json-predicate]: https://tools.ietf.org/id/draft-snell-json-test-01.html
