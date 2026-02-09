# JSON Predicate Operations

This document covers [JSON Predicate][json-predicate] operations for conditional testing and validation:

- `test` - Test equality with optional negation
- `contains` - Check if arrays contain values or strings contain substrings
- `defined` - Test if paths exist
- `undefined` - Test if paths don't exist
- `starts` - Test if strings start with specific prefixes
- `ends` - Test if strings end with specific suffixes
- `in` - Check membership in arrays
- `less` - Numeric less-than comparison
- `more` - Numeric greater-than comparison
- `matches` - Regular expression matching
- `type` - Type validation (single type)
- `test_type` - Type validation (single or multiple types)
- `test_string` - Position-based string testing
- `test_string_len` - String length validation
- `and` - Logical AND operation
- `or` - Logical OR operation
- `not` - Logical NOT operation

## Basic Usage

```go
import "github.com/kaptinlin/jsonpatch"

doc := map[string]any{
    "user": map[string]any{
        "name":   "Alice",
        "email":  "alice@example.com",
        "age":    25,
        "active": true,
    },
    "tags": []any{"admin", "user"},
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

### Defined / Undefined Operations

Check if paths exist.

```go
// Check if path exists
{Op: "defined", Path: "/user/email"}

// Check if path doesn't exist
{Op: "undefined", Path: "/user/phone"}
```

### Type Operation

Check the type of a value (single type string).

```go
{Op: "type", Path: "/user/age", Value: "number"}
```

Supported types: `"string"`, `"number"`, `"boolean"`, `"object"`, `"array"`, `"null"`, `"integer"`

### Test Type Operation

Check the type of a value with support for multiple types.

```go
// Single type
{Op: "test_type", Path: "/user/age", Type: "number"}

// Multiple types
{Op: "test_type", Path: "/data", Type: []any{"string", "number"}}
```

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
{Op: "in", Path: "/user/role", Value: []any{"admin", "moderator", "user"}}
```

### Test String Operation

Test a substring at a specific position in a string. Supports `Not` flag.

```go
// Test substring at position
{Op: "test_string", Path: "/text", Pos: 0, Str: "Hello"}

// Test negation
{Op: "test_string", Path: "/text", Pos: 0, Str: "Goodbye", Not: true}
```

### Test String Length Operation

Validate the length of a string. Supports `Not` flag.

```go
// Test exact length
{Op: "test_string_len", Path: "/name", Len: 5}

// Test that length is not equal
{Op: "test_string_len", Path: "/name", Len: 10, Not: true}
```

### Matches Operation

Regular expression matching. Requires custom matcher configuration.

```go
import "regexp"

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
    Op:   "and",
    Path: "",
    Apply: []jsonpatch.Operation{
        {Op: "defined", Path: "/user/email"},
        {Op: "type", Path: "/user/age", Value: "number"},
        {Op: "more", Path: "/user/age", Value: 18},
    },
}
```

### Or Operation

At least one condition must be true.

```go
{
    Op:   "or",
    Path: "",
    Apply: []jsonpatch.Operation{
        {Op: "contains", Path: "/tags", Value: "admin"},
        {Op: "contains", Path: "/tags", Value: "moderator"},
    },
}
```

### Not Operation

Invert a condition.

```go
{
    Op:   "not",
    Path: "",
    Apply: []jsonpatch.Operation{
        {Op: "contains", Path: "/tags", Value: "banned"},
    },
}
```

## Common Patterns

### User Validation

```go
patch := []jsonpatch.Operation{
    {
        Op:   "and",
        Path: "",
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
        Op:   "or",
        Path: "",
        Apply: []jsonpatch.Operation{
            {Op: "contains", Path: "/roles", Value: "admin"},
            {
                Op:   "and",
                Path: "",
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
        Op:   "not",
        Path: "",
        Apply: []jsonpatch.Operation{
            {Op: "test", Path: "/required_field", Value: ""},
        },
    },
}
```

[json-predicate]: https://tools.ietf.org/id/draft-snell-json-test-01.html
