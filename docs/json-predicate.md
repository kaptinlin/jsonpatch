# JSON Predicate Operations

This document covers all [JSON Predicate][json-predicate] operations implemented in this library:

- `test` - Test equality with optional negation
- `contains` - Check if arrays contain values or strings contain substrings
- `defined` - Test if paths exist and are not undefined
- `undefined` - Test if paths don't exist or are undefined
- `ends` - Test if strings end with specific suffixes
- `starts` - Test if strings start with specific prefixes
- `in` - Check membership in arrays
- `less` - Numeric less-than comparison
- `more` - Numeric greater-than comparison
- `matches` - Regular expression matching
- `type` - Type validation
- `and` - Logical AND operation
- `or` - Logical OR operation
- `not` - Logical NOT operation

Only a subset of types supported by the `type` operation are implemented.

By default, `ValidateOperation` does not allow the `matches` operation, as it uses regular expressions and can be exploited using ReDoS attacks. You can allow it through options.

## Basic Usage

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "regexp"
    
    "github.com/kaptinlin/jsonpatch"
)

func main() {
    doc := map[string]interface{}{
        "user": map[string]interface{}{
            "name":   "Alice",
            "email":  "alice@example.com",
            "age":    25,
            "active": true,
        },
        "tags": []interface{}{"admin", "user"},
    }
    
    // Use predicate operations for conditional checks
    patch := []jsonpatch.Operation{
        {
            "op":   "defined",
            "path": "/user/name",
        },
        {
            "op":    "type",
            "path":  "/user/age",
            "value": "number",
        },
        {
            "op":    "contains",
            "path":  "/tags",
            "value": "admin",
        },
    }
    
    // Create regex matcher for matches operation
    options := jsonpatch.ApplyPatchOptions{
        JsonPatchOptions: jsonpatch.JsonPatchOptions{
            CreateMatcher: func(pattern string, ignoreCase bool) jsonpatch.RegexMatcher {
                var flags string
                if ignoreCase {
                    flags = "(?i)"
                }
                re := regexp.MustCompile(flags + pattern)
                return func(value string) bool {
                    return re.MatchString(value)
                }
            },
        },
    }
    
    result, err := jsonpatch.ApplyPatch(doc, patch, options)
    if err != nil {
        log.Fatalf("Predicate operation failed: %v", err)
    }
    
    output, _ := json.MarshalIndent(result.Doc, "", "  ")
    fmt.Println(string(output))
}
```

## First-Order Predicate Operations

### Test Operation

Test if a value equals the expected value.

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
        "not":   true,
    },
}
```

### Defined Operation

Check if a path exists and is not undefined.

```go
patch := []jsonpatch.Operation{
    {
        "op":   "defined",
        "path": "/user/email",
    },
}
```

### Undefined Operation

Check if a path doesn't exist or is undefined.

```go
patch := []jsonpatch.Operation{
    {
        "op":   "undefined",
        "path": "/user/phone",
    },
}
```

### Type Operation

Check the type of a value.

```go
patch := []jsonpatch.Operation{
    {
        "op":    "type",
        "path":  "/user/age",
        "value": "number",
    },
}

// Supported types: "string", "number", "boolean", "object", "array", "null", "integer"
patch = []jsonpatch.Operation{
    {
        "op":    "type",
        "path":  "/user/preferences",
        "value": "object",
    },
}

// Multiple types can be specified
patch = []jsonpatch.Operation{
    {
        "op":    "type",
        "path":  "/data",
        "value": []string{"string", "number"},
    },
}
```

### Contains Operation

Check if arrays contain specified values, or if strings contain substrings.

```go
// Array contains check
patch := []jsonpatch.Operation{
    {
        "op":    "contains",
        "path":  "/tags",
        "value": "admin",
    },
}

// String contains check
patch = []jsonpatch.Operation{
    {
        "op":    "contains",
        "path":  "/user/email",
        "value": "@example.com",
    },
}

// Case-insensitive string check
patch = []jsonpatch.Operation{
    {
        "op":           "contains",
        "path":         "/user/name",
        "value":        "ALICE",
        "ignore_case":  true,
    },
}
```

### Starts Operation

Check if strings start with specified prefixes.

```go
patch := []jsonpatch.Operation{
    {
        "op":    "starts",
        "path":  "/user/email",
        "value": "alice",
    },
}

// Case-insensitive check
patch = []jsonpatch.Operation{
    {
        "op":          "starts",
        "path":        "/user/name",
        "value":       "al",
        "ignore_case": true,
    },
}
```

### Ends Operation

Check if strings end with specified suffixes.

```go
patch := []jsonpatch.Operation{
    {
        "op":    "ends",
        "path":  "/user/email",
        "value": "@example.com",
    },
}

// Case-insensitive check
patch = []jsonpatch.Operation{
    {
        "op":          "ends",
        "path":        "/filename",
        "value":       ".TXT",
        "ignore_case": true,
    },
}
```

### In Operation

Check if a value is in an array.

```go
patch := []jsonpatch.Operation{
    {
        "op":    "in",
        "path":  "/user/role",
        "value": []interface{}{"admin", "moderator", "user"},
    },
}
```

### Less Operation

Check if a numeric value is less than the specified value.

```go
patch := []jsonpatch.Operation{
    {
        "op":    "less",
        "path":  "/user/age",
        "value": 30,
    },
}
```

### More Operation

Check if a numeric value is greater than the specified value.

```go
patch := []jsonpatch.Operation{
    {
        "op":    "more",
        "path":  "/user/score",
        "value": 100,
    },
}
```

### Matches Operation

Check if strings match regular expressions. **Note**: This operation must be explicitly allowed.

```go
func matchesExample() {
    doc := map[string]interface{}{
        "email": "user@example.com",
        "phone": "+1-555-123-4567",
    }
    
    patch := []jsonpatch.Operation{
        {
            "op":    "matches",
            "path":  "/email",
            "value": `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`,
        },
        {
            "op":          "matches",
            "path":        "/phone",
            "value":       `^\+\d{1,3}-\d{3}-\d{3}-\d{4}$`,
            "ignore_case": false,
        },
    }
    
    // Custom regex matcher
    options := jsonpatch.ApplyPatchOptions{
        JsonPatchOptions: jsonpatch.JsonPatchOptions{
            CreateMatcher: func(pattern string, ignoreCase bool) jsonpatch.RegexMatcher {
                flags := ""
                if ignoreCase {
                    flags = "(?i)"
                }
                re, _ := regexp.Compile(flags + pattern)
                return func(value string) bool {
                    return re.MatchString(value)
                }
            },
        },
    }
    
    result, err := jsonpatch.ApplyPatch(doc, patch, options)
    if err != nil {
        log.Printf("Matches operation failed: %v", err)
        return
    }
    
    fmt.Printf("Validation passed: %+v\n", result.Doc)
}
```

## Second-Order Predicate Operations

### And Operation

Logical AND - all conditions must be true.

```go
patch := []jsonpatch.Operation{
    {
        "op":   "and",
        "path": "/user",
        "apply": []jsonpatch.Operation{
            {
                "op":    "defined",
                "path":  "/user/name",
            },
            {
                "op":    "type",
                "path":  "/user/age",
                "value": "number",
            },
            {
                "op":    "more",
                "path":  "/user/age",
                "value": 18,
            },
        },
    },
}
```

### Or Operation

Logical OR - at least one condition must be true.

```go
patch := []jsonpatch.Operation{
    {
        "op":   "or",
        "path": "/contact",
        "apply": []jsonpatch.Operation{
            {
                "op":   "defined",
                "path": "/contact/email",
            },
            {
                "op":   "defined",
                "path": "/contact/phone",
            },
        },
    },
}
```

### Not Operation

Logical NOT - condition must be false.

```go
patch := []jsonpatch.Operation{
    {
        "op":   "not",
        "path": "/user",
        "apply": []jsonpatch.Operation{
            {
                "op":    "test",
                "path":  "/user/status",
                "value": "banned",
            },
        },
    },
}
```

## Advanced Usage Examples

### Complex Validation

```go
func complexValidation() {
    doc := map[string]interface{}{
        "user": map[string]interface{}{
            "name":   "Alice",
            "email":  "alice@example.com",
            "age":    25,
            "roles":  []interface{}{"user", "admin"},
            "status": "active",
        },
    }
    
    // Complex validation with nested conditions
    patch := []jsonpatch.Operation{
        {
            "op":   "and",
            "path": "/user",
            "apply": []jsonpatch.Operation{
                // Must have name and email
                {
                    "op":   "defined",
                    "path": "/user/name",
                },
                {
                    "op":   "defined",
                    "path": "/user/email",
                },
                // Age must be valid
                {
                    "op":   "and",
                    "path": "/user",
                    "apply": []jsonpatch.Operation{
                        {
                            "op":    "type",
                            "path":  "/user/age",
                            "value": "number",
                        },
                        {
                            "op":    "more",
                            "path":  "/user/age",
                            "value": 13,
                        },
                        {
                            "op":    "less",
                            "path":  "/user/age",
                            "value": 120,
                        },
                    },
                },
                // Must be active or pending
                {
                    "op":    "in",
                    "path":  "/user/status",
                    "value": []interface{}{"active", "pending"},
                },
                // Must have admin role
                {
                    "op":    "contains",
                    "path":  "/user/roles",
                    "value": "admin",
                },
            },
        },
    }
    
    options := jsonpatch.ApplyPatchOptions{Mutate: false}
    result, err := jsonpatch.ApplyPatch(doc, patch, options)
    if err != nil {
        log.Printf("Validation failed: %v", err)
        return
    }
    
    fmt.Printf("Complex validation passed: %+v\n", result.Doc)
}
```

### Conditional Processing

```go
func conditionalProcessing() {
    doc := map[string]interface{}{
        "order": map[string]interface{}{
            "status":   "pending",
            "total":    150.00,
            "customer": "premium",
        },
    }
    
    // Process only if conditions are met
    patch := []jsonpatch.Operation{
        // Check if order is eligible for processing
        {
            "op":   "and",
            "path": "/order",
            "apply": []jsonpatch.Operation{
                {
                    "op":    "test",
                    "path":  "/order/status",
                    "value": "pending",
                },
                {
                    "op":    "more",
                    "path":  "/order/total",
                    "value": 100,
                },
                {
                    "op":    "in",
                    "path":  "/order/customer",
                    "value": []interface{}{"premium", "gold"},
                },
            },
        },
        // If validation passes, update the order
        {
            "op":    "replace",
            "path":  "/order/status",
            "value": "processing",
        },
    }
    
    options := jsonpatch.ApplyPatchOptions{Mutate: false}
    result, err := jsonpatch.ApplyPatch(doc, patch, options)
    if err != nil {
        log.Printf("Conditional processing failed: %v", err)
        return
    }
    
    fmt.Printf("Order processed: %+v\n", result.Doc)
}
```

### Data Integrity Checks

```go
func dataIntegrityChecks() {
    doc := map[string]interface{}{
        "product": map[string]interface{}{
            "name":        "Laptop",
            "price":       999.99,
            "category":    "electronics",
            "description": "High-performance laptop",
            "tags":        []interface{}{"computer", "portable"},
        },
    }
    
    // Comprehensive data integrity validation
    patch := []jsonpatch.Operation{
        {
            "op":   "and",
            "path": "/product",
            "apply": []jsonpatch.Operation{
                // Required fields must exist
                {
                    "op":   "defined",
                    "path": "/product/name",
                },
                {
                    "op":   "defined",
                    "path": "/product/price",
                },
                {
                    "op":   "defined",
                    "path": "/product/category",
                },
                // Name must be non-empty string
                {
                    "op":   "and",
                    "path": "/product",
                    "apply": []jsonpatch.Operation{
                        {
                            "op":    "type",
                            "path":  "/product/name",
                            "value": "string",
                        },
                        {
                            "op":   "not",
                            "path": "/product",
                            "apply": []jsonpatch.Operation{
                                {
                                    "op":    "test",
                                    "path":  "/product/name",
                                    "value": "",
                                },
                            },
                        },
                    },
                },
                // Price must be positive number
                {
                    "op":   "and",
                    "path": "/product",
                    "apply": []jsonpatch.Operation{
                        {
                            "op":    "type",
                            "path":  "/product/price",
                            "value": "number",
                        },
                        {
                            "op":    "more",
                            "path":  "/product/price",
                            "value": 0,
                        },
                    },
                },
                // Category must be valid
                {
                    "op":    "in",
                    "path":  "/product/category",
                    "value": []interface{}{"electronics", "clothing", "books", "home"},
                },
            },
        },
    }
    
    options := jsonpatch.ApplyPatchOptions{Mutate: false}
    result, err := jsonpatch.ApplyPatch(doc, patch, options)
    if err != nil {
        log.Printf("Data integrity check failed: %v", err)
        return
    }
    
    fmt.Printf("Data integrity verified: %+v\n", result.Doc)
}
```

## Error Handling

```go
func predicateErrorHandling() {
    doc := map[string]interface{}{
        "user": map[string]interface{}{
            "name": "Alice",
            "age":  "twenty-five", // Invalid type
        },
    }
    
    patch := []jsonpatch.Operation{
        {
            "op":    "type",
            "path":  "/user/age",
            "value": "number",
        },
    }
    
    options := jsonpatch.ApplyPatchOptions{Mutate: false}
    result, err := jsonpatch.ApplyPatch(doc, patch, options)
    if err != nil {
        // Predicate operations return specific error types
        switch {
        case strings.Contains(err.Error(), "test operation failed"):
            log.Printf("Predicate test failed: %v", err)
        case strings.Contains(err.Error(), "path not found"):
            log.Printf("Invalid path: %v", err)
        default:
            log.Printf("Predicate operation error: %v", err)
        }
        return
    }
    
    fmt.Printf("Predicate passed: %+v\n", result.Doc)
}
```

[json-predicate]: https://tools.ietf.org/id/draft-snell-json-test-01.html 
