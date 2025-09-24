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
    {Op: "defined", Path: "/user/name"},
    {Op: "type", Path: "/user/age", Value: "number"},
    {Op: "contains", Path: "/tags", Value: "admin"},
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
{Op: "test", Path: "/status", Value: "active"}

// Test inequality (inverted)
{Op: "test", Path: "/status", Value: "inactive", Not: true}
```

### Defined/Undefined Operations

Check if paths exist.

```go
// Check if path exists
{Op: "defined", Path: "/user/email"}

// Check if path doesn't exist
{Op: "undefined", Path: "/user/phone"}
```

### Type Operation

Check the type of a value.

```go
// Single type
{Op: "type", Path: "/user/age", Value: "number"}

// Multiple types
{Op: "type", Path: "/data", Value: ["string", "number"]}
```

Supported types: `"string"`, `"number"`, `"boolean"`, `"object"`, `"array"`, `"null"`, `"integer"`

### Contains Operation

Check if arrays contain values or strings contain substrings.

```go
// Array contains
{Op: "contains", Path: "/tags", Value: "admin"}

// String contains
{Op: "contains", Path: "/user/email", Value: "@example.com"}
```

### String Operations

Test string prefixes and suffixes.

```go
// Starts with
{Op: "starts", Path: "/user/email", Value: "alice"}

// Ends with
{Op: "ends", Path: "/user/email", Value: ".com"}

// Case-insensitive
{Op: "starts", Path: "/user/name", Value: "ALICE", IgnoreCase: true}
```

### Numeric Comparisons

Compare numeric values.

```go
// Less than
{Op: "less", Path: "/user/age", Value: 30}

// Greater than
{Op: "more", Path: "/user/age", Value: 18}
```

### In Operation

Check if a value is in an array.

```go
{Op: "in", Path: "/user/role", Value: ["admin", "moderator", "user"]}
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
    {Op: "matches", Path: "/user/email", Value: `^[^@]+@[^@]+\.[^@]+$`},
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
    Op: "and",
    Apply: [
        {Op: "defined", Path: "/user/email"},
        {Op: "type", Path: "/user/age", Value: "number"},
        {Op: "more", Path: "/user/age", Value: 18}
    ]
}
```

### Or Operation

At least one condition must be true.

```go
{
    Op: "or",
    Apply: [
        {Op: "contains", Path: "/tags", Value: "admin"},
        {Op: "contains", Path: "/tags", Value: "moderator"}
    ]
}
```

### Not Operation

Invert a condition.

```go
{
    Op: "not",
    Apply: [
        {Op: "contains", Path: "/tags", Value: "banned"}
    ]
}
```

## Common Patterns

### User Validation

```go
patch := []jsonpatch.Operation{
    {
        Op: "and",
        Apply: []jsonpatch.Operation{
            {Op: "defined", Path: "/user/name"},
            {Op: "defined", Path: "/user/email"},
            {Op: "type", Path: "/user/age", Value: "number"},
            {Op: "more", Path: "/user/age", Value: 17},
            {Op: "test", Path: "/user/active", Value: true},
        },
    },
}
```

### Permission Check

```go
patch := []jsonpatch.Operation{
    {
        Op: "or",
        Apply: []jsonpatch.Operation{
            {Op: "contains", Path: "/roles", Value: "admin"},
            {
                Op: "and",
                Apply: []jsonpatch.Operation{
                    {Op: "contains", Path: "/roles", Value: "user"},
                    {Op: "contains", Path: "/permissions", Value: "write"},
                },
            },
        },
    },
}
```

### Data Validation

```go
patch := []jsonpatch.Operation{
    {Op: "defined", Path: "/required_field"},
    {Op: "type", Path: "/required_field", Value: "string"},
    {
        Op: "not",
        Apply: []jsonpatch.Operation{
            {Op: "test", Path: "/required_field", Value: ""},
        },
    },
}
```

[json-predicate]: https://tools.ietf.org/id/draft-snell-json-test-01.html
