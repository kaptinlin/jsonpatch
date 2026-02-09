# CLAUDE.md

This file provides guidance to Claude Code when working with code in this repository.

## Project Overview

This is a comprehensive Go implementation of JSON Patch (RFC 6902), JSON Predicate, and extended operations for JSON document manipulation with full type safety and generic support. It's a Go port of json-joy/json-patch with 95%+ behavioral compatibility.

## API Usage

The library provides a clean struct-based API for JSON Patch operations:

```go
// Create operations using struct literals
patch := []jsonpatch.Operation{
    {Op: "add", Path: "/name", Value: "John"},
    {Op: "inc", Path: "/age", Inc: 1},
    {Op: "str_ins", Path: "/bio", Pos: 0, Str: "Hello "},
}

// Apply with type-safe generic API
result, err := jsonpatch.ApplyPatch(doc, patch)
```

## Development Commands

### Essential Commands
```bash
# Run all tests with race detection
make test

# Run benchmarks
make bench

# Run linter (golangci-lint v2.7.2 required)
make lint

# Format code
make fmt

# Run go vet
make vet

# Full verification (deps, format, vet, lint, test)
make verify

# Clean build artifacts and caches
make clean

# Test with coverage report
make test-coverage

# Run tests with verbose output
make test-verbose
```

### Testing Individual Operations
```bash
# Test specific operation implementations
go test -race ./op -run TestAddOp
go test -race ./op -run TestRemoveOp
go test -race ./op -run TestReplaceOp

# Test codec implementations
go test -race ./codec/json/...
go test -race ./codec/compact/...
go test -race ./codec/binary/...

# Run benchmarks for specific packages
go test -bench=. -benchmem ./codec/json
```

## Architecture & Key Components

### Core Package Structure
- **`/`** - Main package with generic API (`ApplyPatch`, `ApplyOp`, `ApplyOps`)
- **`/op`** - All operation implementations (29 operations total)
- **`/internal`** - Shared interfaces and types (`Op`, `PredicateOp`, `SecondOrderPredicateOp`, `Codec`)
- **`/codec`** - Encoding/decoding formats:
  - `/codec/json` - Standard JSON Patch format (RFC 6902)
  - `/codec/compact` - Compact array format for size optimization
  - `/codec/binary` - Binary format using MessagePack

### Operation Interface Hierarchy

All operations implement the `internal.Op` interface:
```go
type Op interface {
    Op() OpType
    Code() int
    Path() []string
    Apply(doc any) (OpResult[any], error)
    ToJSON() (Operation, error)
    ToCompact() (CompactOperation, error)
    Validate() error
}
```

Predicate operations also implement `internal.PredicateOp`:
```go
type PredicateOp interface {
    Op
    Test(doc any) (bool, error)
    Not() bool
}
```

Second-order predicates implement `internal.SecondOrderPredicateOp`:
```go
type SecondOrderPredicateOp interface {
    PredicateOp
    Ops() []PredicateOp
}
```

### Type-Safe Generic API
The main API uses Go generics to maintain type safety:
```go
func ApplyPatch[T Document](doc T, patch []Operation, opts ...Option) (*PatchResult[T], error)
```
This ensures the result maintains the same type as the input document.

### Operation struct
The `Operation` struct is used as the input format for all patch operations:
```go
type Operation struct {
    Op    string `json:"op"`
    Path  string `json:"path"`
    Value any    `json:"value,omitempty"`
    From  string `json:"from,omitempty"`

    // Extended operation fields
    Inc float64 `json:"inc"`
    Pos int     `json:"pos"`
    Str string  `json:"str"`
    Len int     `json:"len"`

    // Predicate fields
    Not        bool        `json:"not,omitempty"`
    Type       any         `json:"type,omitempty"`
    IgnoreCase bool        `json:"ignore_case,omitempty"`
    Apply      []Operation `json:"apply,omitempty"`

    // Special fields
    Props      map[string]any `json:"props,omitempty"`
    DeleteNull bool           `json:"deleteNull,omitempty"`
    OldValue   any            `json:"oldValue,omitempty"`
}
```

## Development Guidelines

### Core Principles
- All comments and documentation must be in English
- Follow Go conventions and idioms (Google Go Style Guide)
- Maintain json-joy behavioral compatibility
- Prioritize correctness over performance

### Testing Standards
- Table-driven tests with clear test cases
- Always include error cases
- Use `testify/assert` for assertions
- Benchmark critical operations with `testing.B.Loop()`

#### Error Testing Best Practices
- **NEVER** compare error message content with `assert.Contains(t, err.Error(), "message")`
- **USE** type-safe error checking with `assert.ErrorIs(t, err, ErrSpecificType)`
- **PREFER** sentinel errors defined in `op/errors.go` for consistent error types
- **PATTERN**:
  ```go
  // Good: Type-safe error checking
  assert.Error(t, err)
  assert.ErrorIs(t, err, ErrPathNotFound)

  // Bad: Fragile message checking
  assert.Contains(t, err.Error(), "NOT_FOUND")
  ```
- **AVAILABLE** assertions:
  - `assert.ErrorIs(t, err, ErrType)` - Check specific error type
  - `assert.ErrorAs(t, err, &targetType)` - Check error implements interface
  - `assert.Error(t, err)` - Just verify error occurred

### Error Handling
- Use json-joy compatible error messages
- Static errors: Return predefined error constants
- Dynamic errors: Use `fmt.Errorf("%w: context", baseError, ...)`
- Error checking: Use `errors.Is()` for type-safe checks

## Key Dependencies
- `github.com/kaptinlin/jsonpointer` - JSON Pointer path handling
- `github.com/kaptinlin/deepclone` - Deep cloning for immutable operations
- `github.com/go-json-experiment/json` - JSON encoding/decoding (v2)
- `github.com/tinylib/msgp` - MessagePack for binary codec
- `github.com/stretchr/testify` - Testing assertions

## Operation Categories

### Standard JSON Patch (RFC 6902)
- `add`, `remove`, `replace`, `move`, `copy`, `test`

### JSON Predicate Operations
- Type checks: `defined`, `undefined`, `type`, `test_type`
- String operations: `starts`, `ends`, `contains`, `matches`
- String tests: `test_string`, `test_string_len`
- Comparisons: `less`, `more`, `in`

### Extended Operations
- `flip` - Toggle boolean values
- `inc` - Increment numeric values
- `str_ins`/`str_del` - String insertion/deletion
- `split` - Split values at position
- `merge` - Merge adjacent array elements
- `extend` - Extend objects with properties

### Second-Order Predicates
- `and`, `or`, `not` - Logical combinations of predicates

## Important Implementation Notes

### json-joy Compatibility
1. **Negation Pattern**: Only `test`, `test_string`, `test_string_len` support direct `not` field
2. **Second-Order Predicates**: Use `{op: "not", apply: [...]}` for negating other predicates
3. **Empty Path Format**: Second-order predicates use `path: ""` with absolute paths in `apply` operations
4. **Type Coercion**: Follows JavaScript Number() semantics for numeric operations

### Core Features
1. **Path Handling**: All paths use JSON Pointer format with proper escaping
2. **Immutability**: Operations are immutable by default; use `WithMutate(true)` for in-place modification
3. **Error Propagation**: Operations return detailed errors with path context
4. **Type Conversion**: The library handles conversion between document types automatically
5. **Validation**: All operations validate their parameters before execution
