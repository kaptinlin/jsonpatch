# CLAUDE.md

This file provides guidance to Claude Code when working with code in this repository.

## Project Overview

**jsonpatch** is a comprehensive Go implementation of JSON Patch (RFC 6902), JSON Predicate, and extended operations for JSON document manipulation. It's a Go port of [json-joy/json-patch](https://github.com/streamich/json-joy/tree/master/src/json-patch) with 95%+ behavioral compatibility, providing full type safety through Go generics.

**Module:** `github.com/kaptinlin/jsonpatch`
**Go Version:** 1.26

## Design Philosophy

### Core Principles
- **Type Safety First**: Generic API eliminates `interface{}` casting and provides compile-time type safety
- **json-joy Compatibility**: Maintain 95%+ behavioral compatibility with the TypeScript reference implementation
- **Immutability by Default**: Operations create deep copies unless `WithMutate(true)` is explicitly used
- **Correctness Over Performance**: Prioritize correct behavior over micro-optimizations
- **Library Purity**: No logging dependencies; return errors and let callers handle logging

### Engineering Standards
- Follow KISS, DRY, YAGNI, and SOLID principles
- Use Go 1.26 modern features: generics, `testing.B.Loop()`, `slices`/`maps` packages
- English-only comments and documentation
- Table-driven tests with `t.Parallel()` where safe

## Commands

```bash
# Run all tests with race detection
task test

# Run benchmarks
task bench

# Run linter (golangci-lint v2.9.0)
task lint

# Format code
task fmt

# Run go vet
task vet

# Full verification (deps, format, vet, lint, test)
task verify

# Clean build artifacts and caches
task clean

# Test with coverage report
task test-coverage

# Run tests with verbose output
task test-verbose
```

### Testing Specific Components
```bash
# Test specific operation implementations
go test -race ./op -run TestAddOp
go test -race ./op -run TestRemoveOp

# Test codec implementations
go test -race ./codec/json/...
go test -race ./codec/compact/...
go test -race ./codec/binary/...

# Benchmark specific packages
go test -bench=. -benchmem ./codec/json
```

## Architecture

### Package Structure
```
jsonpatch/
├── index.go              # Main API: ApplyPatch, ApplyOp, ApplyOps
├── jsonpatch.go          # Generic entry points
├── validate.go           # Operation validation
├── op/                   # 29 operation implementations
│   ├── base.go           # BaseOp shared functionality
│   ├── add.go            # RFC 6902 operations
│   ├── remove.go
│   ├── replace.go
│   ├── move.go
│   ├── copy.go
│   ├── test.go
│   ├── flip.go           # Extended operations
│   ├── inc.go
│   ├── str_ins.go
│   ├── str_del.go
│   ├── split.go
│   ├── merge.go
│   ├── extend.go
│   ├── contains.go       # Predicate operations
│   ├── defined.go
│   ├── undefined.go
│   ├── starts.go
│   ├── ends.go
│   ├── in.go
│   ├── less.go
│   ├── more.go
│   ├── matches.go
│   ├── type.go
│   ├── test_type.go
│   ├── test_string.go
│   ├── test_string_len.go
│   ├── and.go            # Second-order predicates
│   ├── or.go
│   ├── not.go
│   └── errors.go         # Sentinel errors
├── internal/             # Shared interfaces and types
│   ├── interfaces.go     # Op, PredicateOp, SecondOrderPredicateOp, Codec
│   ├── types.go          # Operation, OpResult, PatchResult, Document
│   ├── options.go        # JSONPatchOptions and functional options
│   ├── constants.go      # Operation type and code constants
│   ├── classify.go       # Operation classification functions
│   └── jsonpatch_type.go # JSONPatchType constants and helpers
└── codec/                # Encoding/decoding formats
    ├── json/             # Standard RFC 6902 JSON format
    ├── compact/          # Array-based format (~35% space savings)
    └── binary/           # MessagePack binary format
```

### Operation Interface Hierarchy

All operations implement `internal.Op`:
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

The result type matches the input document type automatically.

## Key Types and Interfaces

### Operation Struct
Input format for all patch operations:
```go
type Operation struct {
    Op    string `json:"op"`
    Path  string `json:"path"`
    Value any    `json:"value,omitempty"`
    From  string `json:"from,omitempty"`

    // Extended operation fields (no omitempty - 0 is valid)
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

### Document Constraint
```go
type Document interface {
    ~[]byte | ~string | map[string]any | any
}
```

### Result Types
```go
type OpResult[T Document] struct {
    Doc T   // Result document with preserved type
    Old any // Previous value at the path
}

type PatchResult[T Document] struct {
    Doc T             // The patched document
    Res []OpResult[T] // Results for each operation
}
```

## Coding Rules

### json-joy Compatibility Requirements

**1. Negation Pattern**
- Only `test`, `test_string`, `test_string_len` support direct `not` field
- All other predicates use second-order `not` operation:
```go
// Direct negation (only for test, test_string, test_string_len)
{Op: "test", Path: "/value", Value: 42, Not: true}

// Second-order negation (for all other predicates)
{
    Op:   "not",
    Path: "",
    Apply: []jsonpatch.Operation{
        {Op: "starts", Path: "/name", Value: "John"},
    },
}
```

**2. Second-Order Predicate Path Format**
- Use `path: ""` with absolute paths in `apply` operations
- Logical combinations: `and`, `or`, `not`

**3. Type Coercion**
- Follow JavaScript `Number()` semantics for numeric operations
- String-to-number conversion for `inc`, `less`, `more` operations

### Error Handling

**Type-Safe Error Checking (REQUIRED):**
```go
// GOOD: Use errors.Is() with sentinel errors
if err != nil {
    if errors.Is(err, ErrPathNotFound) {
        // Handle specific error
    }
    return err
}

// BAD: Never check error message strings
if strings.Contains(err.Error(), "NOT_FOUND") { ... }
```

**Error Patterns:**
- Static errors: Return predefined sentinel errors from `op/errors.go`
- Dynamic errors: Use `fmt.Errorf("%w: context", baseError, ...)`
- All sentinel errors are defined in `op/errors.go`

### Path Handling
- All paths use JSON Pointer format (RFC 6901)
- Proper escaping: `~0` for `~`, `~1` for `/`
- Array index `-` means append
- Empty path `[]string{}` means root document

### Immutability
- Operations are immutable by default (deep copy via `deepclone.Clone()`)
- Use `WithMutate(true)` option for in-place modification when performance is critical
- Always clone user-provided values before storing

## Testing

### Test Structure
- Use table-driven tests with clear test case names
- Always include error cases
- Use `t.Parallel()` in top-level tests and subtests where safe
- Use `testing.B.Loop()` for benchmarks (Go 1.24+)

### Error Testing Pattern
```go
func TestOperation_Error(t *testing.T) {
    t.Parallel()

    tests := []struct {
        name    string
        doc     any
        op      Operation
        wantErr error
    }{
        {
            name:    "path not found",
            doc:     map[string]any{"a": 1},
            op:      Operation{Op: "remove", Path: "/b"},
            wantErr: ErrPathNotFound,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()

            _, err := ApplyPatch(tt.doc, []Operation{tt.op})
            if err == nil {
                t.Fatal("expected error")
            }
            if !errors.Is(err, tt.wantErr) {
                t.Errorf("got %v, want %v", err, tt.wantErr)
            }
        })
    }
}
```

### Comparison Tools
- Use `github.com/google/go-cmp/cmp` for deep comparisons
- Standard library `testing` package (no third-party assertion libraries)

## Dependencies

### Core Dependencies
- `github.com/kaptinlin/jsonpointer` - JSON Pointer path handling (RFC 6901)
- `github.com/kaptinlin/deepclone` - Deep cloning for immutable operations
- `github.com/go-json-experiment/json` - JSON encoding/decoding (v2 experimental)
- `github.com/tinylib/msgp` - MessagePack for binary codec

### Test Dependencies
- `github.com/google/go-cmp` - Deep comparison in tests
- `github.com/stretchr/testify` - Test assertions

## Operation Categories

### Standard JSON Patch (RFC 6902)
`add`, `remove`, `replace`, `move`, `copy`, `test`

### JSON Predicate Operations
- Type checks: `defined`, `undefined`, `type`, `test_type`
- String operations: `starts`, `ends`, `contains`, `matches`
- String tests: `test_string`, `test_string_len`
- Comparisons: `less`, `more`, `in`

### Extended Operations
- `flip` - Toggle boolean values
- `inc` - Increment/decrement numeric values
- `str_ins`/`str_del` - String insertion/deletion
- `split` - Split values at position
- `merge` - Merge adjacent array elements
- `extend` - Extend objects with properties

### Second-Order Predicates
`and`, `or`, `not` - Logical combinations of predicates

## Performance

### Optimization Guidelines
- Pre-allocate slices when size is known: `make([]T, 0, capacity)`
- Use `strings.Builder` for string concatenation
- Minimize reflection usage
- Profile before optimizing (no premature optimization)

### Mutation Control
```go
// Immutable (default) - safe for concurrent use
result, err := jsonpatch.ApplyPatch(doc, patch)

// Mutable - faster but modifies original
result, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.WithMutate(true))
```

## Agent Skills

Package-local skills available in `.agents/skills/`:
- **agent-md-creating** - Generate CLAUDE.md for Go projects
- **code-simplifying** - Refine code for clarity and consistency
- **committing** - Create conventional commits
- **dependency-selecting** - Select Go dependencies from kaptinlin ecosystem
- **go-best-practices** - Google Go coding best practices
- **linting** - Set up and run golangci-lint v2
- **modernizing** - Go 1.20-1.26 modernization guide
- **ralphy-initializing** - Initialize Ralphy AI coding loop
- **ralphy-todo-creating** - Create Ralphy TODO.yaml task files
- **readme-creating** - Generate README.md for Go libraries
- **releasing** - Guide release process for Go packages
- **testing** - Write Go tests with best practices

Use the `Skill` tool to invoke these when relevant to your task.
