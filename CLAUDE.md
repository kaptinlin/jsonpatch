# CLAUDE.md

This file provides guidance to Claude Code when working with code in this repository.

## Project Overview

This is a comprehensive Go implementation of JSON Patch (RFC 6902), JSON Predicate, and extended operations for JSON document manipulation with full type safety and generic support. It's a Go port of json-joy/json-patch.

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

# Run linter (golangci-lint v2.4.0 required)
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
- **`/op`** - All operation implementations (add, remove, replace, move, copy, test, and extended ops)
- **`/internal`** - Shared interfaces and types (`Op`, `PredicateOp`, `Codec`)
- **`/codec`** - Different encoding/decoding formats:
  - `/codec/json` - Standard JSON Patch format (RFC 6902)
  - `/codec/compact` - Compact array format for size optimization
  - `/codec/binary` - Binary format using MessagePack

### Operation Interface Hierarchy
All operations implement the `internal.Op` interface with methods:
- `Op()` - Returns operation type
- `Path()` - Returns JSON Pointer path
- `Apply()` - Applies operation to document
- `Validate()` - Validates operation parameters

Predicate operations also implement `internal.PredicateOp` with:
- `Test()` - Tests condition on document
- `Not()` - Whether it's a negation predicate

### Type-Safe Generic API
The main API uses Go generics to maintain type safety:
```go
func ApplyPatch[T any](doc T, operations []Operation, opts ...Option) (*Result[T], error)
```
This ensures the result maintains the same type as the input document.

## Development Guidelines (from .cursor/rules.mdc)

### Core Principles
- All comments and documentation must be in English
- Follow Go conventions and idioms
- Prioritize correctness over performance
- Measure before optimize

### Performance Optimization Approach
- Conservative optimization first - use proven safe patterns
- Single change iteration - one optimization at a time
- Immediate rollback for >2% performance regression
- Focus on:
  - String operations with `strings.Builder`
  - Boolean logic simplification
  - Type specialization for simple types
  - Inline simple operations

### Testing Standards
- Table-driven tests with clear test cases
- Always include error cases
- Use `testify/assert` for assertions
- Benchmark critical operations

### Error Handling
- Define base errors as constants in errors.go files
- Static errors: Return predefined error constants
- Dynamic errors: Use `fmt.Errorf("%w: context", baseError, ...)`
- Error checking: Use `errors.Is()` for type-safe checks

## Key Dependencies
- `github.com/kaptinlin/jsonpointer` - JSON Pointer path handling
- `github.com/kaptinlin/deepclone` - Deep cloning for immutable operations
- `github.com/go-json-experiment/json` - JSON encoding/decoding
- `github.com/wapc/tinygo-msgpack` - MessagePack for binary codec
- `github.com/stretchr/testify` - Testing assertions

## Operation Categories

### Standard JSON Patch (RFC 6902)
- `add`, `remove`, `replace`, `move`, `copy`, `test`

### JSON Predicate Operations
- Type checks: `defined`, `undefined`, `type`
- String operations: `starts`, `ends`, `contains`, `matches`
- Comparisons: `less`, `more`, `in`
- String length: `strlen`

### Extended Operations
- `flip` - Toggle boolean values
- `inc` - Increment numeric values
- `strins`/`strdel` - String insertion/deletion
- `split` - Split strings into arrays
- `merge` - Deep merge objects
- `extend` - Extend arrays

### Second-Order Predicates
- `and`, `or`, `not` - Logical combinations of predicates

## Important Implementation Notes

1. **Path Handling**: All paths use JSON Pointer format with proper escaping
2. **Immutability**: Operations can be immutable (default) or mutate in-place (WithMutate option)
3. **Error Propagation**: Operations return detailed errors with path context
4. **Type Conversion**: The library handles conversion between document types automatically
5. **Validation**: All operations validate their parameters before execution
