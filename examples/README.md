# JSON Patch Examples

This directory contains practical examples demonstrating different JSON Patch usage patterns and document types.

## Core Operation Examples

### 1. Basic Operations (`basic-operations/`)
Demonstrates fundamental operations: `add`, `replace`, `remove`, and `test`.

```bash
cd basic-operations && go run main.go
```

### 2. Array Operations (`array-operations/`)
Shows how to work with arrays: adding elements, removing items, and updating array contents.

```bash
cd array-operations && go run main.go
```

### 3. Conditional Operations (`conditional-operations/`)
Illustrates safe updates using `test` operations for validation and optimistic locking.

```bash
cd conditional-operations && go run main.go
```

### 4. Copy and Move Operations (`copy-move-operations/`)
Shows data restructuring and migration using `copy` and `move` operations.

```bash
cd copy-move-operations && go run main.go
```

### 5. String Operations (`string-operations/`)
Illustrates text editing capabilities with string insertion operations.

```bash
cd string-operations && go run main.go
```

## Document Type Examples

### 6. Struct Patch (`struct-patch/`)
Demonstrates patching Go structs with type safety and JSON tag support.

```bash
cd struct-patch && go run main.go
```

## Codec Examples

### 11. Compact Codec (`compact-codec/`)
Demonstrates the compact array-based codec that provides significant space savings over standard JSON format.

```bash
cd compact-codec && go run main.go
```

Shows encoding with both numeric and string opcodes, space savings analysis, and round-trip compatibility testing.

### 7. Map Patch (`map-patch/`)
Shows efficient patching of `map[string]any` documents for dynamic data.

```bash
cd map-patch && go run main.go
```

### 8. JSON Bytes Patch (`json-bytes-patch/`)
Illustrates patching JSON byte data for file processing and API integration.

```bash
cd json-bytes-patch && go run main.go
```

### 9. JSON String Patch (`json-string-patch/`)
Demonstrates patching JSON string data for configuration and API responses.

```bash
cd json-string-patch && go run main.go
```

## Advanced Examples

### 10. Batch Update (`batch-update/`)
Demonstrates efficient batch operations for updating multiple items at once.

```bash
cd batch-update && go run main.go
```

### 11. Error Handling (`error-handling/`)
Demonstrates error handling patterns, validation techniques, and recovery strategies.

```bash
cd error-handling && go run main.go
```

### 12. Mutate Option (`mutate-option/`)
Shows the Mutate option for performance optimization with in-place modifications.

```bash
cd mutate-option && go run main.go
```

## Example Categories

### 🏗️ **Type-Safe Examples**
- `struct-patch/` - Go struct patching with JSON tags
- `map-patch/` - Dynamic map document manipulation

### 📄 **Data Format Examples**  
- `json-bytes-patch/` - Raw JSON byte processing
- `json-string-patch/` - JSON string manipulation

### ⚙️ **Operation Examples**
- `basic-operations/` - Core RFC 6902 operations
- `array-operations/` - Array-specific operations
- `conditional-operations/` - Safe conditional updates
- `copy-move-operations/` - Data restructuring
- `string-operations/` - Text editing operations

### 🚀 **Advanced Examples**
- `batch-update/` - Bulk operations
- `error-handling/` - Error management
- `mutate-option/` - Performance optimization

## Quick Start

Each example is self-contained and can be run independently:

```bash
# Run any example
cd <example-directory>
go run main.go

# For example:
cd struct-patch
go run main.go
```
